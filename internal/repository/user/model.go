package user

import (
	"vk_quests/internal/pkg/time"
	"vk_quests/internal/pkg/types"
	"vk_quests/internal/repository/quest"
)

type User struct {
	ID      types.Id
	Name    string
	Balance uint64
}

type HistoryRecord struct {
	Quest   *quest.Quest
	Created time.FormattedTime
	Balance uint64
}
