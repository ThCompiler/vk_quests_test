package user

import (
	"vk_quests/internal/pkg/time"
	"vk_quests/internal/pkg/types"
	"vk_quests/internal/repository/user"
	"vk_quests/internal/usecase/quest"
)

type User struct {
	ID      types.Id
	Name    string
	Balance uint64
}

func FromRepUser(u *user.User) *User {
	if u == nil {
		return nil
	}

	return &User{
		ID:      u.ID,
		Name:    u.Name,
		Balance: u.Balance,
	}
}

type HistoryRecord struct {
	Quest   *quest.Quest
	Created time.FormattedTime
	Balance uint64
}

func FromRepHistory(hr *user.HistoryRecord) *HistoryRecord {
	return &HistoryRecord{
		Quest:   quest.FromRepQuest(hr.Quest),
		Created: hr.Created,
		Balance: hr.Balance,
	}
}
