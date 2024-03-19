package response

import (
	"vk_quests/internal/pkg/types"
	qu "vk_quests/internal/usecase/quest"
	"vk_quests/pkg/slices"
)

type Quest struct {
	ID          types.Id        `json:"id" swaggertype:"integer" format:"uint64" example:"5"`
	Name        string          `json:"name" swaggertype:"string" example:"Task"`
	Description string          `json:"description" swaggertype:"string" example:"Random quest"`
	Cost        types.Cost      `json:"cost" swaggertype:"integer" format:"uint8" example:"9" minimum:"0" maximum:"1000"`
	Type        types.QuestType `json:"type" swaggertype:"string" enums:"usual,random" example:"random"`
}

func FromUsQuests(quests []qu.Quest) []Quest {
	return slices.Map(quests, func(quest qu.Quest) Quest {
		return *FromUsQuest(&quest)
	})
}

func FromUsQuest(quest *qu.Quest) *Quest {
	if quest == nil {
		return nil
	}

	return &Quest{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		Cost:        quest.Cost,
		Type:        quest.Type,
	}
}
