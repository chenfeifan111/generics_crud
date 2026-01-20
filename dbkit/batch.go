package dbkit

import (
	"gorm.io/gorm"
)

// BatchCreate 批量创建记录
func BatchCreate[T any](db *gorm.DB, entities []T, batchSize int) error {
	if len(entities) == 0 {
		return nil
	}

	if batchSize <= 0 {
		batchSize = 100 // 默认批次大小
	}

	return db.CreateInBatches(entities, batchSize).Error
}

// BatchUpdateByID 根据ID批量更新不同的值
type BatchUpdateItem struct {
	ID      interface{}            `json:"id"`
	Updates map[string]interface{} `json:"updates"`
}

func BatchUpdateByID[T any](db *gorm.DB, items []BatchUpdateItem) (int64, error) {
	if len(items) == 0 {
		return 0, nil
	}

	var totalAffected int64
	var model T

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if len(item.Updates) == 0 {
				continue
			}

			result := tx.Model(&model).Where("id = ?", item.ID).Updates(item.Updates)
			if result.Error != nil {
				return result.Error
			}
			totalAffected += result.RowsAffected
		}
		return nil
	})

	return totalAffected, err
}

// BatchDelete 批量删除（根据ID列表）
func BatchDelete[T any](db *gorm.DB, ids []interface{}) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	var model T
	result := db.Unscoped().Where("id IN ?", ids).Delete(&model)
	return result.RowsAffected, result.Error
}

// BatchUpdateByFilters 根据不同的过滤条件批量更新
type BatchUpdateByFilterItem struct {
	Filters interface{}            `json:"filters"`
	Updates map[string]interface{} `json:"updates"`
}

func BatchUpdateByFilters[T any](db *gorm.DB, items []BatchUpdateByFilterItem) (int64, error) {
	if len(items) == 0 {
		return 0, nil
	}

	var totalAffected int64
	var model T

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if !HasAnyFilter(item.Filters) || len(item.Updates) == 0 {
				continue
			}

			qb := NewQueryBuilder(tx.Model(&model))
			qb.ApplyFilters(item.Filters)

			result := qb.GetDB().Updates(item.Updates)
			if result.Error != nil {
				return result.Error
			}
			totalAffected += result.RowsAffected
		}
		return nil
	})

	return totalAffected, err
}
