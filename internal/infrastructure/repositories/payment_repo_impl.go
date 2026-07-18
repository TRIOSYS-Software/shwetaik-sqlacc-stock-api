package repositories

import (
	"fmt"
	"time"

	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"

	"gorm.io/gorm"
)

// newRecordSeed is the negative primary-key value SQL Account uses to mark a
// row as a draft/pending document not yet renumbered by the desktop client.
const newRecordSeed = -1

// PaymentRepositoryImpl writes payments directly into the live GL_CB /
// GL_CBDTL / GL_TRANS tables, bypassing the vendor REST API entirely. This
// mirrors shwetaik-sql-acc-backend-api's PaymentRepo.Create, which SQL
// Account's desktop client already relies on to post payments — the
// negative-key allocation and GL_TRANS double-entry logic here must stay
// behaviorally identical to that reference, since any divergence risks
// corrupting the accounting ledger.
type PaymentRepositoryImpl struct {
	db                *gorm.DB
	paymentMethodRepo repositories.PaymentMethodRepository
}

func NewPaymentRepository(db *gorm.DB, paymentMethodRepo repositories.PaymentMethodRepository) repositories.PaymentRepository {
	return &PaymentRepositoryImpl{db: db, paymentMethodRepo: paymentMethodRepo}
}

func (r *PaymentRepositoryImpl) Create(payment *entities.Payment) (err error) {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
			err = fmt.Errorf("payment create: recovered from panic: %v", rec)
		}
	}()

	var firstPayment entities.Payment
	if err := tx.First(&firstPayment).Error; err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	var firstDetail entities.PaymentDetail
	if err := tx.First(&firstDetail).Error; err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	// Matches the reference implementation: the payment-method lookup runs
	// on the plain connection, not the transaction.
	paymentMethod, err := r.paymentMethodRepo.GetPaymentMethodByCode(payment.PaymentMethod)
	if err != nil {
		tx.Rollback()
		return err
	}
	currencyCode := ""
	if paymentMethod.CurrencyCode != nil {
		currencyCode = *paymentMethod.CurrencyCode
	}
	journal := ""
	if paymentMethod.Journal != nil {
		journal = *paymentMethod.Journal
	}

	var currencyRate float64
	if err := tx.Raw("SELECT BUYINGRATE FROM CURRENCY WHERE CODE = ?", currencyCode).Scan(&currencyRate).Error; err != nil {
		tx.Rollback()
		return err
	}

	if firstPayment.DocKey > 0 {
		payment.DocKey = newRecordSeed
	} else {
		payment.DocKey = firstDetail.DtlKey - 1
	}

	payment.CurrencyRate = currencyRate
	payment.Journal = journal
	payment.CurrencyCode = currencyCode

	if firstPayment.GLTransID > 0 {
		payment.GLTransID = newRecordSeed
	} else {
		payment.GLTransID = firstPayment.GLTransID - 1
	}

	payment.LastModified = uint(time.Now().Unix())

	total := 0.0
	for i := range payment.Details {
		payment.Details[i].DtlKey = payment.DocKey - (i + 1)
		payment.Details[i].Seq = uint(i+1) * 1000
		total += payment.Details[i].Amount
		payment.Details[i].DocKey = payment.DocKey
		payment.Details[i].CurrencyCode = currencyCode
		payment.Details[i].CurrencyRate = currencyRate

		if err := tx.Create(&payment.Details[i]).Error; err != nil {
			tx.Rollback()
			return err
		}

		detailDescription := ""
		if payment.Details[i].Description != nil {
			detailDescription = *payment.Details[i].Description
		}
		if err := r.createGLTrans(tx, payment.Details[i].Code, detailDescription, payment, payment.Details[i].Amount, 0, payment.Details[i].CurrencyAmount, 0, "S", payment.Details[i].DtlKey); err != nil {
			tx.Rollback()
			return err
		}
	}
	payment.DocAmt = total
	payment.LocalDocAmt = total

	paymentDescription := ""
	if payment.Description != nil {
		paymentDescription = *payment.Description
	}
	if err := r.createGLTrans(tx, payment.PaymentMethod, paymentDescription, payment, 0, payment.DocAmt, 0, payment.LocalDocAmt, "M", payment.DocKey); err != nil {
		tx.Rollback()
		return err
	}

	payment.Details = nil
	if err := tx.Save(payment).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *PaymentRepositoryImpl) createGLTrans(
	tx *gorm.DB,
	code string,
	description string,
	payment *entities.Payment,
	dr float64,
	cr float64,
	localDR float64,
	localCR float64,
	tableType string,
	fromKey int,
) error {
	var firstGLTrans entities.GLTrans
	if err := tx.First(&firstGLTrans).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	glTrans := entities.GLTrans{
		GLTransID:    int64(payment.GLTransID),
		Code:         code,
		Area:         payment.Area,
		Agent:        payment.Agent,
		Project:      payment.Project,
		Journal:      payment.Journal,
		CurrencyCode: payment.CurrencyCode,
		CurrencyRate: payment.CurrencyRate,
		Description:  description,
		DR:           dr,
		CR:           cr,
		LocalDR:      localDR,
		LocalCR:      localCR,
		Ref1:         payment.DocNo,
		FromDocType:  payment.DocType,
		FromKey:      fromKey,
		TableType:    tableType,
		Cancelled:    payment.Cancelled,
	}

	if firstGLTrans.DocKey > 0 {
		glTrans.DocKey = newRecordSeed
	} else {
		glTrans.DocKey = firstGLTrans.DocKey - 1
	}

	return tx.Create(&glTrans).Error
}
