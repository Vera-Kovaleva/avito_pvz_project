package container

import (
	"context"

	"github.com/gin-gonic/gin"
)

// аааааа ??

type Container struct {
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) HTTPServer(ctx context.Context) *gin.Engine {
	router := gin.Default()

	return router
}
