package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// Category service 接口
type CategoryService interface {
	// 保存或修改
	SaveOrUpdate(category *model.Category) error

	// 根据id查询
	GetByID(id uint64) *model.Category

	// 根据userId查询
	GetByUserID(userId uint64) []*model.Category

	// 通过userId获取该用户的商品分类
	GetCategoryByUserID(userId uint64) []*model.Category

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll(shopId uint64) []*model.Category

	// 分页查询
	GetPage(page int, pageSize int, user *model.Category) *helper.PageBean
}

// category service 结构体
type categoryService struct {
	/** 存储对象 */
	repo repositories.CategoryRepository
}

func (cs *categoryService) SaveOrUpdate(category *model.Category) error {
	if category == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if category.ID == 0 {
		// 添加
		return cs.repo.Insert(category)
	} else {
		// 修改
		persist := cs.GetByID(category.ID)
		if persist == nil || category.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		category.ID = persist.ID
		return cs.repo.Update(category)
	}
}

func (cs *categoryService) GetByID(id uint64) *model.Category {
	if id == 0 {
		return nil
	}
	category := cs.repo.FindOne(id).(*model.Category)
	return category
}

func (cs *categoryService) GetByUserID(userId uint64) []*model.Category {
	category := cs.repo.FindSingle("user_id = ?", userId)
	if category == nil {
		return nil
	}

	return category.([]*model.Category)
}

func (cs *categoryService) GetCategoryByUserID(userId uint64) []*model.Category {
	category := cs.repo.FindSingle("user_id = ?", userId)
	if category == nil {
		return nil
	}
	data := category.([]*model.Category)
	return data
}

func (*categoryService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (cs *categoryService) GetAll(shopId uint64) []*model.Category {
	var category []*model.Category

	if shopId != 0 {
		category = cs.repo.FindMore("shop_id=?", shopId).([]*model.Category)
	} else {
		category = cs.repo.FindMore("1=1").([]*model.Category)
	}

	return category
}

func (*categoryService) GetPage(page int, pageSize int, user *model.Category) *helper.PageBean {
	panic("implement me")
}

var categoryServiceIns = &categoryService{}

func CategoryServiceInstance(repo repositories.CategoryRepository) CategoryService {
	categoryServiceIns.repo = repo
	return categoryServiceIns
}
