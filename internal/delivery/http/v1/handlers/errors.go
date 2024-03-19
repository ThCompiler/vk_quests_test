package handlers

import "github.com/pkg/errors"

var (
	ErrorCannotReadBody       = errors.New("can't read body")
	ErrorIncorrectBodyContent = errors.New("incorrect body content")
	ErrorUnknownError         = errors.New("unknown error, try again later")
	ErrorIncorrectQueryParam  = errors.New("invalid query parameter")

	ErrorUserAlreadyCompleteQuest = errors.New("user already complete quest")
	ErrorQuestNameAlreadyExists   = errors.New("quest with this name already exists")
	ErrorQuestNotFound            = errors.New("quest not found")
	ErrorUserNotFound             = errors.New("user not found")
)
