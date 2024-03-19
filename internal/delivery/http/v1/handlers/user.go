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
	ur "vk_quests/internal/repository/user"
	uu "vk_quests/internal/usecase/user"
	"vk_quests/pkg/operate"
)

const UserIdField = "user_id"

type UserHandlers struct {
	users uu.Usecase
}

func NewUserHandlers(users uu.Usecase) *UserHandlers {
	return &UserHandlers{users: users}
}

// CreateUser
//
//	@Summary		Добавление пользователя.
//	@Description	Добавляет пользователя включая его имя. Баланс пользователя при создании 0.
//	@Tags			user
//	@Accept			json
//	@Param			request	body	request.User	true	"Информация о добавляемом пользователе"
//	@Produce		json
//	@Success		201	{object}	response.User		"Пользователь успешно добавлен в базу"
//	@Failure		400	{object}	operate.ModelError	"В теле запроса ошибка"
//	@Failure		500	{object}	operate.ModelError	"Ошибка сервера"
//	@Router			/user [post]
func (uh *UserHandlers) CreateUser(c *gin.Context) {
	l := middleware.GetLogger(c)

	// Получение значения тела запроса
	var createUser request.User
	if code, err := parseRequestBody(c.Request.Body, &createUser, request.ValidateUser, l); err != nil {
		operate.SendError(c, err, code, l)
		return
	}

	createdUser, err := uh.users.CreateUser(createUser.Name)
	if err != nil {
		operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
		l.Error(errors.Wrapf(err, "can't create film"))
		return
	}

	operate.SendStatus(c, http.StatusCreated, response.FromUsUser(createdUser), l)
}

// DeleteUser
//
//	@Summary		Удаление пользователя.
//	@Description	Удаляет информацию об пользователе по его id.
//	@Tags			user
//	@Param			user_id	path	uint64	true	"Уникальный идентификатор пользователя"
//	@Produce		json
//	@Success		200	{object}	response.User		"Пользователь успешно удалён"
//	@Failure		400	{object}	operate.ModelError	"В пути запросе ошибка"
//	@Failure		404	{object}	operate.ModelError	"Пользователь с указанным id не найден"
//	@Failure		500	{object}	operate.ModelError	"Ошибка сервера"
//	@Router			/user/{user_id} [delete]
func (uh *UserHandlers) DeleteUser(c *gin.Context) {
	l := middleware.GetLogger(c)

	// Получение уникального идентификатора
	id, err := strconv.ParseUint(c.Param(UserIdField), 10, 64)
	if err != nil {
		operate.SendError(c, errors.Wrapf(err, "try get user id"), http.StatusBadRequest, l)
		return
	}

	user, err := uh.users.DeleteUser(types.Id(id))
	if err != nil {
		if errors.Is(err, ur.ErrorUserNotFound) {
			operate.SendError(c, ErrorUserNotFound, http.StatusNotFound, l)
			return
		}
		operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
		l.Error(errors.Wrapf(err, "can't delete user"))
		return
	}

	operate.SendStatus(c, http.StatusOK, response.FromUsUser(user), l)
}

// UpdateUser
//
//	@Summary		Обновление данных об пользователе.
//	@Description	Обновляет имя пользователя по его id.
//	@Tags			user
//	@Accept			json
//	@Param			user_id	path	uint64			true	"Уникальный идентификатор пользователя"
//	@Param			request	body	request.User	true	"Информация об обновлении"
//	@Produce		json
//	@Success		200	{object}	response.User		"Пользователь успешно обновлен в базе"
//	@Failure		400	{object}	operate.ModelError	"В теле запроса ошибка"
//	@Failure		404	{object}	operate.ModelError	"Пользователь с указанным id не найден"
//	@Failure		500	{object}	operate.ModelError	"Ошибка сервера"
//	@Router			/user/{user_id} [put]
func (uh *UserHandlers) UpdateUser(c *gin.Context) {
	l := middleware.GetLogger(c)

	// Получение уникального идентификатора
	id, err := strconv.ParseUint(c.Param(UserIdField), 10, 64)
	if err != nil {
		operate.SendError(c, errors.Wrapf(err, "try get user id"), http.StatusBadRequest, l)
		return
	}

	// Получение значения тела запроса
	var updateUser request.User
	if code, err := parseRequestBody(c.Request.Body, &updateUser, request.ValidateUser, l); err != nil {
		operate.SendError(c, err, code, l)
		return
	}

	updatedUser, err := uh.users.UpdateUser(types.Id(id), updateUser.Name)

	if err != nil {
		if errors.Is(err, ur.ErrorUserNotFound) {
			operate.SendError(c, ErrorUserNotFound, http.StatusNotFound, l)
			return
		}
		operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
		l.Error(errors.Wrapf(err, "can't update user"))
		return
	}

	operate.SendStatus(c, http.StatusOK, response.FromUsUser(updatedUser), l)
}

