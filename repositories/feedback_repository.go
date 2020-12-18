package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type FeedbackRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

type feedbackRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var feedbackRepoIns = &feedbackRepository{}

// 实例化存储对象
func FeedbackRepositoryInstance(db *gorm.DB) FeedbackRepository {
	feedbackRepoIns.db = db
	return feedbackRepoIns
}

// 新增
func (r *feedbackRepository) Insert(feedback interface{}) error {
	err := r.db.Create(feedback).Error
	return err
}

// 更新
func (r *feedbackRepository) Update(feedback interface{}) error {
	err := r.db.Save(feedback).Error
	return err
}

// 删除
func (r *feedbackRepository) Delete(feedback interface{}) error {
	err := r.db.Delete(feedback).Error
	return err
}

// 根据 id 查询
func (r *feedbackRepository) FindOne(id uint64) interface{} {
	var feedback model.Feedback
	r.db.Where("feedback_id = ?", id).First(&feedback)
	if feedback.ID == 0 {
		return nil
	}
	return &feedback
}

// 条件查询返回单值
func (r *feedbackRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var feedback model.Feedback
	r.db.Where(condition, params...).First(&feedback)
	if feedback.ID == 0 {
		return nil
	}
	return &feedback
}

// 条件查询返回多值
func (r *feedbackRepository) FindMore(condition string, params ...interface{}) interface{} {
	feedback := make([]*model.Feedback, 0)
	r.db.Where(condition, params...).Find(&feedback)
	return feedback
}

// 分页查询
func (r *feedbackRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.Feedback, 0)
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
