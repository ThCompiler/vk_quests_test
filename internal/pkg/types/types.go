package types

type Id uint64

type Cost uint32

type QuestType string

const (
	RANDOM QuestType = "random"
	USUAL  QuestType = "usual"
)

type ContextField string
