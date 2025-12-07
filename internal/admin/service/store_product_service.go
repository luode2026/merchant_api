package service

import (
	"context"
	"errors"
	"fmt"
	"merchant_api/internal/dao"
	"merchant_api/internal/model"
	"merchant_api/pkg/database"
	"time"

	"gorm.io/gorm"
)

type StoreProductService struct {
	ctx context.Context
}

func NewStoreProductService(ctx context.Context) *StoreProductService {
	dao.SetDefault(database.GetDB())
	return &StoreProductService{ctx: ctx}
}

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	StoreName     string                `json:"store_name" binding:"required"`
	StoreInfo     *string               `json:"store_info"`
	Keyword       string                `json:"keyword"`
	CateID        int32                 `json:"cate_id" binding:"required"`
	UnitName      string                `json:"unit_name" binding:"required"`
	Sort          int32                 `json:"sort"`
	Price         *float64              `json:"price"`
	Cost          *float64              `json:"cost"`
	OtPrice       *float64              `json:"ot_price"`
	IsGood        bool                  `json:"is_good"`
	ProductType   int32                 `json:"product_type"`
	Image         string                `json:"image" binding:"required"`
	SliderImage   string                `json:"slider_image"`
	RefundSwitch  *int32                `json:"refund_switch"`
	BarCodeNumber *string               `json:"bar_code_number"`
	Content       *string               `json:"content"`
	Skus          []CreateProductSkuReq `json:"skus" binding:"required,min=1"`
}

// CreateProductSkuReq SKU请求
type CreateProductSkuReq struct {
	ProductSkuID *int32   `json:"product_sku_id"` // SKU ID，更新时传入，创建时不传
	AttrName     *string  `json:"attr_name"`
	Price        *float64 `json:"price" binding:"required"`
	Cost         *float64 `json:"cost"`
	OtPrice      *float64 `json:"ot_price"`
	Image        *string  `json:"image"`
}

// ProductDetailResponse 商品详情响应
type ProductDetailResponse struct {
	*model.MerStoreProduct
	Category *model.MerStoreCategory       `json:"category"`
	Content  *model.MerStoreProductContent `json:"content"`
	Skus     []*model.MerStoreProductSku   `json:"skus"`
}

