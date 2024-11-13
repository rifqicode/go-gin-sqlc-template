package server_config

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"gitlab.gift.id/lv2/loyalty/internal/db"
)

// Handler handles CRUD operations for server configurations
type Handler struct {
	queries *db.Queries
}

// ServerConfig represents the server configuration stored in the database
type ServerConfigParams struct {
	ConfigName  string `json:"config_name"`
	Status      string `json:"status"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

// NewHandler creates a new ServerConfigHandler
func NewHandler(query *db.Queries) *Handler {
	return &Handler{
		queries: query,
	}
}

// Create godoc
// @Summary Create a new server configuration
// @Schemes http https
// @Description Create a new server configuration entry
// @Tags server-config
// @Accept json
// @Produce json
// @Param config body postgres.ServerConfig true "Server Configuration"
// @Success 201 {object} postgres.ServerConfig
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /server-config [post]
func (h *Handler) Create(c *gin.Context) {
	var config ServerConfigParams
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdConfig, err := h.queries.CreateServerConfig(context.Background(), db.CreateServerConfigParams{
		ConfigName:  config.ConfigName,
		Status:      config.Status,
		Value:       config.Value,
		Description: pgtype.Text{String: config.Description, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create server config"})
		return
	}

	c.JSON(http.StatusCreated, createdConfig)
}

// Get godoc
// @Summary Get a server configuration by ID
// @Schemes http https
// @Description Retrieve a server configuration entry by its ID
// @Tags server-config
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Success 200 {object} postgres.ServerConfig
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /server-config/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	config, err := h.queries.GetServerConfig(context.Background(), pgtype.UUID{Bytes: parsedID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Server config not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve server config"})
		}
		return
	}

	c.JSON(http.StatusOK, config)
}

// Update godoc
// @Summary Update a server configuration
// @Schemes http https
// @Description Update an existing server configuration entry
// @Tags server-config
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Param config body postgres.ServerConfig true "Updated Server Configuration"
// @Success 200 {object} postgres.ServerConfig
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /server-config/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var config ServerConfigParams
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	updatedConfig, err := h.queries.UpdateServerConfig(context.Background(), db.UpdateServerConfigParams{
		ID:          pgtype.UUID{Bytes: parsedID, Valid: true},
		ConfigName:  config.ConfigName,
		Status:      config.Status,
		Value:       config.Value,
		Description: pgtype.Text{String: config.Description, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Server config not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update server config"})
			return
		}
	}

	c.JSON(http.StatusOK, updatedConfig)
}

// Delete godoc
// @Summary Delete a server configuration
// @Schemes http https
// @Description Delete a server configuration entry by its ID
// @Tags server-config
// @Accept json
// @Produce json
// @Param id path string true "Configuration ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /server-config/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	_, err = h.queries.DeleteServerConfig(context.Background(), pgtype.UUID{Bytes: parsedID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Server config not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete server config"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// SetupRoutes sets up the routes for server configuration
func SetupRoutes(router *gin.RouterGroup, query *db.Queries) {
	handler := NewHandler(query)

	group := router.Group("/server-config")
	{
		group.POST("", handler.Create)
		group.GET("/:id", handler.Get)
		group.PUT("/:id", handler.Update)
		group.DELETE("/:id", handler.Delete)
	}
}
