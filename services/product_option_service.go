package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// ProductOption service 接口
type ProductOptionService interface {
	// 保存或修改
	SaveOrUpdate(productOption *model.ProductOption) error

	// 根据id查询
	GetByID(id uint64) *model.ProductOption

	// 根据userId查询
	GetByUserID(userId uint64) []*model.ProductOption

	// 通过userId获取该用户的商品分类
	GetProductOptionByUserID(userId uint64) []*model.ProductOption

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll(productId uint64) []*model.ProductOption

	// 查询所有,并按类型归类
	GetAllGroupByType(productId uint64) map[uint8][]*model.ProductOption

	// 分页查询
	GetPage(page int, pageSize int, user *model.ProductOption) *helper.PageBean
}

// productOption service 结构体
type productOptionService struct {
	/** 存储对象 */
	repo repositories.ProductOptionRepository
}

func (pos *productOptionService) SaveOrUpdate(productOption *model.ProductOption) error {
	if productOption == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if productOption.ID == 0 {
		// 添加
		return pos.repo.Insert(productOption)
	} else {
		// 修改
		persist := pos.GetByID(productOption.ID)
		if persist == nil || productOption.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		productOption.ID = persist.ID
		return pos.repo.Update(productOption)
	}
}

func (pos *productOptionService) GetByID(id uint64) *model.ProductOption {
	if id == 0 {
		return nil
	}
	productOption := pos.repo.FindOne(id).(*model.ProductOption)
	return productOption
}

func (pos *productOptionService) GetByUserID(userId uint64) []*model.ProductOption {
	productOption := pos.repo.FindSingle("user_id = ?", userId)
	if productOption == nil {
		return nil
	}

	return productOption.([]*model.ProductOption)
}

func (pos *productOptionService) GetProductOptionByUserID(userId uint64) []*model.ProductOption {
	productOption := pos.repo.FindSingle("user_id = ?", userId)
	if productOption == nil {
		return nil
	}
	data := productOption.([]*model.ProductOption)
	return data
}

func (*productOptionService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (pos *productOptionService) GetAll(productId uint64) []*model.ProductOption {
	var productOption []*model.ProductOption

	if productId != 0 {
		productOption = pos.repo.FindMore("product_id=?", productId).([]*model.ProductOption)
	} else {
		productOption = pos.repo.FindMore("1=1").([]*model.ProductOption)
	}

	return productOption
}

func (pos *productOptionService) GetAllGroupByType(productId uint64) map[uint8][]*model.ProductOption {
	var productOption []*model.ProductOption

	if productId != 0 {
		productOption = pos.repo.FindMore("product_id=?", productId).([]*model.ProductOption)
	} else {
		productOption = pos.repo.FindMore("1=1").([]*model.ProductOption)
	}

	optionList := make(map[uint8][]*model.ProductOption)

	for _, v := range productOption {
		if _, ok := optionList[v.Type]; ok {
			temp := optionList[v.Type]
			temp = append(temp, v)
			optionList[v.Type] = temp
		} else {
			var temp []*model.ProductOption
			temp = append(temp, v)
			optionList[v.Type] = temp
		}
	}

	return optionList
}

func (*productOptionService) GetPage(page int, pageSize int, user *model.ProductOption) *helper.PageBean {
	panic("implement me")
}

var productOptionServiceIns = &productOptionService{}

func ProductOptionServiceInstance(repo repositories.ProductOptionRepository) ProductOptionService {
	productOptionServiceIns.repo = repo
	return productOptionServiceIns
}
