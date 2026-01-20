package request

type GroupQueryRequest struct {
	GroupBy []string `json:"group_by" binding:"required"` // 分组字段，如 ["department"] 或 ["department", "name"]
	Select  []string `json:"select"`                      // 自定义查询字段，如 ["department", "COUNT(*) AS count"]
	Having  string   `json:"having"`                      // HAVING条件，如 "COUNT(*) > 1"
}
