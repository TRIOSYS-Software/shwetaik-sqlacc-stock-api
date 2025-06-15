package entities

import (
	"time"
)

type STItem struct {
	DocKey             int           `json:"dockey" gorm:"column:DOCKEY;primaryKey;not null"`
	Code               string        `json:"code" gorm:"column:CODE;size:30;not null"`
	Description        *string       `json:"description,omitempty" gorm:"column:DESCRIPTION;size:200"`
	Description2       *string       `json:"description2,omitempty" gorm:"column:DESCRIPTION2;size:200"`
	Description3       []byte        `json:"description3,omitempty" gorm:"column:DESCRIPTION3"`
	StockGroup         string        `json:"stock_group" gorm:"column:STOCKGROUP;size:20;not null"`
	StockControl       *bool         `json:"stock_control,omitempty" gorm:"column:STOCKCONTROL"`
	CostingMethod      int16         `json:"costing_method" gorm:"column:COSTINGMETHOD;not null"`
	SerialNumber       *bool         `json:"serial_number,omitempty" gorm:"column:SERIALNUMBER"`
	Remark1            *string       `json:"remark1,omitempty" gorm:"column:REMARK1;size:200"`
	Remark2            *string       `json:"remark2,omitempty" gorm:"column:REMARK2;size:200"`
	MinQty             *float64      `json:"min_qty,omitempty" gorm:"column:MINQTY;precision:18;scale:4"`
	MaxQty             *float64      `json:"max_qty,omitempty" gorm:"column:MAXQTY;precision:18;scale:4"`
	ReorderLevel       *float64      `json:"reorder_level,omitempty" gorm:"column:REORDERLEVEL;precision:18;scale:4"`
	ReorderQty         *float64      `json:"reorder_qty,omitempty" gorm:"column:REORDERQTY;precision:18;scale:4"`
	Shelf              *string       `json:"shelf,omitempty" gorm:"column:SHELF;size:40"`
	SUOM               *string       `json:"suom,omitempty" gorm:"column:SUOM;size:10"`
	ItemType           string        `json:"item_type" gorm:"column:ITEMTYPE;size:1;not null"`
	LeadTime           *int          `json:"lead_time,omitempty" gorm:"column:LEADTIME"`
	BOMLeadTime        *int          `json:"bom_lead_time,omitempty" gorm:"column:BOM_LEADTIME"`
	BOMASMCost         *float64      `json:"bom_asm_cost,omitempty" gorm:"column:BOM_ASMCOST;precision:18;scale:8"`
	SLTax              *string       `json:"sl_tax,omitempty" gorm:"column:SLTAX;size:10"`
	PHTax              *string       `json:"ph_tax,omitempty" gorm:"column:PHTAX;size:10"`
	Tariff             *string       `json:"tariff,omitempty" gorm:"column:TARIFF;size:20"`
	IRBMClassification *string       `json:"irbm_classification,omitempty" gorm:"column:IRBM_CLASSIFICATION;size:3"`
	StockMatrix        *string       `json:"stock_matrix,omitempty" gorm:"column:STOCKMATRIX;size:15"`
	DefUOMST           *string       `json:"def_uom_st,omitempty" gorm:"column:DEFUOM_ST;size:10"`
	DefUOMSL           *string       `json:"def_uom_sl,omitempty" gorm:"column:DEFUOM_SL;size:10"`
	DefUOMPH           *string       `json:"def_uom_ph,omitempty" gorm:"column:DEFUOM_PH;size:10"`
	ScriptCode         *string       `json:"script_code,omitempty" gorm:"column:SCRIPTCODE;size:15"`
	IsActive           *bool         `json:"is_active,omitempty" gorm:"column:ISACTIVE"`
	BalSQty            *float64      `json:"bals_qty,omitempty" gorm:"column:BALSQTY;precision:18;scale:4"`
	BalSUOMQty         *float64      `json:"bals_uom_qty,omitempty" gorm:"column:BALSUOMQTY;precision:18;scale:4"`
	CreationDate       *time.Time    `json:"creation_date,omitempty" gorm:"column:CREATIONDATE;type:date"`
	Picture            []byte        `json:"picture,omitempty" gorm:"column:PICTURE"`
	PictureClass       *string       `json:"picture_class,omitempty" gorm:"column:PICTURECLASS;size:10"`
	Attachments        []byte        `json:"attachments,omitempty" gorm:"column:ATTACHMENTS"`
	Note               []byte        `json:"note,omitempty" gorm:"column:NOTE"`
	LastModified       *int64        `json:"last_modified,omitempty" gorm:"column:LASTMODIFIED"`
	STItemPrices       []STItemPrice `json:"st_item_prices,omitempty" gorm:"foreignKey:Code;references:Code"`
}

func (STItem) TableName() string {
	return "ST_ITEM"
}
