package quest

import (
	"database/sql"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/pkg/errors"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"vk_quests/internal/pkg/types"
)

var testError = errors.New("test error")

type QuestRepositorySuite struct {
	suite.Suite
	QuestRepository *PostgresQuest
	mock            sqlxmock.Sqlmock
}

func (qrs *QuestRepositorySuite) BeforeEach(t provider.T) {
	db, mock, err := sqlxmock.Newx(sqlxmock.QueryMatcherOption(sqlxmock.QueryMatcherEqual))
	t.Require().NoError(err)
	qrs.QuestRepository = NewPostgresQuest(db)
	qrs.mock = mock
}

func (qrs *QuestRepositorySuite) AfterEach(t provider.T) {
	t.Require().NoError(qrs.mock.ExpectationsWereMet())
}

func (qrs *QuestRepositorySuite) TestCreateFunction(t provider.T) {
	t.Title("CreateQuest function of Quest repository")
	t.NewStep("Init test data")
	quest := &Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	questColumns := []string{
		"id", "name", "description", "cost", "type", "exists",
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(createQuery).
			WithArgs(quest.Name, quest.Description, quest.Cost, quest.Type).
			WillReturnRows(sqlxmock.NewRows(questColumns).
				AddRow(quest.ID, quest.Name, quest.Description, quest.Cost, quest.Type, 0),
			)

		t.NewStep("Check result")
		qst, err := qrs.QuestRepository.CreateQuest(quest)
		t.Require().NoError(err)
		t.Require().EqualValues(quest, qst)
	})

	t.WithNewStep("Conflict name exists execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(createQuery).
			WithArgs(quest.Name, quest.Description, quest.Cost, quest.Type).
			WillReturnRows(sqlxmock.NewRows(questColumns).
				AddRow(quest.ID, quest.Name, quest.Description, quest.Cost, quest.Type, 1),
			)

		t.NewStep("Check result")
		_, err := qrs.QuestRepository.CreateQuest(quest)
		t.Require().ErrorIs(err, ErrorQuestNameAlreadyExists)
	})

	t.WithNewStep("Postgres error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(createQuery).
			WithArgs(quest.Name, quest.Description, quest.Cost, quest.Type).
			WillReturnError(testError)

		t.NewStep("Check result")
		_, err := qrs.QuestRepository.CreateQuest(quest)
		t.Require().ErrorIs(err, testError)
	})
}

func (qrs *QuestRepositorySuite) TestDeleteFunction(t provider.T) {
	t.Title("DeleteQuest function of Quest repository")
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
		qrs.mock.ExpectExec(deleteQuest).
			WithArgs(quest.ID).
			WillReturnResult(sqlxmock.NewResult(0, 1))

		t.NewStep("Check result")
		err := qrs.QuestRepository.DeleteQuest(quest.ID)
		t.Require().NoError(err)
	})

	t.WithNewStep("Postgres error for deleteQuest query", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectExec(deleteQuest).
			WithArgs(quest.ID).WillReturnError(testError)

		t.NewStep("Check result")
		err := qrs.QuestRepository.DeleteQuest(quest.ID)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Row affected error of deleteQuest query", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectExec(deleteQuest).
			WithArgs(quest.ID).
			WillReturnResult(sqlxmock.NewErrorResult(testError))

		t.NewStep("Check result")
		err := qrs.QuestRepository.DeleteQuest(quest.ID)
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Error not found quest", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectExec(deleteQuest).
			WithArgs(quest.ID).
			WillReturnResult(sqlxmock.NewResult(2, 0))

		t.NewStep("Check result")
		err := qrs.QuestRepository.DeleteQuest(quest.ID)
		t.Require().ErrorIs(err, ErrorQuestNotFound)
	})
}

