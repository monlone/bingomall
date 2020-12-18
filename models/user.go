package model

import (
	helper "bingomall/helpers"
	"bingomall/helpers/regex"
	"gorm.io/gorm"
	"strings"
	"time"
)

type User struct {
	/** 主键 */
	Model

	/** 用户ID */
	//UserID uint64 `gorm:"type:bigint;column:user_id;not null" json:"userId" form:"userId"`

	/** 姓名 */
	Username string `gorm:"type:varchar(32);unique" json:"username" form:"username" binding:"required"`

	Nickname string `gorm:"type:varchar(32);" json:"nickname" form:"nickname"`

	/** 密码  */
	Password string `gorm:"type:varchar(64)" json:"-" form:"password" binding:"required"`

	/** 电话 */
	Phone string `gorm:"type:varchar(11);unique" json:"phone" form:"phone" binding:"required"`

	/** 头像 */
	Avatar string `gorm:"type:varchar(500)" json:"avatar" form:"avatar"`

	/** 邮件 */
	Email string `gorm:"type:varchar(64);unique" json:"email" form:"email"`

	/** 商户号 */
	MerchantNo string `gorm:"type:varchar(32)" json:"merchant_no" form:"merchant_no"`

	/** 商户名称 */
	//MerchantName string `gorm:"-"`

	/** 标志 1 表示这个账号是由管理方为商户添加的账号 */
	Flag int `json:"-"`

	/** 登陆次数 */
	LogonCount int `json:"-"`

	/** 状态  0 正常  */
	Status int `json:"status"`

	/** 最后一次登陆时间 */
	LoginTime time.Time `gorm:"default:null" json:"loginTime"`

	/** 增删改的时间 */
	CrudTime

	/** 用户对应的角色 */
	//Role *Role `gorm:"foreignkey:RoleID;save_associations:false"`

	/** 用户类型，1：管理员，2：平台员工，3：商户，4：普通用户*/
	Type uint8 `gorm:"tinyint(3)；column:type;default:2" json:"type"`

	/** 外键 */
	RoleId *string `gorm:"type:varchar(36)" form:"role_id"`

	/** 微信app支付的openid **/
	OpenId string `gorm:"type:varchar(64);column:openid;unique" form:"openid" json:"openid" weChat:"openid"`

	/** 微信商户号的openid **/
	MPOpenId string `gorm:"type:varchar(64);column:mp_openid;unique" json:"mp_openid" weChat:"openid"`

	/** 微信的unionid **/
	UnionId string `gorm:"type:varchar(64);column:unionid;" form:"unionid" json:"unionid" weChat:"unionid"`

	/** 微信的 access token **/
	AccessToken string `gorm:"type:varchar(128)" json:"-" weChat:"access_token"`

	/** 微信的 refresh token **/
	RefreshToken string `gorm:"type:varchar(128)" json:"-" weChat:"refresh_token"`

	/** 上一级用户id **/
	MultiLevel uint64 `gorm:"type:bigint" json:"multi_level"`

	/** 积分 **/
	Score int64 `gorm:"type:bigint(20)" json:"score"`
}

// 表结构初始化
func init() {
	// 创建或更新表结构
	_ = helper.GetUserDB().AutoMigrate(&User{})
	// 生成外键约束
	//helper.SQL.Model(&User{}).AddForeignKey("role_id", "role(id)", "no action", "no action")
}

// 插入前生成主键
func (user *User) BeforeCreate(db *gorm.DB) error {
	//id := uuid.NewV4()
	//db.Set("ID", &id)
	//user.ID = id.String()
	return nil
}

// 校验表单中提交的参数是否合法
func (user *User) Validator() error {
	if ok, err := regex.MatchLetterNumMinAndMax(user.Username, 4, 16, "用户名"); !ok {
		return err
	}
	//if ok, err := regex.MatchStrongPassword(user.Password, 6, 13); !ok && strings.TrimSpace(user.RecordID) == "" {
	//	return err
	//}
	if ok, err := regex.IsPhone(user.Phone); !ok {
		return err
	}
	if ok, err := regex.IsEmail(user.Email); !ok && strings.TrimSpace(user.Email) != "" {
		return err
	}
	return nil
}
