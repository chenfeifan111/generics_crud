package main

import (
	"fmt"
	"github.com/chenfeifan111/generics_crud/config"
	"github.com/chenfeifan111/generics_crud/controller"
	"github.com/chenfeifan111/generics_crud/dbkit"
	"github.com/chenfeifan111/generics_crud/dto"
	"github.com/chenfeifan111/generics_crud/entity"
	"github.com/chenfeifan111/generics_crud/request"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.InitConfig()
	config.InitDB()

	r := gin.Default()

	// ============ 方式1: 使用简化的 Controller 函数 ============
	users := r.Group("/users")
	{
		// 查询列表
		users.POST("/query", controller.QueryUsersV2)
		users.POST("/query2", controller.QueryUsers2V2)

		// 获取单条
		users.POST("/one", controller.GetUserOneV2)

		// 创建
		users.POST("", controller.CreateUser)
		users.POST("/with-transaction", controller.CreateUserWithTransaction)

		// 更新
		users.POST("/update", controller.UpdateUsersV2)

		// 删除
		users.POST("/delete", controller.DeleteUsersV2)

		// 批量操作
		users.POST("/batch", controller.BatchCreateUsers)
		users.POST("/batch-update", controller.BatchUpdateUsers)
		users.POST("/batch-delete", controller.BatchDeleteUsers)

		// 统计
		users.POST("/stats", controller.GetUserStats)

		// 范围查询
		users.POST("/range", controller.QueryUsersWithRange)

		// OR 条件查询
		users.POST("/or", controller.QueryUsersWithOr)
	}

	// ============ 方式2: 直接在路由中使用通用 Handler ============
	usersV3 := r.Group("/users-v3")
	{
		// 查询（返回完整实体）
		usersV3.POST("/query",
			dbkit.GenericQueryHandler[entity.User, request.UserFilters, request.UserOrders](config.DB))

		// 查询（返回DTO）
		usersV3.POST("/query2",
			dbkit.GenericQueryToHandler[entity.User, dto.UserQuery2Item, request.UserFilters, request.UserOrders](config.DB))

		// 获取单条
		usersV3.POST("/one",
			dbkit.GenericGetOneHandler[entity.User, request.UserFilters, request.UserOrders](config.DB))

		// 创建
		usersV3.POST("",
			dbkit.GenericCreateHandler[entity.User](config.DB))

		// 更新
		usersV3.POST("/update",
			dbkit.GenericUpdateHandler[entity.User, request.UserFilters](config.DB))

		// 删除
		usersV3.POST("/delete",
			dbkit.GenericDeleteHandler[entity.User, request.UserFilters](config.DB))

		// 统计
		usersV3.POST("/stats",
			dbkit.GenericStatsHandler[entity.User, request.UserFilters, request.UserOrders](config.DB))

		// 批量创建
		usersV3.POST("/batch",
			dbkit.GenericBatchCreateHandler[entity.User](config.DB, 100))
	}

	// ============ 分组查询示例 ============
	groupExample := r.Group("/group-example")
	{
		groupExample.POST("/group", controller.GroupQuery)
	}

	port := viper.GetInt("server.port")
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
}