// Create 创建商品
func (s *StoreProductService) Create(req *CreateProductRequest, merID int32) (*ProductDetailResponse, error) {
	db := database.GetDB()

	// 验证分类是否存在且属于该商户
	category, err := dao.MerStoreCategory.WithContext(s.ctx).
		Where(dao.MerStoreCategory.StoreCategoryID.Eq(req.CateID)).
		Where(dao.MerStoreCategory.MerID.Eq(merID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("分类不存在或无权访问")
		}
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}

	var result *ProductDetailResponse

	// 使用事务创建商品及关联数据
	err = db.Transaction(func(tx *gorm.DB) error {
		dao.SetDefault(tx)

		// 创建商品主表
		now := time.Now()
		product := &model.MerStoreProduct{
			MerID:         merID,
			StoreName:     req.StoreName,
			StoreInfo:     "",
			Keyword:       req.Keyword,
			IsShow:        1,             // 默认上架
			SaleStatus:    boolPtr(true), // 默认销售中
			CateID:        req.CateID,
			UnitName:      req.UnitName,
			Sort:          req.Sort,
			Sales:         0,
			Price:         req.Price,
			Cost:          req.Cost,
			OtPrice:       req.OtPrice,
			IsGood:        req.IsGood,
			ProductType:   req.ProductType,
			Image:         req.Image,
			SliderImage:   req.SliderImage,
			RefundSwitch:  req.RefundSwitch,
			CreateAt:      now,
			BarCodeNumber: req.BarCodeNumber,
		}

		if req.StoreInfo != nil {
			product.StoreInfo = *req.StoreInfo
		}

		if err := dao.MerStoreProduct.WithContext(s.ctx).Create(product); err != nil {
			return fmt.Errorf("创建商品失败: %w", err)
		}

		// 创建商品详情
		contentStr := ""
		if req.Content != nil {
			contentStr = *req.Content
		}
		content := &model.MerStoreProductContent{
			ProductID: product.ProductID,
			Content:   contentStr,
		}
		if err := dao.MerStoreProductContent.WithContext(s.ctx).Create(content); err != nil {
			return fmt.Errorf("创建商品详情失败: %w", err)
		}

		// 创建商品SKU
		skus := make([]*model.MerStoreProductSku, 0, len(req.Skus))
		for _, skuReq := range req.Skus {
			sku := &model.MerStoreProductSku{
				ProductID: product.ProductID,
				AttrName:  skuReq.AttrName,
				Price:     skuReq.Price,
				Cost:      skuReq.Cost,
				OtPrice:   skuReq.OtPrice,
				Image:     skuReq.Image,
			}
			if err := dao.MerStoreProductSku.WithContext(s.ctx).Create(sku); err != nil {
				return fmt.Errorf("创建商品SKU失败: %w", err)
			}
			skus = append(skus, sku)
		}

		result = &ProductDetailResponse{
			MerStoreProduct: product,
			Category:        category,
			Content:         content,
			Skus:            skus,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update 更新商品
func (s *StoreProductService) Update(productID int32, req *CreateProductRequest, merID int32) error {
	db := database.GetDB()

	// 验证商品是否存在且属于该商户
	product, err := dao.MerStoreProduct.WithContext(s.ctx).
		Where(dao.MerStoreProduct.ProductID.Eq(productID)).
		Where(dao.MerStoreProduct.MerID.Eq(merID)).
		Where(dao.MerStoreProduct.DeleteAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("商品不存在或无权访问")
		}
		return fmt.Errorf("查询商品失败: %w", err)
	}

	// 验证分类是否存在且属于该商户
	_, err = dao.MerStoreCategory.WithContext(s.ctx).
		Where(dao.MerStoreCategory.StoreCategoryID.Eq(req.CateID)).
		Where(dao.MerStoreCategory.MerID.Eq(merID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("分类不存在或无权访问")
		}
		return fmt.Errorf("查询分类失败: %w", err)
	}

	// 使用事务更新商品及关联数据
	return db.Transaction(func(tx *gorm.DB) error {
		dao.SetDefault(tx)

		// 更新商品主表
		now := time.Now()
		updates := map[string]interface{}{
			"store_name":      req.StoreName,
			"keyword":         req.Keyword,
			"cate_id":         req.CateID,
			"unit_name":       req.UnitName,
			"sort":            req.Sort,
			"price":           req.Price,
			"cost":            req.Cost,
			"ot_price":        req.OtPrice,
			"is_good":         req.IsGood,
			"product_type":    req.ProductType,
			"image":           req.Image,
			"slider_image":    req.SliderImage,
			"refund_switch":   req.RefundSwitch,
			"bar_code_number": req.BarCodeNumber,
			"update_at":       now,
		}
		if req.StoreInfo != nil {
			updates["store_info"] = *req.StoreInfo
		}

		_, err := dao.MerStoreProduct.WithContext(s.ctx).
			Where(dao.MerStoreProduct.ProductID.Eq(productID)).
			Updates(updates)
		if err != nil {
			return fmt.Errorf("更新商品失败: %w", err)
		}

		// 更新商品详情
		if req.Content != nil {
			_, err = dao.MerStoreProductContent.WithContext(s.ctx).
				Where(dao.MerStoreProductContent.ProductID.Eq(productID)).
				Updates(map[string]interface{}{
					"content": *req.Content,
				})
			if err != nil {
				return fmt.Errorf("更新商品详情失败: %w", err)
			}
		}

		// 收集请求中的SKU ID
		reqSkuIDs := make(map[int32]bool)
		for _, skuReq := range req.Skus {
			if skuReq.ProductSkuID != nil {
				reqSkuIDs[*skuReq.ProductSkuID] = true
			}
		}

		// 查询当前商品的所有SKU
		existingSkus, err := dao.MerStoreProductSku.WithContext(s.ctx).
			Where(dao.MerStoreProductSku.ProductID.Eq(productID)).
			Find()
		if err != nil {
			return fmt.Errorf("查询现有SKU失败: %w", err)
		}

		// 删除不在请求中的SKU
		for _, existingSku := range existingSkus {
			if !reqSkuIDs[existingSku.ProductSkuID] {
				_, err = dao.MerStoreProductSku.WithContext(s.ctx).
					Where(dao.MerStoreProductSku.ProductSkuID.Eq(existingSku.ProductSkuID)).
					Delete()
				if err != nil {
					return fmt.Errorf("删除旧SKU失败: %w", err)
				}
			}
		}

		// 更新或创建SKU
		for _, skuReq := range req.Skus {
			if skuReq.ProductSkuID != nil {
				// 更新现有SKU
				_, err = dao.MerStoreProductSku.WithContext(s.ctx).
					Where(dao.MerStoreProductSku.ProductSkuID.Eq(*skuReq.ProductSkuID)).
					Where(dao.MerStoreProductSku.ProductID.Eq(productID)).
					Updates(map[string]interface{}{
						"attr_name": skuReq.AttrName,
						"price":     skuReq.Price,
						"cost":      skuReq.Cost,
						"ot_price":  skuReq.OtPrice,
						"image":     skuReq.Image,
					})
				if err != nil {
					return fmt.Errorf("更新SKU失败: %w", err)
				}
			} else {
				// 创建新SKU
				sku := &model.MerStoreProductSku{
					ProductID: product.ProductID,
					AttrName:  skuReq.AttrName,
					Price:     skuReq.Price,
					Cost:      skuReq.Cost,
					OtPrice:   skuReq.OtPrice,
					Image:     skuReq.Image,
				}
				if err := dao.MerStoreProductSku.WithContext(s.ctx).Create(sku); err != nil {
					return fmt.Errorf("创建新SKU失败: %w", err)
				}
			}
		}

		return nil
	})
}

// Delete 删除商品（软删除）
func (s *StoreProductService) Delete(productID int32, merID int32) error {
	// 验证商品是否存在且属于该商户
	_, err := dao.MerStoreProduct.WithContext(s.ctx).
		Where(dao.MerStoreProduct.ProductID.Eq(productID)).
		Where(dao.MerStoreProduct.MerID.Eq(merID)).
		Where(dao.MerStoreProduct.DeleteAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("商品不存在或无权访问")
		}
		return fmt.Errorf("查询商品失败: %w", err)
	}

	// 软删除
	now := int32(time.Now().Unix())
	_, err = dao.MerStoreProduct.WithContext(s.ctx).
		Where(dao.MerStoreProduct.ProductID.Eq(productID)).
		Updates(map[string]interface{}{
			"delete_at": now,
		})
	return err
}

// Get 获取商品详情
func (s *StoreProductService) Get(productID int32, merID int32) (*ProductDetailResponse, error) {
	// 查询商品
	product, err := dao.MerStoreProduct.WithContext(s.ctx).
		Where(dao.MerStoreProduct.ProductID.Eq(productID)).
		Where(dao.MerStoreProduct.MerID.Eq(merID)).
		Where(dao.MerStoreProduct.DeleteAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("商品不存在或无权访问")
		}
		return nil, fmt.Errorf("查询商品失败: %w", err)
	}

	// 查询分类
	category, err := dao.MerStoreCategory.WithContext(s.ctx).
		Where(dao.MerStoreCategory.StoreCategoryID.Eq(product.CateID)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}

	// 查询商品详情
	content, err := dao.MerStoreProductContent.WithContext(s.ctx).
		Where(dao.MerStoreProductContent.ProductID.Eq(productID)).
		First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("查询商品详情失败: %w", err)
	}

	// 查询SKU列表
	skus, err := dao.MerStoreProductSku.WithContext(s.ctx).
		Where(dao.MerStoreProductSku.ProductID.Eq(productID)).
		Find()
	if err != nil {
		return nil, fmt.Errorf("查询SKU失败: %w", err)
	}

	return &ProductDetailResponse{
		MerStoreProduct: product,
		Category:        category,
		Content:         content,
		Skus:            skus,
	}, nil
}

