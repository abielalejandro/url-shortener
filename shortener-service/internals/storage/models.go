package storage

import "time"

type Url struct {
	Short       string
	Long        string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	LastVisited time.Time
}
