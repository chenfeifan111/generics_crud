package dbkit

import (
	"gorm.io/gorm"
)

// OrCondition 表示 OR 条件组
type OrCondition struct {
	Conditions []interface{} `json:"conditions"`
}

// ApplyOrConditions 应用 OR 条件
func (qb *QueryBuilder) ApplyOrConditions(orConditions []interface{}) *QueryBuilder {
	if len(orConditions) == 0 {
		return qb
	}

	// 构建 OR 条件
	for i, condition := range orConditions {
		tempDB := qb.db.Session(&gorm.Session{NewDB: true})
		tempQB := &QueryBuilder{db: tempDB}
		tempQB.ApplyFilters(condition)

		if i == 0 {
			qb.db = qb.db.Where(tempQB.db)
		} else {
			qb.db = qb.db.Or(tempQB.db)
		}
	}

	return qb
}

// ExtendedFilters 扩展的过滤器，支持 OR 条件
type ExtendedFilters struct {
	And interface{}   `json:"and"`
	Or  []interface{} `json:"or"`
}

// ApplyExtendedFilters 应用扩展过滤器
func (qb *QueryBuilder) ApplyExtendedFilters(filters *ExtendedFilters) *QueryBuilder {
	if filters == nil {
		return qb
	}

	// 应用 AND 条件
	if filters.And != nil {
		qb.ApplyFilters(filters.And)
	}

	// 应用 OR 条件
	if len(filters.Or) > 0 {
		qb.ApplyOrConditions(filters.Or)
	}

	return qb
}
