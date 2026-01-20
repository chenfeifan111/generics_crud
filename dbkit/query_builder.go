package dbkit

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type QueryBuilder struct {
	db *gorm.DB
}

func NewQueryBuilder(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{db: db}
}

func (qb *QueryBuilder) ApplyFilters(filters interface{}) *QueryBuilder {
	if filters == nil {
		return qb
	}

	filtersValue := reflect.ValueOf(filters)
	if filtersValue.Kind() == reflect.Ptr {
		filtersValue = filtersValue.Elem()
	}

	if filtersValue.Kind() != reflect.Struct {
		return qb
	}

	filtersType := filtersValue.Type()

	for i := 0; i < filtersValue.NumField(); i++ {
		field := filtersValue.Field(i)
		fieldType := filtersType.Field(i)

		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		columnName := strings.Split(jsonTag, ",")[0]
		filterTag := fieldType.Tag.Get("filter")

		var value interface{}
		if field.Kind() == reflect.Ptr {
			value = field.Elem().Interface()
		} else {
			value = field.Interface()
		}

		qb.applyFilter(columnName, filterTag, value)
	}

	return qb
}

func (qb *QueryBuilder) applyFilter(column, operator string, value interface{}) {
	switch operator {
	case "eq", "":
		qb.db = qb.db.Where(fmt.Sprintf("%s = ?", column), value)
	case "ne", "neq":
		qb.db = qb.db.Where(fmt.Sprintf("%s != ?", column), value)
	case "gt":
		qb.db = qb.db.Where(fmt.Sprintf("%s > ?", column), value)
	case "gte":
		qb.db = qb.db.Where(fmt.Sprintf("%s >= ?", column), value)
	case "lt":
		qb.db = qb.db.Where(fmt.Sprintf("%s < ?", column), value)
	case "lte":
		qb.db = qb.db.Where(fmt.Sprintf("%s <= ?", column), value)
	case "like":
		qb.db = qb.db.Where(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%v%%", value))
	case "in":
		qb.db = qb.db.Where(fmt.Sprintf("%s IN ?", column), value)
	case "not_in":
		qb.db = qb.db.Where(fmt.Sprintf("%s NOT IN ?", column), value)
	case "is_null":
		qb.db = qb.db.Where(fmt.Sprintf("%s IS NULL", column))
	case "is_not_null":
		qb.db = qb.db.Where(fmt.Sprintf("%s IS NOT NULL", column))
	case "between":
		// 处理范围查询
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return
			}
			v = v.Elem()
		}

		if v.Kind() != reflect.Struct {
			return
		}

		minField := v.FieldByName("Min")
		maxField := v.FieldByName("Max")

		if minField.IsValid() && minField.Kind() == reflect.Ptr && !minField.IsNil() {
			qb.db = qb.db.Where(fmt.Sprintf("%s >= ?", column), minField.Elem().Interface())
		}

		if maxField.IsValid() && maxField.Kind() == reflect.Ptr && !maxField.IsNil() {
			qb.db = qb.db.Where(fmt.Sprintf("%s <= ?", column), maxField.Elem().Interface())
		}
	}
}

func (qb *QueryBuilder) ApplyOrders(orders interface{}) *QueryBuilder {
	if orders == nil {
		return qb
	}

	ordersValue := reflect.ValueOf(orders)
	if ordersValue.Kind() == reflect.Ptr {
		ordersValue = ordersValue.Elem()
	}

	if ordersValue.Kind() != reflect.Struct {
		return qb
	}

	ordersType := ordersValue.Type()

	for i := 0; i < ordersValue.NumField(); i++ {
		field := ordersValue.Field(i)
		fieldType := ordersType.Field(i)

		if field.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}

		var directionRaw string
		switch field.Kind() {
		case reflect.Ptr:
			if field.Elem().Kind() != reflect.String {
				qb.db.AddError(fmt.Errorf("order field %s must be string", fieldType.Name))
				continue
			}
			directionRaw = field.Elem().String()
		case reflect.String:
			directionRaw = field.String()
		default:
			qb.db.AddError(fmt.Errorf("order field %s must be string", fieldType.Name))
			continue
		}

		direction := strings.ToLower(strings.TrimSpace(directionRaw))
		if direction == "" {
			qb.db.AddError(fmt.Errorf("order field %s must be 'asc' or 'desc'", fieldType.Name))
			continue
		}

		if direction != "asc" && direction != "desc" {
			qb.db.AddError(fmt.Errorf("order field %s must be 'asc' or 'desc'", fieldType.Name))
			continue
		}

		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		columnName := strings.Split(jsonTag, ",")[0]
		qb.db = qb.db.Order(fmt.Sprintf("%s %s", columnName, strings.ToUpper(direction)))
	}

	return qb
}

func (qb *QueryBuilder) ApplyPagination(page interface{}) *QueryBuilder {
	if page == nil {
		return qb
	}

	if p, ok := page.(Pageable); ok && p.IsValid() {
		qb.db = qb.db.Offset(p.GetOffset()).Limit(p.GetLimit())
	}

	return qb
}

func (qb *QueryBuilder) ApplyGroupBy(columns ...string) *QueryBuilder {
	if len(columns) == 0 {
		return qb
	}
	qb.db = qb.db.Group(strings.Join(columns, ", "))
	return qb
}

func (qb *QueryBuilder) ApplyHaving(condition string, args ...interface{}) *QueryBuilder {
	if condition == "" {
		return qb
	}
	qb.db = qb.db.Having(condition, args...)
	return qb
}

func (qb *QueryBuilder) ApplySelect(columns ...string) *QueryBuilder {
	if len(columns) == 0 {
		return qb
	}
	qb.db = qb.db.Select(columns)
	return qb
}

func (qb *QueryBuilder) Query(result interface{}) error {
	return qb.db.Find(result).Error
}

func (qb *QueryBuilder) QueryWithCount(result interface{}) (int64, error) {
	var count int64

	countDB := qb.db.Session(&gorm.Session{})
	if err := countDB.Count(&count).Error; err != nil {
		return 0, err
	}

	if err := qb.db.Find(result).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (qb *QueryBuilder) GetDB() *gorm.DB {
	return qb.db
}
