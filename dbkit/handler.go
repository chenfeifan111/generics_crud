package dbkit

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// generateUUID32 生成32位UUID（去掉横杠）
func generateUUID32() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// GenericQueryHandler 通用查询处理器
func GenericQueryHandler[T any, F any, O any](db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BaseQueryRequest[F, O]
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Error("Invalid request: "+err.Error()))
			return
		}

		results, count, err := Query[T](db, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessWithPage(results, req.Page, count))
	}
}

// GenericQueryToHandler 通用查询处理器（映射到DTO）
func GenericQueryToHandler[T any, R any, F any, O any](db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BaseQueryRequest[F, O]
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Error("Invalid request: "+err.Error()))
			return
		}

		results, count, err := QueryTo[T, R](db, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessWithPage(results, req.Page, count))
	}
}

// GenericCreateHandler 通用创建处理器
func GenericCreateHandler[T any](db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var entity T
		if err := c.ShouldBindJSON(&entity); err != nil {
			c.JSON(http.StatusBadRequest, Error("Invalid request: "+err.Error()))
			return
		}

		if err := Create(db, &entity); err != nil {
			c.JSON(http.StatusInternalServerError, Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Success(entity))
	}
}

// GenericCreateWithIDHandler 通用创建处理器（自动生成32位UUID）
func GenericCreateWithIDHandler[T any](db *gorm.DB, idSetter func(*T, string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var entity T
		if err := c.ShouldBindJSON(&entity); err != nil {
			c.JSON(http.StatusBadRequest, Error("Invalid request: "+err.Error()))
			return
		}

		// 生成32位UUID（去掉横杠）
		id := generateUUID32()
		idSetter(&entity, id)

		if err := Create(db, &entity); err != nil {
			c.JSON(http.StatusInternalServerError, Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Success(entity))
	}
}

// GenericUpdateHandler 通用更新处理器
func GenericUpdateHandler[T any, F any](db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Filters F                      `json:"filters"`
			Updates map[string]interface{} `json:"updates"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Error("Invalid request: "+err.Error()))
			return
		}

		if len(req.Updates) == 0 {
			c.JSON(http.StatusBadRequest, Error("No fields to update"))
			return
		}

		affected, err := Update[T](db, req.Filters, req.Updates)
		if err != nil {
			if err == ErrFilterRequired {
				c.JSON(http.StatusBadRequest, Error(err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Success(map[string]interface{}{
			"affected": affected,
		}))
	}
}

// GenericDeleteHandler 通用删除处理器
func GenericDeleteHandler[T any, F any](db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Filters F `json:"filters"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Error("Invalid request: "+err.Error()))
			return
		}

		affected, err := Delete[T](db, req.Filters)
		if err != nil {
			if err == ErrFilterRequired {
				c.JSON(http.StatusBadRequest, Error(err.Error()))
				return
			}
			c.JSON(http.StatusInternalServerError, Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Success(map[string]interface{}{
			"affected": affected,
		}))
	}
}

// GenericGetOneHandler 通用获取单条记录处理器
func GenericGetOneHandler[T any, F any, O any](db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BaseQueryRequest[F, O]
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Error("Invalid request: "+err.Error()))
			return
		}

		result, err := First[T](db, &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Success(result))
	}
}

// GenericStatsHandler 通用统计处理器
func GenericStatsHandler[T any, F any, O any](db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			BaseQueryRequest[F, O]
			StatsConfig StatsConfig `json:"stats_config"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, Error("Invalid request: "+err.Error()))
			return
		}

		stats, err := Stats[T](db, &req.BaseQueryRequest, req.StatsConfig)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Success(stats))
	}
}

// GenericBatchCreateHandler 通用批量创建处理器
func GenericBatchCreateHandler[T any](db *gorm.DB, batchSize int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var entities []T
		if err := c.ShouldBindJSON(&entities); err != nil {
			c.JSON(http.StatusBadRequest, Error("Invalid request: "+err.Error()))
			return
		}

		if err := BatchCreate(db, entities, batchSize); err != nil {
			c.JSON(http.StatusInternalServerError, Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, Success(map[string]interface{}{
			"created": len(entities),
		}))
	}
}
