package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"vk_quests/internal/delivery/middleware"
	"vk_quests/internal/pkg/types"
	"vk_quests/pkg/logger"
)

var testError = errors.New("test error")

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, testError
}

func (errReader) Close() error {
	return testError
}

type emptyLogger struct{}

func (*emptyLogger) Debug(_ any, _ ...any)                          {}
func (*emptyLogger) Info(_ any, _ ...any)                           {}
func (*emptyLogger) Warn(_ any, _ ...any)                           {}
func (*emptyLogger) Error(_ any, _ ...any)                          {}
func (*emptyLogger) Panic(_ any, _ ...any)                          {}
func (*emptyLogger) Fatal(_ any, _ ...any)                          {}
func (el *emptyLogger) With(_ logger.Field, _ any) logger.Interface { return el }

func initRequest(method string, path string, body io.Reader, contextValues map[types.ContextField]any) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)

	if err != nil {
		return nil, err
	}

	ctx := req.Context()
	for k, v := range contextValues {
		ctx = context.WithValue(ctx, k, v)
	}

	return req.WithContext(ctx), nil
}

func addEmptyLogger(handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(string(middleware.LoggerField), logger.Interface(&emptyLogger{}))
		handlerFunc(ctx)
	}
}
