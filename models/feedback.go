package model

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Feedback struct {
	/** 主键 */
	Model

	//FeedbackID string `gorm:"type:varchar(36);column:feedback_id;not null" json:"feedback_id" form:"feedback_id"`

	/** 提交反馈的userId*/
	UserID uint64 `gorm:"type:bigint;column:user_id;not null" json:"user_id" form:"user_id"`

	Username string `gorm:"type:varchar(36);column:username;" json:"username" form:"username"`

	Phone string `gorm:"type:varchar(20);column:phone;" json:"phone" form:"phone"`

	Mail string `gorm:"type:varchar(36);column:mail;" json:"mail" form:"mail"`

	/**反馈的title*/
	Title string `gorm:"type:varchar(36);column:Title;" json:"title" form:"title"`

	/**反馈描述*/
	Description string `gorm:"type:varchar(32)" json:"description" form:"description" binding:"required"`

	/** 状态  0,未处理，1：正在处理中，3：已完成，4：重新打开  */
	Status uint `gorm:"type:tinyint(1);default:0" json:"status"  form:"status"`

	/**类别：0：建议， 1：投诉，2：其他*/
	Type uint `gorm:"type:tinyint(1)" json:"type"  form:"type"`

	/** 增删改的时间 */
	CrudTime
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetDBByName(constant.DBMerchant).AutoMigrate(&Feedback{})
}

// 插入前生成主键
func (feedback *Feedback) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//feedback.FeedbackID = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (feedback *Feedback) Validator() error {
	if feedback.Description == "" {
		return errors.New("描述不能为空")
	}

	if feedback.Title == "" {
		return errors.New("标题不能为空")
	}

	return nil
}
