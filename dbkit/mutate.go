package dbkit

import (
	"errors"
	"reflect"

	"gorm.io/gorm"
)

var ErrFilterRequired = errors.New("filters required")

type QueryRequest interface {
	GetPage() interface{}
	GetFilters() interface{}
	GetOrders() interface{}
}

type BaseQueryRequest[F any, O any] struct {
	Page    *Page `json:"page"`
	Filters F     `json:"filters"`
	Orders  O     `json:"orders"`
}

func (r *BaseQueryRequest[F, O]) GetPage() interface{} {
	return r.Page
}

func (r *BaseQueryRequest[F, O]) GetFilters() interface{} {
	return r.Filters
}

func (r *BaseQueryRequest[F, O]) GetOrders() interface{} {
	return r.Orders
}

func Query[T any](db *gorm.DB, req QueryRequest) ([]T, int64, error) {
	var results []T
	var model T

	qb := NewQueryBuilder(db)
	qb.db = qb.db.Model(&model)
	qb.ApplyFilters(req.GetFilters())
	qb.ApplyOrders(req.GetOrders())
	qb.ApplyPagination(req.GetPage())

	count, err := qb.QueryWithCount(&results)
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func QueryTo[T any, R any](db *gorm.DB, req QueryRequest) ([]R, int64, error) {
	var results []R
	var model T

	qb := NewQueryBuilder(db)
	qb.db = qb.db.Model(&model)
	qb.ApplyFilters(req.GetFilters())
	qb.ApplyOrders(req.GetOrders())
	qb.ApplyPagination(req.GetPage())

	count, err := qb.QueryWithCount(&results)
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func First[T any](db *gorm.DB, req QueryRequest) (*T, error) {
	var out T
	var model T

	qb := NewQueryBuilder(db)
	qb.db = qb.db.Model(&model)
	qb.ApplyFilters(req.GetFilters())
	qb.ApplyOrders(req.GetOrders())

	if err := qb.GetDB().First(&out).Error; err != nil {
		return nil, err
	}

	return &out, nil
}

func Create[T any](db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func Update[T any](db *gorm.DB, filters interface{}, updates map[string]interface{}) (int64, error) {
	if !HasAnyFilter(filters) {
		return 0, ErrFilterRequired
	}

	var model T
	qb := NewQueryBuilder(db)
	qb.db = qb.db.Model(&model)
	qb.ApplyFilters(filters)

	res := qb.GetDB().Updates(updates)
	return res.RowsAffected, res.Error
}

func Delete[T any](db *gorm.DB, filters interface{}) (int64, error) {
	if !HasAnyFilter(filters) {
		return 0, ErrFilterRequired
	}

	var model T
	qb := NewQueryBuilder(db)
	qb.db = qb.db.Model(&model)
	qb.ApplyFilters(filters)

	res := qb.GetDB().Unscoped().Delete(&model)
	return res.RowsAffected, res.Error
}

func HasAnyFilter(filters interface{}) bool {
	if filters == nil {
		return false
	}

	v := reflect.ValueOf(filters)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return true
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.IsValid() {
			continue
		}

		if f.Kind() == reflect.Ptr {
			if !f.IsNil() {
				return true
			}
			continue
		}

		zero := reflect.Zero(f.Type())
		if !reflect.DeepEqual(f.Interface(), zero.Interface()) {
			return true
		}
	}

	return false
}
