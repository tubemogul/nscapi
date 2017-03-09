package main

import (
	"reflect"
	"testing"
)

var lookupUseCase = []struct {
	custom      *customFields
	key         *fieldClassifier
	initialMap  map[string]interface{}
	expectedMap map[string]interface{}
}{
	{custom: &customFields{}, key: &fieldClassifier{hostgroup: "host1", service: "serviceA"}, initialMap: map[string]interface{}{}, expectedMap: map[string]interface{}{}},
	{custom: &customFields{}, key: &fieldClassifier{hostgroup: "host1", service: "serviceA"}, initialMap: map[string]interface{}{"foo": "Bar"}, expectedMap: map[string]interface{}{"foo": "Bar"}},
	{
		custom:      &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "host1", service: "serviceA"}: map[string]interface{}{"foo": "Bar"}}},
		key:         &fieldClassifier{hostgroup: "host1", service: "serviceA"},
		initialMap:  map[string]interface{}{},
		expectedMap: map[string]interface{}{"foo": "Bar"},
	},
	{
		custom:      &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "host1", service: "serviceA"}: map[string]interface{}{"foo": "NewVal"}}},
		key:         &fieldClassifier{hostgroup: "host1", service: "serviceA"},
		initialMap:  map[string]interface{}{"foo": "Bar"},
		expectedMap: map[string]interface{}{"foo": "NewVal"},
	},
	{
		custom:      &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "host1", service: "serviceB"}: map[string]interface{}{"foo": "Unrelated val"}}},
		key:         &fieldClassifier{hostgroup: "host1", service: "serviceA"},
		initialMap:  map[string]interface{}{"foo": "Bar"},
		expectedMap: map[string]interface{}{"foo": "Bar"},
	},
}

func TestLookup(t *testing.T) {
	for _, uc := range lookupUseCase {
		uc.custom.lookup(uc.initialMap, uc.key)
		if !reflect.DeepEqual(uc.initialMap, uc.expectedMap) {
			t.Fatalf("Expecting map to be %v after lookup. Got: %v\n", uc.expectedMap, uc.initialMap)
		}
	}
}
