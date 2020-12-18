package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

type userRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var userRepoIns = &userRepository{}

// 实例化存储对象
func UserRepositoryInstance(db *gorm.DB) UserRepository {
	userRepoIns.db = db
	return userRepoIns
}

// 新增
func (r *userRepository) Insert(user interface{}) error {
	err := r.db.Create(user).Error
	return err
}

// 更新
func (r *userRepository) Update(user interface{}) error {
	err := r.db.Save(user).Error
	return err
}

// 删除
func (r *userRepository) Delete(user interface{}) error {
	err := r.db.Delete(user).Error
	return err
}

// 根据 id 查询
func (r *userRepository) FindOne(id uint64) interface{} {
	var user model.User
	r.db.Where("id = ?", id).First(&user)
	if user.ID == 0 {
		return nil
	}
	return &user
}

// 条件查询返回单值
func (r *userRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var user model.User
	r.db.Where(condition, params...).First(&user)
	if user.ID == 0 {
		return nil
	}
	return &user
}

// 条件查询返回多值
func (r *userRepository) FindMore(condition string, params ...interface{}) interface{} {
	users := make([]*model.User, 0)
	r.db.Where(condition, params...).Find(&users)
	return users
}

// 分页查询
func (r *userRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.User, 0)
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
	r.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("login_time desc").Find(&rows).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
