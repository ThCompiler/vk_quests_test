package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"vk_quests/internal/delivery/http/v1/model/request"
	"vk_quests/internal/delivery/http/v1/model/response"
	"vk_quests/internal/delivery/middleware"
	"vk_quests/internal/pkg/types"
	qr "vk_quests/internal/repository/quest"
	qu "vk_quests/internal/usecase/quest"
	"vk_quests/pkg/operate"
)

const (
	QuestIdField = "quest_id"
)

type QuestHandlers struct {
	quests qu.Usecase
}

func NewQuestHandlers(quests qu.Usecase) *QuestHandlers {
	return &QuestHandlers{quests: quests}
}

// CreateQuest
//
//	@Summary		Добавление задание.
//	@Description	Добавляет задание включая его название(уникальное), описание, стоимость и тип. Есть обычная задание, которое выполняется как только вызывается метод, сигнализирующий о выполнении для пользователя задачи. И случайная задача, которая выполняется в с вероятностью 0,5.
//	@Tags			quest
//	@Accept			json
//	@Param			request	body	request.CreateQuest	true	"Информация о добавляемом фильме"
//	@Produce		json
//	@Success		201	{object}	response.Quest		"Задание успешно добавлен в базу"
//	@Failure		400	{object}	operate.ModelError	"В теле запроса ошибка"
//	@Failure		409	{object}	operate.ModelError	"Задача с таким название уже существует"
//	@Failure		500	{object}	operate.ModelError	"Ошибка сервера"
//	@Router			/quest [post]
func (qh *QuestHandlers) CreateQuest(c *gin.Context) {
	l := middleware.GetLogger(c)

	// Получение значения тела запроса
	var createQuest request.CreateQuest
	if code, err := parseRequestBody(c.Request.Body, &createQuest, request.ValidateCreateQuest, l); err != nil {
		operate.SendError(c, err, code, l)
		return
	}

	createdQuest, err := qh.quests.CreateQuest(createQuest.ToUsQuest())
	if err != nil {
		if errors.Is(err, qr.ErrorQuestNameAlreadyExists) {
			operate.SendError(c, ErrorQuestNameAlreadyExists, http.StatusConflict, l)
			l.Info(errors.Wrapf(err, "can't create quest"))
			return
		}
		operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
		l.Error(errors.Wrapf(err, "can't create quest"))
		return
	}

	operate.SendStatus(c, http.StatusCreated, response.FromUsQuest(createdQuest), l)
}

// DeleteQuest
//
//	@Summary		Удаление задания.
//	@Description	Удаляет информацию о задании из системы по его id.
//	@Tags			quest
//	@Param			quest_id	path	uint64	true	"Уникальный идентификатор задания"
//	@Produce		json
//	@Success		200	"Задание успешно удалено"
//	@Failure		400	{object}	operate.ModelError	"В пути запроса ошибка"
//	@Failure		404	{object}	operate.ModelError	"Задание с указанным id не найдено"
//	@Failure		500	{object}	operate.ModelError	"Ошибка сервера"
//	@Router			/quest/{quest_id} [delete]
func (qh *QuestHandlers) DeleteQuest(c *gin.Context) {
	l := middleware.GetLogger(c)

	// Получение уникального идентификатора
	id, err := strconv.ParseUint(c.Param(QuestIdField), 10, 64)
	if err != nil {
		operate.SendError(c, errors.Wrapf(err, "try get quest id"), http.StatusBadRequest, l)
		return
	}

	if err = qh.quests.DeleteQuest(types.Id(id)); err != nil {
		if errors.Is(err, qr.ErrorQuestNotFound) {
			operate.SendError(c, ErrorQuestNotFound, http.StatusNotFound, l)
			return
		}
		operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
		l.Error(errors.Wrapf(err, "can't delete quest"))
		return
	}

	operate.SendStatus(c, http.StatusOK, nil, l)
}

