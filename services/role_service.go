package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// role service 接口
type RoleService interface {
	// 保存或修改
	SaveOrUpdate(role *model.Role) error

	// 根据id
	GetByID(id uint64) *model.Role

	/** 根据角色名称查询 */
	GetByRoleName(roleName string) *model.Role

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有 */
	GetAll() []*model.Role

	/** 分页查询 */
	GetPage(page int, pageSize int, role *model.Role) *helper.PageBean
}

var roleServiceIns = &roleService{}

// 获取 roleService实例
func RoleServiceInstance(repo repositories.RoleRepository) RoleService {
	roleServiceIns.repo = repo
	return roleServiceIns
}

// role service 结构体
type roleService struct {
	// role repository 对象
	repo repositories.RoleRepository
}

func (rs *roleService) SaveOrUpdate(role *model.Role) error {
	if role == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断角色名称是否已存在
	roleByName := rs.GetByRoleName(role.RoleName)
	if role.ID == 0 {
		// 添加
		if roleByName != nil && roleByName.ID != 0 {
			return errors.New(helper.StatusText(helper.ExistSameNameErr))
		}
		return rs.repo.Insert(role)
	} else {
		// 修改
		persist := rs.GetByID(role.ID)
		if persist == nil || persist.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		if roleByName != nil && roleByName.ID != role.ID {
			return errors.New(helper.StatusText(helper.ExistSameNameErr))
		}
		role.ID = persist.ID
		return rs.repo.Update(role)
	}
}

func (rs *roleService) GetByID(id uint64) *model.Role {
	if id == 0 {
		return nil
	}
	role := rs.repo.FindOne(id)
	if role != nil {
		return role.(*model.Role)
	}
	return nil
}

func (rs *roleService) GetByRoleName(roleName string) *model.Role {
	role := rs.repo.FindSingle("role_name = ?", roleName)
	if role != nil {
		return role.(*model.Role)
	}
	return nil
}

func (rs *roleService) DeleteByID(id uint64) error {
	role := rs.repo.FindOne(id).(*model.Role)
	if role == nil || role.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := rs.repo.Delete(role)
	return err
}

func (rs *roleService) GetAll() []*model.Role {
	roles := rs.repo.FindMore("1=1").([]*model.Role)
	return roles
}

func (rs *roleService) GetPage(page int, pageSize int, role *model.Role) *helper.PageBean {
	andCons := make(map[string]interface{})
	if role != nil && role.RoleName != "" {
		andCons["role_name LIKE ?"] = role.RoleName + "%"
	}
	if role != nil && role.RoleKey != "" {
		andCons["role_key LIKE ?"] = role.RoleKey + "%"
	}
	pageBean := rs.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}
