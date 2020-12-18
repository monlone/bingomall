package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// sku repository 接口
type SkuRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

var skuRepoIns = &skuRepository{}

// 实例化 存储对象
func SkuRepositoryInstance(db *gorm.DB) SkuRepository {
	skuRepoIns.db = db
	return skuRepoIns
}

type skuRepository struct {
	db *gorm.DB
}

func (cr *skuRepository) Insert(sku interface{}) error {
	err := cr.db.Create(sku).Error
	return err
}

func (cr *skuRepository) Update(sku interface{}) error {
	err := cr.db.Save(sku).Error
	return err
}

func (cr *skuRepository) Delete(sku interface{}) error {
	err := cr.db.Delete(sku).Error
	return err
}

func (cr *skuRepository) FindOne(id uint64) interface{} {
	var sku model.Sku
	cr.db.Where("id = ?", id).First(&sku)
	return &sku
}

func (cr *skuRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var sku model.Sku
	cr.db.Where(condition, params...).First(&sku)
	return &sku
}

func (cr *skuRepository) FindMore(condition string, params ...interface{}) interface{} {
	categories := make([]*model.Sku, 0)
	cr.db.Where(condition, params...).Find(&categories)
	return categories
}

func (cr *skuRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Sku, 0)
	if andCons != nil && len(andCons) > 0 {
		for k, v := range andCons {
			cr.db = cr.db.Where(k, v)
		}

	}
	if orCons != nil && len(orCons) > 0 {
		for k, v := range orCons {
			cr.db = cr.db.Or(k, v)
		}
	}
	cr.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&rows).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
