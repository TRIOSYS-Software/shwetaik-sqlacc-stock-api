package monitor

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type stockItemChange struct {
	Code         string  `gorm:"column:CODE"`
	Description  *string `gorm:"column:DESCRIPTION"`
	StockGroup   string  `gorm:"column:STOCKGROUP"`
	LastModified int64   `gorm:"column:LASTMODIFIED"`
}

// StartStockItemChangeMonitor polls ST_ITEM.LASTMODIFIED at the given
// interval and logs any rows changed since the last poll. The baseline is
// the highest LASTMODIFIED value present at startup, so only changes made
// while the monitor is running are reported.
func StartStockItemChangeMonitor(db *gorm.DB, interval time.Duration) {
	go func() {
		var watermark int64
		if err := db.Raw("SELECT COALESCE(MAX(LASTMODIFIED), 0) FROM ST_ITEM").Scan(&watermark).Error; err != nil {
			log.Printf("stock item monitor: failed to establish baseline: %v", err)
			return
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			var changes []stockItemChange
			err := db.Raw(`
				SELECT CODE, DESCRIPTION, STOCKGROUP, LASTMODIFIED
				FROM ST_ITEM
				WHERE LASTMODIFIED > ?
				ORDER BY LASTMODIFIED
			`, watermark).Scan(&changes).Error
			if err != nil {
				log.Printf("stock item monitor: query failed: %v", err)
				continue
			}

			for _, change := range changes {
				description := ""
				if change.Description != nil {
					description = *change.Description
				}
				log.Printf("stock item changed: code=%s description=%q stock_group=%s last_modified=%d",
					change.Code, description, change.StockGroup, change.LastModified)
				if change.LastModified > watermark {
					watermark = change.LastModified
				}
			}
		}
	}()
}
