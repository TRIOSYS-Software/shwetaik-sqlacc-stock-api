package entities

import "time"

// GLTrans mirrors GL_TRANS, the general ledger journal. Rows are only ever
// written by PaymentRepository.Create as the double-entry counterpart of a
// GL_CB/GL_CBDTL payment posting — nothing reads this table back through the
// API, so it carries gorm tags only, no json tags.
type GLTrans struct {
	DocKey       int        `gorm:"column:DOCKEY;primaryKey"`
	GLTransID    int64      `gorm:"column:GLTRANSID;primaryKey;autoIncrement"`
	Code         string     `gorm:"column:CODE"`
	DocDate      time.Time  `gorm:"column:DOCDATE;autoCreateTime"`
	PostDate     time.Time  `gorm:"column:POSTDATE;autoUpdateTime"`
	TaxDate      time.Time  `gorm:"column:TAXDATE;autoCreateTime"`
	Area         string     `gorm:"column:AREA;default:----"`
	Agent        string     `gorm:"column:AGENT;default:----"`
	Project      string     `gorm:"column:PROJECT;default:----"`
	Tax          string     `gorm:"column:TAX"`
	Journal      string     `gorm:"column:JOURNAL"`
	CurrencyCode string     `gorm:"column:CURRENCYCODE"`
	CurrencyRate float64    `gorm:"column:CURRENCYRATE"`
	Description  string     `gorm:"column:DESCRIPTION"`
	Description2 string     `gorm:"column:DESCRIPTION2"`
	DR           float64    `gorm:"column:DR"`
	CR           float64    `gorm:"column:CR"`
	LocalDR      float64    `gorm:"column:LOCALDR"`
	LocalCR      float64    `gorm:"column:LOCALCR"`
	Ref1         string     `gorm:"column:REF1"`
	Ref2         string     `gorm:"column:REF2"`
	FromDocType  string     `gorm:"column:FROMDOCTYPE"`
	FromKey      int        `gorm:"column:FROMKEY"`
	TableType    string     `gorm:"column:TABLETYPE"`
	ReconDate    *time.Time `gorm:"column:RECONDATE"`
	Cancelled    bool       `gorm:"column:CANCELLED"`
	AutoPost     bool       `gorm:"column:AUTOPOST"`
	Nonce        string     `gorm:"column:NONCE;default:0000"`
	Digest       string     `gorm:"column:DIGEST"`
}

func (GLTrans) TableName() string {
	return "GL_TRANS"
}
