package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// service 接口
type AreaService interface {
	// 保存或修改
	SaveOrUpdate(area *model.Area) error

	// 根据id查询
	GetByID(id uint64) *model.Area

	// 根据userId查询
	GetByUserID(userId uint64) []*model.Area

	// 通过userId获取该用户的商品分类
	GetAreaByUserID(userId uint64) []*model.Area

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll(cityCode uint64) []*model.Area

	// 分页查询
	GetPage(page int, pageSize int, user *model.Area) *helper.PageBean
}

// area service 结构体
type areaService struct {
	/** 存储对象 */
	repo repositories.AreaRepository
}

func (cs *areaService) SaveOrUpdate(area *model.Area) error {
	if area == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if area.Code == 0 {
		// 添加
		return cs.repo.Insert(area)
	} else {
		// 修改
		persist := cs.GetByID(area.Code)
		if persist == nil || area.Code == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		area.ID = persist.ID
		return cs.repo.Update(area)
	}
}

func (cs *areaService) GetByID(id uint64) *model.Area {
	if id == 0 {
		return nil
	}
	area := cs.repo.FindOne(id).(*model.Area)
	return area
}

func (cs *areaService) GetByUserID(userId uint64) []*model.Area {
	area := cs.repo.FindSingle("user_id = ?", userId)
	if area == nil {
		return nil
	}

	return area.([]*model.Area)
}

func (cs *areaService) GetAreaByUserID(userId uint64) []*model.Area {
	area := cs.repo.FindSingle("user_id = ?", userId)
	if area == nil {
		return nil
	}
	data := area.([]*model.Area)
	return data
}

func (*areaService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (cs *areaService) GetAll(provinceCode uint64) []*model.Area {
	var area []*model.Area

	if provinceCode != 0 {
		area = cs.repo.FindMore("code=?", provinceCode).([]*model.Area)
	} else {
		area = cs.repo.FindMore("1=1").([]*model.Area)
	}

	return area
}

func (*areaService) GetPage(page int, pageSize int, user *model.Area) *helper.PageBean {
	panic("implement me")
}

var areaServiceIns = &areaService{}

func AreaServiceInstance(repo repositories.AreaRepository) AreaService {
	areaServiceIns.repo = repo
	return areaServiceIns
}
