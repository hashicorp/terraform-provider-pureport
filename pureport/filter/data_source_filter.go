package filter

import (
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// Filter holds the associated values to filter on for the specified named field.
type Filter struct {
	Name   string
	Values []string
}

// Function to find the reflection field indices using the nested path name
// Since the field we want to filter on can be nested, we need to use the
// FieldByIndex() on the top level value. This takes the nested field indices
// so we need to traverse the class hierarchy to find the indices.
func findValue(inst interface{}, path string) (out reflect.Value) {

	// Get the Reflection Type
	e := reflect.Indirect(reflect.ValueOf(inst))
	t := reflect.TypeOf(inst)

	indices := []int{}

	for _, seg := range strings.Split(path, ".") {

		// If this is a map, we need to handle it here and return
		// Maps are only handled as the last lookup in the path.
		if t.Kind() == reflect.Map {

			log.Printf("Found Map: '%+v'\n", t)

			top := reflect.Indirect(reflect.ValueOf(inst))
			myMap := top.FieldByIndex(indices)

			iter := myMap.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()

				if k.String() == seg {
					out = v
				}
			}

			// Need to make sure we set this to a Zero value if nothing was found.
			if out == reflect.ValueOf(nil) {
				out = reflect.Zero(reflect.TypeOf(seg))
			}

			return
		}

		s, ok := t.FieldByName(seg)
		if !ok {

			log.Printf("Unable to find the specified path: %s", path)

			// Need to make sure we set this to a Zero value if nothing was found.
			if out == reflect.ValueOf(nil) {
				out = reflect.Zero(reflect.TypeOf(seg))
			}
			return
		}

		field := e.FieldByIndex(s.Index)
		indices = append(indices, s.Index...)

		switch field.Kind() {
		case reflect.Ptr:
			t = field.Type().Elem()
		default:
			t = field.Type()
		}

		e = reflect.Indirect(field)
	}

	top := reflect.Indirect(reflect.ValueOf(inst))
	out = top.FieldByIndex(indices)

	return
}

func FilterType(is []interface{}, filters []*Filter) []interface{} {

	matched := make(map[*interface{}]bool)

	for i, x := range is {

		matched[&is[i]] = true

		// Iterate through all of the filters
		for _, f := range filters {

			field := findValue(x, f.Name)
			log.Printf("Filter='%+v', Values: values='%s', Found:'%+v'\n", f, f.Values, field)

			if field.Interface() != reflect.Zero(field.Type()).Interface() {

				// Since we only need to match one filter, keep track
				// of the match status
				has_match := false

				// Check to see if the field value matches any of our filters
				for _, v := range f.Values {
					r := regexp.MustCompile(v)

					if r.MatchString(field.String()) {
						has_match = true
						log.Printf("Matched!!! v='%s', f='%s'\n", v, field)
						break
					}
				}

				// We didn't match any of the values for this field
				if !has_match {
					matched[&is[i]] = false
				}
			} else {
				matched[&is[i]] = false
			}

		}
	}

	log.Printf("Matched='%+v'\n", matched)

	results := make([]interface{}, 0)
	for k, v := range matched {
		if v {
			results = append(results, *k)
		}
	}

	return results
}

func BuildDataSourceFilters(set *schema.Set) []*Filter {
	var filters []*Filter
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		filters = append(filters, &Filter{
			Name:   m["name"].(string),
			Values: filterValues,
		})
	}
	return filters
}

func DataSourceFiltersSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},

				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}
