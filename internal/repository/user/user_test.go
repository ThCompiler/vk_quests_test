package user

import (
	"github.com/lib/pq"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"vk_quests/internal/pkg/types"
	qr "vk_quests/internal/repository/quest"
)

var testError = errors.New("test error")

type UserRepositorySuite struct {
	suite.Suite
	userRepository *PostgresUser
	mock           sqlxmock.Sqlmock
}

func (urs *UserRepositorySuite) BeforeEach(t provider.T) {
	db, mock, err := sqlxmock.Newx(sqlxmock.QueryMatcherOption(sqlxmock.QueryMatcherEqual))
	t.Require().NoError(err)
	urs.userRepository = NewPostgresUser(db)
	urs.mock = mock
}

func (urs *UserRepositorySuite) AfterEach(t provider.T) {
	t.Require().NoError(urs.mock.ExpectationsWereMet())
}

func (urs *UserRepositorySuite) TestCreateFunction(t provider.T) {
	t.Title("CreateUser function of User repository")
	t.NewStep("Init test data")
	user := &User{
		ID:      1,
		Name:    "user",
		Balance: 200,
	}

	userColumns := []string{
		"id", "name", "balance",
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(createQuery).
			WithArgs(user.Name).
			WillReturnRows(sqlxmock.NewRows(userColumns).
				AddRow(user.ID, user.Name, user.Balance),
			)

		t.NewStep("Check result")
		usr, err := urs.userRepository.CreateUser(user)
		t.Require().NoError(err)
		t.Require().EqualValues(user, usr)
	})

	t.WithNewStep("Postgres error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(createQuery).
			WithArgs(user.Name).
			WillReturnError(testError)

		t.NewStep("Check result")
		_, err := urs.userRepository.CreateUser(user)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Empty result of execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(createQuery).
			WithArgs(user.Name).
			WillReturnRows(sqlxmock.NewRows(userColumns))

		t.NewStep("Check result")
		_, err := urs.userRepository.CreateUser(user)
		t.Require().Error(err)
	})
}

func (urs *UserRepositorySuite) TestDeleteFunction(t provider.T) {
	t.Title("DeleteUser function of User repository")
	t.NewStep("Init test data")
	user := &User{
		ID:      1,
		Name:    "user",
		Balance: 10,
	}

	userColumns := []string{
		"id", "name", "balance",
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(deleteUser).
			WithArgs(user.ID).
			WillReturnRows(sqlxmock.NewRows(userColumns).
				AddRow(user.ID, user.Name, user.Balance),
			)

		t.NewStep("Check result")
		usr, err := urs.userRepository.DeleteUser(user.ID)
		t.Require().NoError(err)
		t.Require().EqualValues(user, usr)
	})

	t.WithNewStep("Postgres error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(deleteUser).
			WithArgs(user.ID).
			WillReturnError(testError)

		t.NewStep("Check result")
		_, err := urs.userRepository.DeleteUser(user.ID)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Empty result of execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(deleteUser).
			WithArgs(user.ID).
			WillReturnRows(sqlxmock.NewRows(userColumns))

		t.NewStep("Check result")
		_, err := urs.userRepository.DeleteUser(user.ID)
		t.Require().ErrorIs(err, ErrorUserNotFound)
	})
}

