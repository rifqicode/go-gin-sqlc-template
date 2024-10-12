package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ServerConfig represents the server configuration stored in the database
type ServerConfig struct {
	ID          uuid.UUID `db:"id"`
	ConfigName  string    `db:"config_name"`
	Status      string    `db:"status"`
	Value       string    `db:"value"`
	Description string    `db:"description"`
	TimestampFields
}

// ServerConfigRepository handles database operations for server configurations
type ServerConfigRepository struct {
	db *pgxpool.Pool
}

// NewServerConfigRepository creates a new ServerConfigRepository
func NewServerConfigRepository(db *pgxpool.Pool) *ServerConfigRepository {
	return &ServerConfigRepository{db: db}
}

// Create inserts a new server configuration into the database
func (r *ServerConfigRepository) Create(ctx context.Context, config *ServerConfig) (*ServerConfig, error) {
	query := `
		INSERT INTO server_configs (config_name, status, value, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		config.ConfigName, config.Status, config.Value, config.Description).
		Scan(&config.ID, &config.CreatedAt, &config.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// GetByID retrieves a server configuration by its ID
func (r *ServerConfigRepository) GetByID(ctx context.Context, id uuid.UUID) (*ServerConfig, error) {
	var config ServerConfig
	query := `SELECT * FROM server_configs WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&config.ID, &config.ConfigName, &config.Status, &config.Value, &config.Description,
		&config.CreatedAt, &config.UpdatedAt, &config.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Update updates an existing server configuration
func (r *ServerConfigRepository) Update(ctx context.Context, config *ServerConfig) (*ServerConfig, error) {
	query := `
		UPDATE server_configs
		SET config_name = $2, status = $3, value = $4, description = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		config.ID, config.ConfigName, config.Status, config.Value, config.Description).
		Scan(&config.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Delete soft-deletes a server configuration
func (r *ServerConfigRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE server_configs SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// List retrieves all active server configurations
func (r *ServerConfigRepository) List(ctx context.Context) ([]*ServerConfig, error) {
	var configs []*ServerConfig
	query := `SELECT * FROM server_configs WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var config ServerConfig
		err := rows.Scan(
			&config.ID, &config.ConfigName, &config.Status, &config.Value, &config.Description,
			&config.CreatedAt, &config.UpdatedAt, &config.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &config)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}
