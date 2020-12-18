package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// category repository 接口
type CategoryRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

var categoryRepoIns = &categoryRepository{}

// 实例化 存储对象
func CategoryRepositoryInstance(db *gorm.DB) CategoryRepository {
	categoryRepoIns.db = db
	return categoryRepoIns
}

type categoryRepository struct {
	db *gorm.DB
}

func (cr *categoryRepository) Insert(category interface{}) error {
	err := cr.db.Create(category).Error
	return err
}

func (cr *categoryRepository) Update(category interface{}) error {
	err := cr.db.Save(category).Error
	return err
}

func (cr *categoryRepository) Delete(category interface{}) error {
	err := cr.db.Delete(category).Error
	return err
}

func (cr *categoryRepository) FindOne(id uint64) interface{} {
	var category model.Category
	cr.db.Where("id = ?", id).First(&category)
	return &category
}

func (cr *categoryRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var category model.Category
	cr.db.Where(condition, params...).First(&category)
	return &category
}

func (cr *categoryRepository) FindMore(condition string, params ...interface{}) interface{} {
	categories := make([]*model.Category, 0)
	cr.db.Where(condition, params...).Find(&categories)
	return categories
}

func (cr *categoryRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Category, 0)
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
