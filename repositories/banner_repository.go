package repositories

import (
	helper "bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

type BannerRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

type bannerRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var bannerRepoIns = &bannerRepository{}

// 实例化存储对象
func BannerRepositoryInstance(db *gorm.DB) BannerRepository {
	bannerRepoIns.db = db
	return bannerRepoIns
}

// 新增
func (r *bannerRepository) Insert(banner interface{}) error {
	err := r.db.Create(banner).Error
	return err
}

// 更新
func (r *bannerRepository) Update(banner interface{}) error {
	err := r.db.Save(banner).Error
	return err
}

// 删除
func (r *bannerRepository) Delete(banner interface{}) error {
	err := r.db.Delete(banner).Error
	return err
}

// 根据 id 查询
func (r *bannerRepository) FindOne(id uint64) interface{} {
	var banner model.Banner
	r.db.Where("banner_id = ?", id).First(&banner)
	if banner.ID == 0 {
		return nil
	}
	return &banner
}

// 根据 id 查询
func (r *bannerRepository) FindByShopId(shopId string) interface{} {
	var banner model.Banner
	r.db.Where("shop_id = ?", shopId).Find(&banner)
	if banner.ID == 0 {
		return nil
	}
	return &banner
}

// 条件查询返回单值
func (r *bannerRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var banner model.Banner
	r.db.Where(condition, params...).First(&banner)
	if banner.ID == 0 {
		return nil
	}
	return &banner
}

// 条件查询返回多值
func (r *bannerRepository) FindMore(condition string, params ...interface{}) interface{} {
	banners := make([]*model.Banner, 0)
	r.db.Where(condition, params...).Find(&banners)
	return banners
}

// 分页查询
func (r *bannerRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.Banner, 0)
	rows2 := make([]model.Banner, 0)
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
	r.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("begin_at desc").Find(&rows)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
