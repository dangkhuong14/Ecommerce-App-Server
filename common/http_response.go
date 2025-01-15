package common

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/viettranx/service-context/core"
)

type DebugCarrier interface {
	WithDebug(debugMessage string) *core.DefaultError
}

func WriteErrorResponse(c *gin.Context, err error) {
	if errSt, ok := err.(core.StatusCodeCarrier); ok {
		fmt.Printf("Type of errSt: %s\n", reflect.TypeOf(errSt))
		// If current environment is not debugging: delete debug message
		if !gin.IsDebugging() {
			if errSt, ok := errSt.(DebugCarrier); ok{
				errWithoutDebug := errSt.WithDebug("")
				c.JSON(errWithoutDebug.StatusCode(), errWithoutDebug)
				return
			}
		}

		c.JSON(errSt.StatusCode(), errSt)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}