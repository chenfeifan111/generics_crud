package dbkit

// Range 表示范围查询
type Range[T any] struct {
	Min *T `json:"min"`
	Max *T `json:"max"`
}

// IsValid 检查范围是否有效
func (r *Range[T]) IsValid() bool {
	return r != nil && (r.Min != nil || r.Max != nil)
}
