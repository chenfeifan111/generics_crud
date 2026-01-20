package controller

import (
	"database/sql"
	"github.com/chenfeifan111/generics_crud/config"
	"github.com/chenfeifan111/generics_crud/dbkit"
	"github.com/chenfeifan111/generics_crud/request"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GroupQuery(c *gin.Context) {
	var req request.GroupQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dbkit.Error("Invalid request: "+err.Error()))
		return
	}

	if len(req.GroupBy) == 0 {
		c.JSON(http.StatusBadRequest, dbkit.Error("group_by is required"))
		return
	}

	var results []map[string]interface{}

	db := config.DB.Table("group_example")

	if len(req.Select) > 0 {
		db = db.Select(req.Select)
	}

	db = db.Group(strings.Join(req.GroupBy, ", "))

	if req.Having != "" {
		db = db.Having(req.Having)
	}

	rows, err := db.Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dbkit.Error(err.Error()))
		return
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]sql.RawBytes, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			c.JSON(http.StatusInternalServerError, dbkit.Error(err.Error()))
			return
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			if values[i] == nil {
				row[col] = nil
			} else {
				str := string(values[i])
				if intVal, err := strconv.ParseInt(str, 10, 64); err == nil {
					row[col] = intVal
				} else if floatVal, err := strconv.ParseFloat(str, 64); err == nil {
					row[col] = floatVal
				} else {
					row[col] = str
				}
			}
		}
		results = append(results, row)
	}

	c.JSON(http.StatusOK, dbkit.Success(results))
}
