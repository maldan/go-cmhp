package cmhp_sqllite

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func typeToSqlType(t string, size string) string {
	switch t {
	case "int8":
		return "TINYINT"
	case "uint8":
		return "TINYINT UNSIGNED"
	case "int16":
		return "SMALLINT"
	case "uint16":
		return "SMALLINT UNSIGNED"
	case "int32", "int":
		return "INTEGER"
	case "uint32", "uint":
		return "INTEGER UNSIGNED"
	case "string":
		if size == "" {
			return "TEXT"
		} else {
			return "VARCHAR(" + size + ")"
		}
	case "Time":
		return "DATETIME"
	default:
		panic("unknown type " + t)
	}
	return ""
}

func getValueFieldNames[T any](v T, useQuotes bool) []string {
	typeOf := reflect.TypeOf(v)
	out := make([]string, 0)

	for i := 0; i < typeOf.NumField(); i++ {
		fieldName := typeOf.Field(i).Name
		if typeOf.Field(i).Tag.Get("json") != "" {
			fieldName = typeOf.Field(i).Tag.Get("json")
		}
		if useQuotes {
			out = append(out, "`"+fieldName+"`")
		} else {
			out = append(out, fieldName)
		}
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
		fieldType := typeToSqlType(typeOf.Field(i).Type.Name(), typeOf.Field(i).Tag.Get("len"))
		opts := ""

		if typeOf.Field(i).Name == "Id" {
			fieldType = "INTEGER"
			opts += "PRIMARY KEY AUTOINCREMENT "
		} else {
			if typeOf.Field(i).Type.Name() == "string" {
				opts += "DEFAULT \"\" "
			} else {
				opts += "DEFAULT 0 "
			}
		}
		opts += "NOT NULL "

		if typeOf.Field(i).Tag.Get("json") != "" {
			fieldName = typeOf.Field(i).Tag.Get("json")
		}

		out += fmt.Sprintf("\t\t`%v` %v %v", fieldName, fieldType, opts)

		if i != typeOf.NumField()-1 {
			out += ","
		}
		out += "\n"
	}

	out += ");\n"
	fmt.Printf("%v", out)
	// Prepare
	statement, err := db.Prepare(out)
	if err != nil {
		return err
	}

	// Execute
	_, err = statement.Exec()

	return err
}

func SelectOne[T any](db *sql.DB, from string, where string, values ...any) (T, error) {
	out := *new(T)
	outType := reflect.TypeOf(&out).Elem()

	fields := getValueFieldNames(out, false)
	query := fmt.Sprintf("SELECT %v FROM %v WHERE %v LIMIT 1", strings.Join(fields, ","), from, where)

	destForScan := make([]any, len(fields))
	for i := 0; i < len(fields); i++ {
		ptr := unsafe.Add(unsafe.Pointer(&out), outType.Field(i).Offset)

		if outType.Field(i).Type.Name() == "Time" {
			destForScan[i] = (*time.Time)(ptr)
		} else if outType.Field(i).Type.Kind() == reflect.String {
			destForScan[i] = (*string)(ptr)
		} else if outType.Field(i).Type.Kind() == reflect.Uint32 {
			destForScan[i] = (*uint32)(ptr)
		} else if outType.Field(i).Type.Kind() == reflect.Int {
			destForScan[i] = (*int)(ptr)
		} else if outType.Field(i).Type.Kind() == reflect.Int8 {
			destForScan[i] = (*int8)(ptr)
		} else {
			panic("unsupported type " + outType.Field(i).Type.Name())
		}
	}

	// Prepare
	statement, err := db.Prepare(query)
	defer statement.Close()
	if err != nil {
		return out, err
	}

	// Execute statement
	rows, err := statement.Query(values...)
	defer rows.Close()
	if err != nil {
		return out, err
	}

	// Scan rows
	for rows.Next() {
		err2 := rows.Scan(destForScan...)
		if err2 != nil {
			return out, err2
		}
	}

	return out, err
}

func Insert[T any](db *sql.DB, table string, value T) error {
	fields := getValueFieldNames(value, true)
	values := getValues(value)
	valuesQ := make([]string, len(values))

	query := fmt.Sprintf("INSERT INTO '%v'(\n", table)

	for i := 0; i < len(fields); i++ {
		if fields[i] == "`id`" && reflect.ValueOf(values[i]).IsZero() {
			values[i] = nil
		}

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

	fmt.Printf("%v\n", query)
	fmt.Printf("%v\n", values[0])

	// Execute statement
	_, err = statement.Exec(values...)
	if err != nil {
		return err
	}

	return nil
}

func InsertMany[T any](db *sql.DB, table string, vs []T) error {
	fields := getValueFieldNames(vs[0], true)
	// v0 := getValues(vs[0])
	valuesQ := make([]string, len(getValues(vs[0])))
	valuesOut := make([]any, 0)
	query := fmt.Sprintf("INSERT INTO '%v'(\n", table)

	for i := 0; i < len(vs); i++ {
		vals := getValues(vs[i])
		for j := 0; j < len(fields); j++ {
			if fields[j] == "`id`" && reflect.ValueOf(vs[i]).Field(j).IsZero() {
				valuesOut = append(valuesOut, nil)
			} else {
				valuesOut = append(valuesOut, vals[j])
			}

			valuesQ[j] = "?"
		}
	}

	query += "" + strings.Join(fields, ",\n ") + "\n) \n"
	query += "VALUES"

	for i := 0; i < len(vs); i++ {
		query += "(\n" + strings.Join(valuesQ, ",\n") + "\n)"
		if i <= len(vs)-2 {
			query += ","
		}
	}

	// Prepare
	statement, err := db.Prepare(query)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", query)
	fmt.Print(valuesOut)

	// Execute statement
	_, err = statement.Exec(valuesOut...)
	if err != nil {
		return err
	}

	return nil
}

func Raw(db *sql.DB, query string, values ...any) error {
	// Prepare
	statement, err := db.Prepare(query)
	defer statement.Close()
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