// GetUsers
//
//	@Summary		Получение списка пользователь.
//	@Description	Формирует список всех пользователей в системы.
//	@Tags			user
//	@Produce		json
//	@Success		200	{array}		response.User		"Список пользователей успешно сформирован"
//	@Failure		500	{object}	operate.ModelError	"Ошибка сервера"
//	@Router			/user/list [get]
func (uh *UserHandlers) GetUsers(c *gin.Context) {
	l := middleware.GetLogger(c)

	users, err := uh.users.GetUsers()
	if err != nil {
		operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
		l.Error(errors.Wrapf(err, "can't get users"))
		return
	}

	operate.SendStatus(c, http.StatusOK, response.FromUsUsers(users), l)
}

// GetUserHistory
//
//	@Summary		Получение истории выполнения заданий пользователем.
//	@Description	Формирует список выполненных заданий пользователя по его id. Если задача была удалена, то информация о ней не будет выводиться.
//	@Tags			user
//	@Param			user_id	path	uint64	true	"Уникальный идентификатор пользователя"
//	@Produce		json
//	@Success		200	{array}		response.HistoryRecord	"Список выполненных заданий пользователя сформирован"
//	@Failure		400	{object}	operate.ModelError		"В пути запроса ошибка"
//	@Failure		500	{object}	operate.ModelError		"Ошибка сервера"
//	@Router			/user/{user_id}/history [get]
func (uh *UserHandlers) GetUserHistory(c *gin.Context) {
	l := middleware.GetLogger(c)

	// Получение уникального идентификатора
	id, err := strconv.ParseUint(c.Param(UserIdField), 10, 64)
	if err != nil {
		operate.SendError(c, errors.Wrapf(err, "try get user id"), http.StatusBadRequest, l)
		return
	}

	history, err := uh.users.GetUserHistory(types.Id(id))
	if err != nil {
		operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
		l.Error(errors.Wrapf(err, "can't get users"))
		return
	}

	operate.SendStatus(c, http.StatusOK, response.FromUsHistory(history), l)
}

// CompleteQuest
//
//	@Summary		Сообщение о выполнение условии для определённого пользователя определённого задания.
//	@Description	Обрабатывает информацию о выполнение условии для определённого пользователя определённого задания по их идентификаторам.
//	@Tags			user
//	@Param			user_id		query	uint64	true	"Уникальный идентификатор пользователи"
//	@Param			quest_id	query	uint64	true	"Уникальный идентификатор задачи"
//	@Produce		json
//	@Success		200	{array}		response.StatusApplyCost	"Результат применения задания к пользователю. Если 'success' - то задача засчитана пользователю, иначе не засчитана"
//	@Failure		400	{object}	operate.ModelError			"В параметрах запроса ошибка"
//	@Failure		404	{object}	operate.ModelError			"Пользователь или задача не найдены"
//	@Failure		409	{object}	operate.ModelError			"Данную задачу пользователь уже выполнил"
//	@Failure		500	{object}	operate.ModelError			"Ошибка сервера"
//	@Router			/user/complete [post]
func (uh *UserHandlers) CompleteQuest(c *gin.Context) {
	l := middleware.GetLogger(c)

	// Получение уникального идентификатора
	userId, err := strconv.ParseUint(c.Query(UserIdField), 10, 64)
	if err != nil {
		operate.SendError(c, ErrorIncorrectQueryParam, http.StatusBadRequest, l)
		l.Error(errors.Wrapf(err, "try get user id"))
		return
	}

	// Получение уникального идентификатора
	questId, err := strconv.ParseUint(c.Query(QuestIdField), 10, 64)
	if err != nil {
		operate.SendError(c, ErrorIncorrectQueryParam, http.StatusBadRequest, l)
		l.Error(errors.Wrapf(err, "try get quest id"))
		return
	}

	if err = uh.users.ApplyQuests(types.Id(questId), types.Id(userId)); err != nil {
		switch {
		case errors.Is(err, uu.QuestNotApplied):
			operate.SendStatus(c, http.StatusOK, &response.StatusApplyCost{Status: response.Failure}, l)
		case errors.Is(err, qr.ErrorQuestNotFound):
			operate.SendError(c, ErrorQuestNotFound, http.StatusNotFound, l)
		case errors.Is(err, ur.ErrorUserNotFound):
			operate.SendError(c, ErrorUserNotFound, http.StatusNotFound, l)
		case errors.Is(err, ur.ErrorUserAlreadyCompleteQuest):
			operate.SendError(c, ErrorUserAlreadyCompleteQuest, http.StatusConflict, l)
		default:
			operate.SendError(c, ErrorUnknownError, http.StatusInternalServerError, l)
			l.Error(errors.Wrapf(err, "can't apply quest with id %d to user with id %d", questId, userId))
			return
		}
		l.Warn(errors.Wrapf(err, "can't apply quest with id %d to user with id %d", questId, userId))
		return
	}

	operate.SendStatus(c, http.StatusOK, &response.StatusApplyCost{Status: response.Success}, l)
}
