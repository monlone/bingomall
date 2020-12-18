package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
	"time"
)

// moneyLog_service 接口
type MoneyLogService interface {
	/** 保存或修改 */
	SaveOrUpdate(moneyLog *model.MoneyLog) error

	Save(moneyLog *model.MoneyLog) error

	Update(moneyLog *model.MoneyLog) error

	/** 根据 moneyLog_id 查询 */
	GetByMoneyLogID(moneyLogId uint64) *model.MoneyLog

	/** */
	GetAllByShopId(shopId uint64) []*model.MoneyLog

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	DeleteByMoneyLogID(moneyLogId uint64) error

	/** 查询所有  */
	GetAll() []*model.MoneyLog

	/** 分页查询 */
	GetPage(page int, pageSize int, moneyLog *model.MoneyLog) *helper.PageBean

	GetPageByMonth(pageOjb *model.PageObject, moneyLog *model.MoneyLog) *helper.PageBean

	GetPageByMonthWithMerchant(pageOjb *model.PageObject, moneyLog *model.MoneyLog) *helper.PageBean
}

var moneyLogServiceIns = &moneyLogService{}

// 获取 moneyLogService 实例
func MoneyLogServiceInstance(repo repositories.MoneyLogRepository) MoneyLogService {
	moneyLogServiceIns.repo = repo
	return moneyLogServiceIns
}

// 结构体
type moneyLogService struct {
	/** 存储对象 */
	repo repositories.MoneyLogRepository
}

func (mls *moneyLogService) GetByMoneyLogOpenId(moneyLogId uint64) *model.MoneyLog {
	moneyLog := mls.repo.FindSingle("id = ?", moneyLogId)
	if moneyLog != nil {
		return moneyLog.(*model.MoneyLog)
	}
	return nil
}

func (mls *moneyLogService) SaveOrUpdate(moneyLog *model.MoneyLog) error {
	if moneyLog == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	if moneyLog.ID == 0 {
		// 添加
		return mls.repo.Insert(moneyLog)
	} else {
		// 修改
		persist := mls.GetByMoneyLogID(moneyLog.ID)
		if persist == nil || persist.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}

		moneyLog.ID = persist.ID
		return mls.repo.Update(moneyLog)
	}
}

func (mls *moneyLogService) Save(moneyLog *model.MoneyLog) error {
	if moneyLog == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return mls.repo.Insert(moneyLog)
}

func (mls *moneyLogService) Update(moneyLog *model.MoneyLog) error {
	if moneyLog == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	persist := mls.GetByMoneyLogID(moneyLog.ID)
	if persist == nil || persist.ID == 0 {
		return errors.New(helper.StatusText(helper.UpdateObjIsNil))
	}

	moneyLog.ID = persist.ID
	return mls.repo.Update(moneyLog)
}

func (mls *moneyLogService) GetAll() []*model.MoneyLog {
	moneyLogs := mls.repo.FindMore("1=1").([]*model.MoneyLog)
	return moneyLogs
}

func (mls *moneyLogService) GetAllByShopId(shopId uint64) []*model.MoneyLog {
	if shopId == 0 {
		return nil
	}
	moneyLog := mls.repo.FindMore("shop_id = ?", shopId).([]*model.MoneyLog)
	return moneyLog
}

func (mls *moneyLogService) GetByMoneyLogID(moneyLogId uint64) *model.MoneyLog {
	if moneyLogId == 0 {
		return nil
	}
	moneyLog := mls.repo.FindSingle("id = ?", moneyLogId).(*model.MoneyLog)
	return moneyLog
}

func (mls *moneyLogService) DeleteByID(id uint64) error {
	moneyLog := mls.repo.FindOne(id).(*model.MoneyLog)
	if moneyLog == nil || moneyLog.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := mls.repo.Delete(moneyLog)
	return err
}

func (mls *moneyLogService) DeleteByMoneyLogID(moneyLogId uint64) error {
	if moneyLogId == 0 {
		return nil
	}
	moneyLog := mls.repo.FindSingle("moneyLog_id = ?", moneyLogId).(*model.MoneyLog)
	if moneyLog == nil || moneyLog.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := mls.repo.Delete(moneyLog)
	return err
}

func (mls *moneyLogService) GetPage(page int, pageSize int, moneyLog *model.MoneyLog) *helper.PageBean {
	andCons := make(map[string]interface{})

	if moneyLog != nil && moneyLog.RelationUserID != 0 {
		andCons["relation_user_id = ?"] = moneyLog.RelationUserID
	}
	pageBean := mls.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}

func (mls *moneyLogService) GetPageByMonth(pageObj *model.PageObject, moneyLog *model.MoneyLog) *helper.PageBean {
	andCons := make(map[string]interface{})

	//if moneyLog != nil && moneyLog.Status != 0 {
	//	andCons["status = ?"] = moneyLog.Status
	//}

	if moneyLog != nil && moneyLog.UserID != 0 {
		andCons["user_id = ?"] = moneyLog.UserID
	}

	if pageObj.Month != "" {
		monthBegin, _ := time.Parse("2006-01-02", pageObj.Month+"-01")
		andCons["updated_at >= ?"] = monthBegin
		monthEnd := monthBegin.AddDate(0, 1, 0)
		andCons["updated_at < ?"] = monthEnd
	}

	pageBean := mls.repo.FindPage(pageObj.Page, pageObj.PageSize, andCons, nil)
	return pageBean
}

func (mls *moneyLogService) GetPageByMonthWithMerchant(pageObj *model.PageObject, moneyLog *model.MoneyLog) *helper.PageBean {
	andCons := make(map[string]interface{})

	//if moneyLog != nil && moneyLog.Status != 0 {
	//	andCons["status = ?"] = moneyLog.Status
	//}

	if moneyLog != nil && moneyLog.UserID != 0 {
		andCons["user_id = ?"] = moneyLog.UserID
	}

	if pageObj.Month != "" {
		monthBegin, _ := time.Parse("2006-01-02", pageObj.Month+"-01")
		andCons["updated_at >= ?"] = monthBegin
		monthEnd := monthBegin.AddDate(0, 1, 0)
		andCons["updated_at < ?"] = monthEnd
	}

	pageBean := mls.repo.FindPageWithMerchant(pageObj.Page, pageObj.PageSize, andCons, nil)
	return pageBean
}
