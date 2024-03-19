package response

import (
	"vk_quests/internal/pkg/time"
	"vk_quests/internal/pkg/types"
	uu "vk_quests/internal/usecase/user"
	"vk_quests/pkg/slices"
)

type User struct {
	ID      types.Id `json:"id" swaggertype:"integer" format:"uint64" example:"5"`
	Name    string   `json:"name" swaggertype:"string" example:"User"`
	Balance uint64   `json:"balance" swaggertype:"integer" format:"uint64"  example:"25"`
}

func FromUsUsers(users []uu.User) []User {
	return slices.Map(users, func(user uu.User) User {
		return *FromUsUser(&user)
	})
}

func FromUsUser(user *uu.User) *User {
	return &User{
		ID:      user.ID,
		Name:    user.Name,
		Balance: user.Balance,
	}
}

type HistoryRecord struct {
	Quest   *Quest             `json:"quest,omitempty"`
	Created time.FormattedTime `json:"created" swaggertype:"integer" format:"uint64" example:"5"`
	Balance uint64             `json:"balance" swaggertype:"integer" format:"uint64" example:"5"`
}

func FromUsHistoryRecord(record *uu.HistoryRecord) *HistoryRecord {
	return &HistoryRecord{
		Quest:   FromUsQuest(record.Quest),
		Created: record.Created,
		Balance: record.Balance,
	}
}

func FromUsHistory(history []uu.HistoryRecord) []HistoryRecord {
	return slices.Map(history, func(record uu.HistoryRecord) HistoryRecord {
		return *FromUsHistoryRecord(&record)
	})
}

type Status string

const (
	Success Status = "success"
	Failure Status = "failure"
)

type StatusApplyCost struct {
	Status Status `json:"status" swaggertype:"string"  enums:"success,failure" example:"success"`
}
