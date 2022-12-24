package dbstruct

import (
	"reflect"

	sq "github.com/Masterminds/squirrel"
)

func GetColumns(structure interface{}) []string {
	var columns []string

	refElem := reflect.TypeOf(structure).Elem()
	for i := 0; i < refElem.NumField(); i++ {
		columns = append(columns, refElem.Field(i).Tag.Get("db"))
	}

	return columns
}
func GetValues(structure interface{}) []interface{} {
	var values []interface{}

	ref := reflect.ValueOf(structure).Elem()
	for i := 0; i < ref.NumField(); i++ {
		values = append(values, ref.Field(i).Interface())
	}

	return values
}
func SetValues(builder sq.UpdateBuilder, structure interface{}) sq.UpdateBuilder {
	columns := GetColumns(structure)
	values := GetValues(structure)

	for i := range columns {
		builder = builder.Set(columns[i], values[i])
	}

	return builder
}
