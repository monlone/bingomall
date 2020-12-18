package repositories

import (
	"gorm.io/gorm"
)

// base repository 接口
type BaseRepository interface {
	/** 基础 repository 提供最基础的增删改查 */
}

var baseRepoIns = &baseRepository{}

// 实例化 存储对象
func BaseRepositoryInstance(db *gorm.DB) BaseRepository {
	baseRepoIns.db = db
	return baseRepoIns
}

type baseRepository struct {
	db *gorm.DB
}
