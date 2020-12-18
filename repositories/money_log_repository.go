package repositories

import (
	helper "bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type MoneyLogRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository

	FindPageWithMerchant(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean)
}

type moneyLogRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var moneyLogRepoIns = &moneyLogRepository{}

// 实例化存储对象
func MoneyLogRepositoryInstance(db *gorm.DB) MoneyLogRepository {
	moneyLogRepoIns.db = db
	return moneyLogRepoIns
}

// 新增
func (r *moneyLogRepository) Insert(moneyLog interface{}) error {
	err := r.db.Create(moneyLog).Error
	return err
}

// 更新
func (r *moneyLogRepository) Update(moneyLog interface{}) error {
	err := r.db.Save(moneyLog).Error
	return err
}

// 删除
func (r *moneyLogRepository) Delete(moneyLog interface{}) error {
	err := r.db.Delete(moneyLog).Error
	return err
}

// 根据 id 查询
func (r *moneyLogRepository) FindOne(id uint64) interface{} {
	var moneyLog model.MoneyLog
	r.db.Where("moneyLog_id = ?", id).First(&moneyLog)
	if moneyLog.ID == 0 {
		return nil
	}
	return &moneyLog
}

// 根据 id 查询
func (r *moneyLogRepository) FindByShopId(shopId string) interface{} {
	var moneyLog model.MoneyLog
	r.db.Where("shop_id = ?", shopId).Find(&moneyLog)
	if moneyLog.ID == 0 {
		return nil
	}
	return &moneyLog
}

// 条件查询返回单值
func (r *moneyLogRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var moneyLog model.MoneyLog
	r.db.Where(condition, params...).First(&moneyLog)
	if moneyLog.ID == 0 {
		return nil
	}
	return &moneyLog
}

// 条件查询返回多值
func (r *moneyLogRepository) FindMore(condition string, params ...interface{}) interface{} {
	moneyLogs := make([]*model.MoneyLog, 0)
	r.db.Where(condition, params...).Find(&moneyLogs)
	return moneyLogs
}

// 分页查询
func (r *moneyLogRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.MoneyLog, 0)
	rows2 := make([]model.MoneyLog, 0)
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
	r.db.Find(&rows2).Count(&total)
	r.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("updated_at desc").Find(&rows)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}

func (r *moneyLogRepository) FindPageWithMerchant(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.MoneyLog, 0)
	rows2 := make([]model.MoneyLog, 0)
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
	r.db.Find(&rows2).Count(&total)
	r.db.Preload("MerchantSummary").Limit(pageSize).Offset((page - 1) * pageSize).Order("updated_at desc").Find(&rows)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
