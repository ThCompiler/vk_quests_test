package quest

import (
	"vk_quests/internal/pkg/types"
)

type Quest struct {
	ID          types.Id
	Name        string
	Description string
	Cost        types.Cost
	Type        types.QuestType
}

type UpdateQuest struct {
	ID          types.Id
	Description *string
	Cost        *types.Cost
	Type        *types.QuestType
}
