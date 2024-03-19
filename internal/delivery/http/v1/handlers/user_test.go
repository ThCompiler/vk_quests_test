package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"

	"vk_quests/internal/delivery/http/v1/model/response"
	"vk_quests/internal/pkg/types"
	qr "vk_quests/internal/repository/quest"
	ur "vk_quests/internal/repository/user"
	qu "vk_quests/internal/usecase/quest"
	uu "vk_quests/internal/usecase/user"
	muu "vk_quests/internal/usecase/user/mocks"
)

type UserHandlersSuite struct {
	suite.Suite
	handlers *UserHandlers
	mockUser *muu.UserUsecase
	gmc      *gomock.Controller
}

func (uhs *UserHandlersSuite) BeforeEach(t provider.T) {
	uhs.gmc = gomock.NewController(t)
	uhs.mockUser = muu.NewUserUsecase(uhs.gmc)
	uhs.handlers = NewUserHandlers(uhs.mockUser)
}

func (uhs *UserHandlersSuite) AfterEach(t provider.T) {
	uhs.gmc.Finish()
}

func (uhs *UserHandlersSuite) TestGetUsersHandler(t provider.T) {
	t.Title("GetUsers handler of user handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/", addEmptyLogger(uhs.handlers.GetUsers))

	t.NewStep("Init test data")
	user := &uu.User{
		ID:      1,
		Name:    "User",
		Balance: 10,
	}
	users := []uu.User{*user, *user, *user}

	responseUser := &response.User{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}

	responseUsers := []response.User{*responseUser, *responseUser, *responseUser}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().GetUsers().Return(users, nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
		var usrs []response.User
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&usrs))
		t.Require().EqualValues(responseUsers, usrs)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().GetUsers().Return(nil, testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})
}

func (uhs *UserHandlersSuite) TestDeleteUserHandler(t provider.T) {
	t.Title("DeleteUser handler of user handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/:"+UserIdField, addEmptyLogger(uhs.handlers.DeleteUser))

	t.NewStep("Init test data")
	user := &uu.User{
		ID:      1,
		Name:    "User",
		Balance: 20,
	}

	responseUser := &response.User{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().DeleteUser(user.ID).Return(user, nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
		var usr response.User
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&usr))
		t.Require().EqualValues(responseUser, &usr)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().DeleteUser(user.ID).Return(nil, testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})

	t.WithNewStep("User not found error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().DeleteUser(user.ID).Return(nil, ur.ErrorUserNotFound).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusNotFound, recorder.Code)
	})

	t.WithNewStep("Incorrect path param error execute", func(t provider.StepCtx) {
		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/sus", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusBadRequest, recorder.Code)
	})
}

func (uhs *UserHandlersSuite) TestCreateUserHandler(t provider.T) {
	t.Title("CreateUser handler of user handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/", addEmptyLogger(uhs.handlers.CreateUser))

	t.NewStep("Init test data")
	user := &uu.User{
		ID:      1,
		Name:    "User",
		Balance: 0,
	}

	body := `
		{
			"name": "User"
		}
	`

	responseUser := &response.User{
		ID:      user.ID,
		Name:    user.Name,
		Balance: 0,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().CreateUser(user.Name).Return(user, nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusCreated, recorder.Code)
		var usr response.User
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&usr))
		t.Require().EqualValues(responseUser, &usr)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().CreateUser(user.Name).Return(nil, testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})

	t.WithNewStep("Incorrect body error execute", func(t provider.StepCtx) {
		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/", errReader(1), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})
}

func (uhs *UserHandlersSuite) TestUpdateUserHandler(t provider.T) {
	t.Title("UpdateUser handler of user handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/:"+UserIdField, addEmptyLogger(uhs.handlers.UpdateUser))

	t.NewStep("Init test data")
	user := &uu.User{
		ID:      1,
		Name:    "User",
		Balance: 25,
	}

	body := `
		{
			"name": "User"
		}
	`

	responseUser := &response.User{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().UpdateUser(user.ID, user.Name).Return(user, nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
		var usr response.User
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&usr))
		t.Require().EqualValues(responseUser, &usr)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().UpdateUser(user.ID, user.Name).Return(nil, testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})

	t.WithNewStep("User not found error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().UpdateUser(user.ID, user.Name).Return(nil, ur.ErrorUserNotFound).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusNotFound, recorder.Code)
	})

	t.WithNewStep("Incorrect query param execute", func(t provider.StepCtx) {
		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/qwerty", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusBadRequest, recorder.Code)
	})

	t.WithNewStep("Incorrect body error execute", func(t provider.StepCtx) {
		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", errReader(1), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})
}

func (uhs *UserHandlersSuite) TestGetHistoryHandler(t provider.T) {
	t.Title("GetUserHistory handler of user handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/:"+UserIdField, addEmptyLogger(uhs.handlers.GetUserHistory))

	t.NewStep("Init test data")
	userId := types.Id(1)
	history := []uu.HistoryRecord{
		{
			Quest: &qu.Quest{
				ID:   1,
				Name: "Not null",
			},
			Balance: 30,
		},
		{
			Quest:   nil,
			Balance: 25,
		},
	}

	responseHistory := []response.HistoryRecord{
		{
			Quest: &response.Quest{
				ID:   history[0].Quest.ID,
				Name: history[0].Quest.Name,
			},
			Balance: history[0].Balance,
		},
		{
			Quest:   nil,
			Balance: history[1].Balance,
		},
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().GetUserHistory(userId).Return(history, nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
		var hist []response.HistoryRecord
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&hist))
		t.Require().EqualValues(responseHistory, hist)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().GetUserHistory(userId).Return(nil, testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})

	t.WithNewStep("Incorrect path param execute", func(t provider.StepCtx) {
		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/qwerty", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusBadRequest, recorder.Code)
	})
}

