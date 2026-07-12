package monitor

import (
	"log"
	"time"

	"gorm.io/gorm"

	"shwetaik-sqlacc-stock-api/internal/infrastructure/webhook"
)

type stockItemChange struct {
	Code         string `gorm:"column:CODE"`
	LastModified int64  `gorm:"column:LASTMODIFIED"`
}

type stockItemCode struct {
	Code string `gorm:"column:CODE"`
}

// StartStockItemChangeMonitor polls ST_ITEM at the given interval and fires
// webhooks for updates and deletions.
//
// Startup is non-blocking: only MAX(LASTMODIFIED) is fetched up front (a
// fast single-row query), so the server isn't held up. The full set of
// known codes needed for delete detection is loaded lazily on the first
// poll tick instead.
//
// Update detection polls WHERE LASTMODIFIED > watermark for exact changed
// codes. Delete detection tracks known codes in memory and only pays for a
// full code scan when COUNT(*) drops, then diffs against the known set to
// find exactly which codes vanished.
func StartStockItemChangeMonitor(db *gorm.DB, webhookClient *webhook.Client, interval time.Duration) {
	go func() {
		var watermark int64
		if err := db.Raw("SELECT COALESCE(MAX(LASTMODIFIED), 0) FROM ST_ITEM").Scan(&watermark).Error; err != nil {
			log.Printf("stock item monitor: failed to establish baseline: %v", err)
			return
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		var knownCodes map[string]struct{}
		var lastCount int
		codesLoaded := false

		for range ticker.C {
			if !codesLoaded {
				var ok bool
				knownCodes, lastCount, ok = loadKnownCodes(db)
				if ok {
					codesLoaded = true
					log.Printf("stock item monitor: loaded %d known codes", len(knownCodes))
				}
			}

			var changedCodes []string
			changedCodes, watermark = checkUpdatedStockItems(db, watermark)
			if len(changedCodes) > 0 {
				for _, code := range changedCodes {
					knownCodes[code] = struct{}{}
				}
				go webhookClient.Send("item.updated", changedCodes)
			}

			var deletedCodes []string
			deletedCodes, knownCodes, lastCount = checkDeletedStockItems(db, knownCodes, lastCount)
			if len(deletedCodes) > 0 {
				go webhookClient.Send("item.deleted", deletedCodes)
			}
		}
	}()
}

// loadKnownCodes returns ok=false on any query failure so the caller can
// retry on a later tick instead of permanently treating the load as done
// with an empty/incomplete code set.
func loadKnownCodes(db *gorm.DB) (codes map[string]struct{}, count int, ok bool) {
	var total int64
	if err := db.Raw("SELECT COUNT(*) FROM ST_ITEM").Scan(&total).Error; err != nil {
		log.Printf("stock item monitor: failed to count items: %v", err)
		return nil, 0, false
	}

	var rows []stockItemCode
	if err := db.Raw("SELECT CODE FROM ST_ITEM").Scan(&rows).Error; err != nil {
		log.Printf("stock item monitor: failed to load known codes: %v", err)
		return nil, 0, false
	}

	known := make(map[string]struct{}, len(rows))
	for _, r := range rows {
		known[r.Code] = struct{}{}
	}
	return known, int(total), true
}

func checkUpdatedStockItems(db *gorm.DB, watermark int64) ([]string, int64) {
	var changes []stockItemChange
	err := db.Raw(`
		SELECT CODE, LASTMODIFIED
		FROM ST_ITEM
		WHERE LASTMODIFIED > ?
		ORDER BY LASTMODIFIED DESC
	`, watermark).Scan(&changes).Error
	if err != nil {
		log.Printf("stock item monitor: updated-poll query failed: %v", err)
		return nil, watermark
	}
	if len(changes) == 0 {
		return nil, watermark
	}

	codes := make([]string, len(changes))
	for i, change := range changes {
		codes[i] = change.Code
	}
	// changes is ordered DESC by LASTMODIFIED, so the first row carries the
	// new high watermark.
	newWatermark := changes[0].LastModified
	log.Printf("stock item monitor: %d item(s) updated: %v", len(codes), codes)
	return codes, newWatermark
}

func checkDeletedStockItems(db *gorm.DB, knownCodes map[string]struct{}, lastCount int) ([]string, map[string]struct{}, int) {
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM ST_ITEM").Scan(&count).Error; err != nil {
		log.Printf("stock item monitor: delete-check count query failed: %v", err)
		return nil, knownCodes, lastCount
	}
	newCount := int(count)

	if newCount >= lastCount {
		return nil, knownCodes, newCount
	}

	// Count dropped — pay for a full scan to find exactly which codes
	// vanished, rather than doing this on every poll tick.
	log.Printf("stock item monitor: item count dropped %d -> %d, scanning for deleted codes", lastCount, newCount)
	var codes []stockItemCode
	if err := db.Raw("SELECT CODE FROM ST_ITEM").Scan(&codes).Error; err != nil {
		log.Printf("stock item monitor: delete scan failed: %v", err)
		return nil, knownCodes, lastCount
	}

	currentCodes := make(map[string]struct{}, len(codes))
	for _, c := range codes {
		currentCodes[c.Code] = struct{}{}
	}

	var deleted []string
	for code := range knownCodes {
		if _, ok := currentCodes[code]; !ok {
			deleted = append(deleted, code)
		}
	}

	if len(deleted) > 0 {
		log.Printf("stock item monitor: %d item(s) deleted: %v", len(deleted), deleted)
	}
	return deleted, currentCodes, newCount
}
