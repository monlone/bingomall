package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// Demo service 接口
type DemoService interface {
	// 保存或修改
	SaveOrUpdate(demo *model.Demo) error

	// 根据id查询
	GetByID(id uint64) *model.Demo

	// 根据userId查询
	GetByUserID(userId uint64) []*model.Demo

	// 通过userId获取该用户的商品分类
	GetDemoByUserID(userId uint64) []*model.Demo

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll(shopId uint64) []*model.Demo

	// 分页查询
	GetPage(page int, pageSize int, user *model.Demo) *helper.PageBean
}

// demo service 结构体
type demoService struct {
	/** 存储对象 */
	repo repositories.DemoRepository
}

func (cs *demoService) SaveOrUpdate(demo *model.Demo) error {
	if demo == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if demo.ID == 0 {
		// 添加
		return cs.repo.Insert(demo)
	} else {
		// 修改
		persist := cs.GetByID(demo.ID)
		if persist == nil || demo.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		demo.ID = persist.ID
		return cs.repo.Update(demo)
	}
}

func (cs *demoService) GetByID(id uint64) *model.Demo {
	if id == 0 {
		return nil
	}
	demo := cs.repo.FindOne(id).(*model.Demo)
	return demo
}

func (cs *demoService) GetByUserID(userId uint64) []*model.Demo {
	demo := cs.repo.FindSingle("user_id = ?", userId)
	if demo == nil {
		return nil
	}

	return demo.([]*model.Demo)
}

func (cs *demoService) GetDemoByUserID(userId uint64) []*model.Demo {
	demo := cs.repo.FindSingle("user_id = ?", userId)
	if demo == nil {
		return nil
	}
	data := demo.([]*model.Demo)
	return data
}

func (*demoService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (cs *demoService) GetAll(shopId uint64) []*model.Demo {
	var demo []*model.Demo

	if shopId != 0 {
		demo = cs.repo.FindMore("id=?", shopId).([]*model.Demo)
	} else {
		demo = cs.repo.FindMore("1=1").([]*model.Demo)
	}

	return demo
}

func (*demoService) GetPage(page int, pageSize int, user *model.Demo) *helper.PageBean {
	panic("implement me")
}

var demoServiceIns = &demoService{}

func DemoServiceInstance(repo repositories.DemoRepository) DemoService {
	demoServiceIns.repo = repo
	return demoServiceIns
}