func (qrs *QuestRepositorySuite) TestGetQuestFunction(t provider.T) {
	t.Title("GetQuest function of Quest repository")
	t.NewStep("Init test data")

	quest := &Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	questColumns := []string{
		"id", "name", "description", "cost", "type",
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(getQuest).
			WithArgs(quest.ID).
			WillReturnRows(sqlxmock.NewRows(questColumns).
				AddRow(quest.ID, quest.Name, quest.Description, quest.Cost, quest.Type),
			)

		t.NewStep("Check result")
		qst, err := qrs.QuestRepository.GetQuest(quest.ID)
		t.Require().NoError(err)
		t.Require().EqualValues(quest, qst)
	})

	t.WithNewStep("Error no quests found execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(getQuest).
			WithArgs(quest.ID).
			WillReturnRows(sqlxmock.NewRows(questColumns))

		t.NewStep("Check result")
		_, err := qrs.QuestRepository.GetQuest(quest.ID)
		t.Require().ErrorIs(err, ErrorQuestNotFound)
	})

	t.WithNewStep("Postgres error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(getQuest).
			WithArgs(quest.ID).
			WillReturnError(testError)

		t.NewStep("Check result")
		_, err := qrs.QuestRepository.GetQuest(quest.ID)
		t.Require().ErrorIs(err, testError)
	})
}

func (qrs *QuestRepositorySuite) TestUpdateFunction(t provider.T) {
	t.Title("UpdateQuest function of Quest repository")
	t.NewStep("Init test data")
	quest := &Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	questColumns := []string{
		"id", "name", "description", "cost", "type",
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(updateQuest).
			WithArgs(quest.ID,
				getNullString(&quest.Description),
				sql.NullInt64{Valid: true, Int64: int64(quest.Cost)},
				getNullString((*string)(&quest.Type)),
			).
			WillReturnRows(sqlxmock.NewRows(questColumns).AddRow(
				quest.ID, quest.Name, quest.Description, quest.Cost, quest.Type,
			))

		t.NewStep("Check result")
		updatedQuest, err := qrs.QuestRepository.UpdateQuest(&UpdateQuest{
			ID:          quest.ID,
			Description: &quest.Description,
			Cost:        &quest.Cost,
			Type:        &quest.Type,
		})
		t.Require().NoError(err)
		t.Require().EqualValues(quest, updatedQuest)
	})

	t.WithNewStep("Correct only description execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(updateQuest).
			WithArgs(quest.ID,
				getNullString(&quest.Description),
				sql.NullInt64{Valid: false},
				getNullString(nil),
			).
			WillReturnRows(sqlxmock.NewRows(questColumns).AddRow(
				quest.ID, quest.Name, quest.Description, quest.Cost, quest.Type,
			))

		t.NewStep("Check result")
		updatedQuest, err := qrs.QuestRepository.UpdateQuest(&UpdateQuest{
			ID:          quest.ID,
			Description: &quest.Description,
			Cost:        nil,
			Type:        nil,
		})
		t.Require().NoError(err)
		t.Require().EqualValues(quest, updatedQuest)
	})

	t.WithNewStep("Correct only cost execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(updateQuest).
			WithArgs(quest.ID,
				getNullString(nil),
				sql.NullInt64{Valid: true, Int64: int64(quest.Cost)},
				getNullString(nil),
			).
			WillReturnRows(sqlxmock.NewRows(questColumns).AddRow(
				quest.ID, quest.Name, quest.Description, quest.Cost, quest.Type,
			))

		t.NewStep("Check result")
		updatedQuest, err := qrs.QuestRepository.UpdateQuest(&UpdateQuest{
			ID:          quest.ID,
			Description: nil,
			Cost:        &quest.Cost,
			Type:        nil,
		})
		t.Require().NoError(err)
		t.Require().EqualValues(quest, updatedQuest)
	})

	t.WithNewStep("Correct only type execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(updateQuest).
			WithArgs(quest.ID,
				getNullString(nil),
				sql.NullInt64{Valid: false},
				getNullString((*string)(&quest.Type)),
			).
			WillReturnRows(sqlxmock.NewRows(questColumns).AddRow(
				quest.ID, quest.Name, quest.Description, quest.Cost, quest.Type,
			))

		t.NewStep("Check result")
		updatedQuest, err := qrs.QuestRepository.UpdateQuest(&UpdateQuest{
			ID:          quest.ID,
			Description: nil,
			Cost:        nil,
			Type:        &quest.Type,
		})
		t.Require().NoError(err)
		t.Require().EqualValues(quest, updatedQuest)
	})

	t.WithNewStep("Error quest not found execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(updateQuest).
			WithArgs(quest.ID,
				getNullString(&quest.Description),
				sql.NullInt64{Valid: true, Int64: int64(quest.Cost)},
				getNullString((*string)(&quest.Type)),
			).
			WillReturnRows(sqlxmock.NewRows(questColumns))

		t.NewStep("Check result")
		_, err := qrs.QuestRepository.UpdateQuest(&UpdateQuest{
			ID:          quest.ID,
			Description: &quest.Description,
			Cost:        &quest.Cost,
			Type:        &quest.Type,
		})
		t.Require().ErrorIs(err, ErrorQuestNotFound)
	})

	t.WithNewStep("Postgres error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(updateQuest).
			WithArgs(quest.ID,
				getNullString(&quest.Description),
				sql.NullInt64{Valid: true, Int64: int64(quest.Cost)},
				getNullString((*string)(&quest.Type)),
			).WillReturnError(testError)

		t.NewStep("Check result")
		_, err := qrs.QuestRepository.UpdateQuest(&UpdateQuest{
			ID:          quest.ID,
			Description: &quest.Description,
			Cost:        &quest.Cost,
			Type:        &quest.Type,
		})
		t.Require().ErrorIs(err, testError)
	})
}

