package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurentpoirierfr/ora-cdc-go/internal/controller"

	docs "github.com/laurentpoirierfr/ora-cdc-go/api"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/penglongli/gin-metrics/ginmetrics"

	"github.com/laurentpoirierfr/ora-cdc-go/pkg/config"
)

func NewRouter(context context.Context, conf config.Config, app *gin.Engine) {

	// ==========================================================================================
	// Monitoring Prometheus
	// ==========================================================================================
	m := ginmetrics.GetMonitor()
	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/ops/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin
	m.Use(app)
	// ==========================================================================================
	// Swagger
	// ==========================================================================================
	app.GET("/swagger/*any", func(context *gin.Context) {
		docs.SwaggerInfo.Host = context.Request.Host
		docs.SwaggerInfo.Description = conf.GetPropertyString("application.description")
		docs.SwaggerInfo.Title = conf.GetPropertyString("application.name")

		ginSwagger.WrapHandler(swaggerfiles.Handler)(context)
	})

	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	indexController := controller.NewIndexController(context, conf)

	// ==========================================================================================
	// API
	// ==========================================================================================
	api := app.Group("/api")
	{
		api.GET("/demo", indexController.Demo)
	}

	// ==========================================================================================
	// OPS
	// ==========================================================================================
	opsCtrl := controller.NewOpsController(context, conf)

	ops := app.Group("/ops")
	{
		ops.GET("/info", opsCtrl.Info)
		ops.GET("/liveness", opsCtrl.Liveness)
		ops.GET("/readiness", opsCtrl.Readiness)
	}
	// ==========================================================================================
}
