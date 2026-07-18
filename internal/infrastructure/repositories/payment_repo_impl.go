package repositories

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"shwetaik-sqlacc-stock-api/internal/domain/entities"
	"shwetaik-sqlacc-stock-api/internal/domain/repositories"

	"gorm.io/gorm"
)

// newRecordSeed is the negative primary-key value SQL Account uses to mark a
// row as a draft/pending document not yet renumbered by the desktop client.
const newRecordSeed = -1

// maxCreateAttempts bounds the retry loop in Create — see its comment for
// why a retry is needed at all.
const maxCreateAttempts = 3

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

// Create allocates negative keys from the current minimum row in each table
// (see createOnce), which races if two requests read the same "first" row
// concurrently — both then try to insert the same key. Rather than locking
// (unverified whether the Firebird GORM dialect honors FOR UPDATE), this
// retries on what looks like a primary/unique-key collision, re-reading the
// current minimum fresh each attempt.
func (r *PaymentRepositoryImpl) Create(payment *entities.Payment) error {
	// createOnce nils payment.Details out once the detail rows are inserted
	// (see its comment), so a retry needs the original slice restored —
	// otherwise a failure after that point would retry with zero detail
	// lines.
	originalDetails := payment.Details

	var lastErr error
	for attempt := 1; attempt <= maxCreateAttempts; attempt++ {
		payment.Details = originalDetails
		err := r.createOnce(payment)
		if err == nil {
			return nil
		}
		if !isKeyCollision(err) {
			return err
		}
		lastErr = err
		time.Sleep(20 * time.Millisecond)
	}
	return fmt.Errorf("payment create: gave up after %d attempts on key collision: %w", maxCreateAttempts, lastErr)
}

// isKeyCollision reports whether err looks like a primary/unique-key
// constraint violation — the retryable case in Create. gorm.ErrDuplicatedKey
// only fires if the dialect implements GORM's error-translator interface;
// flylink888/gorm-firebird is a small community driver that may not, so this
// also falls back to matching Firebird's own violation message text. That
// text match hasn't been verified against a live Firebird error, so treat it
// as best-effort — worth confirming against real driver output.
func isKeyCollision(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "violation of primary or unique key constraint")
}

func (r *PaymentRepositoryImpl) createOnce(payment *entities.Payment) (err error) {
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

	var firstGLTrans entities.GLTrans
	if err := tx.First(&firstGLTrans).Error; err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}
	// nextGLTransDocKey tracks the next negative GL_TRANS.DOCKEY to use.
	// Re-derived once here instead of re-querying inside createGLTrans for
	// every row: within this one transaction nothing else can insert
	// between our own calls, so the row we just inserted is always the new
	// minimum — decrementing locally is equivalent to re-querying it, without
	// the extra round trip per detail line.
	nextGLTransDocKey := firstGLTrans.DocKey - 1
	if firstGLTrans.DocKey > 0 {
		nextGLTransDocKey = newRecordSeed
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
		if err := r.createGLTrans(tx, &nextGLTransDocKey, payment.Details[i].Code, detailDescription, payment, payment.Details[i].Amount, 0, payment.Details[i].CurrencyAmount, 0, "S", payment.Details[i].DtlKey); err != nil {
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
	if err := r.createGLTrans(tx, &nextGLTransDocKey, payment.PaymentMethod, paymentDescription, payment, 0, payment.DocAmt, 0, payment.LocalDocAmt, "M", payment.DocKey); err != nil {
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
	nextDocKey *int,
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
	glTrans := entities.GLTrans{
		DocKey:       *nextDocKey,
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

	if err := tx.Create(&glTrans).Error; err != nil {
		return err
	}
	*nextDocKey--
	return nil
}
