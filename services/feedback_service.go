package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// feedback_service 接口
type FeedbackService interface {
	/** 保存或修改 */
	SaveOrUpdate(feedback *model.Feedback) error

	Update(feedback *model.Feedback) error

	/** 根据 id 查询 */
	GetByID(id uint64) *model.Feedback

	/** 根据 feedbackId 查询 */
	GetByFeedbackID(feedbackId uint64) *model.Feedback

	/** 根据title查询 */
	GetByTitle(title string) *model.Feedback

	GetAllByUserID(userId uint64) []*model.Feedback

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有  */
	GetAll() []*model.Feedback

	/** 分页查询 */
	GetPage(page int, pageSize int, feedback *model.Feedback) *helper.PageBean
}

var feedbackServiceIns = &feedbackService{}

// 获取 feedbackService 实例
func FeedbackServiceInstance(repo repositories.FeedbackRepository) FeedbackService {
	feedbackServiceIns.repo = repo
	return feedbackServiceIns
}

// 结构体
type feedbackService struct {
	/** 存储对象 */
	repo repositories.FeedbackRepository
}

func (fs *feedbackService) GetByTitle(title string) *model.Feedback {
	feedback := fs.repo.FindSingle("title = ?", title)
	if feedback != nil {
		return feedback.(*model.Feedback)
	}
	return nil
}

func (fs *feedbackService) SaveOrUpdate(feedback *model.Feedback) error {
	if feedback == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	if feedback.ID == 0 {
		// 添加
		return fs.repo.Insert(feedback)
	} else {
		// 修改
		persist := fs.GetByFeedbackID(feedback.ID)
		if persist == nil || persist.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		feedback.ID = persist.ID
		return fs.repo.Update(feedback)
	}
}

func (fs *feedbackService) Update(feedback *model.Feedback) error {
	if feedback == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	persist := fs.GetByFeedbackID(feedback.ID)
	if persist == nil || persist.ID == 0 {
		return errors.New(helper.StatusText(helper.UpdateObjIsNil))
	}

	feedback.ID = persist.ID
	return fs.repo.Update(feedback)
}

func (fs *feedbackService) GetAll() []*model.Feedback {
	feedback := fs.repo.FindMore("1=1").([]*model.Feedback)
	return feedback
}

func (fs *feedbackService) GetAllByUserID(userId uint64) []*model.Feedback {
	if userId == 0 {
		return nil
	}
	banner := fs.repo.FindMore("user_id = ?", userId).([]*model.Feedback)
	return banner
}

func (fs *feedbackService) GetByID(id uint64) *model.Feedback {
	if id == 0 {
		return nil
	}
	feedback := fs.repo.FindOne(id).(*model.Feedback)
	return feedback
}

func (fs *feedbackService) GetByFeedbackID(feedbackId uint64) *model.Feedback {
	if feedbackId == 0 {
		return nil
	}
	feedback := fs.repo.FindSingle("feedback_id = ?", feedbackId)
	if feedback == nil {
		return nil
	}

	return feedback.(*model.Feedback)
}

func (fs *feedbackService) DeleteByID(id uint64) error {
	feedback := fs.repo.FindOne(id).(*model.Feedback)
	if feedback == nil || feedback.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := fs.repo.Delete(feedback)
	return err
}

func (fs *feedbackService) GetPage(page int, pageSize int, feedback *model.Feedback) *helper.PageBean {
	andCons := make(map[string]interface{})
	if feedback != nil && feedback.Title != "" {
		andCons["title LIKE ?"] = feedback.Title + "%"
	}
	if feedback != nil && feedback.Description != "" {
		andCons["description LIKE ?"] = feedback.Description + "%"
	}
	pageBean := fs.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}
