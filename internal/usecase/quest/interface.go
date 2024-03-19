package quest

import "vk_quests/internal/pkg/types"

//go:generate mockgen -destination=mocks/usecase.go -package=mu -mock_names=Usecase=QuestUsecase . Usecase

type Usecase interface {
	CreateQuest(quest *Quest) (*Quest, error)
	DeleteQuest(id types.Id) error
	UpdateQuest(id types.Id, quest *UpdateQuest) (*Quest, error)
	GetQuests() ([]Quest, error)
	GetQuest(id types.Id) (*Quest, error)
}
