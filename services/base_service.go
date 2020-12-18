package service

import (
	"bingomall/models"
	"bingomall/repositories"
)

// base service 接口
type BaseService interface {
	IsPlatformOfficial(shop *model.Shop) bool
}

// function service 结构体
type baseService struct {
	/** 存储对象 */
	repo repositories.BaseRepository
}

var baseServiceIns = &baseService{}

func BaseServiceInstance(repo repositories.BaseRepository) BaseService {
	baseServiceIns.repo = repo
	return baseServiceIns
}

// 平台官方自己运营的店，或者商户代平台运营的店，商户要拿平台给的提成
func (us *baseService) IsPlatformOfficial(shop *model.Shop) bool {
	if shop.Type == 1 || shop.Type == 2 {
		return true
	}
	return false
}
