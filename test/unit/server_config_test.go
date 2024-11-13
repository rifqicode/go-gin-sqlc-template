package test

import (
	"context"
	"net/http"
	"testing"

	. "github.com/franela/goblin"
	"github.com/onsi/gomega"
	"gitlab.gift.id/lv2/loyalty/internal/db"
	"gitlab.gift.id/lv2/loyalty/internal/routes/server_config"
	"gitlab.gift.id/lv2/loyalty/test"
)

func TestServerConfig(t *testing.T) {
	g := Goblin(t)
	var app *http.Server
	var database *db.Queries

	g.Describe("Server Config Test", func() {
		g.Before(func() {
			app, database = test.InitTest(g)
			database.TruncateForTesting(context.Background())
		})

		g.After(func() {
			database.TruncateForTesting(context.Background())
		})

		g.It("Create Config - Success ", func() {
			serverConfigParams := server_config.ServerConfigParams{
				ConfigName:  "Testing",
				Status:      "active",
				Value:       "TestingValue",
				Description: "Description",
			}

			res := test.Request(app.Handler, "POST", "/server-config", serverConfigParams)

			gomega.Expect(res.StatusCode).To(gomega.Equal(http.StatusCreated))
		})

		g.It("Create Config - Failed ", func() {
			serverConfigParams := server_config.ServerConfigParams{
				ConfigName: "Testing",
				Status:     "active",
				Value:      "TestingValue",
			}

			res := test.Request(app.Handler, "POST", "/server-config", serverConfigParams)

			gomega.Expect(res.StatusCode).To(gomega.Equal(http.StatusInternalServerError))
		})
	})
}