// GetQuest
//
//	@Summary		Получение задания.
//	@Description	Позволяет информацию о задании по его id.
//	@Tags			quest
//	@Param			quest_id	path	uint64	true	"Уникальный идентификатор задания"
//	@Produce		json
//	@Success		200	{array}		response.Quest		"Полученное задание"
//	@Failure		400	{object}	operate.ModelError	"В пути запросе ошибка"
//	@Failure		404	{object}	operate.ModelError	"Задание с указанным id не найдено"
//	@Failure		500	{object}	operate.ModelError	"Ошибка сервера"
//	@Router			/quest/{quest_id} [get]
func (qh *QuestHandlers) GetQuest(c *gin.Context) {
	l := middleware.GetLogger(c)

	// Получение уникального идентификатора
	id, err := strconv.ParseUint(c.Param(QuestIdField), 10, 64)
	if err != nil {
		operate.SendError(c, errors.Wrapf(err, "try get quest id"), http.StatusBadRequest, l)
		return
	}

	quests, err := qh.quests.GetQuest(types.Id(id))
	if err != nil {
		if errors.Is(err, qr.ErrorQuestNotFound) {
			operate.SendError(c, ErrorQuestNotFound, http.StatusNotFound, l)
			l.Error(errors.Wrapf(err, "can't get quests"))
			return
		}
		operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
		l.Error(errors.Wrapf(err, "can't get quests"))
		return
	}

	operate.SendStatus(c, http.StatusOK, response.FromUsQuest(quests), l)
}

// UpdateQuest
//
//	@Summary		Обновление данных об задании.
//	@Description	Обновляет данные об задании. Все переданные поля будут обновлены. Отсутствующие поля будут оставлены без изменений.
//	@Tags			quest
//	@Accept			json
//	@Param			quest_id	path	uint64				true	"Уникальный идентификатор задания"
//	@Param			request		body	request.UpdateQuest	true	"Информация об обновлении"
//	@Produce		json
//	@Success		200	{object}	response.Quest		"Задание успешно обновлено в базе"
//	@Failure		400	{object}	operate.ModelError	"В теле запроса ошибка"
//	@Failure		404	{object}	operate.ModelError	"Задание с указанным id не найден"
//	@Failure		500	{object}	operate.ModelError	"Ошибка сервера"
//	@Router			/quest/{quest_id} [put]
func (qh *QuestHandlers) UpdateQuest(c *gin.Context) {
	l := middleware.GetLogger(c)

	// Получение уникального идентификатора
	id, err := strconv.ParseUint(c.Param(QuestIdField), 10, 64)
	if err != nil {
		operate.SendError(c, errors.Wrapf(err, "try get quest id"), http.StatusBadRequest, l)
		return
	}

	// Получение значения тела запроса
	var updateQuest request.UpdateQuest
	if code, err := parseRequestBody(c.Request.Body, &updateQuest, request.ValidateUpdateQuest, l); err != nil {
		operate.SendError(c, err, code, l)
		return
	}

	updatedQuest, err := qh.quests.UpdateQuest(types.Id(id), updateQuest.ToUsUpdateQuest())
	if err != nil {
		if errors.Is(err, qr.ErrorQuestNotFound) {
			operate.SendError(c, ErrorQuestNotFound, http.StatusNotFound, l)
			return
		}

		operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
		l.Error(errors.Wrapf(err, "can't update quest"))
		return
	}

	operate.SendStatus(c, http.StatusOK, response.FromUsQuest(updatedQuest), l)
}

// GetQuests
//
//	@Summary		Получение списка заданий.
//	@Description	Позволяет получить список заданий.
//	@Tags			quest
//	@Produce		json
//	@Success		200	{array}		response.Quest		"Список заданий успешно сформирован"
//	@Failure		500	{object}	operate.ModelError	"Ошибка сервера"
//	@Router			/quest/list [get]
func (qh *QuestHandlers) GetQuests(c *gin.Context) {
	l := middleware.GetLogger(c)

	quests, err := qh.quests.GetQuests()
	if err != nil {
		operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
		l.Error(errors.Wrapf(err, "can't get quests"))
		return
	}

	operate.SendStatus(c, http.StatusOK, response.FromUsQuests(quests), l)
}
