package controller

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/laurentpoirierfr/ora-cdc-go/internal/models"
	"github.com/laurentpoirierfr/ora-cdc-go/pkg/config"
)

// IndexController is the default controller
type IndexController struct {
	context context.Context
	conf    config.Config
}

func NewIndexController(context context.Context, conf config.Config) *IndexController {
	return &IndexController{
		conf:    conf,
		context: context,
	}
}

// @Summary	Get demo
// @Tags		api
// @Accept		json
// @Produce	json
// @Success 200 {object} models.Health
// @Router		/api/demo [get]
func (ctrl *IndexController) Demo(c *gin.Context) {
	c.JSON(200, models.Health{
		Status: "UP",
	})
}