// ListRequest 列表请求
type ListRequest struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=20"`
	CateID     *int32 `form:"cate_id"`
	IsShow     *int32 `form:"is_show"`
	SaleStatus *bool  `form:"sale_status"`
	Keyword    string `form:"keyword"`
}

// GetList 获取商品列表
func (s *StoreProductService) GetList(merID int32, req *ListRequest) ([]*ProductDetailResponse, int64, error) {
	p := dao.MerStoreProduct

	query := p.WithContext(s.ctx).
		Where(p.MerID.Eq(merID)).
		Where(p.DeleteAt.IsNull())

	// 分类筛选
	if req.CateID != nil {
		query = query.Where(p.CateID.Eq(*req.CateID))
	}

	// 上架状态筛选
	if req.IsShow != nil {
		query = query.Where(p.IsShow.Eq(*req.IsShow))
	}

	// 售完状态筛选
	if req.SaleStatus != nil {
		query = query.Where(p.SaleStatus.Is(*req.SaleStatus))
	}

	// 关键字搜索
	if req.Keyword != "" {
		query = query.Where(p.StoreName.Like("%" + req.Keyword + "%"))
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, 0, fmt.Errorf("查询商品总数失败: %w", err)
	}

	// 获取列表
	products, err := query.
		Order(p.Sort.Desc(), p.ProductID.Desc()).
		Limit(req.PageSize).
		Offset((req.Page - 1) * req.PageSize).
		Find()
	if err != nil {
		return nil, 0, fmt.Errorf("查询商品列表失败: %w", err)
	}

	// 组装详情数据
	result := make([]*ProductDetailResponse, 0, len(products))
	for _, product := range products {
		// 查询分类
		category, _ := dao.MerStoreCategory.WithContext(s.ctx).
			Where(dao.MerStoreCategory.StoreCategoryID.Eq(product.CateID)).
			First()

		// 查询商品详情
		content, _ := dao.MerStoreProductContent.WithContext(s.ctx).
			Where(dao.MerStoreProductContent.ProductID.Eq(product.ProductID)).
			First()

		// 查询SKU列表
		skus, _ := dao.MerStoreProductSku.WithContext(s.ctx).
			Where(dao.MerStoreProductSku.ProductID.Eq(product.ProductID)).
			Find()

		result = append(result, &ProductDetailResponse{
			MerStoreProduct: product,
			Category:        category,
			Content:         content,
			Skus:            skus,
		})
	}

	return result, total, nil
}

