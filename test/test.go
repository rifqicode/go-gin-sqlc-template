package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/onsi/gomega"
	"gitlab.gift.id/lv2/loyalty/internal/db"
	"gitlab.gift.id/lv2/loyalty/internal/routes/ping"
	"gitlab.gift.id/lv2/loyalty/internal/routes/server_config"
)

// Request creates a request for testing purposes with chi and sqlc setup
func Request(app http.Handler, method, endpoint string, optional ...interface{}) *http.Response {
	var requestBody *bytes.Buffer
	var req *http.Request
	var err error

	// If there's a request body, marshal it into JSON
	if len(optional) > 0 && optional[0] != nil {
		jsonBody, _ := json.Marshal(optional[0])
		requestBody = bytes.NewBuffer(jsonBody)
		req, err = http.NewRequest(method, endpoint, requestBody)
	} else {
		req, err = http.NewRequest(method, endpoint, nil)
	}
	if err != nil {
		log.Fatal("Failed to create request:", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	res := rec.Result()

	return res
}

func InitTest(g *goblin.G) (*http.Server, *db.Queries) {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Println("Unable to load ENVIRONMENT VARIABLES:", err.Error())
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("db_url"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	queries := db.New(conn)

	router := gin.Default()

	router.GET("/ping", ping.GetHandler(queries))

	server_config.SetupRoutes(&router.RouterGroup, queries)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("port")), // Port di mana server akan mendengarkan
		Handler:      router,                                // Attach the mux to handle routes
		ReadTimeout:  5 * time.Second,                       // Timeout untuk membaca request
		WriteTimeout: 10 * time.Second,                      // Timeout untuk menulis response
		IdleTimeout:  15 * time.Second,                      // Timeout untuk koneksi idle
	}

	gomega.RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	return httpServer, queries
}
