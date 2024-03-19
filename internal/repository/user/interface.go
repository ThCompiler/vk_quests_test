package user

import (
	"github.com/pkg/errors"
	"vk_quests/internal/pkg/types"
	"vk_quests/internal/repository/quest"
)

var (
	ErrorUserNotFound             = errors.New("user with id not found")
	ErrorUserAlreadyCompleteQuest = errors.New("user already complete quest")
)

//go:generate mockgen -destination=mocks/repository.go -package=mr -mock_names=Repository=UserRepository . Repository

type Repository interface {
	// CreateUser
	// Returns Error:
	//   - SQLError
	CreateUser(user *User) (*User, error)

	// UpdateUser
	// Returns Error:
	//   - SQLError
	//   - ErrorUserNotFound
	UpdateUser(user *User) (*User, error)

	// DeleteUser
	// Returns Error:
	//   - SQLError
	//   - ErrorUserNotFound
	DeleteUser(id types.Id) (*User, error)

	// GetUsers
	// Returns Error:
	//   - SQLError
	GetUsers() ([]User, error)

	// HasUser
	// Returns Error:
	//   - SQLError
	//   - ErrorUserNotFound
	HasUser(userId types.Id) error

	// GetHistory
	// Returns Error:
	//   - SQLError
	//   - ErrorUserNotFound
	GetHistory(id types.Id) ([]HistoryRecord, error)

	// ApplyCost
	// Returns Error:
	//   - SQLError
	//   - ErrorUserNotFound
	//   - quest.ErrorQuestNotFound
	//   - ErrorUserAlreadyCompleteQuest
	ApplyCost(user *User, quest *quest.Quest) error

	// IsCompletedQuest
	// Returns Error:
	//   - SQLError
	//   - ErrorUserNotFound
	//   - quest.ErrorQuestNotFound
	IsCompletedQuest(user *User, quest *quest.Quest) error
}
