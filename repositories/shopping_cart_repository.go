package repositories

import (
	"fmt"
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type ShoppingCartRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

type shoppingCartRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var shoppingCartRepoIns = &shoppingCartRepository{}

// 实例化存储对象
func ShoppingCartRepositoryInstance(db *gorm.DB) ShoppingCartRepository {
	shoppingCartRepoIns.db = db
	return shoppingCartRepoIns
}

// 新增
func (r *shoppingCartRepository) Insert(shoppingCart interface{}) error {
	err := r.db.Create(shoppingCart).Error
	return err
}

// 更新
func (r *shoppingCartRepository) Update(shoppingCart interface{}) error {
	fmt.Println(helper.Json(shoppingCart))
	err := r.db.Save(shoppingCart).Error
	return err
}

// 删除
func (r *shoppingCartRepository) Delete(shoppingCart interface{}) error {
	err := r.db.Delete(shoppingCart).Error
	return err
}

// 根据 id 查询
func (r *shoppingCartRepository) FindOne(id uint64) interface{} {
	var shoppingCart model.ShoppingCart
	r.db.Where("id = ?", id).First(&shoppingCart)
	if shoppingCart.ID == 0 {
		return nil
	}
	return &shoppingCart
}

// 条件查询返回单值
func (r *shoppingCartRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var shoppingCart model.ShoppingCart
	//r.db.Preload("Product").Where(condition, params...).First(&shoppingCart)
	r.db.Where(condition, params...).First(&shoppingCart)
	if shoppingCart.ID == 0 {
		return nil
	}
	return &shoppingCart
}

// 条件查询返回多值
func (r *shoppingCartRepository) FindMore(condition string, params ...interface{}) interface{} {
	shoppingCarts := make([]*model.ShoppingCart, 0)
	//r.db.Preload("OptionList").Preload("Product").Preload("Sku").Where(condition, params...).Find(&shoppingCarts)
	r.db.Preload("OptionList").Preload("Product").Preload("Sku").Preload("OptionList.Option").Where(condition, params...).Find(&shoppingCarts)
	return shoppingCarts
}

// 分页查询
func (r *shoppingCartRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.ShoppingCart, 0)
	rows2 := make([]model.ShoppingCart, 0)
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
