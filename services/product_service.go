package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// product_service 接口
type ProductService interface {
	/** 保存或修改 */
	SaveOrUpdate(product *model.Product) error

	Save(product *model.Product) error

	Update(product *model.Product) error

	/** 根据 id 查询 */
	GetByID(id uint64) *model.Product

	/** 根据 product_id 查询 */
	GetByProductId(productId uint64) *model.Product

	GetListByProductIds(productId []uint64) []*model.Product

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有  */
	GetAll() []*model.Product

	/** 分页查询 */
	GetPage(page int, pageSize int, product *model.Product) *helper.PageBean

	GetListByShopId(page int, pageSize int, shopId string) *helper.PageBean
}

var productServiceIns = &productService{}

// 获取 productService 实例
func ProductServiceInstance(repo repositories.ProductRepository) ProductService {
	productServiceIns.repo = repo
	return productServiceIns
}

// 结构体
type productService struct {
	/** 存储对象 */
	repo repositories.ProductRepository
}

func (us *productService) GetByProductOpenId(productId uint64) *model.Product {
	product := us.repo.FindSingle("id = ?", productId)
	if product != nil {
		return product.(*model.Product)
	}
	return nil
}

func (us *productService) SaveOrUpdate(product *model.Product) error {
	if product == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	if product.ID == 0 {
		// 添加
		return us.repo.Insert(product)
	} else {
		// 修改
		persist := us.GetByProductId(product.ID)
		if persist == nil {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		product.ID = persist.ID
		return us.repo.Update(product)
	}
}

func (us *productService) Save(product *model.Product) error {
	if product == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return us.repo.Insert(product)
}

func (us *productService) Update(product *model.Product) error {
	if product == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	persist := us.GetByProductId(product.ID)
	if persist == nil {
		return errors.New(helper.StatusText(helper.UpdateObjIsNil))
	}
	product.ID = persist.ID
	return us.repo.Update(product)
}

func (us *productService) GetAll() []*model.Product {
	products := us.repo.FindMore("1=1").([]*model.Product)
	return products
}

func (us *productService) GetByID(id uint64) *model.Product {
	if id == 0 {
		return nil
	}
	product := us.repo.FindOne(id)
	if product == nil {
		return nil
	}
	return product.(*model.Product)
}

func (us *productService) GetByProductId(productId uint64) *model.Product {
	if productId == 0 {
		return nil
	}
	product := us.repo.FindSingle("id = ?", productId)
	if product == nil {
		return nil
	}
	return product.(*model.Product)
}

func (us *productService) GetListByProductIds(productIds []uint64) []*model.Product {
	product := us.repo.FindMore("id in (?)", productIds)
	if product == nil {
		return nil
	}
	return product.([]*model.Product)
}

func (us *productService) GetListByProductIdsWithOption(productIds []uint64) []*model.Product {
	product := us.repo.FindMore("id in (?)", productIds)
	if product == nil {
		return nil
	}
	return product.([]*model.Product)
}

func (us *productService) DeleteByID(productId uint64) error {
	product := us.repo.FindSingle("id = ?", productId).(*model.Product)
	if product == nil || product.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := us.repo.Delete(product)
	return err
}

func (us *productService) GetPage(page int, pageSize int, product *model.Product) *helper.PageBean {
	andCons := make(map[string]interface{})

	if product != nil && product.ShopId != 0 {
		andCons["shop_id = ?"] = product.ShopId
	}
	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}

func (us *productService) GetListByShopId(page int, pageSize int, shopId string) *helper.PageBean {
	andCons := make(map[string]interface{})

	if shopId != "" {
		andCons["shop_id = ?"] = shopId
	}

	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)

	return pageBean
}
