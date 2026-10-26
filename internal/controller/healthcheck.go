package controller

import (
	"chatsmith-mini/pkg/response"

	"github.com/gin-gonic/gin"
)

const healthStatusOK = "ok"

type HealthCheck struct{}

func NewHealthCheck() *HealthCheck {
	return &HealthCheck{}
}

func (h *HealthCheck) Check(ctx *gin.Context) {
	response.Success(ctx, healthStatusOK)
}
