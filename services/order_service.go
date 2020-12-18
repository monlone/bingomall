package service

import (
	"errors"
	"bingomall/constant"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
	"time"
)

// order_service 接口
type OrderService interface {
	/** 保存或修改 */
	SaveOrUpdate(order *model.Order) error

	Save(order *model.Order) error

	Update(order *model.Order) error

	/** 根据 id 查询 */
	GetByOrderId(id uint64) *model.Order

	GetByOrderIdAndUserId(id uint64, userId uint64) *model.Order

	GetLevelOrder(userId uint64) *model.Order

	GetByOutTradNo(outTradNo string) *model.Order

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有  */
	GetAll() []*model.Order

	GetByUserId(userId uint64) []*model.Order

	GetByUserIdAndStatus(userId uint64, status int8) []*model.Order

	Statistics(userId uint64) []*model.Order

	/** 分页查询 */
	GetPage(page int, pageSize int, order *model.Order) *helper.PageBean

	/**  按月份取 order 列表*/
	GetPageByMonth(page int, pageSize int, order *model.Order, month string) *helper.PageBean
}

var orderServiceIns = &orderService{}

// 获取 orderService 实例
func OrderServiceInstance(repo repositories.OrderRepository) OrderService {
	orderServiceIns.repo = repo
	return orderServiceIns
}

// 结构体
type orderService struct {
	/** 存储对象 */
	repo repositories.OrderRepository
}

func (os *orderService) GetByOrderOpenId(orderId uint64) *model.Order {
	order := os.repo.FindSingle("id = ?", orderId)
	if order != nil {
		return order.(*model.Order)
	}
	return nil
}

func (os *orderService) GetLevelOrder(userId uint64) *model.Order {
	order := os.repo.FindSingle("user_id = ? and type = ?", userId, constant.OrderLevel)
	if order != nil {
		return order.(*model.Order)
	}
	return nil
}

func (os *orderService) SaveOrUpdate(order *model.Order) error {
	if order == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	if order.ID == 0 {
		return os.repo.Insert(order)
	} else {
		persist := os.GetByOrderId(order.ID)
		if persist == nil || persist.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}

		order.ID = persist.ID
		return os.repo.Update(order)
	}
}

func (os *orderService) Save(order *model.Order) error {
	if order == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return os.repo.Insert(order)
}

func (os *orderService) Update(order *model.Order) error {
	if order == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	if order.ID == 0 {
		return os.repo.Insert(order)
	} else {
		persist := os.GetByOrderId(order.ID)
		if persist == nil || persist.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}

		order.ID = persist.ID
		return os.repo.Update(order)
	}
}

func (os *orderService) GetAll() []*model.Order {
	orders := os.repo.FindMore("1=1")
	if orders != nil {
		return orders.([]*model.Order)
	}

	return nil
}

func (os *orderService) Statistics(userId uint64) []*model.Order {
	orders := os.repo.Statistics("user_id = ?", userId)

	return orders
}

func (os *orderService) GetByUserId(userId uint64) []*model.Order {
	orders := os.repo.FindMore("user_id = ?", userId)
	if orders != nil {
		return orders.([]*model.Order)
	}

	return nil
}

func (os *orderService) GetByUserIdAndStatus(userId uint64, status int8) []*model.Order {
	orders := os.repo.FindMore("user_id = ? AND status = ?", userId, status)
	if orders != nil {
		return orders.([]*model.Order)
	}

	return nil
}

func (os *orderService) GetByOrderId(orderId uint64) *model.Order {
	if orderId == 0 {
		return nil
	}
	order := os.repo.FindSingle("id = ?", orderId)
	if order != nil {
		return order.(*model.Order)
	}

	return nil
}

func (os *orderService) GetByOrderIdAndUserId(orderId uint64, userId uint64) *model.Order {
	if orderId == 0 {
		return nil
	}
	order := os.repo.FindSingle("id = ? AND user_id = ?", orderId, userId)
	if order != nil {
		return order.(*model.Order)
	}

	return nil
}

func (os *orderService) GetByOutTradNo(outTradeNo string) *model.Order {
	if outTradeNo == "" {
		return nil
	}
	order := os.repo.FindSingle("out_trade_no = ?", outTradeNo).(*model.Order)
	return order
}

func (os *orderService) DeleteByID(id uint64) error {
	order := os.repo.FindOne(id).(*model.Order)
	if order == nil || order.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := os.repo.Delete(order)
	return err
}

func (os *orderService) GetPage(page int, pageSize int, order *model.Order) *helper.PageBean {
	andCons := make(map[string]interface{})

	if order != nil && order.ID != 0 {
		andCons["id = ?"] = order.ID
	}
	pageBean := os.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}

func (os *orderService) GetPageByMonth(page int, pageSize int, order *model.Order, month string) *helper.PageBean {
	andCons := make(map[string]interface{})

	if order != nil {
		if order.UserID != 0 {
			andCons["user_id = ?"] = order.UserID
		}
		if order.Status != -1 {
			andCons["status = ?"] = order.Status
		}
	}

	if month != "" {
		monthBegin, _ := time.Parse("2006-01-02", month+"-01")
		andCons["order.updated_at >= ?"] = monthBegin
		monthEnd := monthBegin.AddDate(0, 1, 0)
		//timeUnix := monthEnd.Unix() + 86400 - 1
		//end := time.Unix(timeUnix, 0)
		//fmt.Println("end:", end)
		andCons["order.updated_at < ?"] = monthEnd
	}
	
	pageBean := os.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}
