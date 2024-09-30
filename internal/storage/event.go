package storage

import "time"

type Event struct {
	Id          uint64 `gorm:"primaryKey;autoIncrement:true;"`
	OpCode      uint64
	Data        string
	CreatedAt   time.Time
	ProcessedAt time.Time
}
