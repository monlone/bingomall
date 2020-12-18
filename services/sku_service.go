package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// Sku service 接口
type SkuService interface {
	// 保存或修改
	SaveOrUpdate(sku *model.Sku) error

	// 根据id查询
	GetByID(id uint64) *model.Sku

	// 根据userId查询
	GetByUserID(userId uint64) []*model.Sku

	// 通过productIdCombineId获取该商品sku
	GetByProductIdCombineId(productId uint64, combineId string) *model.Sku

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll(shopId uint64) []*model.Sku

	// 分页查询
	GetPage(page int, pageSize int, user *model.Sku) *helper.PageBean
}

// sku service 结构体
type skuService struct {
	/** 存储对象 */
	repo repositories.SkuRepository
}

func (ss *skuService) SaveOrUpdate(sku *model.Sku) error {
	if sku == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if sku.ID == 0 {
		// 添加
		return ss.repo.Insert(sku)
	} else {
		// 修改
		persist := ss.GetByID(sku.ID)
		if persist == nil || sku.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		sku.ID = persist.ID
		return ss.repo.Update(sku)
	}
}

func (ss *skuService) GetByID(id uint64) *model.Sku {
	if id == 0 {
		return nil
	}
	sku := ss.repo.FindOne(id).(*model.Sku)
	return sku
}

func (ss *skuService) GetByUserID(userId uint64) []*model.Sku {
	sku := ss.repo.FindSingle("user_id = ?", userId)
	if sku == nil {
		return nil
	}

	return sku.([]*model.Sku)
}

func (ss *skuService) GetByProductIdCombineId(productId uint64, combineId string) *model.Sku {
	sku := ss.repo.FindSingle("product_id = ? AND combine_id = ?", productId, combineId)
	if sku == nil {
		return nil
	}
	data := sku.(*model.Sku)
	return data
}

func (*skuService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (ss *skuService) GetAll(productId uint64) []*model.Sku {
	var sku []*model.Sku

	if productId != 0 {
		sku = ss.repo.FindMore("product_id=?", productId).([]*model.Sku)
	} else {
		sku = ss.repo.FindMore("1=1").([]*model.Sku)
	}

	return sku
}

func (*skuService) GetPage(page int, pageSize int, user *model.Sku) *helper.PageBean {
	panic("implement me")
}

var skuServiceIns = &skuService{}

func SkuServiceInstance(repo repositories.SkuRepository) SkuService {
	skuServiceIns.repo = repo
	return skuServiceIns
}
