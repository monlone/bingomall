package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
	"time"
)

// userShop_service 接口
type UserShopService interface {
	/** 保存或修改 */
	SaveOrUpdate(userShop *model.UserShop) error

	Save(userShop *model.UserShop) error

	Update(userShop *model.UserShop) error

	GetByUserIDShopId(userId uint64, shopId string) *model.UserShop

	GetByShopId(id uint64) []*model.UserShop

	GetByUserID(id uint64) []*model.UserShop

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有  */
	GetAll() []*model.UserShop

	/** 分页查询 */
	GetPage(page int, pageSize int, userShop *model.UserShop) *helper.PageBean

	/**  按月份取 userShop 列表*/
	GetPageByMonth(page int, pageSize int, userShop *model.UserShop, month string) *helper.PageBean
}

var userShopServiceIns = &userShopService{}

// 获取 userShopService 实例
func UserShopServiceInstance(repo repositories.UserShopRepository) UserShopService {
	userShopServiceIns.repo = repo
	return userShopServiceIns
}

// 结构体
type userShopService struct {
	/** 存储对象 */
	repo repositories.UserShopRepository
}

func (us *userShopService) GetByUserIDShopId(userId uint64, shopId string) *model.UserShop {
	userShop := us.repo.FindSingle("user_id = ? AND shop_id", userId, shopId)
	if userShop != nil {
		return userShop.(*model.UserShop)
	}
	return nil
}

func (us *userShopService) SaveOrUpdate(userShop *model.UserShop) error {
	if userShop == nil || userShop.UserID == 0 || userShop.ShopId == 0 {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return us.repo.Insert(userShop)
}

func (us *userShopService) Save(userShop *model.UserShop) error {
	if userShop == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return us.repo.Insert(userShop)
}

func (us *userShopService) Update(userShop *model.UserShop) error {
	if userShop == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return us.repo.Update(userShop)
}

func (us *userShopService) GetAll() []*model.UserShop {
	userShops := us.repo.FindMore("1=1").([]*model.UserShop)
	return userShops
}

func (us *userShopService) GetByUserID(userId uint64) []*model.UserShop {
	if userId == 0 {
		return nil
	}
	userShop := us.repo.FindMore("user_id = ?", userId).([]*model.UserShop)
	return userShop
}

func (us *userShopService) GetByShopId(shopId uint64) []*model.UserShop {
	if shopId == 0 {
		return nil
	}
	userShop := us.repo.FindMore("shop_id = ?", shopId).([]*model.UserShop)
	return userShop
}

func (us *userShopService) DeleteByID(id uint64) error {
	userShop := us.repo.FindOne(id).(*model.UserShop)
	if userShop == nil || userShop.UserID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := us.repo.Delete(userShop)
	return err
}

func (us *userShopService) GetPage(page int, pageSize int, userShop *model.UserShop) *helper.PageBean {
	andCons := make(map[string]interface{})

	if userShop != nil && userShop.ShopId != 0 {
		andCons["shop_id = ?"] = userShop.ShopId
	}
	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}

func (us *userShopService) GetPageByMonth(page int, pageSize int, userShop *model.UserShop, month string) *helper.PageBean {
	andCons := make(map[string]interface{})

	if userShop != nil && userShop.ShopId != 0 {
		andCons["shop_id = ?"] = userShop.ShopId
	}
	if month != "" {
		monthBegin, _ := time.Parse("2006-01-02", month+"-01")
		andCons["userShop.updated_at >= ?"] = monthBegin
		monthEnd := monthBegin.AddDate(0, 1, 0)
		//timeUnix := monthEnd.Unix() + 86400 - 1
		//end := time.Unix(timeUnix, 0)
		//fmt.Println("end:", end)
		andCons["userShop.updated_at < ?"] = monthEnd
	}

	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}
