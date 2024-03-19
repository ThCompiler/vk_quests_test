package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vk_quests/internal/delivery/http/v1/model/response"
	"vk_quests/internal/pkg/types"
	qr "vk_quests/internal/repository/quest"
	qu "vk_quests/internal/usecase/quest"
	muq "vk_quests/internal/usecase/quest/mocks"
)

type QuestHandlersSuite struct {
	suite.Suite
	handlers  *QuestHandlers
	mockQuest *muq.QuestUsecase
	gmc       *gomock.Controller
}

func (qhs *QuestHandlersSuite) BeforeEach(t provider.T) {
	qhs.gmc = gomock.NewController(t)
	qhs.mockQuest = muq.NewQuestUsecase(qhs.gmc)
	qhs.handlers = NewQuestHandlers(qhs.mockQuest)
}

func (qhs *QuestHandlersSuite) AfterEach(t provider.T) {
	qhs.gmc.Finish()
}

func (qhs *QuestHandlersSuite) TestGetQuestsHandler(t provider.T) {
	t.Title("GetQuests handler of quest handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/", addEmptyLogger(qhs.handlers.GetQuests))

	t.NewStep("Init test data")
	quest := &qu.Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}
	quests := []qu.Quest{*quest, *quest, *quest}

	responseQuest := &response.Quest{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		Cost:        quest.Cost,
		Type:        quest.Type,
	}

	responseQuests := []response.Quest{*responseQuest, *responseQuest, *responseQuest}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().GetQuests().Return(quests, nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
		var qsts []response.Quest
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&qsts))
		t.Require().EqualValues(responseQuests, qsts)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().GetQuests().Return(nil, testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})
}

func (qhs *QuestHandlersSuite) TestGetQuestHandler(t provider.T) {
	t.Title("GetQuest handler of quest handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/:"+QuestIdField, addEmptyLogger(qhs.handlers.GetQuest))

	t.NewStep("Init test data")
	quest := &qu.Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	responseQuest := &response.Quest{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		Cost:        quest.Cost,
		Type:        quest.Type,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().GetQuest(quest.ID).Return(quest, nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
		var qst response.Quest
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&qst))
		t.Require().EqualValues(responseQuest, &qst)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().GetQuest(quest.ID).Return(nil, testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})

	t.WithNewStep("Quest not found error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().GetQuest(quest.ID).Return(nil, qr.ErrorQuestNotFound).Times(1)

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

func (qhs *QuestHandlersSuite) TestDeleteQuestHandler(t provider.T) {
	t.Title("DeleteQuest handler of quest handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/:"+QuestIdField, addEmptyLogger(qhs.handlers.DeleteQuest))

	t.NewStep("Init test data")
	quest := &qu.Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().DeleteQuest(quest.ID).Return(nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().DeleteQuest(quest.ID).Return(testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", nil, nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})

	t.WithNewStep("Quest not found error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().DeleteQuest(quest.ID).Return(qr.ErrorQuestNotFound).Times(1)

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

func (qhs *QuestHandlersSuite) TestCreateQuestHandler(t provider.T) {
	t.Title("CreateQuest handler of quest handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/", addEmptyLogger(qhs.handlers.CreateQuest))

	t.NewStep("Init test data")
	quest := &qu.Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	newQuest := &qu.Quest{
		Name:        quest.Name,
		Description: quest.Description,
		Cost:        quest.Cost,
		Type:        quest.Type,
	}

	body := `
		{
			"name": "Quest",
			"description": "good Quest",
			"cost": 10,
			"type": "usual"
		}
	`

	responseQuest := &response.Quest{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		Cost:        quest.Cost,
		Type:        quest.Type,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().CreateQuest(newQuest).Return(quest, nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusCreated, recorder.Code)
		var qst response.Quest
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&qst))
		t.Require().EqualValues(responseQuest, &qst)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().CreateQuest(newQuest).Return(nil, testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})

	t.WithNewStep("Quest name already exists error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().CreateQuest(newQuest).Return(nil, qr.ErrorQuestNameAlreadyExists).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusConflict, recorder.Code)
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

func (qhs *QuestHandlersSuite) TestUpdateQuestHandler(t provider.T) {
	t.Title("UpdateQuest handler of quest handlers")
	t.NewStep("Init gin routes")
	r := gin.New()
	r.POST("/:"+QuestIdField, addEmptyLogger(qhs.handlers.UpdateQuest))

	t.NewStep("Init test data")
	quest := &qu.Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	updateQuest := &qu.UpdateQuest{
		Description: &quest.Description,
		Cost:        &quest.Cost,
		Type:        &quest.Type,
	}

	body := `
		{
			"description": "good Quest",
			"cost": 10,
			"type": "usual"
		}
	`

	nilUpdateQuest := &qu.UpdateQuest{
		Description: nil,
		Cost:        nil,
		Type:        nil,
	}

	nilBody := `
		{
		}
	`

	responseQuest := &response.Quest{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		Cost:        quest.Cost,
		Type:        quest.Type,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().UpdateQuest(quest.ID, updateQuest).Return(quest, nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
		var qst response.Quest
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&qst))
		t.Require().EqualValues(responseQuest, &qst)
	})

	t.WithNewStep("Correct no changes execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().UpdateQuest(quest.ID, nilUpdateQuest).Return(quest, nil).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", strings.NewReader(nilBody), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusOK, recorder.Code)
		var qst response.Quest
		dec := json.NewDecoder(recorder.Body)
		t.Require().NoError(dec.Decode(&qst))
		t.Require().EqualValues(responseQuest, &qst)
	})

	t.WithNewStep("Usecase error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().UpdateQuest(quest.ID, updateQuest).Return(nil, testError).Times(1)

		t.NewStep("Init http")
		req, err := initRequest(http.MethodPost, "/1", strings.NewReader(body), nil)
		t.Require().NoError(err)

		recorder := httptest.NewRecorder()

		t.NewStep("Check result")
		r.ServeHTTP(recorder, req)

		t.Require().Equal(http.StatusInternalServerError, recorder.Code)
	})

	t.WithNewStep("Quest not found error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qhs.mockQuest.EXPECT().UpdateQuest(quest.ID, updateQuest).Return(nil, qr.ErrorQuestNotFound).Times(1)

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

func TestRunQuestHandlersSuite(t *testing.T) {
	suite.RunSuite(t, new(QuestHandlersSuite))
}
