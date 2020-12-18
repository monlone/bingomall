package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// service 接口
type CityService interface {
	// 保存或修改
	SaveOrUpdate(city *model.City) error

	// 根据id查询
	GetByID(id uint64) *model.City

	// 根据userId查询
	GetByUserID(userId uint64) []*model.City

	// 通过userId获取该用户的商品分类
	GetCityByUserID(userId uint64) []*model.City

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll(provinceCode uint64) []*model.City

	// 分页查询
	GetPage(page int, pageSize int, user *model.City) *helper.PageBean
}

// city service 结构体
type cityService struct {
	/** 存储对象 */
	repo repositories.CityRepository
}

func (cs *cityService) SaveOrUpdate(city *model.City) error {
	if city == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if city.Code == 0 {
		// 添加
		return cs.repo.Insert(city)
	} else {
		// 修改
		persist := cs.GetByID(city.Code)
		if persist == nil || city.Code == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		city.ID = persist.ID
		return cs.repo.Update(city)
	}
}

func (cs *cityService) GetByID(id uint64) *model.City {
	if id == 0 {
		return nil
	}
	city := cs.repo.FindOne(id).(*model.City)
	return city
}

func (cs *cityService) GetByUserID(userId uint64) []*model.City {
	city := cs.repo.FindSingle("user_id = ?", userId)
	if city == nil {
		return nil
	}

	return city.([]*model.City)
}

func (cs *cityService) GetCityByUserID(userId uint64) []*model.City {
	city := cs.repo.FindSingle("user_id = ?", userId)
	if city == nil {
		return nil
	}
	data := city.([]*model.City)
	return data
}

func (*cityService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (cs *cityService) GetAll(provinceCode uint64) []*model.City {
	var city []*model.City

	if provinceCode != 0 {
		city = cs.repo.FindMore("province_code=?", provinceCode).([]*model.City)
	} else {
		city = cs.repo.FindMore("1=1").([]*model.City)
	}

	return city
}

func (*cityService) GetPage(page int, pageSize int, user *model.City) *helper.PageBean {
	panic("implement me")
}

var cityServiceIns = &cityService{}

func CityServiceInstance(repo repositories.CityRepository) CityService {
	cityServiceIns.repo = repo
	return cityServiceIns
}
