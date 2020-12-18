package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// option repository 接口
type OptionRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

var optionRepoIns = &optionRepository{}

// 实例化 存储对象
func OptionRepositoryInstance(db *gorm.DB) OptionRepository {
	optionRepoIns.db = db
	return optionRepoIns
}

type optionRepository struct {
	db *gorm.DB
}

func (cr *optionRepository) Insert(option interface{}) error {
	err := cr.db.Create(option).Error
	return err
}

func (cr *optionRepository) Update(option interface{}) error {
	err := cr.db.Save(option).Error
	return err
}

func (cr *optionRepository) Delete(option interface{}) error {
	err := cr.db.Delete(option).Error
	return err
}

func (cr *optionRepository) FindOne(id uint64) interface{} {
	var option model.Option
	cr.db.Where("id = ?", id).First(&option)
	return &option
}

func (cr *optionRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var option model.Option
	cr.db.Where(condition, params...).First(&option)
	return &option
}

func (cr *optionRepository) FindMore(condition string, params ...interface{}) interface{} {
	categories := make([]*model.Option, 0)
	cr.db.Where(condition, params...).Find(&categories)
	return categories
}

func (cr *optionRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Option, 0)
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
