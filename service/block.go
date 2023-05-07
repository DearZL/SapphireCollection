package service

import "P/repository"

type BlockServiceInterface interface {
}

type BlockService struct {
	BlockRepo repository.BlockRepoInterface
}
