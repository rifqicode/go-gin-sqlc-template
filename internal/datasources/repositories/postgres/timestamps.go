package postgres

import "time"

// TimestampFields represents common timestamp fields for database records
type TimestampFields struct {
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
