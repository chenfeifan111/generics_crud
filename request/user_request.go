package request

import "github.com/chenfeifan111/generics_crud/dbkit"

type UserFilters struct {
	ID   *string `json:"id" filter:"eq"`
	Age  *int    `json:"age" filter:"gte"`
	Name *string `json:"name" filter:"like"`
}

type UserOrders struct {
	Age *string `json:"age"`
	ID  *string `json:"id"`
}

// type UserQueryRequest = dbkit.BaseQueryRequest[UserFilters, struct{}]//如果不需要排序可以使用struct{}代替
type UserQueryRequest = dbkit.BaseQueryRequest[UserFilters, UserOrders]

type UserGetOneRequest struct {
	Filters struct {
		ID   *string `json:"id" filter:"eq"`
		Age  *int    `json:"age" filter:"gte"`
		Name *string `json:"name" filter:"like"`
	} `json:"filters"`
	Orders struct {
		Age *string `json:"age"`
		ID  *string `json:"id"`
	} `json:"orders"`
}

type UserUpdateByFiltersRequest struct {
	Filters struct {
		ID   *string `json:"id" filter:"eq"`
		Age  *int    `json:"age" filter:"gte"`
		Name *string `json:"name" filter:"like"`
	} `json:"filters"`
	Updates struct {
		Name *string `json:"name"`
		Age  *int    `json:"age"`
	} `json:"updates"`
}

type UserDeleteByFiltersRequest struct {
	Filters struct {
		ID   *string `json:"id" filter:"eq"`
		Age  *int    `json:"age" filter:"gte"`
		Name *string `json:"name" filter:"like"`
	} `json:"filters"`
}

type UserCreateRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required"`
}
