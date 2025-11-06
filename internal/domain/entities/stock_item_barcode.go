package entities

type STItemBarcode struct {
	AutoKey int    `gorm:"column:AUTOKEY" json:"autokey"`
	Code    string `gorm:"column:CODE" json:"code"`
	Barcode string `gorm:"column:BARCODE" json:"barcode"`
	UOM     string `gorm:"column:UOM" json:"uom"`
}

func (STItemBarcode) TableName() string {
	return "ST_ITEM_BARCODE"
}
