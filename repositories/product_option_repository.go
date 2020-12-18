package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// productOption repository 接口
type ProductOptionRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

var productOptionRepoIns = &productOptionRepository{}

// 实例化 存储对象
func ProductOptionRepositoryInstance(db *gorm.DB) ProductOptionRepository {
	productOptionRepoIns.db = db
	return productOptionRepoIns
}

type productOptionRepository struct {
	db *gorm.DB
}

func (cr *productOptionRepository) Insert(productOption interface{}) error {
	err := cr.db.Create(productOption).Error
	return err
}

func (cr *productOptionRepository) Update(productOption interface{}) error {
	err := cr.db.Save(productOption).Error
	return err
}

func (cr *productOptionRepository) Delete(productOption interface{}) error {
	err := cr.db.Delete(productOption).Error
	return err
}

func (cr *productOptionRepository) FindOne(id uint64) interface{} {
	var productOption model.ProductOption
	cr.db.Where("id = ?", id).First(&productOption)
	return &productOption
}

func (cr *productOptionRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var productOption model.ProductOption
	cr.db.Where(condition, params...).First(&productOption)
	return &productOption
}

func (cr *productOptionRepository) FindMore(condition string, params ...interface{}) interface{} {
	categories := make([]*model.ProductOption, 0)
	cr.db.Where(condition, params...).Find(&categories)
	return categories
}

func (cr *productOptionRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.ProductOption, 0)
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