func (urs *UserRepositorySuite) TestGetFunction(t provider.T) {
	t.Title("GetUsers function of User repository")
	t.NewStep("Init test data")
	user := &User{
		ID:      1,
		Name:    "user",
		Balance: 20,
	}

	userColumns := []string{
		"id", "name", "balance",
	}

	usersRows := func() *sqlxmock.Rows {
		return sqlxmock.NewRows(userColumns).
			AddRow(user.ID, user.Name, user.Balance).
			AddRow(user.ID, user.Name, user.Balance).
			AddRow(user.ID, user.Name, user.Balance)
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getUsers).WillReturnRows(usersRows())

		t.NewStep("Check result")
		users, err := urs.userRepository.GetUsers()
		t.Require().NoError(err)
		t.Require().EqualValues([]User{*user, *user, *user}, users)
	})

	t.WithNewStep("Postgres error query", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getUsers).WillReturnError(testError)

		t.NewStep("Check result")
		_, err := urs.userRepository.GetUsers()
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Rows error query", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getUsers).WillReturnRows(usersRows().RowError(1, testError))

		t.NewStep("Check result")
		_, err := urs.userRepository.GetUsers()
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Incorrect field in row of getUsers query", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getUsers).WillReturnRows(usersRows().AddRow(1, 1, "top"))

		t.NewStep("Check result")
		_, err := urs.userRepository.GetUsers()
		t.Require().Error(err)
	})

	t.WithNewStep("Rows close error on getUsers query", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getUsers).WillReturnRows(usersRows().CloseError(testError))

		t.NewStep("Check result")
		_, err := urs.userRepository.GetUsers()
		t.Require().ErrorIs(err, testError)
	})
}

func (urs *UserRepositorySuite) TestHasUserFunction(t provider.T) {
	t.Title("HasUser function of User repository")
	t.NewStep("Init test data")
	userId := types.Id(1)

	userColumns := []string{
		"id",
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(hasUser).
			WithArgs(userId).
			WillReturnRows(sqlxmock.NewRows(userColumns).
				AddRow(userId),
			)

		t.NewStep("Check result")
		err := urs.userRepository.HasUser(userId)
		t.Require().NoError(err)
	})

	t.WithNewStep("Postgres error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(hasUser).
			WithArgs(userId).
			WillReturnError(testError)

		t.NewStep("Check result")
		err := urs.userRepository.HasUser(userId)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Empty result of execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(hasUser).
			WithArgs(userId).
			WillReturnRows(sqlxmock.NewRows(userColumns))

		t.NewStep("Check result")
		err := urs.userRepository.HasUser(userId)
		t.Require().ErrorIs(err, ErrorUserNotFound)
	})
}

func (urs *UserRepositorySuite) TestUpdateFunction(t provider.T) {
	t.Title("UpdateUser function of User repository")
	t.NewStep("Init test data")
	user := &User{
		ID:      1,
		Name:    "user",
		Balance: 10,
	}

	userColumns := []string{
		"id", "name", "balance",
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(updateUser).
			WithArgs(user.ID, user.Name).
			WillReturnRows(sqlxmock.NewRows(userColumns).
				AddRow(user.ID, user.Name, user.Balance),
			)

		t.NewStep("Check result")
		usr, err := urs.userRepository.UpdateUser(user)
		t.Require().NoError(err)
		t.Require().EqualValues(user, usr)
	})

	t.WithNewStep("Postgres error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(updateUser).
			WithArgs(user.ID, user.Name).
			WillReturnError(testError)

		t.NewStep("Check result")
		_, err := urs.userRepository.UpdateUser(user)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Empty result of execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(updateUser).
			WithArgs(user.ID, user.Name).
			WillReturnRows(sqlxmock.NewRows(userColumns))

		t.NewStep("Check result")
		_, err := urs.userRepository.UpdateUser(user)
		t.Require().ErrorIs(err, ErrorUserNotFound)
	})
}

