package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// merchant_service 接口
type MerchantService interface {
	/** 保存或修改 */
	SaveOrUpdate(merchant *model.Merchant) error

	Save(merchant *model.Merchant) error

	Update(merchant *model.Merchant) error

	/** 根据 merchant_id 查询 */
	GetByMerchantId(merchantId uint64) *model.Merchant

	/** 根据 user_id 查询 */
	GetByUserID(userId uint64) []*model.Merchant

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有  */
	GetAll() []*model.Merchant

	/** 分页查询 */
	GetPage(page int, pageSize int, merchant *model.Merchant) *helper.PageBean

	ShopList(page int, pageSize int, merchantId string) *helper.PageBean
}

var merchantServiceIns = &merchantService{}

// 获取 merchantService 实例
func MerchantServiceInstance(repo repositories.MerchantRepository) MerchantService {
	merchantServiceIns.repo = repo
	return merchantServiceIns
}

// 结构体
type merchantService struct {
	/** 存储对象 */
	repo repositories.MerchantRepository
}

func (us *merchantService) GetByMerchantOpenId(merchantId uint64) *model.Merchant {
	merchant := us.repo.FindSingle("id = ?", merchantId)
	if merchant != nil {
		return merchant.(*model.Merchant)
	}
	return nil
}

func (us *merchantService) SaveOrUpdate(merchant *model.Merchant) error {
	if merchant == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	if merchant.ID == 0 {
		// 添加
		return us.repo.Insert(merchant)
	} else {
		// 修改
		persist := us.GetByMerchantId(merchant.ID)
		if persist == nil || persist.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}

		merchant.ID = persist.ID
		return us.repo.Update(merchant)
	}
}

func (us *merchantService) Save(merchant *model.Merchant) error {
	if merchant == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return us.repo.Insert(merchant)
}

func (us *merchantService) Update(merchant *model.Merchant) error {
	if merchant == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return us.repo.Update(merchant)
}

func (us *merchantService) GetAll() []*model.Merchant {
	merchants := us.repo.FindMore("1=1").([]*model.Merchant)
	return merchants
}

func (us *merchantService) GetByMerchantId(merchantId uint64) *model.Merchant {
	if merchantId == 0 {
		return nil
	}
	merchant := us.repo.FindSingle("id = ?", merchantId)
	if merchant == nil {
		return nil
	}

	return merchant.(*model.Merchant)
}

func (us *merchantService) GetByUserID(userId uint64) []*model.Merchant {
	if userId == 0 {
		return nil
	}
	merchant := us.repo.FindMore("user_id = ?", userId).([]*model.Merchant)
	return merchant
}

func (us *merchantService) DeleteByID(id uint64) error {
	merchant := us.repo.FindOne(id).(*model.Merchant)
	if merchant == nil || merchant.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := us.repo.Delete(merchant)
	return err
}

func (us *merchantService) GetPage(page int, pageSize int, merchant *model.Merchant) *helper.PageBean {
	andCons := make(map[string]interface{})

	if merchant != nil && merchant.ID != 0 {
		andCons["id = ?"] = merchant.ID
	}
	if merchant != nil && merchant.Title != "" {
		andCons["title like ?"] = "%" + merchant.Title + "%"
	}
	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}

func (us *merchantService) ShopList(page int, pageSize int, merchantId string) *helper.PageBean {
	pageBean := us.repo.ShopList(page, pageSize, merchantId)
	return pageBean
}
