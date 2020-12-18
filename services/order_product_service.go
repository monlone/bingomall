package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// OrderProduct service 接口
type OrderProductService interface {
	// 保存或修改
	SaveOrUpdate(orderProduct *model.OrderProduct) error

	// 根据id查询
	GetByID(id uint64) *model.OrderProduct

	GetByOrderProductId(orderId uint64) *model.OrderProduct

	UpdateOrderProductByOrderId(orderId uint64) error

	Update(order *model.OrderProduct) error

	GetOrderItemsByOrderId(orderProductId uint64) []*model.OrderProduct

	BatchInsert(data []model.OrderProduct) error

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll(shopId uint64) []*model.OrderProduct

	// 分页查询
	GetPage(page int, pageSize int, user *model.OrderProduct) *helper.PageBean
}

// orderProduct service 结构体
type orderProductService struct {
	/** 存储对象 */
	repo repositories.OrderProductRepository
}

func (ops *orderProductService) SaveOrUpdate(orderProduct *model.OrderProduct) error {
	if orderProduct == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if orderProduct.ID == 0 {
		// 添加
		return ops.repo.Insert(orderProduct)
	} else {
		// 修改
		persist := ops.GetByID(orderProduct.ID)
		if persist == nil || orderProduct.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		orderProduct.ID = persist.ID
		return ops.repo.Update(orderProduct)
	}
}

func (ops *orderProductService) GetByID(id uint64) *model.OrderProduct {
	if id == 0 {
		return nil
	}
	orderProduct := ops.repo.FindOne(id).(*model.OrderProduct)
	return orderProduct
}

func (ops *orderProductService) GetByOrderProductId(orderProductId uint64) *model.OrderProduct {
	if orderProductId == 0 {
		return nil
	}
	orderProduct := ops.repo.FindSingle("id = ?", orderProductId).(*model.OrderProduct)
	return orderProduct
}

func (ops *orderProductService) GetByOrderIdAndUserID(orderId uint64, userID uint64) []*model.OrderProduct {
	orderProduct := ops.repo.FindMore("order_id = ? AND user_id = ?", orderId, userID)
	if orderProduct == nil {
		return nil
	}

	return orderProduct.([]*model.OrderProduct)
}

func (ops *orderProductService) Update(orderProduct *model.OrderProduct) error {
	if orderProduct == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	if orderProduct.ID == 0 {
		return ops.repo.Insert(orderProduct)
	} else {
		persist := ops.GetByOrderProductId(orderProduct.ID)
		if persist == nil || persist.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}

		orderProduct.ID = persist.ID
		return ops.repo.Update(orderProduct)
	}
}

func (ops *orderProductService) GetOrderItemsByOrderId(orderId uint64) []*model.OrderProduct {
	if orderId == 0 {
		return nil
	}
	product := ops.repo.FindMore("order_id = ?", orderId)
	if product == nil {
		return nil
	}
	return product.([]*model.OrderProduct)
}

// TODO 这个要实现
func (ops *orderProductService) UpdateOrderProductByOrderId(orderProductId uint64) error {
	panic("implement me")
	return nil
}

func (*orderProductService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (ops *orderProductService) GetAll(shopId uint64) []*model.OrderProduct {
	var orderProduct []*model.OrderProduct

	if shopId != 0 {
		orderProduct = ops.repo.FindMore("shop_id=?", shopId).([]*model.OrderProduct)
	} else {
		orderProduct = ops.repo.FindMore("1=1").([]*model.OrderProduct)
	}

	return orderProduct
}

func (*orderProductService) GetPage(page int, pageSize int, user *model.OrderProduct) *helper.PageBean {
	panic("implement me")
}

func (ops *orderProductService) BatchInsert(data []model.OrderProduct) error {
	err := ops.repo.BatchInsert(data)

	return err
}

var orderProductServiceIns = &orderProductService{}

func OrderProductServiceInstance(repo repositories.OrderProductRepository) OrderProductService {
	orderProductServiceIns.repo = repo
	return orderProductServiceIns
}
