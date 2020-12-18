package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// Address service 接口
type AddressService interface {
	// 保存或修改
	SaveOrUpdate(address *model.Address) error

	// 根据id查询
	GetByID(id uint64) *model.Address

	// 根据userId查询
	GetByUserID(userId uint64) []*model.Address

	// 通过userId获取该用户的地址
	GetAddressByAddressId(addressId uint64, UserID uint64) *model.Address

	SetDefaultAddress(addressId uint64, UserID uint64) error

	GetDefaultAddressByUserID(userId uint64) *model.Address

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll(userId uint64) []*model.Address

	// 分页查询
	GetPage(page int, pageSize int, user *model.Address) *helper.PageBean
}

// address service 结构体
type addressService struct {
	/** 存储对象 */
	repo repositories.AddressRepository
}

func (as *addressService) SaveOrUpdate(address *model.Address) error {
	if address == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if address.ID == 0 {
		// 添加
		return as.repo.Insert(address)
	} else {
		// 修改
		persist := as.GetByID(address.ID)
		if persist == nil || address.ID == 0 || persist.UserID != address.UserID {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		address.ID = persist.ID
		return as.repo.Update(address)
	}
}

func (as *addressService) GetByID(id uint64) *model.Address {
	if id == 0 {
		return nil
	}
	address := as.repo.FindOne(id).(*model.Address)
	return address
}

func (as *addressService) GetByUserID(userId uint64) []*model.Address {
	address := as.repo.FindMore("user_id = ?", userId)
	if address == nil {
		return nil
	}

	return address.([]*model.Address)
}

func (as *addressService) GetAddressByAddressId(addressId uint64, userID uint64) *model.Address {
	address := as.repo.FindSingle("id=? AND user_id=?", addressId, userID)
	if address == nil {
		return nil
	}

	return address.(*model.Address)
}

func (as *addressService) GetDefaultAddressByUserID(userId uint64) *model.Address {
	address := as.repo.FindSingle("user_id = ? AND is_default = ?", userId, true)
	if address == nil {
		return nil
	}

	return address.(*model.Address)
}

func (*addressService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (as *addressService) GetAll(shopId uint64) []*model.Address {
	var address []*model.Address

	if shopId != 0 {
		address = as.repo.FindMore("id=?", shopId).([]*model.Address)
	} else {
		address = as.repo.FindMore("1=1").([]*model.Address)
	}

	return address
}

func (*addressService) GetPage(page int, pageSize int, user *model.Address) *helper.PageBean {
	panic("implement me")
}

func (as *addressService) SetDefaultAddress(userId uint64, addressId uint64) error {
	res := as.repo.SetDefaultAddress(userId, addressId)
	return res
}

var addressServiceIns = &addressService{}

func AddressServiceInstance(repo repositories.AddressRepository) AddressService {
	addressServiceIns.repo = repo
	return addressServiceIns
}
