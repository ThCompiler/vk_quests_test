package user

import (
	"github.com/pkg/errors"
	"vk_quests/internal/pkg/types"
)

//go:generate mockgen -destination=mocks/usecase.go -package=mu -mock_names=Usecase=UserUsecase . Usecase

var QuestNotApplied = errors.New("quest not applied")

type Usecase interface {
	CreateUser(name string) (*User, error)
	DeleteUser(id types.Id) (*User, error)
	UpdateUser(id types.Id, name string) (*User, error)
	GetUsers() ([]User, error)
	GetUserHistory(id types.Id) ([]HistoryRecord, error)
	ApplyQuests(questId, userId types.Id) error
}
