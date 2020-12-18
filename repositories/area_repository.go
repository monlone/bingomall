package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// area repository 接口
type AreaRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

var areaRepoIns = &areaRepository{}

// 实例化 存储对象
func AreaRepositoryInstance(db *gorm.DB) AreaRepository {
	areaRepoIns.db = db
	return areaRepoIns
}

type areaRepository struct {
	db *gorm.DB
}

func (pr *areaRepository) Insert(area interface{}) error {
	err := pr.db.Create(area).Error
	return err
}

func (pr *areaRepository) Update(area interface{}) error {
	err := pr.db.Save(area).Error
	return err
}

func (pr *areaRepository) Delete(area interface{}) error {
	err := pr.db.Delete(area).Error
	return err
}

func (pr *areaRepository) FindOne(id uint64) interface{} {
	var area model.Area
	pr.db.Where("id = ?", id).First(&area)
	return &area
}

func (pr *areaRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var area model.Area
	pr.db.Where(condition, params...).First(&area)
	return &area
}

func (pr *areaRepository) FindMore(condition string, params ...interface{}) interface{} {
	categories := make([]*model.Area, 0)
	pr.db.Where(condition, params...).Find(&categories)
	return categories
}

func (pr *areaRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Area, 0)
	if andCons != nil && len(andCons) > 0 {
		for k, v := range andCons {
			pr.db = pr.db.Where(k, v)
		}

	}
	if orCons != nil && len(orCons) > 0 {
		for k, v := range orCons {
			pr.db = pr.db.Or(k, v)
		}
	}
	pr.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&rows).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
