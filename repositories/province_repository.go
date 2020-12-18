package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// province repository 接口
type ProvinceRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

var provinceRepoIns = &provinceRepository{}

// 实例化 存储对象
func ProvinceRepositoryInstance(db *gorm.DB) ProvinceRepository {
	provinceRepoIns.db = db
	return provinceRepoIns
}

type provinceRepository struct {
	db *gorm.DB
}

func (pr *provinceRepository) Insert(province interface{}) error {
	err := pr.db.Create(province).Error
	return err
}

func (pr *provinceRepository) Update(province interface{}) error {
	err := pr.db.Save(province).Error
	return err
}

func (pr *provinceRepository) Delete(province interface{}) error {
	err := pr.db.Delete(province).Error
	return err
}

func (pr *provinceRepository) FindOne(id uint64) interface{} {
	var province model.Province
	pr.db.Where("id = ?", id).First(&province)
	return &province
}

func (pr *provinceRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var province model.Province
	pr.db.Where(condition, params...).First(&province)
	return &province
}

func (pr *provinceRepository) FindMore(condition string, params ...interface{}) interface{} {
	categories := make([]*model.Province, 0)
	pr.db.Where(condition, params...).Find(&categories)
	return categories
}

func (pr *provinceRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Province, 0)
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
