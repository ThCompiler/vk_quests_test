package request

import (
	"github.com/miladibra10/vjson"
	"vk_quests/internal/pkg/evjson"
	"vk_quests/internal/pkg/types"
	qu "vk_quests/internal/usecase/quest"
)

type CreateQuest struct {
	Name        string          `json:"name" swaggertype:"string" example:"Task"`
	Description string          `json:"description" swaggertype:"string" example:"Random quest"`
	Cost        types.Cost      `json:"cost" swaggertype:"integer" format:"uint8" example:"9" minimum:"0" maximum:"1000"`
	Type        types.QuestType `json:"type" swaggertype:"string" enums:"usual,random" example:"random"`
}

func (c *CreateQuest) ToUsQuest() *qu.Quest {
	return &qu.Quest{
		Name:        c.Name,
		Description: c.Description,
		Cost:        c.Cost,
		Type:        c.Type,
	}
}

func ValidateCreateQuest(data []byte) error {
	schema := evjson.NewSchema(
		vjson.String("name").Required(),
		vjson.String("description").Required(),
		vjson.Integer("cost").Range(0, 1000).Required(),
		vjson.String("type").Choices(string(types.USUAL), string(types.RANDOM)).Required(),
	)
	return schema.ValidateBytes(data)
}

type UpdateQuest struct {
	Description *string          `json:"description,omitempty" swaggertype:"string" example:"Random quest"`
	Cost        *types.Cost      `json:"cost,omitempty" swaggertype:"integer" format:"uint8" example:"9" minimum:"0" maximum:"1000"`
	Type        *types.QuestType `json:"type,omitempty" swaggertype:"string" enums:"usual,random" example:"random"`
}

func (u *UpdateQuest) ToUsUpdateQuest() *qu.UpdateQuest {
	return &qu.UpdateQuest{
		Description: u.Description,
		Cost:        u.Cost,
		Type:        u.Type,
	}
}

func ValidateUpdateQuest(data []byte) error {
	schema := evjson.NewSchema(
		vjson.String("name"),
		vjson.String("description"),
		vjson.Integer("cost").Range(0, 1000),
		vjson.String("type").Choices(string(types.USUAL), string(types.RANDOM)),
	)

	return schema.ValidateBytes(data)
}
