package user

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	"go.uber.org/mock/gomock"
	"math/rand"
	"testing"
	"vk_quests/internal/pkg/types"
	qr "vk_quests/internal/repository/quest"
	mrq "vk_quests/internal/repository/quest/mocks"
	ur "vk_quests/internal/repository/user"
	mru "vk_quests/internal/repository/user/mocks"
	qu "vk_quests/internal/usecase/quest"
)

var testError = errors.New("test error")

type UserUsecaseSuite struct {
	suite.Suite
	userUsecase *UserUsecase
	mockQuest   *mrq.QuestRepository
	mockUser    *mru.UserRepository
	gmc         *gomock.Controller
}

func (uus *UserUsecaseSuite) BeforeEach(t provider.T) {
	uus.gmc = gomock.NewController(t)
	uus.mockQuest = mrq.NewQuestRepository(uus.gmc)
	uus.mockUser = mru.NewUserRepository(uus.gmc)
	uus.userUsecase = NewUserUsecase(uus.mockUser, uus.mockQuest)
}

func (uus *UserUsecaseSuite) AfterEach(t provider.T) {
	uus.gmc.Finish()
}

func (uus *UserUsecaseSuite) TestCreateUserFunction(t provider.T) {
	t.Title("CreateUser function of user usecase")
	t.NewStep("Init test data")
	user := &User{
		ID:      1,
		Name:    "User",
		Balance: 30,
	}

	repositoryUser := &ur.User{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}

	repositoryOnlyNameUser := &ur.User{
		Name: user.Name,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockUser.EXPECT().CreateUser(repositoryOnlyNameUser).Return(repositoryUser, nil).Times(1)

		t.NewStep("Check result")
		usr, err := uus.userUsecase.CreateUser(user.Name)
		t.Require().NoError(err)
		t.Require().Equal(user, usr)
	})

	t.WithNewStep("Repository error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockUser.EXPECT().CreateUser(repositoryOnlyNameUser).Return(nil, testError).Times(1)

		t.NewStep("Check result")
		_, err := uus.userUsecase.CreateUser(user.Name)
		t.Require().ErrorIs(err, testError)
	})
}

func (uus *UserUsecaseSuite) TestDeleteUserFunction(t provider.T) {
	t.Title("DeleteUser function of user usecase")
	t.NewStep("Init test data")
	user := &User{
		ID:      1,
		Name:    "User",
		Balance: 30,
	}

	repositoryUser := &ur.User{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockUser.EXPECT().DeleteUser(user.ID).Return(repositoryUser, nil).Times(1)

		t.NewStep("Check result")
		usr, err := uus.userUsecase.DeleteUser(user.ID)
		t.Require().NoError(err)
		t.Require().Equal(user, usr)
	})

	t.WithNewStep("Repository error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockUser.EXPECT().DeleteUser(user.ID).Return(nil, testError).Times(1)

		t.NewStep("Check result")
		_, err := uus.userUsecase.DeleteUser(user.ID)
		t.Require().ErrorIs(err, testError)
	})
}

func (uus *UserUsecaseSuite) TestUpdateUserFunction(t provider.T) {
	t.Title("UpdateUser function of user usecase")
	t.NewStep("Init test data")
	user := &User{
		ID:      1,
		Name:    "User",
		Balance: 30,
	}

	repositoryUser := &ur.User{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}

	repositoryWithoutBalanceUser := &ur.User{
		ID:   user.ID,
		Name: user.Name,
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockUser.EXPECT().UpdateUser(repositoryWithoutBalanceUser).Return(repositoryUser, nil).Times(1)

		t.NewStep("Check result")
		usr, err := uus.userUsecase.UpdateUser(user.ID, user.Name)
		t.Require().NoError(err)
		t.Require().Equal(user, usr)
	})

	t.WithNewStep("Repository error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockUser.EXPECT().UpdateUser(repositoryWithoutBalanceUser).Return(nil, testError).Times(1)

		t.NewStep("Check result")
		_, err := uus.userUsecase.UpdateUser(user.ID, user.Name)
		t.Require().ErrorIs(err, testError)
	})
}

func (uus *UserUsecaseSuite) TestGetUsersFunction(t provider.T) {
	t.Title("GetUsers function of user usecase")
	t.NewStep("Init test data")
	user := &User{
		ID:      1,
		Name:    "User",
		Balance: 30,
	}

	repositoryUser := []ur.User{{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	},
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockUser.EXPECT().GetUsers().Return(repositoryUser, nil).Times(1)

		t.NewStep("Check result")
		usr, err := uus.userUsecase.GetUsers()
		t.Require().NoError(err)
		t.Require().Equal([]User{*user}, usr)
	})

	t.WithNewStep("Repository error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockUser.EXPECT().GetUsers().Return(nil, testError).Times(1)

		t.NewStep("Check result")
		_, err := uus.userUsecase.GetUsers()
		t.Require().ErrorIs(err, testError)
	})
}

