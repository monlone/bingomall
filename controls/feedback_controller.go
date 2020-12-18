package control

import (
	"bingomall/constant"
	helper "bingomall/helpers"
	"bingomall/helpers/convention"
	model "bingomall/models"
	"bingomall/repositories"
	service "bingomall/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// 获取 feedback 列表
// @Summary 获取 feedback 列表,传 feedback_id 按 feedback_id 查，不传就查平台的
// @Tags FeedbackController
// @Accept json
// @Produce json
// @Param shop_id        query string false "产品记录id"
// @Success 200 {object} model.Feedback
// @Router /api/app/feedback_list [get]
func FeedbackList(context *gin.Context) {
	var feedback []*model.Feedback
	userId := context.DefaultPostForm("user_id", "")
	feedbackService := service.FeedbackServiceInstance(repositories.FeedbackRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	feedback = feedbackService.GetAllByUserID(convention.StringToUint64(userId))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: feedback,
		})
}

// 获取 feedback 列表
// @Summary 获取 feedback 列表
// @Tags FeedbackController
// @Accept json
// @Produce json
// @Success 200 {object} model.ListsResponse
// @Router /api/feedback_list_all [get]
func FeedbackListAll(context *gin.Context) {
	var feedback *model.Feedback
	feedbackService := service.FeedbackServiceInstance(repositories.FeedbackRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	pageStr := context.DefaultQuery("page", "0")
	page, _ := strconv.Atoi(pageStr)
	pageSizeStr := context.DefaultQuery("page_size", constant.PageSize)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	feedbackList := feedbackService.GetPage(page, pageSize, feedback)

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: feedbackList,
		})
}

// 获取 feedback 详情
// @Summary  获取 feedback 详情
// @Tags FeedbackController
// @Accept json
// @Produce json
// @Param feedback_id query string true "列表的feedback_id"
// @Success 200 {object} model.Feedback
// @Router /api/feedback_detail [get]
func FeedbackDetail(context *gin.Context) {
	feedbackService := service.FeedbackServiceInstance(repositories.FeedbackRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
	feedbackId := context.Query("feedback_id")
	feedback := feedbackService.GetByFeedbackID(convention.StringToUint64(feedbackId))

	context.JSON(http.StatusOK,
		&helper.JsonObject{
			Code:    0,
			Message: "ok",
			Content: feedback,
		})
}

// 编辑 feedback
// @Summary  feedback 编辑
// @Tags FeedbackController
// @Accept json
// @Produce json
// @Success 200 {object} model.Feedback
// @Router /api/feedback/save [post]
func SaveFeedback(context *gin.Context) {
	userId := helper.GetUserID(context)
	if userId == 0 {
		context.JSON(http.StatusOK,
			&helper.JsonObject{
				Code:    4301,
				Message: "用户id非法",
			})
		return
	}
	feedback := &model.Feedback{}
	if err := context.Bind(feedback); err == nil {
		feedback.UserID = userId
		feedbackService := service.FeedbackServiceInstance(repositories.FeedbackRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		err := feedbackService.SaveOrUpdate(feedback)
		if err == nil {
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    0,
					Message: helper.StatusText(helper.SaveStatusOK),
				})
			return
		} else {
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4302,
					Message: err.Error(),
				})
			return
		}
	} else {
		context.JSON(http.StatusUnprocessableEntity, helper.JsonObject{
			Code:    4303,
			Message: helper.StatusText(helper.BindModelErr),
			Content: err,
		})
	}
}

// 编辑 feedback
// @Summary  feedback 编辑
// @Tags FeedbackController
// @Accept json
// @Produce json
// @Success 200 {object} model.Feedback
// @Router /api/web/merchant/register [post]
// 以后再移走吧，临时放这里
func MerchantRegister(context *gin.Context) {
	feedback := &model.Feedback{}
	if err := context.Bind(feedback); err == nil {
		feedback.UserID = 0
		feedbackService := service.FeedbackServiceInstance(repositories.FeedbackRepositoryInstance(helper.GetDBByName(constant.DBMerchant)))
		err := feedbackService.SaveOrUpdate(feedback)
		if err == nil {
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    0,
					Message: helper.StatusText(helper.SaveStatusOK),
				})
			return
		} else {
			context.JSON(http.StatusOK,
				&helper.JsonObject{
					Code:    4302,
					Message: err.Error(),
				})
			return
		}
	} else {
		context.JSON(http.StatusUnprocessableEntity, helper.JsonObject{
			Code:    4304,
			Message: helper.StatusText(helper.BindModelErr),
			Content: err,
		})
	}
}
