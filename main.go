package main

import (
	"fmt"
	"github.com/chenfeifan111/generics_crud/config"
	"github.com/chenfeifan111/generics_crud/controller"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.InitConfig()
	config.InitDB()

	r := gin.Default()

	users := r.Group("/users")
	{
		// 基础 CRUD（使用通用 Handler）
		users.POST("/query", controller.QueryUsersV2)      //查询列表(返回全部字段，支持范围查询)
		users.POST("/query2", controller.QueryUsers2V2)    //查询列表(自定义返回字段示例)
		users.POST("", controller.CreateUserV2)            //新增一个（通用Handler）
		users.POST("/native", controller.CreateUserNative) //新增一个（原生实现对比）
		users.POST("/one", controller.GetUserOneV2)        //获取一条记录
		users.POST("/update", controller.UpdateUsersV2)    //更新
		users.POST("/delete", controller.DeleteUsersV2)    //删除

		// 高级功能
		users.POST("/query-or", controller.QueryUsersWithOr)     //查询列表(OR条件示例)
		users.POST("/batch", controller.BatchCreateUsers)        //批量创建
		users.POST("/batch-update", controller.BatchUpdateUsers) //批量更新
		users.POST("/batch-delete", controller.BatchDeleteUsers) //批量删除
		users.POST("/stats", controller.GetUserStats)            //统计查询
	}

	groupExample := r.Group("/group-example")
	{
		groupExample.POST("/group", controller.GroupQuery) //分组查询示例(自定义分组字段)
	}

	port := viper.GetInt("server.port")
	err := r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
}
