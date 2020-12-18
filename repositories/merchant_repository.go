package repositories

import (
	helper "bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type MerchantRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository

	ShopList(page int, pageSize int, merchantId string) (pageBean *helper.PageBean)
}

type merchantRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var merchantRepoIns = &merchantRepository{}

// 实例化存储对象
func MerchantRepositoryInstance(db *gorm.DB) MerchantRepository {
	merchantRepoIns.db = db
	return merchantRepoIns
}

// 新增
func (r *merchantRepository) Insert(merchant interface{}) error {
	err := r.db.Create(merchant).Error
	return err
}

// 更新
func (r *merchantRepository) Update(merchant interface{}) error {
	err := r.db.Save(merchant).Error
	return err
}

// 删除
func (r *merchantRepository) Delete(merchant interface{}) error {
	err := r.db.Delete(merchant).Error
	return err
}

// 根据 id 查询
func (r *merchantRepository) FindOne(id uint64) interface{} {
	var merchant model.Merchant
	r.db.Where("merchant_id = ?", id).First(&merchant)
	if merchant.ID == 0 {
		return nil
	}
	return &merchant
}

// 根据 id 查询
func (r *merchantRepository) FindByMerchantId(merchantId string) interface{} {
	var merchant model.Merchant
	r.db.Where("merchant_id = ?", merchantId).Find(&merchant)
	if merchant.ID == 0 {
		return nil
	}
	return &merchant
}

// 条件查询返回单值
func (r *merchantRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var merchant model.Merchant
	r.db.Where(condition, params...).First(&merchant)
	if merchant.ID == 0 {
		return nil
	}
	return &merchant
}

// 条件查询返回多值
func (r *merchantRepository) FindMore(condition string, params ...interface{}) interface{} {
	merchants := make([]*model.Merchant, 0)
	r.db.Where(condition, params...).Find(&merchants)
	return merchants
}

// 分页查询
func (r *merchantRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.Merchant, 0)
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
	r.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("updated_at desc").Find(&rows).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}

func (r *merchantRepository) ShopList(page int, pageSize int, merchantId string) (pageBean *helper.PageBean) {
	var merchant model.Merchant
	total := int64(0)

	r.db.Preload("ShopList", func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset((page - 1) * pageSize)
	}).Where("merchant_id = ? ", merchantId).First(&merchant)

	r.db.Model(model.Product{}).Where("merchant_id = ?", merchantId).Count(&total)

	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: merchant}
}
