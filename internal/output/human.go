package output

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

// HumanFormatter outputs data in human-readable format.
type HumanFormatter struct{}

// Print outputs data with pretty formatting for structs.
func (f *HumanFormatter) Print(data interface{}) error {
	if data == nil {
		return nil
	}

	v := reflect.ValueOf(data)

	// Dereference pointer if needed
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		f.printStruct(v, "")
	case reflect.Slice, reflect.Array:
		f.printSlice(v)
	case reflect.Map:
		f.printMap(v)
	default:
		fmt.Println(data)
	}

	return nil
}

// printStruct prints a struct in key: value format.
func (f *HumanFormatter) printStruct(v reflect.Value, indent string) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get field name from json tag if available, otherwise use field name
		name := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				name = parts[0]
			}
		}

		// Handle nested structs
		if value.Kind() == reflect.Struct {
			fmt.Printf("%s%s:\n", indent, name)
			f.printStruct(value, indent+"  ")
		} else if value.Kind() == reflect.Ptr && !value.IsNil() && value.Elem().Kind() == reflect.Struct {
			fmt.Printf("%s%s:\n", indent, name)
			f.printStruct(value.Elem(), indent+"  ")
		} else {
			fmt.Printf("%s%s: %v\n", indent, name, value.Interface())
		}
	}
}

// printSlice prints a slice with numbered items.
func (f *HumanFormatter) printSlice(v reflect.Value) {
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Struct || (item.Kind() == reflect.Ptr && !item.IsNil() && item.Elem().Kind() == reflect.Struct) {
			fmt.Printf("[%d]\n", i+1)
			if item.Kind() == reflect.Ptr {
				f.printStruct(item.Elem(), "  ")
			} else {
				f.printStruct(item, "  ")
			}
			if i < v.Len()-1 {
				fmt.Println()
			}
		} else {
			fmt.Printf("[%d] %v\n", i+1, item.Interface())
		}
	}
}

// printMap prints a map in key: value format.
func (f *HumanFormatter) printMap(v reflect.Value) {
	iter := v.MapRange()
	for iter.Next() {
		fmt.Printf("%v: %v\n", iter.Key().Interface(), iter.Value().Interface())
	}
}

// PrintError outputs an error message to stderr in red.
func (f *HumanFormatter) PrintError(err error) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Fprintln(os.Stderr, red("Error:"), err.Error())
}

// PrintMessage outputs a simple message to stdout.
func (f *HumanFormatter) PrintMessage(msg string) {
	fmt.Println(msg)
}

// PrintTable outputs tabular data with aligned columns.
func (f *HumanFormatter) PrintTable(headers []string, rows [][]string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print headers
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print separator
	separators := make([]string, len(headers))
	for i, h := range headers {
		separators[i] = strings.Repeat("-", len(h))
	}
	fmt.Fprintln(w, strings.Join(separators, "\t"))

	// Print rows
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	w.Flush()
}
