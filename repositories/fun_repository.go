package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// function repository 接口
type FunctionReposotory interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

var funRepoIns = &functionRepository{}

// 实例化 存储对象
func FunRepositoryInstance(db *gorm.DB) FunctionReposotory {
	funRepoIns.db = db
	return funRepoIns
}

type functionRepository struct {
	db *gorm.DB
}

func (fr *functionRepository) Insert(function interface{}) error {
	err := fr.db.Create(function).Error
	return err
}

func (fr *functionRepository) Update(function interface{}) error {
	err := fr.db.Save(function).Error
	return err
}

func (fr *functionRepository) Delete(function interface{}) error {
	err := fr.db.Delete(function).Error
	return err
}

func (fr *functionRepository) FindOne(id uint64) interface{} {
	var function model.Function
	fr.db.Where("id = ?", id).First(&function)
	return &function
}

func (fr *functionRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var function model.Function
	fr.db.Where(condition, params...).First(&function)
	return &function
}

func (fr *functionRepository) FindMore(condition string, params ...interface{}) interface{} {
	functions := make([]*model.Function, 0)
	fr.db.Where(condition, params...).Find(&functions)
	return &functions
}

func (fr *functionRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Function, 0)
	if andCons != nil && len(andCons) > 0 {
		for k, v := range andCons {
			fr.db = fr.db.Where(k, v)
		}

	}
	if orCons != nil && len(orCons) > 0 {
		for k, v := range orCons {
			fr.db = fr.db.Or(k, v)
		}
	}
	fr.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&rows).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
