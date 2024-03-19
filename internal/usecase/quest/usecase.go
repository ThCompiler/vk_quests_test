package quest

import (
	"vk_quests/internal/pkg/types"
	"vk_quests/internal/repository/quest"
	"vk_quests/pkg/slices"
)

type QuestUsecase struct {
	quests quest.Repository
}

func NewQuestUsecase(quests quest.Repository) *QuestUsecase {
	return &QuestUsecase{
		quests: quests,
	}
}

func (qu *QuestUsecase) CreateQuest(qst *Quest) (*Quest, error) {
	createdQst, err := qu.quests.CreateQuest(
		&quest.Quest{
			ID:          qst.ID,
			Name:        qst.Name,
			Description: qst.Description,
			Cost:        qst.Cost,
			Type:        qst.Type,
		},
	)

	return FromRepQuest(createdQst), err
}

func (qu *QuestUsecase) DeleteQuest(id types.Id) error {
	return qu.quests.DeleteQuest(id)
}

func (qu *QuestUsecase) UpdateQuest(id types.Id, qst *UpdateQuest) (*Quest, error) {
	updatedQst, err := qu.quests.UpdateQuest(qst.ToRepUpdateQuest(id))

	return FromRepQuest(updatedQst), err
}

func (qu *QuestUsecase) GetQuests() ([]Quest, error) {
	quests, err := qu.quests.GetQuests()
	if err != nil {
		return nil, err
	}

	return slices.Map(quests, func(q quest.Quest) Quest { return *FromRepQuest(&q) }), nil
}

func (qu *QuestUsecase) GetQuest(id types.Id) (*Quest, error) {
	qst, err := qu.quests.GetQuest(id)

	return FromRepQuest(qst), err
}
