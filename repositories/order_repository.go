package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
	Statistics(condition string, params ...interface{}) []*model.Order
}

type orderRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var orderRepoIns = &orderRepository{}

// 实例化存储对象
func OrderRepositoryInstance(db *gorm.DB) OrderRepository {
	orderRepoIns.db = db
	return orderRepoIns
}

// 新增
func (or *orderRepository) Insert(order interface{}) error {
	err := or.db.Create(order).Error
	return err
}

// 更新
func (or *orderRepository) Update(order interface{}) error {
	err := or.db.Save(order).Error
	return err
}

// 删除
func (or *orderRepository) Delete(order interface{}) error {
	err := or.db.Delete(order).Error
	return err
}

// 根据 id 查询
func (or *orderRepository) FindOne(id uint64) interface{} {
	var order model.Order
	or.db.Where("order_id = ?", id).First(&order)
	if order.ID == 0 {
		return nil
	}
	return &order
}

// 条件查询返回单值
func (or *orderRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var order model.Order
	or.db.Preload("OrderProduct").Preload("OrderProduct.Shop").Where(condition, params...).First(&order)
	if order.ID == 0 {
		return nil
	}
	return &order
}

func (or *orderRepository) Statistics(condition string, params ...interface{}) []*model.Order {
	orders := make([]*model.Order, 0)
	or.db.Where(condition, params...).Select("count(status) as total, user_id, status").Group("status").Find(&orders)
	return orders
}

// 条件查询返回多值
func (or *orderRepository) FindMore(condition string, params ...interface{}) interface{} {
	orders := make([]*model.Order, 0)
	or.db.Where(condition, params...).Find(&orders)
	return orders
}

// 分页查询
func (or *orderRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Order, 0)
	if andCons != nil && len(andCons) > 0 {
		for k, v := range andCons {
			or.db = or.db.Where(k, v)
		}
	}
	if orCons != nil && len(orCons) > 0 {
		for k, v := range orCons {
			or.db = or.db.Or(k, v)
		}
	}
	or.db.Preload("OrderProduct").Preload("OrderProduct.Product").Limit(pageSize).Offset((page - 1) * pageSize).Order("updated_at desc").Find(&rows).Count(&total)

	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
