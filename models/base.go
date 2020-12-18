package model

import "bingomall/helpers/datetime"

type Model struct {
	ID uint64 `gorm:"column:id;type:bigint;primaryKey;" json:"id"`
}

// 定义 增删改时间 结构体
type CrudTime struct {
	/** 创建时间 */
	CreatedAt datetime.DateTime `json:"createdAt"`

	/** 更新时间 */
	UpdatedAt datetime.DateTime `json:"updatedAt"`

	/** 删除时间 */
	//DeletedAt pq.NullTime `json:"-"`
	DeletedAt *datetime.DateTime `json:"-"`
}

type ExpiredTime struct {
	ExpiredAt datetime.DateTime `json:"expiredAt"`
}

type Duration struct {
	/** 开始时间 */
	BeginAt datetime.DateTime `gorm:"type:timestamp;" json:"beginAt"`

	/** 结束时间 */
	EndAt datetime.DateTime `gorm:"type:timestamp;" json:"endAt"`
}

type StringDuration struct {
	/** 开始时间 */
	BeginAt string `gorm:"type:varchar(20);default:NULL" json:"beginAt"`

	/** 结束时间 */
	EndAt string `gorm:"type:varchar(20);default:NULL" json:"endAt"`
}

type PageObject struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"pageSize" form:"pageSize"`
	Month    string `json:"month" form:"month"`
	Status   int8   `json:"status" form:"status"`
}
