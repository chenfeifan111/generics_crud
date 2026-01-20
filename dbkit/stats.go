package dbkit

import (
	"fmt"
	"gorm.io/gorm"
)

// QueryStats 查询统计信息
type QueryStats struct {
	Count int64                  `json:"count"`
	Sum   map[string]float64     `json:"sum,omitempty"`
	Avg   map[string]float64     `json:"avg,omitempty"`
	Min   map[string]interface{} `json:"min,omitempty"`
	Max   map[string]interface{} `json:"max,omitempty"`
}

// StatsConfig 统计配置
type StatsConfig struct {
	SumFields []string `json:"sum_fields"` // 需要求和的字段
	AvgFields []string `json:"avg_fields"` // 需要求平均的字段
	MinFields []string `json:"min_fields"` // 需要求最小值的字段
	MaxFields []string `json:"max_fields"` // 需要求最大值的字段
}

// Stats 执行统计查询
func Stats[T any](db *gorm.DB, req QueryRequest, config StatsConfig) (*QueryStats, error) {
	var model T
	stats := &QueryStats{
		Sum: make(map[string]float64),
		Avg: make(map[string]float64),
		Min: make(map[string]interface{}),
		Max: make(map[string]interface{}),
	}

	qb := NewQueryBuilder(db.Model(&model))
	qb.ApplyFilters(req.GetFilters())

	// 获取总数
	if err := qb.db.Count(&stats.Count).Error; err != nil {
		return nil, err
	}

	// 构建 SELECT 语句
	var selectFields []string

	for _, field := range config.SumFields {
		selectFields = append(selectFields, fmt.Sprintf("SUM(%s) as sum_%s", field, field))
	}

	for _, field := range config.AvgFields {
		selectFields = append(selectFields, fmt.Sprintf("AVG(%s) as avg_%s", field, field))
	}

	for _, field := range config.MinFields {
		selectFields = append(selectFields, fmt.Sprintf("MIN(%s) as min_%s", field, field))
	}

	for _, field := range config.MaxFields {
		selectFields = append(selectFields, fmt.Sprintf("MAX(%s) as max_%s", field, field))
	}

	if len(selectFields) == 0 {
		return stats, nil
	}

	// 执行聚合查询
	var result map[string]interface{}
	qb2 := NewQueryBuilder(db.Model(&model))
	qb2.ApplyFilters(req.GetFilters())

	// 直接在 db 上执行 Select 和 Find
	if err := qb2.db.Select(selectFields).Limit(1).Find(&result).Error; err != nil {
		return nil, err
	}

	// 解析结果
	for _, field := range config.SumFields {
		key := fmt.Sprintf("sum_%s", field)
		if val, ok := result[key]; ok && val != nil {
			if floatVal, ok := val.(float64); ok {
				stats.Sum[field] = floatVal
			}
		}
	}

	for _, field := range config.AvgFields {
		key := fmt.Sprintf("avg_%s", field)
		if val, ok := result[key]; ok && val != nil {
			if floatVal, ok := val.(float64); ok {
				stats.Avg[field] = floatVal
			}
		}
	}

	for _, field := range config.MinFields {
		key := fmt.Sprintf("min_%s", field)
		if val, ok := result[key]; ok && val != nil {
			stats.Min[field] = val
		}
	}

	for _, field := range config.MaxFields {
		key := fmt.Sprintf("max_%s", field)
		if val, ok := result[key]; ok && val != nil {
			stats.Max[field] = val
		}
	}

	return stats, nil
}

// SimpleStats 简单统计（只返回总数、总和、平均值）
func SimpleStats[T any](db *gorm.DB, req QueryRequest, sumFields []string) (*QueryStats, error) {
	return Stats[T](db, req, StatsConfig{
		SumFields: sumFields,
		AvgFields: sumFields,
	})
}
