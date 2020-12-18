package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
	"time"
)

// merchantOrderDetail_service 接口
type MerchantOrderDetailService interface {
	/** 保存或修改 */
	SaveOrUpdate(merchantOrderDetail *model.MerchantOrderDetail) error

	Save(merchantOrderDetail *model.MerchantOrderDetail) error

	Update(merchantOrderDetail *model.MerchantOrderDetail) error

	/** 根据 id 查询 */
	GetByMerchantOrderDetailID(id uint64) *model.MerchantOrderDetail

	GetByOrderId(outTradNo uint64) *model.MerchantOrderDetail

	Liquidate(userId, shopIds []uint64) error

	LiquidateTotal(userId []uint64, shopIds []uint64) (*model.TotalResult, error)

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有  */
	GetAll() []*model.MerchantOrderDetail

	/** 分页查询 */
	GetPage(page int, pageSize int, merchantOrderDetail *model.MerchantOrderDetail) *helper.PageBean

	/**  按商户id merchantOrderDetail 列表*/
	GetPageByMerchantIds(page int, pageSize int, merchantOrderDetail *model.MerchantOrderDetail) *helper.PageBean

	GetPageByShopIds(merchantOrderDetailPage *model.MerchantOrderDetailPage, merchantOrderDetail *model.MerchantOrderDetail, shopIDs []uint64) *helper.PageBean

	GetPageByCheckUserID(page int, pageSize int, merchantOrderDetail *model.MerchantOrderDetail) *helper.PageBean
}

var merchantOrderDetailServiceIns = &merchantOrderDetailService{}

// 获取 merchantOrderDetailService 实例
func MerchantOrderDetailServiceInstance(repo repositories.MerchantOrderDetailRepository) MerchantOrderDetailService {
	merchantOrderDetailServiceIns.repo = repo
	return merchantOrderDetailServiceIns
}

// 结构体
type merchantOrderDetailService struct {
	/** 存储对象 */
	repo repositories.MerchantOrderDetailRepository
}

func (us *merchantOrderDetailService) GetByMerchantOrderDetailOpenId(merchantOrderDetailId uint64) *model.MerchantOrderDetail {
	merchantOrderDetail := us.repo.FindSingle("merchantOrderDetail_id = ?", merchantOrderDetailId)
	if merchantOrderDetail != nil {
		return merchantOrderDetail.(*model.MerchantOrderDetail)
	}
	return nil
}

func (us *merchantOrderDetailService) SaveOrUpdate(merchantOrderDetail *model.MerchantOrderDetail) error {
	if merchantOrderDetail == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	if merchantOrderDetail.OrderId == 0 {
		return us.repo.Insert(merchantOrderDetail)
	} else {
		persist := us.GetByMerchantOrderDetailID(merchantOrderDetail.OrderId)
		if persist == nil || persist.OrderId == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}

		merchantOrderDetail.ID = persist.ID
		return us.repo.Update(merchantOrderDetail)
	}
}

func (us *merchantOrderDetailService) Save(merchantOrderDetail *model.MerchantOrderDetail) error {
	if merchantOrderDetail == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return us.repo.Insert(merchantOrderDetail)
}

func (us *merchantOrderDetailService) Update(merchantOrderDetail *model.MerchantOrderDetail) error {
	if merchantOrderDetail == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	persist := us.GetByOrderId(merchantOrderDetail.OrderId)
	if persist == nil || persist.ShopId == 0 {
		return errors.New(helper.StatusText(helper.UpdateObjIsNil))
	}

	merchantOrderDetail.ID = persist.ID
	return us.repo.Update(merchantOrderDetail)
}

func (us *merchantOrderDetailService) GetAll() []*model.MerchantOrderDetail {
	merchantOrderDetails := us.repo.FindMore("1=1").([]*model.MerchantOrderDetail)
	return merchantOrderDetails
}

