package entities

import (
	"time"
)

type STItemPrice struct {
	DtlKey     int        `gorm:"column:DTLKEY;primaryKey;not null" json:"dtl_key"`
	Code       string     `gorm:"column:CODE;type:VARCHAR(30);not null" json:"code"`
	TagType    string     `gorm:"column:TAGTYPE;type:CHAR(1);not null" json:"tag_type"`
	Company    *string    `gorm:"column:COMPANY;type:VARCHAR(10)" json:"company"`
	Seq        int        `gorm:"column:SEQ;not null" json:"seq"`
	PriceTag   *string    `gorm:"column:PRICETAG;type:VARCHAR(10)" json:"price_tag"`
	UOM        string     `gorm:"column:UOM;type:VARCHAR(10);not null" json:"uom"`
	Qty        *float64   `gorm:"column:QTY;type:DECIMAL(18,4)" json:"qty"`
	StockValue *float64   `gorm:"column:STOCKVALUE;type:DECIMAL(18,8)" json:"stock_value"`
	Discount   *string    `gorm:"column:DISCOUNT;type:VARCHAR(20)" json:"discount"`
	DateFrom   *time.Time `gorm:"column:DATEFROM;type:DATE" json:"date_from"`
	DateTo     *time.Time `gorm:"column:DATETO;type:DATE" json:"date_to"`
	Note       *[]byte    `gorm:"column:NOTE;type:BLOB" json:"note"`
}

// TableName specifies the table name for GORM
func (STItemPrice) TableName() string {
	return "ST_ITEM_PRICE"
}

// func (s *STItemPrice) BeforeCreate(db *gorm.DB) error {
// 	var maxSeq int
// 	err := db.Model(s).Select("MAX(SEQ)").Where("CODE = ?", s.Code).Scan(&maxSeq).Error
// 	if err != nil {
// 		return err
// 	}
// 	s.Seq = maxSeq + 1000
// 	return nil
// }
