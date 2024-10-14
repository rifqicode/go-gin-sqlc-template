package shared

import "github.com/jackc/pgx/v5/pgxpool"

type SharedMeta struct {
	DB *pgxpool.Pool
}
