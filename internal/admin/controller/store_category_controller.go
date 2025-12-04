package controller

import (
	"errors"
	"merchant_api/internal/admin/service"
	"merchant_api/internal/model"
	"merchant_api/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StoreCategoryController struct{}

func NewStoreCategoryController() *StoreCategoryController { return &StoreCategoryController{} }

// CreateRequest 创建请求
type CreateCategoryRequest struct {
	CateName string `json:"cate_name" binding:"required"`
	Pic      string `json:"pic"`
	Sort     int32  `json:"sort"`
}

// Create 创建分类
func (ctrl *StoreCategoryController) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	merId, err := getMerId(c)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	merIdInt32 := int32(merId)
	category := &model.MerStoreCategory{
		CateName: req.CateName,
		Pic:      req.Pic,
		Sort:     req.Sort,
		MerID:    &merIdInt32,
		Level:    1, // 默认为1级
	}

	svc := service.NewStoreCategoryService(c.Request.Context())
	if err := svc.Create(category); err != nil {
		response.InternalServerError(c, "创建失败："+err.Error())
		return
	}

	response.Success(c, category)
}

// UpdateRequest 更新请求
type UpdateCategoryRequest struct {
	CateName string `json:"cate_name"`
	Pic      string `json:"pic"`
	Sort     int32  `json:"sort"`
}

// Update 更新分类
func (ctrl *StoreCategoryController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "ID无效")
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	merId, err := getMerId(c)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	category := &model.MerStoreCategory{
		CateName: req.CateName,
		Pic:      req.Pic,
		Sort:     req.Sort,
	}

	svc := service.NewStoreCategoryService(c.Request.Context())
	if err := svc.Update(int32(id), int32(merId), category); err != nil {
		response.InternalServerError(c, "更新失败："+err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除分类
func (ctrl *StoreCategoryController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "ID无效")
		return
	}

	merId, err := getMerId(c)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	svc := service.NewStoreCategoryService(c.Request.Context())
	if err := svc.Delete(int32(id), int32(merId)); err != nil {
		response.InternalServerError(c, "删除失败："+err.Error())
		return
	}

	response.Success(c, nil)
}

// ListRequest 列表请求
type ListCategoryRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

// List 获取列表
func (ctrl *StoreCategoryController) List(c *gin.Context) {
	var req ListCategoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误："+err.Error())
		return
	}

	merId, err := getMerId(c)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	svc := service.NewStoreCategoryService(c.Request.Context())
	list, total, err := svc.GetList(int32(merId), req.Page, req.PageSize)
	if err != nil {
		response.InternalServerError(c, "获取列表失败："+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// Get 获取单个
func (ctrl *StoreCategoryController) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, "ID无效")
		return
	}

	merId, err := getMerId(c)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	svc := service.NewStoreCategoryService(c.Request.Context())
	category, err := svc.Get(int32(id), int32(merId))
	if err != nil {
		response.InternalServerError(c, "获取失败："+err.Error())
		return
	}

	response.Success(c, category)
}

// Helper function to get mer_id from context safely
func getMerId(c *gin.Context) (uint, error) {
	merIdValue, exists := c.Get("mer_id")
	if !exists {
		return 0, errors.New("商户ID不存在")
	}
	merIdUint, ok := merIdValue.(uint)
	if !ok {
		return 0, errors.New("商户ID类型错误")
	}
	return merIdUint, nil
}
