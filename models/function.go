package model

import (
	helper "bingomall/helpers"
	"gorm.io/gorm"
)

// 功能菜单结构体
type Function struct {
	/** 主键 id */
	Model

	//FunctionID string `gorm:"type:varchar(36);" form:"functionId"`

	/** 功能名称 */
	FunName string

	/** 访问路径 */
	FunUrl string

	/** 权限功能等级 */
	funLevel int

	/** 是否生成菜单  */
	IsMenu bool

	/** 图标 */
	FunIcon string

	/** 序号 */
	Seq int

	/** 父功能 id */
	PId *string

	/** 父功能 */
	ParentFunction *Function `gorm:"foreignkey:PId;save_associations:false" `

	/** 子功能 */
	ChildFunctions []*Function `gorm:"foreignkey:ID"`

	/** 对应的角色 */
	Roles []*Role `gorm:"many2many:role_functions;" json:"-"`

	/** 增删改的时间 */
	CrudTime
}

// 插入前生成主键
func (function *Function) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//function.FunctionID = id.String()
	return nil
}

func init() {
	// 创建或更新表结构
	_ = helper.GetUserDB().AutoMigrate(&Function{})
}
