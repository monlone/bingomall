package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type UserOrderDetailRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

type userOrderDetailRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var userOrderDetailRepoIns = &userOrderDetailRepository{}

// 实例化存储对象
func UserOrderDetailRepositoryInstance(db *gorm.DB) UserOrderDetailRepository {
	userOrderDetailRepoIns.db = db
	return userOrderDetailRepoIns
}

// 新增
func (r *userOrderDetailRepository) Insert(userOrderDetail interface{}) error {
	err := r.db.Create(userOrderDetail).Error
	return err
}

// 更新
func (r *userOrderDetailRepository) Update(userOrderDetail interface{}) error {
	err := r.db.Save(userOrderDetail).Error
	return err
}

// 删除
func (r *userOrderDetailRepository) Delete(userOrderDetail interface{}) error {
	err := r.db.Delete(userOrderDetail).Error
	return err
}

// 根据 id 查询
func (r *userOrderDetailRepository) FindOne(id uint64) interface{} {
	var userOrderDetail model.UserOrderDetail
	r.db.Where("user_order_detail_id = ?", id).First(&userOrderDetail)
	if userOrderDetail.OrderId == 0 {
		return nil
	}
	return &userOrderDetail
}

// 条件查询返回单值
func (r *userOrderDetailRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var userOrderDetail model.UserOrderDetail
	r.db.Preload("ShopDetail").Where(condition, params...).First(&userOrderDetail)
	if userOrderDetail.OrderId == 0 {
		return nil
	}
	return &userOrderDetail
}

// 条件查询返回多值
func (r *userOrderDetailRepository) FindMore(condition string, params ...interface{}) interface{} {
	userOrderDetails := make([]*model.UserOrderDetail, 0)
	r.db.Where(condition, params...).Find(&userOrderDetails)
	return userOrderDetails
}

// 分页查询
func (r *userOrderDetailRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.UserOrderDetail, 0)
	rows2 := make([]model.UserOrderDetail, 0)
	if andCons != nil && len(andCons) > 0 {
		for k, v := range andCons {
			r.db = r.db.Where(k, v)
		}
	}
	if orCons != nil && len(orCons) > 0 {
		for k, v := range orCons {
			r.db = r.db.Or(k, v)
		}
	}
	r.db.Preload("ShopDetail").Preload("ProductDetail").Limit(pageSize).Offset((page - 1) * pageSize).Order("updated_at desc").Find(&rows)
	r.db.Find(&rows2).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
