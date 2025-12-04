package service

import (
	"context"
	"merchant_api/internal/dao"
	"merchant_api/internal/model"
	"merchant_api/pkg/database"
)

type StoreCategoryService struct {
	ctx context.Context
}

func NewStoreCategoryService(ctx context.Context) *StoreCategoryService {
	dao.SetDefault(database.GetDB())
	return &StoreCategoryService{ctx: ctx}
}

// Create 创建分类
func (s *StoreCategoryService) Create(req *model.MerStoreCategory) error {
	return dao.MerStoreCategory.WithContext(s.ctx).Create(req)
}

// Update 更新分类
func (s *StoreCategoryService) Update(id int32, merId int32, req *model.MerStoreCategory) error {
	c := dao.MerStoreCategory

	// 确保只能更新自己商户的分类
	_, err := c.WithContext(s.ctx).
		Where(c.StoreCategoryID.Eq(id), c.MerID.Eq(merId)).
		Updates(req)
	return err
}

// Delete 删除分类
func (s *StoreCategoryService) Delete(id int32, merId int32) error {
	c := dao.MerStoreCategory

	// 确保只能删除自己商户的分类
	_, err := c.WithContext(s.ctx).
		Where(c.StoreCategoryID.Eq(id), c.MerID.Eq(merId)).
		Delete()
	return err
}

// GetList 获取分类列表
func (s *StoreCategoryService) GetList(merId int32, page, pageSize int) ([]*model.MerStoreCategory, int64, error) {
	c := dao.MerStoreCategory

	query := c.WithContext(s.ctx).Where(c.MerID.Eq(merId))

	total, err := query.Count()
	if err != nil {
		return nil, 0, err
	}

	list, err := query.Order(c.Sort.Desc(), c.StoreCategoryID.Desc()).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find()

	return list, total, err
}

// Get 获取单个分类
func (s *StoreCategoryService) Get(id int32, merId int32) (*model.MerStoreCategory, error) {
	c := dao.MerStoreCategory

	return c.WithContext(s.ctx).
		Where(c.StoreCategoryID.Eq(id), c.MerID.Eq(merId)).
		First()
}
