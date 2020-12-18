package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type UserShopRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

type userShopRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var userShopRepoIns = &userShopRepository{}

// 实例化存储对象
func UserShopRepositoryInstance(db *gorm.DB) UserShopRepository {
	userShopRepoIns.db = db
	return userShopRepoIns
}

// 新增
func (r *userShopRepository) Insert(userShop interface{}) error {
	err := r.db.Create(userShop).Error
	return err
}

// 更新
func (r *userShopRepository) Update(userShop interface{}) error {
	err := r.db.Save(userShop).Error
	return err
}

// 删除
func (r *userShopRepository) Delete(userShop interface{}) error {
	err := r.db.Delete(userShop).Error
	return err
}

// 根据 id 查询
func (r *userShopRepository) FindOne(id uint64) interface{} {
	var userShop model.UserShop
	r.db.Where("user_id = ?", id).First(&userShop)
	if userShop.UserID == 0 {
		return nil
	}
	return &userShop
}

// 条件查询返回单值
func (r *userShopRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var userShop model.UserShop
	r.db.Where(condition, params...).First(&userShop)
	if userShop.UserID == 0 {
		return nil
	}
	return &userShop
}

// 条件查询返回多值
func (r *userShopRepository) FindMore(condition string, params ...interface{}) interface{} {
	userShops := make([]*model.UserShop, 0)
	r.db.Where(condition, params...).Find(&userShops)
	return userShops
}

// 分页查询
func (r *userShopRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.UserShop, 0)
	rows2 := make([]model.UserShop, 0)
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
	r.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("updated_at desc").Find(&rows)
	r.db.Find(&rows2).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
