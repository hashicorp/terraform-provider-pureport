package filter

import (
	"flag"
	"sort"
	"testing"

	"github.com/pureport/pureport-sdk-go/pureport/client"
)

var (
	items = []interface{}{
		client.Account{Name: "Testing 1", Description: "First Test Account"},
		client.Account{Name: "Testing 2", Description: "Second Test Account"},
		client.Account{Name: "Testing 3", Description: "Third Test Account"},
	}

	nested_items = []interface{}{
		client.Connection{Name: "TestConnection 1", Location: &client.Link{Title: "Raleigh"}},
		client.Connection{Name: "TestConnection 2", Location: &client.Link{Title: "San Jose"}},
		client.Connection{Name: "TestConnection 3", Location: &client.Link{Title: "Seattle"}},
	}

	tagged_items = []interface{}{
		client.Account{
			Name:        "Testing 1",
			Description: "First Test Account",
			Tags:        map[string]string{"some_name": "value1"},
		},
		client.Account{
			Name:        "Testing 2",
			Description: "Second Test Account",
			Tags:        map[string]string{"some_name": "value2"},
		},
		client.Account{
			Name:        "Testing 3",
			Description: "Third Test Account",
			Tags:        map[string]string{"some_name": "value3"},
		},
	}
)

func init() {
	var _ *string = flag.String("sweep", "", "Eat the sweep for unit tests")
}

func TestFilterNestedItems(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Location.Title", Values: []string{"Raleigh"}},
	}

	results := FilterType(nested_items, filter)

	if len(results) != 1 {
		t.Errorf("Name filter failed: expected '%d', got: '%d'", 1, len(results))
	}
}

func TestFilterNameNoItems(t *testing.T) {

	no_items := []interface{}{
		client.Account{Name: "Blah", Description: "Blah Blah Blah"},
	}

	filter := []*Filter{
		&Filter{Name: "Name", Values: []string{"Testing"}},
	}

	results := FilterType(no_items, filter)

	if len(results) != 0 {
		t.Errorf("Name filter failed: expected '%d', got: '%d'", 0, len(results))
	}
}

func TestFilterInvalid(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Blah", Values: []string{"Testing"}},
	}

	results := FilterType(items, filter)

	if len(results) != 0 {
		t.Errorf("Name filter failed: expected '%d', got: '%d'", 0, len(results))
	}
}

func TestFilterNameGeneric(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Name", Values: []string{"Testing"}},
	}

	results := FilterType(items, filter)

	if len(results) != 3 {
		t.Errorf("Name filter failed: expected '%d', got: '%d'", 3, len(results))
	}
}

func TestFilterNameSpecific1(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Name", Values: []string{"ting 1"}},
	}

	results := FilterType(items, filter)

	if len(results) != 1 {
		t.Errorf("Name filter failed: expected: '%d', got: '%d'", 1, len(results))
		return
	}

	result := results[0].(client.Account)

	if result.Name != "Testing 1" {
		t.Errorf("Invalid name: expected: 'Testing 1', got '%s'", result.Name)
	}
}

func TestFilterNameSpecific2(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Name", Values: []string{"Testing 2"}},
	}

	results := FilterType(items, filter)

	if len(results) != 1 {
		t.Errorf("Name filter failed: expected: '%d', got: '%d'", 1, len(results))
		return
	}

	result := results[0].(client.Account)

	if result.Name != "Testing 2" {
		t.Errorf("Invalid name: expected: 'Testing 2', got '%s'", result.Name)
	}
}

func TestFilterDescriptionGeneric(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Description", Values: []string{"Test Account"}},
	}

	results := FilterType(items, filter)

	if len(results) != 3 {
		t.Errorf("Description filter failed: expected '%d', got: '%d'", 3, len(results))
	}
}

func TestFilterDescriptionSpecific1(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Description", Values: []string{"Third"}},
	}

	results := FilterType(items, filter)

	if len(results) != 1 {
		t.Errorf("Description filter failed: expected: '%d', got: '%d'", 1, len(results))
		return
	}

	result := results[0].(client.Account)

	if result.Name != "Testing 3" {
		t.Errorf("Invalid name: expected: 'Testing 1', got '%s'", result.Name)
	}
}

func TestFilterDescriptionSpecific2(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Description", Values: []string{"First", "Second"}},
	}

	results := FilterType(items, filter)

	if len(results) != 2 {
		t.Errorf("Description filter failed: expected: '%d', got: '%d'", 2, len(results))
		return
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].(client.Account).Name < results[j].(client.Account).Name
	})

	result := results[0].(client.Account)
	if result.Name != "Testing 1" {
		t.Errorf("Invalid name: expected: 'Testing 1', got '%s'", result.Name)
	}

	result = results[1].(client.Account)
	if result.Name != "Testing 2" {
		t.Errorf("Invalid name: expected: 'Testing 2', got '%s'", result.Name)
	}
}

func TestFilterMap1(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Tags.some_name", Values: []string{"value1"}},
	}

	results := FilterType(tagged_items, filter)

	if len(results) != 1 {
		t.Errorf("Description filter failed: expected: '%d', got: '%d'", 1, len(results))
		return
	}

	result := results[0].(client.Account)
	if result.Name != "Testing 1" {
		t.Errorf("Invalid name: expected: 'Testing 1', got '%s'", result.Name)
	}
}

func TestFilterMap2(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Tags.some_name", Values: []string{"value2", "value3"}},
	}

	results := FilterType(tagged_items, filter)

	if len(results) != 2 {
		t.Errorf("Description filter failed: expected: '%d', got: '%d'", 2, len(results))
		return
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].(client.Account).Name < results[j].(client.Account).Name
	})

	result := results[0].(client.Account)
	if result.Name != "Testing 2" {
		t.Errorf("Invalid name: expected: 'Testing 2', got '%s'", result.Name)
	}

	result = results[1].(client.Account)
	if result.Name != "Testing 3" {
		t.Errorf("Invalid name: expected: 'Testing 3', got '%s'", result.Name)
	}
}

func TestFilterMapError(t *testing.T) {

	filter := []*Filter{
		&Filter{Name: "Tags.some_name", Values: []string{"value1"}},
	}

	results := FilterType(items, filter)

	if len(results) != 0 {
		t.Errorf("Description filter failed: expected: '%d', got: '%d'", 1, len(results))
	}
}
