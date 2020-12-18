package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// shop_service 接口
type ShopService interface {
	/** 保存或修改 */
	SaveOrUpdate(shop *model.Shop) error

	Save(shop *model.Shop) error

	Update(shop *model.Shop) error

	/** 根据 shop_id 查询 */
	GetByShopId(shopId uint64) *model.Shop

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有  */
	GetAll() []*model.Shop

	GetShopsByMerchantIds(merchantIds []uint64) []*model.Shop

	GetShopsByShopIds(shopIds []uint64) []*model.Shop

	/** 分页查询 */
	GetPage(page int, pageSize int, shop *model.Shop) *helper.PageBean

	ShopListNearby(page int, pageSize int, shop *model.ShopDetailDistance) *helper.PageBean

	ProductList(page int, pageSize int, shopId string) *helper.PageBean

	RawSqlInsert(shop *model.Shop) error

	RawSqlUpdate(shop *model.Shop) error

	AddAttention(shop *model.Shop) error
}

var shopServiceIns = &shopService{}

// 获取 shopService 实例
func ShopServiceInstance(repo repositories.ShopRepository) ShopService {
	shopServiceIns.repo = repo
	return shopServiceIns
}

// 结构体
type shopService struct {
	/** 存储对象 */
	repo repositories.ShopRepository
}

func (ss *shopService) GetByShopOpenId(shopId string) *model.Shop {
	shop := ss.repo.FindSingle("shop_id = ?", shopId)
	if shop != nil {
		return shop.(*model.Shop)
	}
	return nil
}

func (ss *shopService) SaveOrUpdate(shop *model.Shop) error {
	if shop == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	if shop.ID == 0 {
		// 添加
		return ss.repo.RawSqlInsert(shop)
	} else {
		// 修改
		persist := ss.GetByShopId(shop.ID)
		if persist == nil {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}

		shop.ID = persist.ID
		_ = ss.repo.Update(shop)
		return ss.repo.RawSqlUpdate(shop)
	}
}

func (ss *shopService) Save(shop *model.Shop) error {
	if shop == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return ss.repo.Insert(shop)
}

func (ss *shopService) Update(shop *model.Shop) error {
	if shop == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 修改
	persist := ss.GetByShopId(shop.ID)
	if persist == nil {
		return errors.New(helper.StatusText(helper.UpdateObjIsNil))
	}

	shop.ID = persist.ID
	_ = ss.repo.Update(shop)

	return ss.repo.Update(shop)
}

func (ss *shopService) GetAll() []*model.Shop {
	shops := ss.repo.FindMore("1=1").([]*model.Shop)
	return shops
}

func (ss *shopService) GetShopsByMerchantIds(merchantIds []uint64) []*model.Shop {
	shops := ss.repo.FindMore("merchant_id in (?)", merchantIds)
	if shops != nil {
		return shops.([]*model.Shop)
	}
	return nil
}
func (ss *shopService) GetShopsByShopIds(shopIds []uint64) []*model.Shop {
	shops := ss.repo.FindMore("shop_id in (?)", shopIds)
	if shops != nil {
		return shops.([]*model.Shop)
	}
	return nil
}
func (ss *shopService) GetByShopId(shopId uint64) *model.Shop {
	if shopId == 0 {
		return nil
	}
	shop := ss.repo.FindSingle("id = ?", shopId)
	if shop != nil {
		return shop.(*model.Shop)
	}

	return nil
}

func (ss *shopService) DeleteByID(id uint64) error {
	shop := ss.repo.FindOne(id).(*model.Shop)
	if shop == nil || shop.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := ss.repo.Delete(shop)
	return err
}

func (ss *shopService) GetPage(page int, pageSize int, shop *model.Shop) *helper.PageBean {
	andCons := make(map[string]interface{})

	if shop != nil && shop.ID != 0 {
		andCons["shop_id = ?"] = shop.ID
	}
	if shop != nil && shop.Title != "" {
		andCons["title like ?"] = "%" + shop.Title + "%"
	}
	andCons["status = ?"] = 1
	pageBean := ss.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}

func (ss *shopService) ProductList(page int, pageSize int, shopId string) *helper.PageBean {
	pageBean := ss.repo.ProductList(page, pageSize, shopId)
	return pageBean
}

func (ss *shopService) RawSqlInsert(shop *model.Shop) error {
	err := ss.repo.RawSqlInsert(shop)
	return err
}
func (ss *shopService) RawSqlUpdate(shop *model.Shop) error {
	err := ss.repo.RawSqlUpdate(shop)
	return err
}

func (ss *shopService) ShopListNearby(page int, pageSize int, shop *model.ShopDetailDistance) *helper.PageBean {
	err := ss.repo.ShopListNearby(page, pageSize, shop)
	return err
}

func (ss *shopService) AddAttention(shop *model.Shop) error {
	persist := ss.GetByShopId(shop.ID)
	if persist == nil || persist.ID == 0 {
		return errors.New(helper.StatusText(helper.UpdateObjIsNil))
	}

	persist.AttentionNum = persist.AttentionNum + 1
	err := ss.repo.Update(persist)
	_ = ss.repo.RawSqlUpdate(persist)
	return err
}
