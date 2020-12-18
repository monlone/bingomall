package repositories

import (
	helper "bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

type productRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var productRepoIns = &productRepository{}

// 实例化存储对象
func ProductRepositoryInstance(db *gorm.DB) ProductRepository {
	productRepoIns.db = db
	return productRepoIns
}

// 新增
func (r *productRepository) Insert(product interface{}) error {
	err := r.db.Create(product).Error
	return err
}

// 更新
func (r *productRepository) Update(product interface{}) error {
	err := r.db.Save(product).Error
	return err
}

// 删除
func (r *productRepository) Delete(product interface{}) error {
	err := r.db.Delete(product).Error
	return err
}

// 根据 id 查询
func (r *productRepository) FindOne(id uint64) interface{} {
	var product model.Product
	r.db.Where("product_id = ?", id).First(&product)
	if product.ID == 0 {
		return nil
	}
	return &product
}

// 根据 id 查询
func (r *productRepository) FindByProductId(productId uint64) interface{} {
	var product model.Product
	r.db.Where("product_id = ?", productId).Find(&product)
	if product.ID == 0 {
		return nil
	}
	return &product
}

// 条件查询返回单值
func (r *productRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var product model.Product
	r.db.Where(condition, params...).First(&product)
	if product.ID == 0 {
		return nil
	}
	return &product
}

// 条件查询返回多值
func (r *productRepository) FindMore(condition string, params ...interface{}) interface{} {
	products := make([]*model.Product, 0)
	r.db.Preload("OptionList").Where(condition, params...).Find(&products)
	return products
}

// 分页查询
func (r *productRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.Product, 0)
	rows2 := make([]model.Product, 0)
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
	r.db.Find(&rows2).Count(&total)
	r.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("id desc").Find(&rows)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
