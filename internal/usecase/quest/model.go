package quest

import (
	"vk_quests/internal/pkg/types"
	"vk_quests/internal/repository/quest"
)

type Quest struct {
	ID          types.Id
	Name        string
	Description string
	Cost        types.Cost
	Type        types.QuestType
}

func FromRepQuest(q *quest.Quest) *Quest {
	if q == nil {
		return nil
	}

	return &Quest{
		ID:          q.ID,
		Name:        q.Name,
		Description: q.Description,
		Cost:        q.Cost,
		Type:        q.Type,
	}
}

type UpdateQuest struct {
	Description *string
	Cost        *types.Cost
	Type        *types.QuestType
}

func (uq *UpdateQuest) ToRepUpdateQuest(id types.Id) *quest.UpdateQuest {
	return &quest.UpdateQuest{
		ID:          id,
		Description: uq.Description,
		Cost:        uq.Cost,
		Type:        uq.Type,
	}
}