func (qrs *QuestRepositorySuite) TestGetQuestsFunction(t provider.T) {
	t.Title("GetQuests function of Quest repository")
	t.NewStep("Init test data")

	quest := &Quest{
		ID:          1,
		Name:        "Quest",
		Description: "good Quest",
		Cost:        10,
		Type:        types.USUAL,
	}

	questColumns := []string{
		"id", "name", "description", "cost", "type",
	}

	questRows := func() *sqlxmock.Rows {
		return sqlxmock.NewRows(questColumns).
			AddRow(quest.ID, quest.Name, quest.Description, quest.Cost, quest.Type).
			AddRow(quest.ID, quest.Name, quest.Description, quest.Cost, quest.Type).
			AddRow(quest.ID, quest.Name, quest.Description, quest.Cost, quest.Type)
	}

	t.WithNewStep("Correct execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(getQuests).WillReturnRows(questRows())

		t.NewStep("Check result")
		qsts, err := qrs.QuestRepository.GetQuests()
		t.Require().NoError(err)
		t.Require().EqualValues([]Quest{
			*quest, *quest, *quest,
		}, qsts)
	})

	t.WithNewStep("Postgres error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(getQuests).WillReturnError(testError)

		t.NewStep("Check result")
		_, err := qrs.QuestRepository.GetQuests()
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Row error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(getQuests).WillReturnRows(questRows().RowError(1, testError))

		t.NewStep("Check result")
		_, err := qrs.QuestRepository.GetQuests()
		t.Require().ErrorIs(err, testError)
	})

	t.WithNewStep("Scan error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(getQuests).WillReturnRows(questRows().AddRow(1, 1, 1, 1, 1))

		t.NewStep("Check result")
		_, err := qrs.QuestRepository.GetQuests()
		t.Require().Error(err)
	})

	t.WithNewStep("Close row error execute", func(t provider.StepCtx) {
		t.NewStep("Init mock")
		qrs.mock.ExpectQuery(getQuests).WillReturnRows(questRows().CloseError(testError))

		t.NewStep("Check result")
		_, err := qrs.QuestRepository.GetQuests()
		t.Require().ErrorIs(err, testError)
	})
}

func TestRunQuestRepositorySuite(t *testing.T) {
	suite.RunSuite(t, new(QuestRepositorySuite))
}
