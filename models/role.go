package model

import (
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 角色结构体
type Role struct {
	/** 主键id */
	Model

	//RoleID string `gorm:"type:varchar(36);not null" form:"roleId"`

	/** 角色名称 */
	RoleName string `gorm:"type:varchar(32);unique;not null" form:"role_name" binding:"required"`

	/** 角色类别标识 */
	RoleKey string `gorm:"type:varchar(16);not null" form:"role_key" binding:"required"`

	/** 角色描述 */
	Description string `gorm:"type:varchar(128)" form:"description"`

	/** 角色关联的功能 */
	Functions []*Function `gorm:"many2many:role_functions;" json:"-"`

	/** 增删改的时间 */
	CrudTime
}

// 插入前生成主键
func (role *Role) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//role.RoleID = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	helper.GetUserDB().AutoMigrate(&Role{})
}
