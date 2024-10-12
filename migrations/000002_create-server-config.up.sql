CREATE TABLE server_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    config_name VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    value TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Add indexes for improved query performance
CREATE INDEX idx_server_configs_status ON server_configs(status);
CREATE INDEX idx_server_configs_created_at ON server_configs(created_at);
CREATE INDEX idx_server_configs_updated_at ON server_configs(updated_at);
CREATE INDEX idx_server_configs_deleted_at ON server_configs(deleted_at);
