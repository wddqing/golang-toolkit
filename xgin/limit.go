package xgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeMaxAllowedMiddleWare(n int) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }

	return func(c *gin.Context) {
		acquire()       // before request
		defer release() // after request
		c.Next()
	}
}

func MakeMaxAllowedProtectMiddleWare(n int) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	acquire := func() bool {
		select {
		case sem <- struct{}{}:
			return true
		default:
			return false
		}
	}
	release := func() {
		select {
		case <-sem:
		default:
		}
	}

	return func(c *gin.Context) {
		if acquire() { // before request
			defer release() // after request
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusGatewayTimeout)
		}
	}
}
