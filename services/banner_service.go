package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// banner_service 接口
type BannerService interface {
	/** 保存或修改 */
	SaveOrUpdate(banner *model.Banner) error

	Save(banner *model.Banner) error

	Update(banner *model.Banner) error

	/** 根据 banner_id 查询 */
	GetByBannerID(bannerId uint64) *model.Banner

	/** */
	GetAllByShopId(shopId uint64) []*model.Banner

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	DeleteByBannerID(bannerId uint64) error

	/** 查询所有  */
	GetAll() []*model.Banner

	/** 分页查询 */
	GetPage(page int, pageSize int, banner *model.Banner) *helper.PageBean
}

var bannerServiceIns = &bannerService{}

// 获取 bannerService 实例
func BannerServiceInstance(repo repositories.BannerRepository) BannerService {
	bannerServiceIns.repo = repo
	return bannerServiceIns
}

// 结构体
type bannerService struct {
	/** 存储对象 */
	repo repositories.BannerRepository
}

func (us *bannerService) GetByBannerOpenId(bannerId uint64) *model.Banner {
	banner := us.repo.FindSingle("banner_id = ?", bannerId)
	if banner != nil {
		return banner.(*model.Banner)
	}
	return nil
}

func (us *bannerService) SaveOrUpdate(banner *model.Banner) error {
	if banner == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	if banner.ID == 0 {
		// 添加
		return us.repo.Insert(banner)
	} else {
		// 修改
		persist := us.GetByBannerID(banner.ID)
		if persist == nil {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}

		banner.ID = persist.ID
		return us.repo.Update(banner)
	}
}

func (us *bannerService) Save(banner *model.Banner) error {
	if banner == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return us.repo.Insert(banner)
}

func (us *bannerService) Update(banner *model.Banner) error {
	if banner == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}

	persist := us.GetByBannerID(banner.ID)
	if persist == nil || persist.ID == 0 {
		return errors.New(helper.StatusText(helper.UpdateObjIsNil))
	}

	banner.ID = persist.ID
	return us.repo.Update(banner)
}

func (us *bannerService) GetAll() []*model.Banner {
	banners := us.repo.FindMore("1=1").([]*model.Banner)
	return banners
}

func (us *bannerService) GetAllByShopId(shopId uint64) []*model.Banner {
	if shopId == 0 {
		return nil
	}
	banner := us.repo.FindMore("id = ?", shopId).([]*model.Banner)
	return banner
}

func (us *bannerService) GetByBannerID(bannerId uint64) *model.Banner {
	if bannerId == 0 {
		return nil
	}
	banner := us.repo.FindSingle("id = ?", bannerId).(*model.Banner)
	return banner
}

func (us *bannerService) DeleteByID(id uint64) error {
	banner := us.repo.FindOne(id).(*model.Banner)
	if banner == nil || banner.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := us.repo.Delete(banner)
	return err
}

func (us *bannerService) DeleteByBannerID(bannerId uint64) error {
	if bannerId == 0 {
		return nil
	}
	banner := us.repo.FindSingle("banner_id = ?", bannerId).(*model.Banner)
	if banner == nil {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := us.repo.Delete(banner)
	return err
}

func (us *bannerService) GetPage(page int, pageSize int, banner *model.Banner) *helper.PageBean {
	andCons := make(map[string]interface{})

	if banner != nil && banner.ShopId != "" {
		andCons["shop_id = ?"] = banner.ShopId
	}
	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}
