package controller

import (
	"errors"
	"merchant_api/internal/admin/service"
	"merchant_api/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StoreProductController struct{}

func NewStoreProductController() *StoreProductController {
	return &StoreProductController{}
}

// Create 创建商品
func (ctrl *StoreProductController) Create(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithKey(c, "error.invalid_params", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	merID, err := getMerID(c)
	if err != nil {
		response.InternalServerErrorWithKey(c, "error.merchant_id_missing")
		return
	}

	svc := service.NewStoreProductService(c.Request.Context())
	result, err := svc.Create(&req, int32(merID))
	if err != nil {
		response.BadRequestWithKey(c, "error.product.create_failed", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	response.SuccessWithKey(c, "success.product.created", result)
}

// Update 更新商品
func (ctrl *StoreProductController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequestWithKey(c, "error.invalid_id", nil)
		return
	}

	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithKey(c, "error.invalid_params", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	merID, err := getMerID(c)
	if err != nil {
		response.InternalServerErrorWithKey(c, "error.merchant_id_missing")
		return
	}

	svc := service.NewStoreProductService(c.Request.Context())
	if err := svc.Update(int32(id), &req, int32(merID)); err != nil {
		response.BadRequestWithKey(c, "error.product.update_failed", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	response.SuccessWithKey(c, "success.product.updated", nil)
}

// Delete 删除商品
func (ctrl *StoreProductController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequestWithKey(c, "error.invalid_id", nil)
		return
	}

	merID, err := getMerID(c)
	if err != nil {
		response.InternalServerErrorWithKey(c, "error.merchant_id_missing")
		return
	}

	svc := service.NewStoreProductService(c.Request.Context())
	if err := svc.Delete(int32(id), int32(merID)); err != nil {
		response.BadRequestWithKey(c, "error.product.delete_failed", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	response.SuccessWithKey(c, "success.product.deleted", nil)
}

// Get 获取商品详情
func (ctrl *StoreProductController) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequestWithKey(c, "error.invalid_id", nil)
		return
	}

	merID, err := getMerID(c)
	if err != nil {
		response.InternalServerErrorWithKey(c, "error.merchant_id_missing")
		return
	}

	svc := service.NewStoreProductService(c.Request.Context())
	product, err := svc.Get(int32(id), int32(merID))
	if err != nil {
		response.BadRequestWithKey(c, "error.product.not_found", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	response.Success(c, product)
}

// List 获取商品列表
func (ctrl *StoreProductController) List(c *gin.Context) {
	var req service.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequestWithKey(c, "error.invalid_params", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	merID, err := getMerID(c)
	if err != nil {
		response.InternalServerErrorWithKey(c, "error.merchant_id_missing")
		return
	}

	svc := service.NewStoreProductService(c.Request.Context())
	list, total, err := svc.GetList(int32(merID), &req)
	if err != nil {
		response.BadRequestWithKey(c, "error.product.list_failed", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	response.Success(c, gin.H{
		"list":      list,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// UpdateListingStatusRequest 更新上架状态请求
type UpdateListingStatusRequest struct {
	IsShow int32 `json:"is_show" binding:"required,oneof=0 1"`
}

// UpdateListingStatus 更新上架状态
func (ctrl *StoreProductController) UpdateListingStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequestWithKey(c, "error.invalid_id", nil)
		return
	}

	var req UpdateListingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithKey(c, "error.invalid_params", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	merID, err := getMerID(c)
	if err != nil {
		response.InternalServerErrorWithKey(c, "error.merchant_id_missing")
		return
	}

	svc := service.NewStoreProductService(c.Request.Context())
	if err := svc.UpdateListingStatus(int32(id), int32(merID), req.IsShow); err != nil {
		response.BadRequestWithKey(c, "error.product.update_listing_failed", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	response.SuccessWithKey(c, "success.product.listing_updated", nil)
}

// UpdateSoldOutStatusRequest 更新售完状态请求
type UpdateSoldOutStatusRequest struct {
	SaleStatus bool `json:"sale_status" binding:"required"`
}

// UpdateSoldOutStatus 更新售完状态
func (ctrl *StoreProductController) UpdateSoldOutStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequestWithKey(c, "error.invalid_id", nil)
		return
	}

	var req UpdateSoldOutStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestWithKey(c, "error.invalid_params", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	merID, err := getMerID(c)
	if err != nil {
		response.InternalServerErrorWithKey(c, "error.merchant_id_missing")
		return
	}

	svc := service.NewStoreProductService(c.Request.Context())
	if err := svc.UpdateSoldOutStatus(int32(id), int32(merID), req.SaleStatus); err != nil {
		response.BadRequestWithKey(c, "error.product.update_soldout_failed", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	response.SuccessWithKey(c, "success.product.soldout_updated", nil)
}

// Helper function to get mer_id from context safely
func getMerID(c *gin.Context) (uint, error) {
	merIDValue, exists := c.Get("mer_id")
	if !exists {
		return 0, errors.New("商户ID不存在")
	}
	merIDUint, ok := merIDValue.(uint)
	if !ok {
		return 0, errors.New("商户ID类型错误")
	}
	return merIDUint, nil
}
