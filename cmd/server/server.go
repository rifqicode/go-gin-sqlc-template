package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"gitlab.gift.id/lv2/loyalty/internal/db"
	"gitlab.gift.id/lv2/loyalty/internal/routes/ping"
	"gitlab.gift.id/lv2/loyalty/internal/routes/server_config"
	logger "gitlab.gift.id/lv2/loyalty/pkg/logger/zap"
)

// App represents the structure for the application server
type App struct {
	httpServer   *http.Server
	dbConnection *pgx.Conn
}

// NewApp initializes a new application server instance
func NewApp() (*App, error) {
	// Connect to the database
	conn, err := pgx.Connect(context.Background(), os.Getenv("db_url"))
	if err != nil {
		logger.Error("Unable to connect to database:", err)
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	queries := db.New(conn)

	// Initialize Gin engine with middleware
	router := gin.Default()

	// Set up logging middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Define routes
	router.GET("/ping", ping.GetHandler(queries))
	server_config.SetupRoutes(&router.RouterGroup, queries)

	// Create an instance of http.Server with simple configuration
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("port")), // Server port
		Handler:      router,                                // Attach the router
		ReadTimeout:  5 * time.Second,                       // Read timeout
		WriteTimeout: 10 * time.Second,                      // Write timeout
		IdleTimeout:  15 * time.Second,                      // Idle timeout
	}

	return &App{httpServer: httpServer, dbConnection: conn}, nil
}

// Run starts the application server
func (app *App) Run() error {
	logger.Info("Starting server on port", app.httpServer.Addr)

	// Start the HTTP server
	if err := app.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("could not listen on %s: %v", app.httpServer.Addr, err)
	}

	return nil
}

// Stop manages graceful shutdown of the application server
func (app *App) Stop(ctx context.Context) error {
	logger.Info("Shutting down server gracefully")
	defer func() {
		if err := app.dbConnection.Close(ctx); err != nil {
			logger.Error("Error closing database connection:", err)
		}
	}()

	return app.httpServer.Shutdown(ctx)
}
