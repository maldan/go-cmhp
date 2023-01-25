package cmhp_sqllite

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func typeToSqlType(t string) string {
	switch t {
	case "int8":
		return "TINYINT"
	case "uint8":
		return "TINYINT UNSIGNED"
	case "int16":
		return "SMALLINT"
	case "uint16":
		return "SMALLINT UNSIGNED"
	case "int32":
		return "INTEGER"
	case "uint32":
		return "INTEGER UNSIGNED"
	case "string":
		return "TEXT"
	}
	return ""
}

func getValueFieldNames[T any](v T) []string {
	typeOf := reflect.TypeOf(v)
	out := make([]string, 0)

	for i := 0; i < typeOf.NumField(); i++ {
		fieldName := typeOf.Field(i).Name
		if typeOf.Field(i).Tag.Get("json") != "" {
			fieldName = typeOf.Field(i).Tag.Get("json")
		}
		out = append(out, "'"+fieldName+"'")
	}

	return out
}

func getValues[T any](v T) []any {
	typeOf := reflect.ValueOf(v)
	out := make([]any, 0)

	for i := 0; i < typeOf.NumField(); i++ {
		fieldValue := typeOf.Field(i).Interface()
		out = append(out, fieldValue)
	}

	return out
}

func CreateTable[T any](db *sql.DB, name string) error {
	out := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %v (\n", name)

	typeOf := reflect.TypeOf(*new(T))

	for i := 0; i < typeOf.NumField(); i++ {
		fieldName := typeOf.Field(i).Name
		fieldType := typeToSqlType(typeOf.Field(i).Type.Name())
		defaultValue := "DEFAULT 0"
		isNull := "NOT NULL"

		if typeOf.Field(i).Type.Name() == "string" {
			defaultValue = "DEFAULT \"\""
		}

		if typeOf.Field(i).Tag.Get("json") != "" {
			fieldName = typeOf.Field(i).Tag.Get("json")
		}

		out += fmt.Sprintf("\t\t\"%v\" %v %v %v", fieldName, fieldType, isNull, defaultValue)

		if i != typeOf.NumField()-1 {
			out += ","
		}
		out += "\n"
	}

	out += ");"

	// Prepare
	statement, err := db.Prepare(out)
	if err != nil {
		return err
	}

	// Execute
	_, err = statement.Exec()

	return err
}

func SelectOne[T any](db *sql.DB, from string, where string) string {
	return "SELECT * FROM user WHERE id=1 LIMIT 1"
}

func Insert[T any](db *sql.DB, table string, value T) error {
	fields := getValueFieldNames(value)
	values := getValues(value)
	valuesQ := make([]string, len(values))

	query := fmt.Sprintf("INSERT INTO '%v'(\n", table)

	for i := 0; i < len(fields); i++ {
		fields[i] = "\t" + fields[i]
		valuesQ[i] = "\t?"
	}

	query += "" + strings.Join(fields, ",\n ") + "\n) \n"
	query += "VALUES(\n" + strings.Join(valuesQ, ",\n") + "\n)"

	// Prepare
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}

	// Execute statement
	_, err = statement.Exec(values...)
	if err != nil {
		return err
	}

	return nil
}
