package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// 接口
type RoleRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

// 结构体
type roleRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var roleRepoIns = &roleRepository{}

// 实例化存储对象
func RoleRepositoryInstance(db *gorm.DB) RoleRepository {
	roleRepoIns.db = db
	return roleRepoIns
}

// 新增
func (r *roleRepository) Insert(role interface{}) error {
	err := r.db.Create(role).Error
	return err
}

// 更新
func (r *roleRepository) Update(role interface{}) error {
	err := r.db.Save(role).Error
	return err
}

// 删除
func (r *roleRepository) Delete(role interface{}) error {
	err := r.db.Delete(role).Error
	return err
}

// 根据id查询
func (r *roleRepository) FindOne(id uint64) interface{} {
	var role model.Role
	r.db.Where("id = ?", id).First(&role)
	if role.ID == 0 {
		return nil
	}
	return &role
}

func (r *roleRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var role model.Role
	r.db.Where(condition, params...).First(&role)
	if role.ID == 0 {
		return nil
	}
	return &role
}

func (r *roleRepository) FindMore(condition string, params ...interface{}) interface{} {
	roles := make([]*model.Role, 0)
	r.db.Where(condition, params...).First(&roles)
	return &roles
}

func (r *roleRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Role, 0)
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
	r.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&rows).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
