package quest

import (
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"go.uber.org/mock/gomock"

	"vk_quests/internal/pkg/types"
	qr "vk_quests/internal/repository/quest"
	mrq "vk_quests/internal/repository/quest/mocks"
)

var testError = errors.New("test error")

type QuestUsecaseSuite struct {
	suite.Suite
	questUsecase *QuestUsecase
	mockQuest    *mrq.QuestRepository
	gmc          *gomock.Controller
}

func (qus *QuestUsecaseSuite) BeforeEach(t provider.T) {
	qus.gmc = gomock.NewController(t)
	qus.mockQuest = mrq.NewQuestRepository(qus.gmc)
	qus.questUsecase = NewQuestUsecase(qus.mockQuest)
}

func (qus *QuestUsecaseSuite) AfterEach(t provider.T) {
	qus.gmc.Finish()
}

func (qus *QuestUsecaseSuite) TestCreateQuestFunction(t provider.T) {
	t.Title("CreateQuest function of quest usecase")
	t.NewStep("Init test data")
	quest := &Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	repositoryQuest := &qr.Quest{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		Cost:        quest.Cost,
		Type:        quest.Type,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qus.mockQuest.EXPECT().CreateQuest(repositoryQuest).Return(repositoryQuest, nil).Times(1)

		t.NewStep("Check result")
		qst, err := qus.questUsecase.CreateQuest(quest)
		t.Require().NoError(err)
		t.Require().Equal(quest, qst)
	})

	t.WithNewStep("Repository error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qus.mockQuest.EXPECT().CreateQuest(repositoryQuest).Return(nil, testError).Times(1)

		t.NewStep("Check result")
		_, err := qus.questUsecase.CreateQuest(quest)
		t.Require().ErrorIs(err, testError)
	})
}

func (qus *QuestUsecaseSuite) TestDeleteQuestFunction(t provider.T) {
	t.Title("DeleteQuest function of quest usecase")
	t.NewStep("Init test data")
	quest := &Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qus.mockQuest.EXPECT().DeleteQuest(quest.ID).Return(nil).Times(1)

		t.NewStep("Check result")
		err := qus.questUsecase.DeleteQuest(quest.ID)
		t.Require().NoError(err)
	})

	t.WithNewStep("Repository error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qus.mockQuest.EXPECT().DeleteQuest(quest.ID).Return(testError).Times(1)

		t.NewStep("Check result")
		err := qus.questUsecase.DeleteQuest(quest.ID)
		t.Require().ErrorIs(err, testError)
	})
}

func (qus *QuestUsecaseSuite) TestUpdateQuestFunction(t provider.T) {
	t.Title("UpdateQuest function of quest usecase")
	t.NewStep("Init test data")
	quest := &Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	repositoryQuest := &qr.Quest{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		Cost:        quest.Cost,
		Type:        quest.Type,
	}

	updateQuest := &UpdateQuest{
		Type:        nil,
		Description: &quest.Description,
		Cost:        nil,
	}

	respositoryUpdateQuest := &qr.UpdateQuest{
		ID:          quest.ID,
		Type:        updateQuest.Type,
		Description: updateQuest.Description,
		Cost:        updateQuest.Cost,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qus.mockQuest.EXPECT().UpdateQuest(respositoryUpdateQuest).Return(repositoryQuest, nil).Times(1)

		t.NewStep("Check result")
		qst, err := qus.questUsecase.UpdateQuest(quest.ID, updateQuest)
		t.Require().NoError(err)
		t.Require().Equal(quest, qst)
	})

	t.WithNewStep("Repository error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qus.mockQuest.EXPECT().UpdateQuest(respositoryUpdateQuest).Return(nil, testError).Times(1)

		t.NewStep("Check result")
		_, err := qus.questUsecase.UpdateQuest(quest.ID, updateQuest)
		t.Require().ErrorIs(err, testError)
	})
}

func (qus *QuestUsecaseSuite) TestGetQuestFunction(t provider.T) {
	t.Title("GetQuest function of quest usecase")
	t.NewStep("Init test data")
	quest := &Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	repositoryQuest := &qr.Quest{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		Cost:        quest.Cost,
		Type:        quest.Type,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qus.mockQuest.EXPECT().GetQuest(quest.ID).Return(repositoryQuest, nil).Times(1)

		t.NewStep("Check result")
		qst, err := qus.questUsecase.GetQuest(quest.ID)
		t.Require().NoError(err)
		t.Require().Equal(quest, qst)
	})

	t.WithNewStep("Repository error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qus.mockQuest.EXPECT().GetQuest(quest.ID).Return(nil, testError).Times(1)

		t.NewStep("Check result")
		_, err := qus.questUsecase.GetQuest(quest.ID)
		t.Require().ErrorIs(err, testError)
	})
}

func (qus *QuestUsecaseSuite) TestGetQuestsFunction(t provider.T) {
	t.Title("GetQuests function of quest usecase")
	t.NewStep("Init test data")
	quest := &Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	repositoryQuest := []qr.Quest{
		{
			ID:          quest.ID,
			Name:        quest.Name,
			Description: quest.Description,
			Cost:        quest.Cost,
			Type:        quest.Type,
		},
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qus.mockQuest.EXPECT().GetQuests().Return(repositoryQuest, nil).Times(1)

		t.NewStep("Check result")
		qst, err := qus.questUsecase.GetQuests()
		t.Require().NoError(err)
		t.Require().Equal([]Quest{*quest}, qst)
	})

	t.WithNewStep("Repository error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qus.mockQuest.EXPECT().GetQuests().Return(nil, testError).Times(1)

		t.NewStep("Check result")
		_, err := qus.questUsecase.GetQuests()
		t.Require().ErrorIs(err, testError)
	})
}

func TestRunQuestUsecaseSuite(t *testing.T) {
	suite.RunSuite(t, new(QuestUsecaseSuite))
}
