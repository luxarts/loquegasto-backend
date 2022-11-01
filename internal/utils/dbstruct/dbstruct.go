package dbstruct

import "reflect"

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
