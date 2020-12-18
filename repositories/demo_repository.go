package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// demo repository 接口
type DemoRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

var demoRepoIns = &demoRepository{}

// 实例化 存储对象
func DemoRepositoryInstance(db *gorm.DB) DemoRepository {
	demoRepoIns.db = db
	return demoRepoIns
}

type demoRepository struct {
	db *gorm.DB
}

func (cr *demoRepository) Insert(demo interface{}) error {
	err := cr.db.Create(demo).Error
	return err
}

func (cr *demoRepository) Update(demo interface{}) error {
	err := cr.db.Save(demo).Error
	return err
}

func (cr *demoRepository) Delete(demo interface{}) error {
	err := cr.db.Delete(demo).Error
	return err
}

func (cr *demoRepository) FindOne(id uint64) interface{} {
	var demo model.Demo
	cr.db.Where("id = ?", id).First(&demo)
	return &demo
}

func (cr *demoRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var demo model.Demo
	cr.db.Where(condition, params...).First(&demo)
	return &demo
}

func (cr *demoRepository) FindMore(condition string, params ...interface{}) interface{} {
	categories := make([]*model.Demo, 0)
	cr.db.Where(condition, params...).Find(&categories)
	return categories
}

func (cr *demoRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Demo, 0)
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
