package service

import (
	"P/model"
	"P/repository"
)

type CommonService struct {
	CommonRepo repository.CommonRepositoryInterface
}
type CommonServiceInterface interface {
	Upload(files []model.File) error
	Download(filename string) (*model.File, error)
}

func (s *CommonService) Upload(files []model.File) error {
	return s.CommonRepo.AddFile(files)
}
func (s *CommonService) Download(filename string) (*model.File, error) {
	return s.CommonRepo.FindFile(filename)
}