func (us *merchantOrderDetailService) GetByMerchantOrderDetailID(merchantOrderDetailId uint64) *model.MerchantOrderDetail {
	if merchantOrderDetailId == 0 {
		return nil
	}
	merchantOrderDetail := us.repo.FindSingle("merchantOrderDetail_id = ?", merchantOrderDetailId).(*model.MerchantOrderDetail)
	return merchantOrderDetail
}

func (us *merchantOrderDetailService) GetByOrderId(orderId uint64) *model.MerchantOrderDetail {
	if orderId == 0 {
		return nil
	}
	merchantOrderDetail := us.repo.FindSingle("order_id = ?", orderId).(*model.MerchantOrderDetail)
	return merchantOrderDetail
}

func (us *merchantOrderDetailService) DeleteByID(id uint64) error {
	merchantOrderDetail := us.repo.FindOne(id).(*model.MerchantOrderDetail)
	if merchantOrderDetail == nil || merchantOrderDetail.OrderId == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := us.repo.Delete(merchantOrderDetail)
	return err
}

func (us *merchantOrderDetailService) GetPage(page int, pageSize int, merchantOrderDetail *model.MerchantOrderDetail) *helper.PageBean {
	andCons := make(map[string]interface{})

	if merchantOrderDetail != nil && merchantOrderDetail.OrderId != 0 {
		andCons["merchantOrderDetail_id = ?"] = merchantOrderDetail.OrderId
	}
	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}

func (us *merchantOrderDetailService) GetPageByMerchantIds(page int, pageSize int, merchantOrderDetail *model.MerchantOrderDetail) *helper.PageBean {
	andCons := make(map[string]interface{})

	if merchantOrderDetail.ShopId > 0 {
		andCons["shop_id = ?"] = merchantOrderDetail.ShopId
	}
	if merchantOrderDetail != nil && merchantOrderDetail.MerchantId != 0 {
		andCons["merchant_id = ?"] = merchantOrderDetail.MerchantId
	}

	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}

func (us *merchantOrderDetailService) GetPageByShopIds(merchantOrderDetailPage *model.MerchantOrderDetailPage, merchantOrderDetail *model.MerchantOrderDetail, shopIDs []uint64) *helper.PageBean {
	andCons := make(map[string]interface{})

	if len(shopIDs) >= 0 {
		andCons["shop_id in (?)"] = shopIDs
	}
	if merchantOrderDetail != nil && merchantOrderDetail.Status > 0 {
		andCons["status = ?"] = merchantOrderDetail.Status
	}

	if merchantOrderDetailPage.Month != "" {
		monthBegin, _ := time.Parse("2006-01-02", merchantOrderDetailPage.Month+"-01")
		andCons["updated_at >= ?"] = monthBegin
		monthEnd := monthBegin.AddDate(0, 1, 0)
		andCons["updated_at < ?"] = monthEnd
	}

	pageBean := us.repo.FindPage(merchantOrderDetailPage.Page, merchantOrderDetailPage.PageSize, andCons, nil)
	return pageBean
}

func (us *merchantOrderDetailService) GetPageByCheckUserID(page int, pageSize int, merchantOrderDetail *model.MerchantOrderDetail) *helper.PageBean {
	andCons := make(map[string]interface{})

	if merchantOrderDetail != nil && merchantOrderDetail.ShopId != 0 {
		andCons["shop_id = ?"] = merchantOrderDetail.ShopId
	}

	if merchantOrderDetail != nil {
		andCons["check_user_id = ?"] = merchantOrderDetail.CheckUserID
	} else {
		andCons["check_user_id = ?"] = ""
	}

	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}

func (us *merchantOrderDetailService) Liquidate(userIds, shopIds []uint64) error {
	err := us.repo.Liquidate(userIds, shopIds)
	return err
}

func (us *merchantOrderDetailService) LiquidateTotal(merchantId, shopIds []uint64) (data *model.TotalResult, err error) {
	data, err = us.repo.LiquidateTotal(merchantId, shopIds)
	return
}