// UpdateListingStatus 更新上架状态
func (s *StoreProductService) UpdateListingStatus(productID int32, merID int32, isShow int32) error {
	// 验证商品是否存在且属于该商户
	_, err := dao.MerStoreProduct.WithContext(s.ctx).
		Where(dao.MerStoreProduct.ProductID.Eq(productID)).
		Where(dao.MerStoreProduct.MerID.Eq(merID)).
		Where(dao.MerStoreProduct.DeleteAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("商品不存在或无权访问")
		}
		return fmt.Errorf("查询商品失败: %w", err)
	}

	// 更新状态
	_, err = dao.MerStoreProduct.WithContext(s.ctx).
		Where(dao.MerStoreProduct.ProductID.Eq(productID)).
		Updates(map[string]interface{}{
			"is_show": isShow,
		})
	return err
}

// UpdateSoldOutStatus 更新售完状态
func (s *StoreProductService) UpdateSoldOutStatus(productID int32, merID int32, saleStatus bool) error {
	// 验证商品是否存在且属于该商户
	_, err := dao.MerStoreProduct.WithContext(s.ctx).
		Where(dao.MerStoreProduct.ProductID.Eq(productID)).
		Where(dao.MerStoreProduct.MerID.Eq(merID)).
		Where(dao.MerStoreProduct.DeleteAt.IsNull()).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("商品不存在或无权访问")
		}
		return fmt.Errorf("查询商品失败: %w", err)
	}

	// 更新状态
	_, err = dao.MerStoreProduct.WithContext(s.ctx).
		Where(dao.MerStoreProduct.ProductID.Eq(productID)).
		Updates(map[string]interface{}{
			"sale_status": saleStatus,
		})
	return err
}

// Helper function
func boolPtr(b bool) *bool {
	return &b
}