func (urs *UserRepositorySuite) TestGetHistoryFunction(t provider.T) {
	t.Title("GetHistory function of User repository")
	t.NewStep("Init test data")
	userId := types.Id(1)

	historyColumns := []string{
		"id", "name", "description", "cost", "type", "created", "balance",
	}

	resHistory := []HistoryRecord{
		{
			Quest: &qr.Quest{
				ID:   1,
				Name: "Not null",
			},
			Balance: 30,
		},
		{
			Quest:   nil,
			Balance: 25,
		},
		{
			Quest:   nil,
			Balance: 26,
		},
	}

	historyRows := func() *sqlxmock.Rows {
		return sqlxmock.NewRows(historyColumns).
			AddRow(resHistory[0].Quest.ID, resHistory[0].Quest.Name, resHistory[0].Quest.Description,
				resHistory[0].Quest.Cost, resHistory[0].Quest.Type, resHistory[0].Created.Time, resHistory[0].Balance).
			AddRow(nil, nil, nil, nil, nil, resHistory[1].Created.Time, resHistory[1].Balance).
			AddRow(resHistory[0].Quest.ID, nil, nil,
				resHistory[0].Quest.Cost, nil, resHistory[2].Created.Time, resHistory[2].Balance)
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getHistory).WillReturnRows(historyRows())

		t.NewStep("Check result")
		hist, err := urs.userRepository.GetHistory(userId)
		t.Require().NoError(err)
		t.Require().EqualValues(resHistory, hist)
	})

	t.WithNewStep("Postgres error query", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getHistory).WillReturnError(testError)

		t.NewStep("Check result")
		_, err := urs.userRepository.GetHistory(userId)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Rows error query", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getHistory).WillReturnRows(historyRows().RowError(1, testError))

		t.NewStep("Check result")
		_, err := urs.userRepository.GetHistory(userId)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Incorrect field in row of getUsers query", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getHistory).WillReturnRows(historyRows().AddRow(1, 1, 1, 1, 1, 1, 1))

		t.NewStep("Check result")
		_, err := urs.userRepository.GetHistory(userId)
		t.Require().Error(err)
	})

	t.WithNewStep("Rows close error on getUsers query", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getHistory).WillReturnRows(historyRows().CloseError(testError))

		t.NewStep("Check result")
		_, err := urs.userRepository.GetHistory(userId)
		t.Require().ErrorIs(err, testError)
	})
}

func (urs *UserRepositorySuite) TestIsCompletedQuestFunction(t provider.T) {
	t.Title("IsCompletedQuest function of User repository")
	t.NewStep("Init test data")

	user := &User{
		ID:      1,
		Name:    "user",
		Balance: 10,
	}

	quest := &qr.Quest{
		ID:   2,
		Name: "Quest",
	}

	userColumns := []string{
		"id",
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getCompleteQuest).
			WithArgs(user.ID, quest.ID).
			WillReturnRows(sqlxmock.NewRows(userColumns))

		t.NewStep("Check result")
		err := urs.userRepository.IsCompletedQuest(user, quest)
		t.Require().NoError(err)
	})

	t.WithNewStep("Postgres error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getCompleteQuest).
			WithArgs(user.ID, quest.ID).
			WillReturnError(testError)

		t.NewStep("Check result")
		err := urs.userRepository.IsCompletedQuest(user, quest)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Quest already completed for user error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectQuery(getCompleteQuest).
			WithArgs(user.ID, quest.ID).
			WillReturnRows(sqlxmock.NewRows(userColumns).AddRow(quest.ID))

		t.NewStep("Check result")
		err := urs.userRepository.IsCompletedQuest(user, quest)
		t.Require().ErrorIs(err, ErrorUserAlreadyCompleteQuest)
	})
}

