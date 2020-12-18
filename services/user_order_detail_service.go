package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
	"time"
)

// merchantUserOrderDetailDetail_service 接口
type UserOrderDetailService interface {
	/** 保存或修改 */
	SaveOrUpdate(merchantUserOrderDetailDetail *model.UserOrderDetail) error

	Save(merchantUserOrderDetailDetail *model.UserOrderDetail) error

	Update(merchantUserOrderDetailDetail *model.UserOrderDetail) error

	/** 根据 id 查询 */
	GetByUserOrderDetailID(id uint64) *model.UserOrderDetail

	GetByOrderId(outTradNo uint64) *model.UserOrderDetail

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有  */
	GetAll() []*model.UserOrderDetail

	/** 分页查询 */
	GetPage(page int, pageSize int, merchantUserOrderDetailDetail *model.UserOrderDetail) *helper.PageBean

	/**  按月份取 merchantUserOrderDetailDetail 列表*/
	GetPageByMonth(page int, pageSize int, merchantUserOrderDetailDetail *model.UserOrderDetail, month string) *helper.PageBean
}

var merchantUserOrderDetailDetailServiceIns = &merchantUserOrderDetailDetailService{}

// 获取 merchantUserOrderDetailDetailService 实例
func UserOrderDetailServiceInstance(repo repositories.UserOrderDetailRepository) UserOrderDetailService {
	merchantUserOrderDetailDetailServiceIns.repo = repo
	return merchantUserOrderDetailDetailServiceIns
}

// 结构体
type merchantUserOrderDetailDetailService struct {
	/** 存储对象 */
	repo repositories.UserOrderDetailRepository
}

func (us *merchantUserOrderDetailDetailService) GetByUserOrderDetailOpenId(merchantUserOrderDetailDetailId string) *model.UserOrderDetail {
	merchantUserOrderDetailDetail := us.repo.FindSingle("merchantUserOrderDetailDetail_id = ?", merchantUserOrderDetailDetailId)
	if merchantUserOrderDetailDetail != nil {
		return merchantUserOrderDetailDetail.(*model.UserOrderDetail)
	}
	return nil
}

func (us *merchantUserOrderDetailDetailService) SaveOrUpdate(merchantUserOrderDetailDetail *model.UserOrderDetail) error {
	if merchantUserOrderDetailDetail == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	if merchantUserOrderDetailDetail.OrderId == 0 {
		return us.repo.Insert(merchantUserOrderDetailDetail)
	} else {
		persist := us.GetByUserOrderDetailID(merchantUserOrderDetailDetail.OrderId)
		if persist == nil || persist.OrderId == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}

		merchantUserOrderDetailDetail.ID = persist.ID
		return us.repo.Update(merchantUserOrderDetailDetail)
	}
}

func (us *merchantUserOrderDetailDetailService) Save(merchantUserOrderDetailDetail *model.UserOrderDetail) error {
	if merchantUserOrderDetailDetail == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return us.repo.Insert(merchantUserOrderDetailDetail)
}

func (us *merchantUserOrderDetailDetailService) Update(merchantUserOrderDetailDetail *model.UserOrderDetail) error {
	if merchantUserOrderDetailDetail == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	persist := us.GetByOrderId(merchantUserOrderDetailDetail.OrderId)
	if persist == nil || persist.ShopId == 0 {
		return errors.New(helper.StatusText(helper.UpdateObjIsNil))
	}

	merchantUserOrderDetailDetail.ID = persist.ID

	return us.repo.Update(merchantUserOrderDetailDetail)
}

func (us *merchantUserOrderDetailDetailService) GetAll() []*model.UserOrderDetail {
	merchantUserOrderDetailDetails := us.repo.FindMore("1=1").([]*model.UserOrderDetail)
	return merchantUserOrderDetailDetails
}

func (us *merchantUserOrderDetailDetailService) GetByUserOrderDetailID(merchantUserOrderDetailDetailId uint64) *model.UserOrderDetail {
	if merchantUserOrderDetailDetailId == 0 {
		return nil
	}
	merchantUserOrderDetailDetail := us.repo.FindSingle("merchantUserOrderDetailDetail_id = ?", merchantUserOrderDetailDetailId).(*model.UserOrderDetail)
	return merchantUserOrderDetailDetail
}

func (us *merchantUserOrderDetailDetailService) GetByOrderId(orderId uint64) *model.UserOrderDetail {
	if orderId == 0 {
		return nil
	}
	merchantUserOrderDetailDetail := us.repo.FindSingle("order_id = ?", orderId).(*model.UserOrderDetail)
	return merchantUserOrderDetailDetail
}

func (us *merchantUserOrderDetailDetailService) DeleteByID(id uint64) error {
	merchantUserOrderDetailDetail := us.repo.FindOne(id).(*model.UserOrderDetail)
	if merchantUserOrderDetailDetail == nil || merchantUserOrderDetailDetail.OrderId == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := us.repo.Delete(merchantUserOrderDetailDetail)
	return err
}

func (us *merchantUserOrderDetailDetailService) GetPage(page int, pageSize int, merchantUserOrderDetailDetail *model.UserOrderDetail) *helper.PageBean {
	andCons := make(map[string]interface{})

	if merchantUserOrderDetailDetail != nil && merchantUserOrderDetailDetail.OrderId != 0 {
		andCons["merchantUserOrderDetailDetail_id = ?"] = merchantUserOrderDetailDetail.OrderId
	}
	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}

func (us *merchantUserOrderDetailDetailService) GetPageByMonth(page int, pageSize int, merchantUserOrderDetailDetail *model.UserOrderDetail, month string) *helper.PageBean {
	andCons := make(map[string]interface{})

	if merchantUserOrderDetailDetail != nil && merchantUserOrderDetailDetail.UserID != 0 {
		andCons["user_id = ?"] = merchantUserOrderDetailDetail.UserID
	}

	if month != "" {
		monthBegin, _ := time.Parse("2006-01-02", month+"-01")
		andCons["updated_at >= ?"] = monthBegin
		monthEnd := monthBegin.AddDate(0, 1, 0)
		//timeUnix := monthEnd.Unix() + 86400 - 1
		//end := time.Unix(timeUnix, 0)
		//fmt.Println("end:", end)
		andCons["updated_at < ?"] = monthEnd
	}

	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}
