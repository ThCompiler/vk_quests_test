package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"vk_quests/internal/pkg/evjson"
	"vk_quests/pkg/logger"
)

func parseRequestBody(b io.ReadCloser, out any, validation func([]byte) error, l logger.Interface) (int, error) {
	body, err := io.ReadAll(b)
	if err != nil {
		l.Error(errors.Wrapf(err, "can't read body"))
		return http.StatusInternalServerError, ErrorCannotReadBody
	}

	// Проверка корректности тела запроса
	if err := validation(body); err != nil {
		if errors.Is(err, evjson.ErrorInvalidJson) {
			l.Warn(errors.Wrapf(err, "try parse body json"))
			return http.StatusBadRequest, ErrorIncorrectBodyContent
		}
		return http.StatusBadRequest, errors.Wrapf(err, "in body error")
	}

	// Получение значения тела запроса
	if err := json.Unmarshal(body, out); err != nil {
		l.Warn(errors.Wrapf(err, "try parse create request entity"))
		return http.StatusBadRequest, ErrorIncorrectBodyContent
	}

	return http.StatusOK, nil
}
