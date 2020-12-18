package helper

// 统一登陆模型
type LoginParams struct {
	/** 用户名或者账号 */
	Username string `json:"username" form:"username"`

	/** 密码 */
	Password string `json:"password" form:"password"`

	/** 验证码 */
	Code string `json:"code" form:"code"`
}

// 统一 json 结构体
type JsonObject struct {
	/** 状态码 */
	Code int `json:"code"`

	/** 内容体 */
	Content interface{} `json:"data"`

	/** 消息 */
	Message string `json:"message"`
}

// 全局分页对象
type PageBean struct {
	/** 当前页  */
	Page int `json:"page"`

	/** 每页显示的最大行数 */
	PageSize int `json:"pageSize"`

	/** 总记录数 */
	Total int64 `json:"total"`

	/** 每行的数据 */
	Rows interface{} `json:"rows"`
}
