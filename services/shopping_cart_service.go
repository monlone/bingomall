package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// shopping_cart_service 接口
type ShoppingCartService interface {
	/** 保存或修改 */
	SaveOrUpdate(shoppingCart *model.ShoppingCart) error

	SaveByApp(shoppingCart *model.ShoppingCart) error

	UpdateByApp(shoppingCart *model.ShoppingCart) error

	/** 根据 id 查询 */
	GetByID(id uint64) *model.ShoppingCart

	/** 根据 shoppingCartId 查询 */
	GetByShoppingCartId(shoppingCartId uint64) *model.ShoppingCart

	GetShoppingCartItem(userId uint64, productId uint64, skuId uint64) *model.ShoppingCart

	/* 根据 UserID 查询*/
	GetShoppingCartByUserID(userId uint64) []*model.ShoppingCart

	GetShoppingCartByUserIDCartIds(userId uint64, cartIds []uint64) []*model.ShoppingCart

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有  */
	GetAll() []*model.ShoppingCart
}

var shoppingCartServiceIns = &shoppingCartService{}

// 获取 shoppingCartService 实例
func ShoppingCartServiceInstance(repo repositories.ShoppingCartRepository) ShoppingCartService {
	shoppingCartServiceIns.repo = repo
	return shoppingCartServiceIns
}

// 结构体
type shoppingCartService struct {
	/** 存储对象 */
	repo repositories.ShoppingCartRepository
}

func (scs *shoppingCartService) GetShoppingCartByUserID(userId uint64) []*model.ShoppingCart {
	shoppingCart := scs.repo.FindMore("user_id = ?", userId)
	if shoppingCart != nil {
		return shoppingCart.([]*model.ShoppingCart)
	}
	return nil
}

func (scs *shoppingCartService) GetShoppingCartByUserIDCartIds(userId uint64, cartIds []uint64) []*model.ShoppingCart {
	shoppingCart := scs.repo.FindMore("user_id = ? AND id IN (?)", userId, cartIds)
	if shoppingCart != nil {
		return shoppingCart.([]*model.ShoppingCart)
	}
	return nil
}

func (scs *shoppingCartService) SaveOrUpdate(shoppingCart *model.ShoppingCart) error {
	if shoppingCart == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	persist := scs.GetShoppingCartItem(shoppingCart.UserID, shoppingCart.ProductId, shoppingCart.SkuId)
	if persist == nil {
		// 添加
		return scs.repo.Insert(shoppingCart)
	} else {
		// 修改
		shoppingCart.ID = persist.ID
		return scs.repo.Update(shoppingCart)
	}
}

func (scs *shoppingCartService) SaveByApp(shoppingCart *model.ShoppingCart) error {
	if shoppingCart == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return scs.repo.Insert(shoppingCart)
}

func (scs *shoppingCartService) UpdateByApp(shoppingCart *model.ShoppingCart) error {
	if shoppingCart == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	persist := scs.GetByShoppingCartId(shoppingCart.ID)
	if persist == nil || persist.ID == 0 {
		return errors.New(helper.StatusText(helper.UpdateObjIsNil))
	}

	shoppingCart.ID = persist.ID
	return scs.repo.Update(shoppingCart)
}

func (scs *shoppingCartService) GetAll() []*model.ShoppingCart {
	shoppingCarts := scs.repo.FindMore("1=1").([]*model.ShoppingCart)
	return shoppingCarts
}

func (scs *shoppingCartService) GetByID(id uint64) *model.ShoppingCart {
	if id == 0 {
		return nil
	}
	shoppingCart := scs.repo.FindOne(id).(*model.ShoppingCart)
	return shoppingCart
}

func (scs *shoppingCartService) GetByShoppingCartId(shoppingCartId uint64) *model.ShoppingCart {
	if shoppingCartId == 0 {
		return nil
	}
	shoppingCart := scs.repo.FindSingle("id = ?", shoppingCartId)
	if shoppingCart == nil {
		return nil
	}

	return shoppingCart.(*model.ShoppingCart)
}

func (scs *shoppingCartService) GetShoppingCartItem(userId uint64, productId uint64, skuId uint64) *model.ShoppingCart {
	shoppingCart := scs.repo.FindSingle("user_id = ? AND product_id = ? AND sku_id = ?", userId, productId, skuId)
	if shoppingCart == nil {
		return nil
	}

	return shoppingCart.(*model.ShoppingCart)
}

func (scs *shoppingCartService) DeleteByID(id uint64) error {
	shoppingCart := scs.repo.FindOne(id).(*model.ShoppingCart)
	if shoppingCart == nil || shoppingCart.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := scs.repo.Delete(shoppingCart)
	return err
}
