package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// service 接口
type ProvinceService interface {
	// 保存或修改
	SaveOrUpdate(province *model.Province) error

	// 根据id查询
	GetByID(id uint64) *model.Province

	// 根据userId查询
	GetByUserID(userId uint64) []*model.Province

	// 通过userId获取该用户的商品分类
	GetProvinceByUserID(userId uint64) []*model.Province

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll(shopId uint64) []*model.Province

	// 分页查询
	GetPage(page int, pageSize int, user *model.Province) *helper.PageBean
}

// province service 结构体
type provinceService struct {
	/** 存储对象 */
	repo repositories.ProvinceRepository
}

func (cs *provinceService) SaveOrUpdate(province *model.Province) error {
	if province == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if province.Code == 0 {
		// 添加
		return cs.repo.Insert(province)
	} else {
		// 修改
		persist := cs.GetByID(province.Code)
		if persist == nil || province.Code == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		province.ID = persist.ID
		return cs.repo.Update(province)
	}
}

func (cs *provinceService) GetByID(code uint64) *model.Province {
	if code == 0 {
		return nil
	}
	province := cs.repo.FindOne(code).(*model.Province)
	return province
}

func (cs *provinceService) GetByUserID(userId uint64) []*model.Province {
	province := cs.repo.FindSingle("user_id = ?", userId)
	if province == nil {
		return nil
	}

	return province.([]*model.Province)
}

func (cs *provinceService) GetProvinceByUserID(userId uint64) []*model.Province {
	province := cs.repo.FindSingle("user_id = ?", userId)
	if province == nil {
		return nil
	}
	data := province.([]*model.Province)
	return data
}

func (*provinceService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (cs *provinceService) GetAll(provinceId uint64) []*model.Province {
	var province []*model.Province

	if provinceId != 0 {
		province = cs.repo.FindMore("province_id=?", provinceId).([]*model.Province)
	} else {
		province = cs.repo.FindMore("1=1").([]*model.Province)
	}

	return province
}

func (*provinceService) GetPage(page int, pageSize int, user *model.Province) *helper.PageBean {
	panic("implement me")
}

var provinceServiceIns = &provinceService{}

func ProvinceServiceInstance(repo repositories.ProvinceRepository) ProvinceService {
	provinceServiceIns.repo = repo
	return provinceServiceIns
}
