package xgin

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

type XGinEngine struct {
	*gin.Engine
}

func (e *XGinEngine) Listen(addr string) error {
	if err := endless.ListenAndServe(addr, e.Engine); err != nil {
		return err
	}
	return nil
}

func NewGin(logger *zap.Logger) *XGinEngine {
	e := gin.New()

	e.Use(gin.Recovery())
	e.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	e.Use(ginzap.RecoveryWithZap(logger, true))

	return &XGinEngine{Engine: e}
}
