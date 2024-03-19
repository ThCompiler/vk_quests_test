package operate

import (
	"github.com/gin-gonic/gin"
	"vk_quests/pkg/logger"
)

type ModelError struct {
	ErrorMessage string `json:"error_message,omitempty"`
}

func SendError(c *gin.Context, err error, code int, l logger.Interface) {
	c.AbortWithStatusJSON(code, ModelError{ErrorMessage: err.Error()})
	l.Info("error %s was sent with status code %d", err, code)
}