func (uus *UserUsecaseSuite) TestGetUserHistoryFunction(t provider.T) {
	t.Title("GetUserHistory function of user usecase")
	t.NewStep("Init test data")
	userId := types.Id(1)
	history := []HistoryRecord{
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

	repositoryHistory := []ur.HistoryRecord{
		{
			Quest: &qr.Quest{
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
		uus.mockUser.EXPECT().GetHistory(userId).Return(repositoryHistory, nil).Times(1)

		t.NewStep("Check result")
		hist, err := uus.userUsecase.GetUserHistory(userId)
		t.Require().NoError(err)
		t.Require().Equal(history, hist)
	})

	t.WithNewStep("Repository error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockUser.EXPECT().GetHistory(userId).Return(nil, testError).Times(1)

		t.NewStep("Check result")
		_, err := uus.userUsecase.GetUserHistory(userId)
		t.Require().ErrorIs(err, testError)
	})
}

func (uus *UserUsecaseSuite) TestApplyQuestsFunction(t provider.T) {
	t.Title("ApplyQuests function of user usecase")
	t.NewStep("Init test data")
	quest := &qu.Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	randomQuest := &qu.Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.RANDOM,
	}

	repositoryQuest := &qr.Quest{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		Cost:        quest.Cost,
		Type:        quest.Type,
	}

	repositoryRandomQuest := &qr.Quest{
		ID:          randomQuest.ID,
		Name:        randomQuest.Name,
		Description: randomQuest.Description,
		Cost:        randomQuest.Cost,
		Type:        randomQuest.Type,
	}

	userId := types.Id(1)

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockQuest.EXPECT().GetQuest(quest.ID).Return(repositoryQuest, nil).Times(1)
		uus.mockUser.EXPECT().HasUser(userId).Return(nil).Times(1)
		uus.mockUser.EXPECT().IsCompletedQuest(&ur.User{ID: userId}, repositoryQuest).Return(nil).Times(1)
		uus.mockUser.EXPECT().ApplyCost(&ur.User{ID: userId}, repositoryQuest).Return(nil).Times(1)

		t.NewStep("Check result")
		err := uus.userUsecase.ApplyQuests(quest.ID, userId)
		t.Require().NoError(err)
	})

	t.WithNewStep("Repository GetQuest method error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockQuest.EXPECT().GetQuest(quest.ID).Return(nil, testError).Times(1)

		t.NewStep("Check result")
		err := uus.userUsecase.ApplyQuests(quest.ID, userId)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Repository HasUser method error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockQuest.EXPECT().GetQuest(quest.ID).Return(repositoryQuest, nil).Times(1)
		uus.mockUser.EXPECT().HasUser(userId).Return(testError).Times(1)

		t.NewStep("Check result")
		err := uus.userUsecase.ApplyQuests(quest.ID, userId)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Repository IsCompletedQuest method error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockQuest.EXPECT().GetQuest(quest.ID).Return(repositoryQuest, nil).Times(1)
		uus.mockUser.EXPECT().HasUser(userId).Return(nil).Times(1)
		uus.mockUser.EXPECT().IsCompletedQuest(&ur.User{ID: userId}, repositoryQuest).Return(testError).Times(1)

		t.NewStep("Check result")
		err := uus.userUsecase.ApplyQuests(quest.ID, userId)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Repository GetQuest method error", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		uus.mockQuest.EXPECT().GetQuest(quest.ID).Return(repositoryQuest, nil).Times(1)
		uus.mockUser.EXPECT().HasUser(userId).Return(nil).Times(1)
		uus.mockUser.EXPECT().IsCompletedQuest(&ur.User{ID: userId}, repositoryQuest).Return(nil).Times(1)
		uus.mockUser.EXPECT().ApplyCost(&ur.User{ID: userId}, repositoryQuest).Return(testError).Times(1)

		t.NewStep("Check result")
		err := uus.userUsecase.ApplyQuests(quest.ID, userId)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Correct random quest failure", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		rnd = rand.New(rand.NewSource(6))
		uus.mockQuest.EXPECT().GetQuest(randomQuest.ID).Return(repositoryRandomQuest, nil).Times(1)
		uus.mockUser.EXPECT().HasUser(userId).Return(nil).Times(1)
		uus.mockUser.EXPECT().IsCompletedQuest(&ur.User{ID: userId}, repositoryRandomQuest).Return(nil).Times(1)
		uus.mockUser.EXPECT().ApplyCost(&ur.User{ID: userId}, repositoryRandomQuest).Return(nil).Times(1)

		t.NewStep("Check result")
		err := uus.userUsecase.ApplyQuests(randomQuest.ID, userId)
		t.Require().ErrorIs(err, QuestNotApplied)
	})
}

func TestRunUserUsecaseSuite(t *testing.T) {
	suite.RunSuite(t, new(UserUsecaseSuite))
}
