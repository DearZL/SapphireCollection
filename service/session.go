package service

import "P/repository"

type SessionServiceInterface interface {
	DropSession(SId string) (string, error)
}

type SessionService struct {
	SessionRepo repository.SessionRepositoryInterface
}

func (s *SessionService) DropSession(SId string) (string, error) {
	return s.SessionRepo.DeleteSession(SId)
}
