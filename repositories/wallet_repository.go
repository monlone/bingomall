package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// wallet repository 接口
type WalletRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
}

var walletRepoIns = &walletRepository{}

// 实例化 存储对象
func WalletRepositoryInstance(db *gorm.DB) WalletRepository {
	walletRepoIns.db = db
	return walletRepoIns
}

type walletRepository struct {
	db *gorm.DB
}

func (wr *walletRepository) Insert(wallet interface{}) error {
	err := wr.db.Create(wallet).Error
	return err
}

func (wr *walletRepository) Update(wallet interface{}) error {
	err := wr.db.Save(wallet).Error
	return err
}

func (wr *walletRepository) Delete(wallet interface{}) error {
	err := wr.db.Delete(wallet).Error
	return err
}

func (wr *walletRepository) FindOne(id uint64) interface{} {
	var wallet model.Wallet
	wr.db.Where("id = ?", id).First(&wallet)
	return &wallet
}

func (wr *walletRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var wallet model.Wallet
	wr.db.Where(condition, params...).First(&wallet)
	return &wallet
}

func (wr *walletRepository) FindMore(condition string, params ...interface{}) interface{} {
	wallets := make([]*model.Wallet, 0)
	wr.db.Where(condition, params...).Find(&wallets)
	return wallets
}

func (wr *walletRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Wallet, 0)
	if andCons != nil && len(andCons) > 0 {
		for k, v := range andCons {
			wr.db = wr.db.Where(k, v)
		}

	}
	if orCons != nil && len(orCons) > 0 {
		for k, v := range orCons {
			wr.db = wr.db.Or(k, v)
		}
	}
	wr.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&rows).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}
