package repositories

import (
	"fmt"
	"bingomall/helpers"
	"bingomall/helpers/datetime"
	"bingomall/models"
	"gorm.io/gorm"
	"strings"
)

// orderProduct repository 接口
type OrderProductRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
	BatchInsert(data []model.OrderProduct) error
}

var orderProductRepoIns = &orderProductRepository{}

// 实例化 存储对象
func OrderProductRepositoryInstance(db *gorm.DB) OrderProductRepository {
	orderProductRepoIns.db = db
	return orderProductRepoIns
}

type orderProductRepository struct {
	db *gorm.DB
}

func (op *orderProductRepository) Insert(orderProduct interface{}) error {
	err := op.db.Create(orderProduct).Error
	return err
}

func (op *orderProductRepository) Update(orderProduct interface{}) error {
	err := op.db.Save(orderProduct).Error
	return err
}

func (op *orderProductRepository) Delete(orderProduct interface{}) error {
	err := op.db.Delete(orderProduct).Error
	return err
}

func (op *orderProductRepository) FindOne(id uint64) interface{} {
	var orderProduct model.OrderProduct
	op.db.Where("id = ?", id).First(&orderProduct)
	return &orderProduct
}

func (op *orderProductRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var orderProduct model.OrderProduct
	op.db.Where(condition, params...).First(&orderProduct)
	return &orderProduct
}

func (op *orderProductRepository) FindMore(condition string, params ...interface{}) interface{} {
	orderProduct := make([]*model.OrderProduct, 0)
	op.db.Where(condition, params...).Find(&orderProduct)
	return orderProduct
}

func (op *orderProductRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.OrderProduct, 0)
	if andCons != nil && len(andCons) > 0 {
		for k, v := range andCons {
			op.db = op.db.Where(k, v)
		}

	}
	if orCons != nil && len(orCons) > 0 {
		for k, v := range orCons {
			op.db = op.db.Or(k, v)
		}
	}
	op.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&rows).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
func (op *orderProductRepository) BatchInsert(data []model.OrderProduct) error {
	sql := "INSERT INTO `order_product` (`order_id`,`shop_id`,`product_id`,`sku_id`, `price`, `score`, `number`, `money`, `created_at`) VALUES "
	// 循环data数组,组合sql语句
	var sqlArr []string
	var sqlStr string
	for _, v := range data {
		sqlTemp := fmt.Sprintf("('%d','%d', '%d', '%d', %d, %d, %d, %d, '%s')",
			v.OrderId, v.ShopId, v.ProductId, v.SkuId, v.Price, v.Score, v.Number, v.Number*v.Price, datetime.String(v.CreatedAt))
		sqlArr = append(sqlArr, sqlTemp)
	}
	sqlStr = strings.Join(sqlArr, ",")
	sql = sql + sqlStr + ";"

	err := op.db.Exec(sql).Error
	return err
}
