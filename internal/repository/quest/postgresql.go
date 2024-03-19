package quest

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"vk_quests/internal/pkg/types"
)

const (
	createQuery = `
		WITH sel AS (
				SELECT id, name, description, cost, type
				FROM quests
				WHERE name = $1 LIMIT 1
		), ins as (
			INSERT INTO quests (name, description, cost, type)
				SELECT $1, $2, $3, $4
			    WHERE not exists (select 1 from sel)
			RETURNING id, name, description, cost, type
		)
		SELECT id, name, description, cost, type, 0
		FROM ins
		UNION ALL
		SELECT id, name, description, cost, type, 1
		FROM sel
	`

	deleteQuest = `
		DELETE FROM quests WHERE id = $1
	`

	updateQuest = `
		UPDATE quests SET description = upd_quest.upd_description, 
		                 cost = upd_quest.upd_cost, type = upd_quest.upd_type
			FROM (
				SELECT COALESCE($2, quests.description) as upd_description, 
					   COALESCE($3, quests.cost) as upd_cost,
					   COALESCE($4, quests.type) as upd_type 
				FROM quests WHERE id = $1
			) as upd_quest
			WHERE id = $1
			RETURNING id, name, description, cost, type
	`

	getQuests = `
		SELECT id, name, description, cost, type FROM quests
	`

	getQuest = `
		SELECT id, name, description, cost, type FROM quests WHERE id = $1
	`
)

type PostgresQuest struct {
	db *sqlx.DB
}

func NewPostgresQuest(db *sqlx.DB) *PostgresQuest {
	return &PostgresQuest{
		db: db,
	}
}

var _ = Repository(&PostgresQuest{})

func (pt *PostgresQuest) CreateQuest(quest *Quest) (*Quest, error) {
	newQuest := &Quest{}
	exists := 0
	if err := pt.db.QueryRowx(createQuery, quest.Name, quest.Description, quest.Cost, quest.Type).
		Scan(
			&newQuest.ID,
			&newQuest.Name,
			&newQuest.Description,
			&newQuest.Cost,
			&newQuest.Type,
			&exists,
		); err != nil {
		return nil, errors.Wrap(err, "can't create quest")
	}

	if exists == 1 {
		return newQuest, ErrorQuestNameAlreadyExists
	}

	return newQuest, nil
}

func getNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: *value}
}

func (pt *PostgresQuest) UpdateQuest(quest *UpdateQuest) (*Quest, error) {
	description := getNullString(quest.Description)
	tp := getNullString((*string)(quest.Type))

	cost := sql.NullInt64{Valid: false}
	if quest.Cost != nil {
		cost = sql.NullInt64{Valid: true, Int64: int64(*quest.Cost)}
	}

	updatedQuest := &Quest{}
	if err := pt.db.QueryRowx(updateQuest, quest.ID, description, cost, tp).
		Scan(
			&updatedQuest.ID,
			&updatedQuest.Name,
			&updatedQuest.Description,
			&updatedQuest.Cost,
			&updatedQuest.Type,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorQuestNotFound
		}

		return nil, errors.Wrapf(err, "can't update quest with id %d", quest.ID)
	}

	return updatedQuest, nil
}

func (pt *PostgresQuest) DeleteQuest(id types.Id) error {
	res, err := pt.db.Exec(deleteQuest, id)
	if err != nil {
		return errors.Wrapf(err, "can't execute deleting query for quest %d", id)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "can't get number affected rows of deleting query for quest %d", id)
	}

	if n != 1 {
		return errors.Wrapf(ErrorQuestNotFound, "with id %d", id)
	}

	return nil
}

func (pt *PostgresQuest) GetQuests() ([]Quest, error) {
	rows, err := pt.db.Queryx(getQuests)
	if err != nil {
		return nil, errors.Wrap(err, "can't execute get quests query")
	}

	quests := make([]Quest, 0)

	for rows.Next() {
		var quest Quest

		err := rows.Scan(
			&quest.ID,
			&quest.Name,
			&quest.Description,
			&quest.Cost,
			&quest.Type,
		)
		if err != nil {
			return nil, errors.Wrap(err, "can't scan get quests query result")
		}

		quests = append(quests, quest)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "can't end scan get quests query result")
	}

	return quests, nil
}

func (pt *PostgresQuest) GetQuest(id types.Id) (*Quest, error) {
	quest := &Quest{}
	if err := pt.db.QueryRowx(getQuest, id).
		Scan(
			&quest.ID,
			&quest.Name,
			&quest.Description,
			&quest.Cost,
			&quest.Type,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorQuestNotFound
		}

		return nil, errors.Wrapf(err, "can't get quest with id %d", quest.ID)
	}

	return quest, nil
}
