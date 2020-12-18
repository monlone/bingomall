package repositories

import (
	"fmt"
	helper "bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
	"time"
)

type ShopRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository

	ProductList(page int, pageSize int, shopId string) (pageBean *helper.PageBean)

	ShopListNearby(page int, pageSize int, shop *model.ShopDetailDistance) (pageBean *helper.PageBean)

	RawSqlInsert(shop *model.Shop) error

	RawSqlUpdate(shop *model.Shop) error
}

type shopRepository struct {
	/** 数据库连接对象 */
	db *gorm.DB
}

var shopRepoIns = &shopRepository{}

// 实例化存储对象
func ShopRepositoryInstance(db *gorm.DB) ShopRepository {
	shopRepoIns.db = db
	return shopRepoIns
}

// 新增
func (r *shopRepository) Insert(shop interface{}) error {
	_ = r.db.Callback().Create().Register("update_created_at", UpdateCreated)
	err := r.db.Create(shop).Error
	return err
}

func UpdateCreated(db *gorm.DB) {
	if db.Migrator().HasColumn(model.Shop{}, "created_at") {
		db.Statement.SetColumn("created_at", time.Now())
	}
}

// 更新
func (r *shopRepository) Update(shop interface{}) error {
	err := r.db.Save(shop).Error
	return err
}

// 删除
func (r *shopRepository) Delete(shop interface{}) error {
	err := r.db.Delete(shop).Error
	return err
}

// 根据 id 查询
func (r *shopRepository) FindOne(id uint64) interface{} {
	var shop model.Shop
	r.db.Where("shop_id = ?", id).First(&shop)
	if shop.ID == 0 {
		return nil
	}
	return &shop
}

// 根据 id 查询
func (r *shopRepository) FindByShopId(shopId string) interface{} {
	var shop model.Shop
	r.db.Where("shop_id = ?", shopId).Find(&shop)
	if shop.ID == 0 {
		return nil
	}
	return &shop
}

// 条件查询返回单值
func (r *shopRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var shop model.Shop
	r.db.Where(condition, params...).First(&shop)
	if shop.ID == 0 {
		return nil
	}
	return &shop
}

// 条件查询返回多值
func (r *shopRepository) FindMore(condition string, params ...interface{}) interface{} {
	shops := make([]*model.Shop, 0)
	r.db.Where(condition, params...).Find(&shops)
	return shops
}

// 分页查询
func (r *shopRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]model.Shop, 0)
	rows2 := make([]model.Shop, 0)
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
	r.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("order_weight desc").Order("begin_at desc").Find(&rows)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}

func (r *shopRepository) ProductList(page int, pageSize int, shopId string) (pageBean *helper.PageBean) {
	var shop model.Shop
	total := int64(0)

	r.db.Preload("ProductList", func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset((page - 1) * pageSize)
	}).Where("shop_id = ? ", shopId).First(&shop)

	r.db.Model(model.Shop{}).Where("shop_id = ?", shopId).Count(&total)

	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: shop}
}

func (r *shopRepository) RawSqlInsert(shop *model.Shop) error {
	r.db.Model(model.Shop{}).Exec("INSERT INTO shop (`id`,`title`,`location`,`merchant_id`,`code`,`star`,"+
		"`business_time`,`attention_num`,`phone`,`address`,`address_description`,`description`,`detail_describe`,"+
		"`logo`,`image_url`,`content`, `longitude`, `latitude`, `created_at`,`updated_at`) "+
		"VALUES (?, ?, GeomFromText(?), ?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		shop.ID, shop.Title, shop.Coordinates.String(), shop.MerchantId, shop.Code, shop.Star, shop.BusinessTime,
		shop.AttentionNum, shop.Phone, shop.Address, shop.AddressDescription, shop.Description, shop.DetailDescribe,
		shop.Logo, shop.ImageUrl, shop.Content, shop.Longitude, shop.Latitude, time.Now(), time.Now())

	r.db.Model(model.Shop{}).Exec("UPDATE  shop SET `coordinates` = NULL WHERE id = ?", shop.ID)

	return nil
}

func (r *shopRepository) RawSqlUpdate(shop *model.Shop) error {
	r.db.Model(model.Shop{}).Exec(
		"UPDATE shop SET `coordinates` = NULL, `location`= GeomFromText(?) WHERE id = ? ",
		shop.Coordinates.String(), shop.ID)

	return nil
}

func (r *shopRepository) ShopListNearby(page int, pageSize int, shop *model.ShopDetailDistance) (pageBean *helper.PageBean) {
	rows, _ := r.db.Raw(
		"SELECT id, title, logo, longitude, latitude,"+
			"ST_Distance_Sphere ( Point(?,?), location ) AS distance"+
			" FROM shop WHERE ST_Distance_Sphere(POINT(?,?),location) < 200000 "+
			" ORDER BY distance asc, order_weight desc LIMIT 100",
		shop.Longitude, shop.Latitude, shop.Longitude, shop.Latitude).Rows()

	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	data := make([]model.ShopDetailDistance, 20)
	for rows.Next() {
		var shopDetail model.ShopDetailDistance
		_ = r.db.ScanRows(rows, &shopDetail)
		fmt.Println("shopDetail:", shopDetail)
		data = append(data, shopDetail)
	}

	return &helper.PageBean{Page: page, PageSize: pageSize, Total: 0, Rows: data}
}
