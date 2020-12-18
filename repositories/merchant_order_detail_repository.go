package repositories

import (
	"bingomall/constant"
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type MerchantOrderDetailRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
	Liquidate(userId, shopIds []uint64) error
	LiquidateTotal(userId, shopIds []uint64) (*model.TotalResult, error)
}

type merchantOrderDetailRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var merchantOrderDetailRepoIns = &merchantOrderDetailRepository{}

// 实例化存储对象
func MerchantOrderDetailRepositoryInstance(db *gorm.DB) MerchantOrderDetailRepository {
	merchantOrderDetailRepoIns.db = db
	return merchantOrderDetailRepoIns
}

// 新增
func (r *merchantOrderDetailRepository) Insert(merchantOrderDetail interface{}) error {
	err := r.db.Create(merchantOrderDetail).Error
	return err
}

// 更新
func (r *merchantOrderDetailRepository) Update(merchantOrderDetail interface{}) error {
	err := r.db.Save(merchantOrderDetail).Error
	return err
}

// 删除
func (r *merchantOrderDetailRepository) Delete(merchantOrderDetail interface{}) error {
	err := r.db.Delete(merchantOrderDetail).Error
	return err
}

// 根据 id 查询
func (r *merchantOrderDetailRepository) FindOne(id uint64) interface{} {
	var merchantOrderDetail model.MerchantOrderDetail
	r.db.Where("merchantOrderDetail_id = ?", id).First(&merchantOrderDetail)
	if merchantOrderDetail.OrderId == 0 {
		return nil
	}
	return &merchantOrderDetail
}

// 条件查询返回单值
func (r *merchantOrderDetailRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var merchantOrderDetail model.MerchantOrderDetail
	r.db.Preload("ShopDetail").Where(condition, params...).First(&merchantOrderDetail)
	if merchantOrderDetail.OrderId == 0 {
		return nil
	}
	return &merchantOrderDetail
}

// 条件查询返回多值
func (r *merchantOrderDetailRepository) FindMore(condition string, params ...interface{}) interface{} {
	merchantOrderDetails := make([]*model.MerchantOrderDetail, 0)
	r.db.Where(condition, params...).Find(&merchantOrderDetails)
	return merchantOrderDetails
}

// 分页查询
func (r *merchantOrderDetailRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.MerchantOrderDetail, 0)
	rows2 := make([]model.MerchantOrderDetail, 0)
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
	r.db.Preload("ShopDetail").Preload("ProductDetail").Limit(pageSize).Offset((page - 1) * pageSize).Order("updated_at desc").Find(&rows)
	r.db.Find(&rows2).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}

func (r *merchantOrderDetailRepository) Liquidate(merchantIds, shopIds []uint64) error {
	db := r.db.Model(model.MerchantOrderDetail{}).Where("merchant_id in (?)", merchantIds).Where("shop_id in (?)", shopIds).
		Where("status = ?", constant.MerchantOrderWaitedForLiquidate).
		Update("status", constant.MerchantOrderLiquidated)

	return db.Error
}

func (r *merchantOrderDetailRepository) LiquidateTotal(merchantIds, shopIds []uint64	) (data *model.TotalResult, err error) {
	var n model.TotalResult
	db := r.db.Model(model.MerchantOrderDetail{}).Where("merchant_id in (?)", merchantIds).Where("shop_id in (?)", shopIds).
		Where("status = ?", constant.MerchantOrderWaitedForLiquidate).Select("sum(money) as total_money").Scan(&n)
	err = db.Error
	data = &n
	return
}