func (uhs *UserHandlersSuite) TestCompleteQuestHandler(t provider.T) {
	t.Title("CompleteQuest handler of user handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/", addEmptyLogger(uhs.handlers.CompleteQuest))

	t.NewStep("Init test data")
	userId := types.Id(1)
	questId := types.Id(2)

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().ApplyQuests(questId, userId).Return(nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/?"+UserIdField+"=1&"+QuestIdField+"=2", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
		var status response.StatusApplyCost
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&status))
		t.Require().EqualValues(response.Success, status.Status)
	})

	t.WithNewStep("Correct execute quest not applied", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().ApplyQuests(questId, userId).Return(uu.QuestNotApplied).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/?"+UserIdField+"=1&"+QuestIdField+"=2", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
		var status response.StatusApplyCost
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&status))
		t.Require().EqualValues(response.Failure, status.Status)
	})

	t.WithNewStep("User not found error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().ApplyQuests(questId, userId).Return(ur.ErrorUserNotFound).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/?"+UserIdField+"=1&"+QuestIdField+"=2", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusNotFound, recorder.Code)
	})

	t.WithNewStep("Quest not found error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().ApplyQuests(questId, userId).Return(qr.ErrorQuestNotFound).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/?"+UserIdField+"=1&"+QuestIdField+"=2", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusNotFound, recorder.Code)
	})

	t.WithNewStep("Quest already complete for user error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().ApplyQuests(questId, userId).Return(ur.ErrorUserAlreadyCompleteQuest).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/?"+UserIdField+"=1&"+QuestIdField+"=2", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusConflict, recorder.Code)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().ApplyQuests(questId, userId).Return(testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/?"+UserIdField+"=1&"+QuestIdField+"=2", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})

	t.WithNewStep("Incorrect user id query param execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().ApplyQuests(questId, userId).Return(testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/?"+UserIdField+"=ar&"+QuestIdField+"=2", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusBadRequest, recorder.Code)
	})

	t.WithNewStep("Incorrect quest id query param execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().ApplyQuests(questId, userId).Return(testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/?"+UserIdField+"=1&"+QuestIdField+"=ar", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusBadRequest, recorder.Code)
	})

	t.WithNewStep("User id query param not presented execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().ApplyQuests(questId, userId).Return(testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/?"+QuestIdField+"=2", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusBadRequest, recorder.Code)
	})

	t.WithNewStep("Quest id query param not presented execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uhs.mockUser.EXPECT().ApplyQuests(questId, userId).Return(testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/?"+UserIdField+"=1", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusBadRequest, recorder.Code)
	})
}

func TestRunUserHandlersSuite(t *testing.T) {
	suite.RunSuite(t, new(UserHandlersSuite))
}
