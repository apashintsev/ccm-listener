package storage

import "time"

type Transaction struct {
	LT          uint64 `gorm:"primaryKey;autoIncrement:false;"`
	ProcessedAt time.Time
}
