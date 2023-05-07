package service

import "P/repository"

type CommodityService struct {
	CommodityRepo repository.CommodityRepoInterface
}
type CommodityServiceInterface interface {
}
