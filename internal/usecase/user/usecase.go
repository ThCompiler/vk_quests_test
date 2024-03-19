package user

import (
	"math/rand"
	"time"

	"vk_quests/internal/pkg/types"
	"vk_quests/internal/repository/quest"
	"vk_quests/internal/repository/user"
	"vk_quests/pkg/slices"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

const CompleteChance = 0.5

type UserUsecase struct {
	users  user.Repository
	quests quest.Repository
}

func NewUserUsecase(users user.Repository, quests quest.Repository) *UserUsecase {
	return &UserUsecase{
		users:  users,
		quests: quests,
	}
}

func (uu *UserUsecase) CreateUser(name string) (*User, error) {
	usr, err := uu.users.CreateUser(&user.User{
		Name: name,
	})

	return FromRepUser(usr), err
}

func (uu *UserUsecase) DeleteUser(id types.Id) (*User, error) {
	usr, err := uu.users.DeleteUser(id)

	return FromRepUser(usr), err
}

func (uu *UserUsecase) UpdateUser(id types.Id, name string) (*User, error) {
	usr, err := uu.users.UpdateUser(&user.User{
		ID:   id,
		Name: name,
	})

	return FromRepUser(usr), err
}

func (uu *UserUsecase) GetUsers() ([]User, error) {
	usrs, err := uu.users.GetUsers()
	if err != nil {
		return nil, err
	}

	return slices.Map(usrs, func(usr user.User) User { return *FromRepUser(&usr) }), nil
}

func (uu *UserUsecase) GetUserHistory(id types.Id) ([]HistoryRecord, error) {
	history, err := uu.users.GetHistory(id)
	if err != nil {
		return nil, err
	}

	return slices.Map(history, func(record user.HistoryRecord) HistoryRecord { return *FromRepHistory(&record) }), nil
}

func (uu *UserUsecase) ApplyQuests(questId, userId types.Id) error {
	qst, err := uu.quests.GetQuest(questId)
	if err != nil {
		return err
	}

	if err := uu.users.HasUser(userId); err != nil {
		return err
	}

	if err := uu.users.IsCompletedQuest(&user.User{ID: userId}, qst); err != nil {
		return err
	}

	if qst.Type == types.USUAL || rnd.Float64() > CompleteChance {
		return uu.users.ApplyCost(&user.User{ID: userId}, qst)
	}

	return QuestNotApplied
}
