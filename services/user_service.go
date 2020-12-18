package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// user_service 接口
type UserService interface {
	/** 保存或修改 */
	SaveOrUpdate(user *model.User) error

	SaveByApp(user *model.User) error

	UpdateByApp(user *model.User) error

	/** 根据 id 查询 */
	GetByID(id uint64) *model.User

	/** 根据 userId 查询 */
	GetByUserID(userId uint64) *model.User

	/** 根据用户名查询 */
	GetByUsername(username string) *model.User

	/* 根据UnionId查询*/
	GetByUserUnionId(UnionId string) *model.User

	GetByUserOpenId(OpenId string) *model.User

	// 根据电话号码查询
	GetByPhone(phone string) *model.User

	/** 根据 id 删除 */
	DeleteByID(id uint64) error

	/** 查询所有  */
	GetAll() []*model.User

	/** 分页查询 */
	GetPage(page int, pageSize int, user *model.User) *helper.PageBean
}

var userServiceIns = &userService{}

// 获取 userService 实例
func UserServiceInstance(repo repositories.UserRepository) UserService {
	userServiceIns.repo = repo
	return userServiceIns
}

// 结构体
type userService struct {
	/** 存储对象 */
	repo repositories.UserRepository
}

func (us *userService) GetByUsername(username string) *model.User {
	user := us.repo.FindSingle("username = ?", username)
	if user != nil {
		return user.(*model.User)
	}
	return nil
}

func (us *userService) GetByUserUnionId(UnionId string) *model.User {
	user := us.repo.FindSingle("unionid = ?", UnionId)
	if user != nil {
		return user.(*model.User)
	}
	return nil
}

func (us *userService) GetByUserOpenId(OpenId string) *model.User {
	user := us.repo.FindSingle("openid = ?", OpenId)
	if user != nil {
		return user.(*model.User)
	}
	return nil
}

func (us *userService) GetByPhone(phone string) *model.User {
	user := us.repo.FindSingle("phone = ?", phone)
	if user != nil {
		return user.(*model.User)
	}
	return nil
}

func (us *userService) SaveOrUpdate(user *model.User) error {
	if user == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 校验用户名是否重复
	userByName := us.GetByUsername(user.Username)

	// 校验手机号码是否重复
	userByPhone := us.GetByPhone(user.Phone)
	if user.ID == 0 {
		// 添加
		if userByPhone != nil && userByPhone.ID != 0 {
			return errors.New(helper.StatusText(helper.ExistSamePhoneErr))
		}
		user.Password = helper.SHA256(user.Password)
		return us.repo.Insert(user)
	} else {
		// 修改
		persist := us.GetByUserID(user.ID)
		if persist == nil || persist.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		if userByName != nil && userByName.ID != user.ID {
			return errors.New(helper.StatusText(helper.ExistSameNameErr))
		}

		if userByPhone != nil && userByPhone.ID != user.ID {
			return errors.New(helper.StatusText(helper.ExistSamePhoneErr))
		}
		user.Password = persist.Password
		user.ID = persist.ID
		return us.repo.Update(user)
	}
}

func (us *userService) SaveByApp(user *model.User) error {
	if user == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	return us.repo.Insert(user)
}

func (us *userService) UpdateByApp(user *model.User) error {
	if user == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	persist := us.GetByUserID(user.ID)
	if persist == nil || persist.ID == 0 {
		return errors.New(helper.StatusText(helper.UpdateObjIsNil))
	}

	user.ID = persist.ID
	return us.repo.Update(user)
}

func (us *userService) GetAll() []*model.User {
	users := us.repo.FindMore("1=1").([]*model.User)
	return users
}

func (us *userService) GetByID(id uint64) *model.User {
	if id == 0 {
		return nil
	}
	user := us.repo.FindOne(id).(*model.User)
	return user
}

func (us *userService) GetByUserID(userId uint64) *model.User {
	if userId == 0 {
		return nil
	}
	user := us.repo.FindSingle("id = ?", userId)
	if user == nil {
		return nil
	}

	return user.(*model.User)
}

func (us *userService) DeleteByID(id uint64) error {
	user := us.repo.FindOne(id).(*model.User)
	if user == nil || user.ID == 0 {
		return errors.New(helper.StatusText(helper.DeleteObjIsNil))
	}
	err := us.repo.Delete(user)
	return err
}

func (us *userService) GetPage(page int, pageSize int, user *model.User) *helper.PageBean {
	andCons := make(map[string]interface{})
	if user != nil && user.Username != "" {
		andCons["username LIKE ?"] = user.Username + "%"
	}
	if user != nil && user.Phone != "" {
		andCons["phone LIKE ?"] = user.Phone + "%"
	}
	pageBean := us.repo.FindPage(page, pageSize, andCons, nil)
	return pageBean
}
