package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// Option service 接口
type OptionService interface {
	// 保存或修改
	SaveOrUpdate(option *model.Option) error

	// 根据id查询
	GetByID(id uint64) *model.Option

	// 根据userId查询
	GetByUserID(userId uint64) []*model.Option

	// 通过userId获取该用户的商品分类
	GetOptionByUserID(userId uint64) []*model.Option

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll(shopId uint64) []*model.Option

	// 分页查询
	GetPage(page int, pageSize int, user *model.Option) *helper.PageBean
}

// option service 结构体
type optionService struct {
	/** 存储对象 */
	repo repositories.OptionRepository
}

func (cs *optionService) SaveOrUpdate(option *model.Option) error {
	if option == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if option.ID == 0 {
		// 添加
		return cs.repo.Insert(option)
	} else {
		// 修改
		persist := cs.GetByID(option.ID)
		if persist == nil || option.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		option.ID = persist.ID
		return cs.repo.Update(option)
	}
}

func (cs *optionService) GetByID(id uint64) *model.Option {
	if id == 0 {
		return nil
	}
	option := cs.repo.FindOne(id).(*model.Option)
	return option
}

func (cs *optionService) GetByUserID(userId uint64) []*model.Option {
	option := cs.repo.FindSingle("user_id = ?", userId)
	if option == nil {
		return nil
	}

	return option.([]*model.Option)
}

func (cs *optionService) GetOptionByUserID(userId uint64) []*model.Option {
	option := cs.repo.FindSingle("user_id = ?", userId)
	if option == nil {
		return nil
	}
	data := option.([]*model.Option)
	return data
}

func (*optionService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (cs *optionService) GetAll(shopId uint64) []*model.Option {
	var option []*model.Option

	if shopId != 0 {
		option = cs.repo.FindMore("shop_id=?", shopId).([]*model.Option)
	} else {
		option = cs.repo.FindMore("1=1").([]*model.Option)
	}

	return option
}

func (*optionService) GetPage(page int, pageSize int, user *model.Option) *helper.PageBean {
	panic("implement me")
}

var optionServiceIns = &optionService{}

func OptionServiceInstance(repo repositories.OptionRepository) OptionService {
	optionServiceIns.repo = repo
	return optionServiceIns
}
