package service

import (
	"ToDoApp/pkg/repository"
	"fmt"
)

type StatusService struct {
	repo repository.Status
}

func NewStatusService(repo repository.Status) *StatusService {
	return &StatusService{repo: repo}
}

func (s *StatusService) GetAllStatuses() ([][]string, error) {
	statuses, err := s.repo.GetUsersStatuses()
	if err != nil {
		return nil, err
	}
	for _, status := range statuses {
		fmt.Println("Name: ", status[0], " Status: ", status[1])
	}
	return statuses, nil
}
