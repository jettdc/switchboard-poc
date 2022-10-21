package u

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpError struct {
	error
	code int
	msg  string
}

func InternalServerError(msg string, params ...any) HttpError {
	return HttpError{
		errors.New(http.StatusText(http.StatusInternalServerError)),
		http.StatusInternalServerError,
		msg,
	}
}

func Err(c *gin.Context, err HttpError) {
	c.JSON(err.code, gin.H{
		"status":  err.code,
		"error":   err.Error(),
		"message": err.msg,
	})
}
