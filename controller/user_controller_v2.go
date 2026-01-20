package controller

import (
	"net/http"
	"strings"

	"github.com/chenfeifan111/generics_crud/config"
	"github.com/chenfeifan111/generics_crud/dbkit"
	"github.com/chenfeifan111/generics_crud/dto"
	"github.com/chenfeifan111/generics_crud/entity"
	"github.com/chenfeifan111/generics_crud/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ============ 使用通用 Handler 的简化版本 ============

// QueryUsersV2 使用通用处理器的查询
func QueryUsersV2(c *gin.Context) {
	dbkit.GenericQueryHandler[entity.User, request.UserFilters, request.UserOrders](config.DB)(c)
}

// QueryUsers2V2 使用通用处理器的查询（映射到DTO）
func QueryUsers2V2(c *gin.Context) {
	dbkit.GenericQueryToHandler[entity.User, dto.UserQuery2Item, request.UserFilters, request.UserOrders](config.DB)(c)
}

// GetUserOneV2 使用通用处理器获取单条记录
func GetUserOneV2(c *gin.Context) {
	dbkit.GenericGetOneHandler[entity.User, request.UserFilters, request.UserOrders](config.DB)(c)
}

// UpdateUsersV2 使用通用处理器更新
func UpdateUsersV2(c *gin.Context) {
	dbkit.GenericUpdateHandler[entity.User, request.UserFilters](config.DB)(c)
}

// DeleteUsersV2 使用通用处理器删除
func DeleteUsersV2(c *gin.Context) {
	dbkit.GenericDeleteHandler[entity.User, request.UserFilters](config.DB)(c)
}

// ============ 新增功能示例 ============

// CreateUserV2 使用通用处理器创建（自动生成ID）
func CreateUserV2(c *gin.Context) {
	dbkit.GenericCreateWithIDHandler[entity.User](config.DB, func(u *entity.User, id string) {
		u.ID = id
	})(c)
}

// CreateUserNative 原生GORM实现示例（完全不使用dbkit，对比用）
func CreateUserNative(c *gin.Context) {
	var req request.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid request: " + err.Error(),
			"data": nil,
		})
		return
	}

	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	user := entity.User{
		ID:   id,
		Name: req.Name,
		Age:  req.Age,
	}

	// 直接使用原生 GORM
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": user,
	})
}

// QueryUsersWithOr OR条件查询示例
func QueryUsersWithOr(c *gin.Context) {
	var req struct {
		Page    *dbkit.Page `json:"page"`
		Filters struct {
			// OR 条件：age >= 30 OR name LIKE '%admin%'
			Or []struct {
				Age  *int    `json:"age" filter:"gte"`
				Name *string `json:"name" filter:"like"`
			} `json:"or"`
		} `json:"filters"`
		Orders struct {
			Age *string `json:"age"`
			ID  *string `json:"id"`
		} `json:"orders"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dbkit.Error("Invalid request: "+err.Error()))
		return
	}

	var users []entity.User
	var model entity.User

	qb := dbkit.NewQueryBuilder(config.DB.Model(&model))

	// 应用 OR 条件
	if len(req.Filters.Or) > 0 {
		orConditions := make([]interface{}, len(req.Filters.Or))
		for i, cond := range req.Filters.Or {
			orConditions[i] = cond
		}
		qb.ApplyOrConditions(orConditions)
	}

	qb.ApplyOrders(req.Orders)
	qb.ApplyPagination(req.Page)

	count, err := qb.QueryWithCount(&users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dbkit.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dbkit.SuccessWithPage(users, req.Page, count))
}

// GetUserStats 获取用户统计信息
func GetUserStats(c *gin.Context) {
	var req request.UserQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dbkit.Error("Invalid request: "+err.Error()))
		return
	}

	stats, err := dbkit.Stats[entity.User](config.DB, &req, dbkit.StatsConfig{
		SumFields: []string{"age"},
		AvgFields: []string{"age"},
		MinFields: []string{"age"},
		MaxFields: []string{"age"},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dbkit.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dbkit.Success(stats))
}

// BatchCreateUsers 批量创建用户
func BatchCreateUsers(c *gin.Context) {
	var users []entity.User
	if err := c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, dbkit.Error("Invalid request: "+err.Error()))
		return
	}

	// 为每个用户生成ID
	for i := range users {
		users[i].ID = strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	// 批量创建，每批100条
	if err := dbkit.BatchCreate(config.DB, users, 100); err != nil {
		c.JSON(http.StatusInternalServerError, dbkit.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dbkit.Success(map[string]interface{}{
		"created": len(users),
		"users":   users,
	}))
}

// BatchUpdateUsers 批量更新用户（不同ID不同值）
func BatchUpdateUsers(c *gin.Context) {
	var items []dbkit.BatchUpdateItem
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, dbkit.Error("Invalid request: "+err.Error()))
		return
	}

	affected, err := dbkit.BatchUpdateByID[entity.User](config.DB, items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dbkit.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dbkit.Success(map[string]interface{}{
		"affected": affected,
	}))
}

// BatchDeleteUsers 批量删除用户
func BatchDeleteUsers(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dbkit.Error("Invalid request: "+err.Error()))
		return
	}

	ids := make([]interface{}, len(req.IDs))
	for i, id := range req.IDs {
		ids[i] = id
	}

	affected, err := dbkit.BatchDelete[entity.User](config.DB, ids)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dbkit.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dbkit.Success(map[string]interface{}{
		"affected": affected,
	}))
}
