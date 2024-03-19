package quest

import (
	"github.com/pkg/errors"
	"vk_quests/internal/pkg/types"
)

var (
	ErrorQuestNotFound          = errors.New("quest with id not found")
	ErrorQuestNameAlreadyExists = errors.New("quest with name already exists")
)

//go:generate mockgen -destination=mocks/repository.go -package=mr -mock_names=Repository=QuestRepository . Repository

type Repository interface {
	// CreateQuest
	// Returns Error:
	//   - SQLError
	//   - ErrorQuestNameAlreadyExists
	CreateQuest(quest *Quest) (*Quest, error)

	// UpdateQuest
	// Returns Error:
	//   - SQLError
	//   - ErrorQuestNotFound
	UpdateQuest(quest *UpdateQuest) (*Quest, error)

	// DeleteQuest
	// Returns Error:
	//   - SQLError
	//   - ErrorQuestNotFound
	DeleteQuest(id types.Id) error

	// GetQuests
	// Returns Error:
	//   - SQLError
	GetQuests() ([]Quest, error)

	// GetQuest
	// Returns Error:
	//   - ErrorQuestNotFound
	GetQuest(id types.Id) (*Quest, error)
}
