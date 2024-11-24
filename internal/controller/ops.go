package controller

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/laurentpoirierfr/ora-cdc-go/internal/models"
	"github.com/laurentpoirierfr/ora-cdc-go/pkg/config"
)

// IndexController is the default controller
type OpsController struct {
	context context.Context
	conf    config.Config
}

func NewOpsController(context context.Context, conf config.Config) *OpsController {
	return &OpsController{
		conf:    conf,
		context: context,
	}
}

// Info micro service
// @Summary Info
// @Schemes
// @Description Informations sur le service
// @Tags ops
// @Accept json
// @Produce json
// @Success 200 {object} models.Info
// @Router /ops/info [get]
func (ctrl *OpsController) Info(c *gin.Context) {
	c.JSON(200, models.Info{
		Version:     ctrl.conf.GetPropertyString("application.version"),
		Name:        ctrl.conf.GetPropertyString("application.name"),
		Description: ctrl.conf.GetPropertyString("application.description"),
	})
}

// Info micro service
// @Summary liveness
// @Schemes
// @Description Informations sur le service
// @Tags ops
// @Accept json
// @Produce json
// @Success 200 {object} models.Health
// @Router /ops/liveness [get]
func (ctrl *OpsController) Liveness(c *gin.Context) {
	c.JSON(200, models.Health{
		Status: "UP",
	})
}

// Info micro service
// @Summary Readiness
// @Schemes
// @Description Informations sur le service
// @Tags ops
// @Accept json
// @Produce json
// @Success 200 {object} models.Health
// @Router /ops/readiness [get]
func (ctrl *OpsController) Readiness(c *gin.Context) {
	c.JSON(200, models.Health{
		Status: "UP",
	})
}

// Info micro service
// @Summary Metrics prometheus
// @Schemes
// @Description Informations sur le service
// @Tags ops
// @Accept json
// @Produce json
// @Success 200
// @Router /ops/metrics [get]
func (ctrl *OpsController) Metrics(c *gin.Context) {
	c.JSON(200, models.Health{
		Status: "UP",
	})
}
