package utils

import (
	"fmt"
	"reflect"
)

func NoIndent(jsLine string) string {
	return fmt.Sprint(jsLine + "\n")
}

func WithIndent(jsLine string, indents int) string {
	indent := ""

	for i := 0; i < indents; i++ {
		indent += "\t"
	}

	return fmt.Sprintf("%s%s%s", indent, jsLine, "\n")
}

func GenerateTSType(field interface{}) string {
	switch field.(type) {
	case string:
		return "string"
	case bool:
		return "boolean"
	default:
		return "number"
	}
}

type DTOFields struct {
	Names []string
	Types []string
}

func GetDTOFields(dto any) DTOFields {
	dtoVal := reflect.ValueOf(dto)
	typeOfDto := dtoVal.Type()

	var fieldNames []string
	var fieldTypes []string

	for i := 0;  i < dtoVal.NumField(); i++ {
		fieldNames = append(fieldNames, typeOfDto.Field(i).Name)
		fieldTypes = append(fieldNames, GenerateTSType(typeOfDto.Field(i).Type))
	}

	return DTOFields{fieldNames, fieldTypes}
} 