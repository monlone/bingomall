package service

import (
	"errors"
	"bingomall/helpers"
	"bingomall/models"
	"bingomall/repositories"
)

// Wallet service 接口
type WalletService interface {
	// 保存或修改
	SaveOrUpdate(wallet *model.Wallet) error

	//初始化个人钱包
	InitMyWallet(userId uint64) error

	// 根据id查询
	GetByID(id uint64) *model.Wallet

	// 根据userId查询
	GetByUserID(userId uint64) *model.Wallet

	// 通过userId获取用户积分
	GetScoreByUserID(userId uint64) uint64

	// 扣除用户积分
	ReduceScoreByUserID(userId uint64, score uint64) error

	AddScoreByUserID(userId uint64, score uint64) error

	ReduceMoneyByUserID(userId uint64, money uint64) error

	AddMoneyByUserID(userId uint64, money uint64) error

	AddGrowthByUserID(userId uint64, money uint64) error

	ReduceGrowthByUserID(userId uint64, money uint64) error

	// 根据 id 删除
	DeleteByID(id uint64) error

	// 查询所有
	GetAll() []*model.Wallet

	// 分页查询
	GetPage(page int, pageSize int, user *model.Wallet) *helper.PageBean
}

// score service 结构体
type walletService struct {
	/** 存储对象 */
	repo repositories.WalletRepository
}

func (ws *walletService) SaveOrUpdate(wallet *model.Wallet) error {
	if wallet == nil {
		return errors.New(helper.StatusText(helper.SaveObjIsNil))
	}
	// 判断 新增还是更新
	if wallet.ID == 0 {
		// 添加
		return ws.repo.Insert(wallet)
	} else {
		// 修改
		persist := ws.GetByUserID(wallet.UserID)
		if persist == nil || wallet.ID == 0 {
			return errors.New(helper.StatusText(helper.UpdateObjIsNil))
		}
		wallet.ID = persist.ID
		return ws.repo.Update(wallet)
	}
}

func (ws *walletService) InitMyWallet(userId uint64) error {
	var wallet model.Wallet
	wallet.UserID = userId
	return ws.repo.Insert(wallet)
}

func (ws *walletService) GetByID(id uint64) *model.Wallet {
	if id == 0 {
		return nil
	}
	wallet := ws.repo.FindOne(id).(*model.Wallet)
	return wallet
}

func (ws *walletService) GetByUserID(userId uint64) *model.Wallet {
	wallet := ws.repo.FindSingle("user_id = ?", userId)
	if wallet == nil {
		return nil
	}

	return wallet.(*model.Wallet)
}

func (ws *walletService) GetScoreByUserID(userId uint64) uint64 {
	wallet := ws.repo.FindSingle("user_id = ?", userId)
	if wallet == nil {
		return 0
	}
	data := wallet.(*model.Wallet)
	return data.Score
}

func (ws *walletService) ReduceScoreByUserID(userId uint64, score uint64) error {
	walletData := ws.repo.FindSingle("user_id = ?", userId)
	if walletData == nil {
		return errors.New("此用户不存在钱包")
	}
	data := walletData.(*model.Wallet)
	if data.Score < score {
		return errors.New("用户积分不足")
	}
	data.Score = data.Score - score
	err := ws.repo.Update(data)

	return err
}

func (ws *walletService) AddScoreByUserID(userId uint64, score uint64) error {
	walletData := ws.repo.FindSingle("user_id = ?", userId)
	if walletData == nil {
		return errors.New("此用户不存在钱包")
	}
	data := walletData.(*model.Wallet)
	data.Score = data.Score + score
	err := ws.repo.Update(data)

	return err
}

func (ws *walletService) ReduceMoneyByUserID(userId uint64, money uint64) error {
	walletData := ws.repo.FindSingle("user_id = ?", userId)
	if walletData == nil {
		return errors.New("此用户不存在钱包")
	}
	data := walletData.(*model.Wallet)
	if data.Balance < money {
		return errors.New("用户金钱不足")
	}
	data.Balance = data.Balance - money

	err := ws.repo.Update(data)

	return err
}

func (ws *walletService) AddMoneyByUserID(userId uint64, money uint64) error {
	walletData := ws.repo.FindSingle("user_id = ?", userId)
	if walletData == nil {
		return errors.New("此用户不存在钱包")
	}

	data := walletData.(*model.Wallet)
	data.Balance = data.Balance + money
	data.Money = data.Money + money
	err := ws.repo.Update(data)

	return err
}

func (ws *walletService) ReduceGrowthByUserID(userId uint64, growth uint64) error {
	walletData := ws.repo.FindSingle("user_id = ?", userId)
	if walletData == nil {
		return errors.New("此用户不存在钱包")
	}
	data := walletData.(*model.Wallet)
	if data.Growth < growth {
		return errors.New("用户成长值不足")
	}
	data.Growth = data.Growth - growth
	err := ws.repo.Update(data)

	return err
}

func (ws *walletService) AddGrowthByUserID(userId uint64, growth uint64) error {
	walletData := ws.repo.FindSingle("user_id = ?", userId)
	if walletData == nil {
		return errors.New("此用户不存在钱包")
	}
	data := walletData.(*model.Wallet)
	data.Growth = data.Growth + growth
	err := ws.repo.Update(data)

	return err
}

func (*walletService) DeleteByID(id uint64) error {
	panic("implement me")
}

func (*walletService) GetAll() []*model.Wallet {
	panic("implement me")
}

func (*walletService) GetPage(page int, pageSize int, user *model.Wallet) *helper.PageBean {
	panic("implement me")
}

var walletServiceIns = &walletService{}

func WalletServiceInstance(repo repositories.WalletRepository) WalletService {
	walletServiceIns.repo = repo
	return walletServiceIns
}