func (urs *UserRepositorySuite) TestApplyCostFunction(t provider.T) {
	t.Title("ApplyCost function of User repository")
	t.NewStep("Init test data")

	user := &User{
		ID:      1,
		Name:    "user",
		Balance: 10,
	}

	quest := &qr.Quest{
		ID:   2,
		Name: "Quest",
	}

	userColumns := []string{
		"id",
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectBegin()
		urs.mock.ExpectQuery(applyCost).
			WithArgs(user.ID, quest.Cost).
			WillReturnRows(sqlxmock.NewRows(userColumns).AddRow(user.ID))
		urs.mock.ExpectExec(createHistory).
			WithArgs(user.ID, quest.ID).
			WillReturnResult(sqlxmock.NewResult(0, 1))
		urs.mock.ExpectCommit()

		t.NewStep("Check result")
		err := urs.userRepository.ApplyCost(user, quest)
		t.Require().NoError(err)
	})

	t.WithNewStep("Postgres error create transaction execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectBegin().WillReturnError(testError)

		t.NewStep("Check result")
		err := urs.userRepository.ApplyCost(user, quest)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Postgres error on applyCost query execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectBegin()
		urs.mock.ExpectQuery(applyCost).
			WithArgs(user.ID, quest.Cost).
			WillReturnError(testError)
		urs.mock.ExpectRollback()

		t.NewStep("Check result")
		err := urs.userRepository.ApplyCost(user, quest)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("No user found on applyCost query execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectBegin()
		urs.mock.ExpectQuery(applyCost).
			WithArgs(user.ID, quest.Cost).
			WillReturnRows(sqlxmock.NewRows(userColumns))
		urs.mock.ExpectRollback()

		t.NewStep("Check result")
		err := urs.userRepository.ApplyCost(user, quest)
		t.Require().ErrorIs(err, ErrorUserNotFound)
	})

	t.WithNewStep("Postgres error on createHistory query execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectBegin()
		urs.mock.ExpectQuery(applyCost).
			WithArgs(user.ID, quest.Cost).
			WillReturnRows(sqlxmock.NewRows(userColumns).AddRow(user.ID))
		urs.mock.ExpectExec(createHistory).
			WithArgs(user.ID, quest.ID).
			WillReturnError(testError)
		urs.mock.ExpectRollback()

		t.NewStep("Check result")
		err := urs.userRepository.ApplyCost(user, quest)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("No user found on createHistory query execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectBegin()
		urs.mock.ExpectQuery(applyCost).
			WithArgs(user.ID, quest.Cost).
			WillReturnRows(sqlxmock.NewRows(userColumns).AddRow(user.ID))
		urs.mock.ExpectExec(createHistory).
			WithArgs(user.ID, quest.ID).
			WillReturnError(&pq.Error{Code: foreignKeyConflictCode, Constraint: userIdConstraintName})
		urs.mock.ExpectRollback()

		t.NewStep("Check result")
		err := urs.userRepository.ApplyCost(user, quest)
		t.Require().ErrorIs(err, ErrorUserNotFound)
	})

	t.WithNewStep("No quest found on createHistory query execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectBegin()
		urs.mock.ExpectQuery(applyCost).
			WithArgs(user.ID, quest.Cost).
			WillReturnRows(sqlxmock.NewRows(userColumns).AddRow(user.ID))
		urs.mock.ExpectExec(createHistory).
			WithArgs(user.ID, quest.ID).
			WillReturnError(&pq.Error{Code: foreignKeyConflictCode, Constraint: questIdConstraintName})
		urs.mock.ExpectRollback()

		t.NewStep("Check result")
		err := urs.userRepository.ApplyCost(user, quest)
		t.Require().ErrorIs(err, qr.ErrorQuestNotFound)
	})

	t.WithNewStep("Quest already completed for user on createHistory query execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectBegin()
		urs.mock.ExpectQuery(applyCost).
			WithArgs(user.ID, quest.Cost).
			WillReturnRows(sqlxmock.NewRows(userColumns).AddRow(user.ID))
		urs.mock.ExpectExec(createHistory).
			WithArgs(user.ID, quest.ID).
			WillReturnError(&pq.Error{Code: uniqueConflictCode, Constraint: uniqueConstraintName})
		urs.mock.ExpectRollback()

		t.NewStep("Check result")
		err := urs.userRepository.ApplyCost(user, quest)
		t.Require().ErrorIs(err, ErrorUserAlreadyCompleteQuest)
	})

	t.WithNewStep("Postgres error commit transaction execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		urs.mock.ExpectBegin()
		urs.mock.ExpectQuery(applyCost).
			WithArgs(user.ID, quest.Cost).
			WillReturnRows(sqlxmock.NewRows(userColumns).AddRow(user.ID))
		urs.mock.ExpectExec(createHistory).
			WithArgs(user.ID, quest.ID).
			WillReturnResult(sqlxmock.NewResult(0, 1))
		urs.mock.ExpectCommit().WillReturnError(testError)

		t.NewStep("Check result")
		err := urs.userRepository.ApplyCost(user, quest)
		t.Require().ErrorIs(err, testError)
	})
}

func TestRunUserRepositorySuite(t *testing.T) {
	suite.RunSuite(t, new(UserRepositorySuite))
}
