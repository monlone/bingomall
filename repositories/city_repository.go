package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// city repository 接口
type CityRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

var cityRepoIns = &cityRepository{}

// 实例化 存储对象
func CityRepositoryInstance(db *gorm.DB) CityRepository {
	cityRepoIns.db = db
	return cityRepoIns
}

type cityRepository struct {
	db *gorm.DB
}

func (pr *cityRepository) Insert(city interface{}) error {
	err := pr.db.Create(city).Error
	return err
}

func (pr *cityRepository) Update(city interface{}) error {
	err := pr.db.Save(city).Error
	return err
}

func (pr *cityRepository) Delete(city interface{}) error {
	err := pr.db.Delete(city).Error
	return err
}

func (pr *cityRepository) FindOne(id uint64) interface{} {
	var city model.City
	pr.db.Where("id = ?", id).First(&city)
	return &city
}

func (pr *cityRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var city model.City
	pr.db.Where(condition, params...).First(&city)
	return &city
}

func (pr *cityRepository) FindMore(condition string, params ...interface{}) interface{} {
	categories := make([]*model.City, 0)
	pr.db.Where(condition, params...).Find(&categories)
	return categories
}

func (pr *cityRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.City, 0)
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
