package httpres

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	errCode = map[string]int{}
)

func HTTPSuccesResponse(ctx *gin.Context, code int, data any) {
	ctx.JSON(code, gin.H{
		"data": data,
	})
}

func HTTPErrorResponse(ctx *gin.Context, err error) {
	code, ok := errCode[err.Error()]
	if !ok {
		code = http.StatusBadRequest
	}

	ctx.JSON(code, gin.H{
		"error": err.Error(),
	})
}
