//go:generate swag init -d .,.. -g main.go -o ../api
package main

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/laurentpoirierfr/ora-cdc-go/pkg/config"
	"github.com/laurentpoirierfr/ora-cdc-go/pkg/helper"
	"github.com/laurentpoirierfr/ora-cdc-go/pkg/logger"

	"github.com/laurentpoirierfr/ora-cdc-go/internal/models"
	"github.com/laurentpoirierfr/ora-cdc-go/internal/router"
)

// @title ms-app
// @version 0.1.0
// @description This is a bff server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host homezone.swagger.io:8080
// @BasePath /
func main() {
	log := logger.NewLogger()

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.LOGGER, log)

	conf, err := config.NewConfig("config.yaml")
	helper.DieIfError(err)

	log.Info("starting " + conf.GetPropertyString("application.name") + " ...")

	// ctx.Value(domain.LOGGER).(*zap.Logger).Info("starting " + conf.GetPropertyString("application.name") + " ...")

	////////////////////////////////////////////////////////////////////////////////////
	// start web app
	////////////////////////////////////////////////////////////////////////////////////
	app := gin.Default()
	router.NewRouter(ctx, conf, app)
	// Listen and Serve
	if err := app.Run(":" + conf.GetPropertyString("server.port")); err != nil {
		helper.DieIfError(err)
	}
	////////////////////////////////////////////////////////////////////////////////////
	log.Info("stopped.")
}
