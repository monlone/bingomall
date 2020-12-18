package repositories

import (
	"bingomall/helpers"
	"bingomall/models"
	"gorm.io/gorm"
)

// address repository 接口
type AddressRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
	Repository
	SetDefaultAddress(userId uint64, addressId uint64) error
}

var addressRepoIns = &addressRepository{}

// 实例化 存储对象
func AddressRepositoryInstance(db *gorm.DB) AddressRepository {
	addressRepoIns.db = db
	return addressRepoIns
}

type addressRepository struct {
	db *gorm.DB
}

func (ar *addressRepository) Insert(address interface{}) error {
	err := ar.db.Create(address).Error
	return err
}

func (ar *addressRepository) Update(address interface{}) error {
	err := ar.db.Save(address).Error
	return err
}

func (ar *addressRepository) Delete(address interface{}) error {
	err := ar.db.Delete(address).Error
	return err
}

func (ar *addressRepository) FindOne(id uint64) interface{} {
	var address model.Address
	ar.db.Where("id = ?", id).First(&address)
	return &address
}

func (ar *addressRepository) FindSingle(condition string, params ...interface{}) interface{} {
	var address model.Address
	ar.db.Where(condition, params...).First(&address)
	return &address
}

func (ar *addressRepository) FindMore(condition string, params ...interface{}) interface{} {
	address := make([]*model.Address, 0)
	ar.db.Where(condition, params...).Find(&address)
	return address
}

func (ar *addressRepository) FindPage(page int, pageSize int, andCons map[string]interface{}, orCons map[string]interface{}) (pageBean *helper.PageBean) {
	total := int64(0)
	rows := make([]*model.Address, 0)
	if andCons != nil && len(andCons) > 0 {
		for k, v := range andCons {
			ar.db = ar.db.Where(k, v)
		}
	}
	if orCons != nil && len(orCons) > 0 {
		for k, v := range orCons {
			ar.db = ar.db.Or(k, v)
		}
	}
	ar.db.Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at desc").Find(&rows).Count(&total)
	return &helper.PageBean{Page: page, PageSize: pageSize, Total: total, Rows: rows}
}

func (ar *addressRepository) SetDefaultAddress(userId uint64, addressId uint64) error {
	var address model.Address
	m := make(map[string]interface{})
	m["is_default"] = false

	//ar.db.Debug().Model(&address).Where("user_id = ?", userId).Update(m)
	ar.db.Model(&address).Where("user_id = ?", userId).Updates(m)

	address.IsDefault = true
	ar.db.Model(&address).Where("user_id = ? and id = ?", userId, addressId).Select("is_default").Updates(address)

	return nil
}
