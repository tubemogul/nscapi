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
	// Emtpy custom fields map & passing empty map
	{custom: &customFields{}, key: &fieldClassifier{hostgroup: "web", service: "serviceA"}, initialMap: map[string]interface{}{}, expectedMap: map[string]interface{}{}},
	// Emtpy custom fields map & passing non empty map
	{custom: &customFields{}, key: &fieldClassifier{hostgroup: "web", service: "serviceA"}, initialMap: map[string]interface{}{"foo": "Bar"}, expectedMap: map[string]interface{}{"foo": "Bar"}},
	// Lookup for an existing classifier with data and passing a empty map
	{
		custom:      &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "web", service: "serviceA"}: map[string]interface{}{"foo": "Bar"}}},
		key:         &fieldClassifier{hostgroup: "web", service: "serviceA"},
		initialMap:  map[string]interface{}{},
		expectedMap: map[string]interface{}{"foo": "Bar"},
	},
	// Lookup for an existing classifier with data and passing a non empty map
	{
		custom:      &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "web", service: "serviceA"}: map[string]interface{}{"foo": "NewVal"}}},
		key:         &fieldClassifier{hostgroup: "web", service: "serviceA"},
		initialMap:  map[string]interface{}{"foo": "Bar"},
		expectedMap: map[string]interface{}{"foo": "NewVal"},
	},
	// Lookup for annon existing classifier with data and passing a non empty map
	{
		custom:      &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "db", service: "serviceB"}: map[string]interface{}{"foo": "Unrelated val"}}},
		key:         &fieldClassifier{hostgroup: "web", service: "serviceA"},
		initialMap:  map[string]interface{}{"foo": "Bar"},
		expectedMap: map[string]interface{}{"foo": "Bar"},
	},
	// Lookup for a partially existing classifier with data and passing a non empty map
	{
		custom:      &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "web", service: "serviceB"}: map[string]interface{}{"foo": "Unrelated val"}}},
		key:         &fieldClassifier{hostgroup: "web", service: "serviceA"},
		initialMap:  map[string]interface{}{"foo": "Bar"},
		expectedMap: map[string]interface{}{"foo": "Bar"},
	},
	// Lookup for a partially existing classifier with data and passing a non empty map
	{
		custom:      &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "db", service: "serviceA"}: map[string]interface{}{"foo": "Unrelated val"}}},
		key:         &fieldClassifier{hostgroup: "web", service: "serviceA"},
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

var getUseCase = []struct {
	custom      *customFields
	hostName    string
	checkName   string
	expectedMap map[string]interface{}
}{
	{custom: &customFields{}, hostName: "web01", checkName: "serviceA", expectedMap: map[string]interface{}{}},
	// One level of hierarchy matching at a time
	{
		custom:   &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "##common##", service: "all"}: map[string]interface{}{"foo": "Bar"}}},
		hostName: "web01", checkName: "serviceA", expectedMap: map[string]interface{}{"foo": "Bar"},
	},
	{
		custom:   &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "web", service: "all"}: map[string]interface{}{"foo": "Bar"}}},
		hostName: "web01", checkName: "serviceA", expectedMap: map[string]interface{}{"foo": "Bar"},
	},
	{
		custom:   &customFields{fields: map[fieldClassifier]map[string]interface{}{fieldClassifier{hostgroup: "web", service: "serviceA"}: map[string]interface{}{"foo": "Bar", "my beautiful field": "Value 2"}}},
		hostName: "web01", checkName: "serviceA", expectedMap: map[string]interface{}{"foo": "Bar", "my beautiful field": "Value 2"},
	},
	// several levels matching at a time
	{
		custom: &customFields{fields: map[fieldClassifier]map[string]interface{}{
			fieldClassifier{hostgroup: "##common##", service: "all"}: map[string]interface{}{"foo": "Low level"},
			fieldClassifier{hostgroup: "web", service: "all"}:        map[string]interface{}{"foo": "Bar"},
			fieldClassifier{hostgroup: "web", service: "serviceB"}:   map[string]interface{}{"foo": "Not Matching"},
		}},
		hostName: "web01", checkName: "serviceA", expectedMap: map[string]interface{}{"foo": "Bar"},
	},
	{
		custom: &customFields{fields: map[fieldClassifier]map[string]interface{}{
			fieldClassifier{hostgroup: "##common##", service: "all"}: map[string]interface{}{"foo": "Low level"},
			fieldClassifier{hostgroup: "web", service: "all"}:        map[string]interface{}{"my beautiful field": "Value 2"},
			fieldClassifier{hostgroup: "web", service: "serviceA"}:   map[string]interface{}{"foo": "Bar"},
			fieldClassifier{hostgroup: "db", service: "all"}:         map[string]interface{}{"Other value": "Non matching"},
			fieldClassifier{hostgroup: "db", service: "serviceA"}:    map[string]interface{}{"foo": "db custom field"},
		}},
		hostName: "web01", checkName: "serviceA", expectedMap: map[string]interface{}{"foo": "Bar", "my beautiful field": "Value 2"},
	},
	{
		custom: &customFields{fields: map[fieldClassifier]map[string]interface{}{
			fieldClassifier{hostgroup: "##common##", service: "all"}: map[string]interface{}{"foo": "Low level"},
			fieldClassifier{hostgroup: "web", service: "all"}:        map[string]interface{}{"foo": "still low level", "my beautiful field": "Value 2"},
			fieldClassifier{hostgroup: "web", service: "serviceA"}:   map[string]interface{}{"foo": "Bar"},
			fieldClassifier{hostgroup: "db", service: "all"}:         map[string]interface{}{"Other value": "Non matching"},
			fieldClassifier{hostgroup: "db", service: "serviceA"}:    map[string]interface{}{"foo": "db custom field"},
		}},
		hostName: "web01", checkName: "serviceA", expectedMap: map[string]interface{}{"foo": "Bar", "my beautiful field": "Value 2"},
	},
}

func TestGet(t *testing.T) {
	for _, uc := range getUseCase {
		returnedMap := uc.custom.get(uc.hostName, uc.checkName)
		if !reflect.DeepEqual(returnedMap, uc.expectedMap) {
			t.Fatalf("Expecting get to return %v. Got: %v\n", uc.expectedMap, returnedMap)
		}
	}
}
