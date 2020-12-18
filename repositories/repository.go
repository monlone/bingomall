package repositories

import (
	"bingomall/helpers"
)

// 基础 repository 接口
type Repository interface {
	// 新增
	Insert(m interface{}) error

	// 更新
	Update(m interface{}) error

	// 删除
	Delete(m interface{}) error

	// 根据 id 查询
	FindOne(id uint64) interface{}

	// 根据条件 查询单条记录
	FindSingle(condition string, params ...interface{}) interface{}

	// 根据条件查询多个结果
	FindMore(condition string, params ...interface{}) interface{}

	/** 分页查询 */
	FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean)
}
