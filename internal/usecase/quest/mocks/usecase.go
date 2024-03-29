// Code generated by MockGen. DO NOT EDIT.
// Source: vk_quests/internal/usecase/quest (interfaces: Usecase)
//
// Generated by this command:
//
//	mockgen -destination=mocks/usecase.go -package=mu -mock_names=Usecase=QuestUsecase . Usecase
//

// Package mu is a generated GoMock package.
package mu

import (
	reflect "reflect"
	types "vk_quests/internal/pkg/types"
	quest "vk_quests/internal/usecase/quest"

	gomock "go.uber.org/mock/gomock"
)

// QuestUsecase is a mock of Usecase interface.
type QuestUsecase struct {
	ctrl     *gomock.Controller
	recorder *QuestUsecaseMockRecorder
}

// QuestUsecaseMockRecorder is the mock recorder for QuestUsecase.
type QuestUsecaseMockRecorder struct {
	mock *QuestUsecase
}

// NewQuestUsecase creates a new mock instance.
func NewQuestUsecase(ctrl *gomock.Controller) *QuestUsecase {
	mock := &QuestUsecase{ctrl: ctrl}
	mock.recorder = &QuestUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *QuestUsecase) EXPECT() *QuestUsecaseMockRecorder {
	return m.recorder
}

// CreateQuest mocks base method.
func (m *QuestUsecase) CreateQuest(arg0 *quest.Quest) (*quest.Quest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQuest", arg0)
	ret0, _ := ret[0].(*quest.Quest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateQuest indicates an expected call of CreateQuest.
func (mr *QuestUsecaseMockRecorder) CreateQuest(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQuest", reflect.TypeOf((*QuestUsecase)(nil).CreateQuest), arg0)
}

// DeleteQuest mocks base method.
func (m *QuestUsecase) DeleteQuest(arg0 types.Id) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteQuest", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteQuest indicates an expected call of DeleteQuest.
func (mr *QuestUsecaseMockRecorder) DeleteQuest(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteQuest", reflect.TypeOf((*QuestUsecase)(nil).DeleteQuest), arg0)
}

// GetQuest mocks base method.
func (m *QuestUsecase) GetQuest(arg0 types.Id) (*quest.Quest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuest", arg0)
	ret0, _ := ret[0].(*quest.Quest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuest indicates an expected call of GetQuest.
func (mr *QuestUsecaseMockRecorder) GetQuest(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuest", reflect.TypeOf((*QuestUsecase)(nil).GetQuest), arg0)
}

// GetQuests mocks base method.
func (m *QuestUsecase) GetQuests() ([]quest.Quest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetQuests")
	ret0, _ := ret[0].([]quest.Quest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetQuests indicates an expected call of GetQuests.
func (mr *QuestUsecaseMockRecorder) GetQuests() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetQuests", reflect.TypeOf((*QuestUsecase)(nil).GetQuests))
}

// UpdateQuest mocks base method.
func (m *QuestUsecase) UpdateQuest(arg0 types.Id, arg1 *quest.UpdateQuest) (*quest.Quest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateQuest", arg0, arg1)
	ret0, _ := ret[0].(*quest.Quest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateQuest indicates an expected call of UpdateQuest.
func (mr *QuestUsecaseMockRecorder) UpdateQuest(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateQuest", reflect.TypeOf((*QuestUsecase)(nil).UpdateQuest), arg0, arg1)
}
