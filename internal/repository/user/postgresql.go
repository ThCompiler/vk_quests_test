package user

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"vk_quests/internal/pkg/types"
	qr "vk_quests/internal/repository/quest"
)

const (
	createQuery = `
		INSERT INTO users (name)
		VALUES ($1)
		RETURNING id, name, balance
	`

	deleteUser = `
		DELETE FROM users WHERE id = $1 RETURNING id, name, balance
	`

	updateUser = `
		UPDATE users SET name = $2 WHERE id = $1
			RETURNING id, name, balance
	`

	getUsers = `
		SELECT id, name, balance FROM users
	`

	applyCost = `
		UPDATE users SET balance = balance + $2 WHERE id = $1 RETURNING id
	`

	createHistory = `
		INSERT INTO balance_history (user_id, quest_id, balance) 
		SELECT $1, $2, users.balance FROM users WHERE id = $1
	`

	getHistory = `
		SELECT quests.id, quests.name, quests.description, quests.cost, quests.type, created, balance 
		FROM balance_history LEFT JOIN quests ON (balance_history.quest_id = quests.id)
		WHERE user_id = $1
	`

	getCompleteQuest = `
		SELECT quest_id FROM balance_history WHERE user_id = $1 and quest_id = $2
	`

	hasUser = `
		SELECT id FROM users WHERE id = $1
	`
)

type PostgresUser struct {
	db *sqlx.DB
}

func NewPostgresUser(db *sqlx.DB) *PostgresUser {
	return &PostgresUser{
		db: db,
	}
}

func (pu *PostgresUser) CreateUser(user *User) (*User, error) {
	newUser := &User{}
	if err := pu.db.QueryRowx(createQuery, user.Name).
		Scan(
			&newUser.ID,
			&newUser.Name,
			&newUser.Balance,
		); err != nil {
		return nil, errors.Wrap(err, "can't create user")
	}

	return newUser, nil
}

func (pu *PostgresUser) UpdateUser(user *User) (*User, error) {
	updatedUser := &User{}
	if err := pu.db.QueryRowx(updateUser, user.ID, user.Name).
		Scan(
			&updatedUser.ID,
			&updatedUser.Name,
			&updatedUser.Balance,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorUserNotFound
		}
		return nil, errors.Wrapf(err, "can't update user with id %d", user.ID)
	}

	return updatedUser, nil
}

func (pu *PostgresUser) DeleteUser(id types.Id) (*User, error) {
	deletedUser := &User{}
	if err := pu.db.QueryRowx(deleteUser, id).
		Scan(
			&deletedUser.ID,
			&deletedUser.Name,
			&deletedUser.Balance,
		); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorUserNotFound
		}
		return nil, errors.Wrapf(err, "can't delete user with id %d", id)
	}

	return deletedUser, nil
}

func (pu *PostgresUser) GetUsers() ([]User, error) {
	rows, err := pu.db.Queryx(getUsers)
	if err != nil {
		return nil, errors.Wrap(err, "can't execute get users query")
	}

	users := make([]User, 0)

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Balance,
		)
		if err != nil {
			return nil, errors.Wrap(err, "can't scan get users query result")
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "can't end scan get users query result")
	}

	return users, nil
}

func (pu *PostgresUser) HasUser(userId types.Id) error {
	id := types.Id(0)
	if err := pu.db.QueryRowx(hasUser, userId).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorUserNotFound
		}
		return errors.Wrapf(err, "can't check user with id %d", userId)
	}

	return nil
}

func (pu *PostgresUser) GetHistory(id types.Id) ([]HistoryRecord, error) {
	rows, err := pu.db.Queryx(getHistory, id)
	if err != nil {
		return nil, errors.Wrapf(err, "can't execute get history query for user with id %d", id)
	}

	history := make([]HistoryRecord, 0)

	for rows.Next() {
		var record HistoryRecord

		questId := sql.Null[types.Id]{}
		name := sql.NullString{}
		description := sql.NullString{}
		cost := sql.Null[types.Cost]{}
		tp := sql.NullString{}

		err := rows.Scan(
			&questId,
			&name,
			&description,
			&cost,
			&tp,
			&record.Created,
			&record.Balance,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "can't scan get history query result for user with id %d", id)
		}

		record.Quest = nil
		if questId.Valid && name.Valid && description.Valid && cost.Valid && tp.Valid {
			record.Quest = &qr.Quest{
				ID:          questId.V,
				Name:        name.String,
				Description: description.String,
				Cost:        cost.V,
				Type:        types.QuestType(tp.String),
			}
		}

		history = append(history, record)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "can't end scan get history query result for user with id %d", id)
	}

	return history, nil
}

func (pu *PostgresUser) IsCompletedQuest(user *User, quest *qr.Quest) error {
	questId := types.Id(0)
	if err := pu.db.QueryRowx(getCompleteQuest, user.ID, quest.ID).Scan(&questId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return errors.Wrapf(err, "can't get quest for user with id %d", user.ID)
	}

	return ErrorUserAlreadyCompleteQuest
}

func (pu *PostgresUser) ApplyCost(user *User, quest *qr.Quest) error {
	tx, err := pu.db.Beginx()
	if err != nil {
		return errors.Wrapf(err,
			"can't begin transaction for apply cost to user with id %d and quest id %d", user.ID, quest.ID)
	}

	id := types.Id(0)
	if err := tx.QueryRowx(applyCost, user.ID, quest.Cost).Scan(&id); err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorUserNotFound
		}
		return errors.Wrapf(err, "can't apply cost to user with id %d and quest id %d", user.ID, quest.ID)
	}

	if _, err := tx.Exec(createHistory, user.ID, quest.ID); err != nil {
		_ = tx.Rollback()
		return errors.Wrapf(
			checkConflictError(err),
			"can't store history for user with id %d and quest id %d", user.ID, quest.ID,
		)
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err,
			"can't commit transaction for apply cost to user with id %d and quest id %d", user.ID, quest.ID)
	}

	return nil
}

const (
	uniqueConflictCode   = "23505"
	uniqueConstraintName = "quest_unique"

	foreignKeyConflictCode = "23503"
	userIdConstraintName   = "balance_history_user_id_fkey"
	questIdConstraintName  = "balance_history_quest_id_fkey"
)

func checkConflictError(err error) error {
	var e *pq.Error

	switch {
	case errors.As(err, &e):
		return checkPgError(e)
	default:
		return err
	}
}

func checkPgError(err *pq.Error) error {
	switch err.Code {
	case foreignKeyConflictCode:
		if err.Constraint == questIdConstraintName {
			return qr.ErrorQuestNotFound
		} else if err.Constraint == userIdConstraintName {
			return ErrorUserNotFound
		}
	case uniqueConflictCode:
		if err.Constraint == uniqueConstraintName {
			return ErrorUserAlreadyCompleteQuest
		}
	}
	return err
}
