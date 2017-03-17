package main

import (
	"reflect"
	"testing"
)

var lookupUseCases = []struct {
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
	// Lookup for a non existing classifier with data and passing a non empty map
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
	for _, uc := range lookupUseCases {
		uc.custom.lookup(uc.initialMap, uc.key)
		if !reflect.DeepEqual(uc.initialMap, uc.expectedMap) {
			t.Fatalf("Expecting map to be %v after lookup. Got: %v\n", uc.expectedMap, uc.initialMap)
		}
	}
}

var getUseCases = []struct {
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
	for _, uc := range getUseCases {
		returnedMap := uc.custom.get(uc.hostName, uc.checkName)
		if !reflect.DeepEqual(returnedMap, uc.expectedMap) {
			t.Fatalf("Expecting get to return %v. Got: %v\n", uc.expectedMap, returnedMap)
		}
	}
}

var processYamlFileUseCases = []struct {
	filePath    string
	expectedMap map[string]interface{}
	expectedErr string
}{
	{filePath: "testData/customFields/common.yaml", expectedMap: map[string]interface{}{"paging": false, "team": []interface{}{"ops"}}, expectedErr: ""},
	{filePath: "testData/customFields/service/db/all.yaml", expectedMap: map[string]interface{}{"paging": true, "team": []interface{}{"dba", "ops"}, "runbook": "https://wiki.example.org/teams/dba/runbooks.html"}, expectedErr: ""},
	{filePath: "testData/customFields/service/db/nonExistingFile.yaml", expectedMap: map[string]interface{}{}, expectedErr: "open testData/customFields/service/db/nonExistingFile.yaml: no such file or directory"},
}

func TestProcessYamlFile(t *testing.T) {
	for _, uc := range processYamlFileUseCases {
		custom := &customFields{}
		returnedMap, err := custom.processYamlFile(uc.filePath)
		if !reflect.DeepEqual(returnedMap, uc.expectedMap) {
			t.Fatalf("Expecting processYamlFile of %s to return %v. Got: %v\n", uc.filePath, uc.expectedMap, returnedMap)
		}

		if err != nil {
			if err.Error() != uc.expectedErr {
				t.Fatalf("Expecting processYamlFile of %s to return the error \"%s\". Got: \"%s\"\n", uc.filePath, uc.expectedErr, err)
			}
		} else {
			if uc.expectedErr != "" {
				t.Fatalf("Expecting processYamlFile of %s to return the error \"%s\". Got: \"%s\"\n", uc.filePath, uc.expectedErr, err)
			}
		}
	}
}

var loadUseCases = []struct {
	filePath    string
	expectedMap *customFields
	expectedErr string
}{
	{filePath: "testData/customFields", expectedMap: &customFields{fields: map[fieldClassifier]map[string]interface{}{
		fieldClassifier{hostgroup: "##common##", service: "all"}:         map[string]interface{}{"paging": false, "team": []interface{}{"ops"}},
		fieldClassifier{hostgroup: "web", service: "all"}:                map[string]interface{}{"team": []interface{}{"webdev"}, "runbook": "https://wiki.example.org/teams/web/runbooks.html"},
		fieldClassifier{hostgroup: "web", service: "apache"}:             map[string]interface{}{"paging": true, "runbook": "https://wiki.example.org/teams/cross/runbooks/apache.html"},
		fieldClassifier{hostgroup: "web", service: "nginx_port"}:         map[string]interface{}{"paging": true},
		fieldClassifier{hostgroup: "db", service: "all"}:                 map[string]interface{}{"paging": true, "team": []interface{}{"dba", "ops"}, "runbook": "https://wiki.example.org/teams/dba/runbooks.html"},
		fieldClassifier{hostgroup: "app", service: "swiftWF_error_rate"}: map[string]interface{}{"paging": true, "team": []interface{}{"stats"}, "runbook": "https://wiki.example.org/teams/cross/runbooks/swif_workflows.html", "alertGroup": "AWS_Swift"},
	}}, expectedErr: ""},
	{filePath: "nonExistentPath", expectedMap: &customFields{fields: map[fieldClassifier]map[string]interface{}{}}, expectedErr: "stat nonExistentPath: no such file or directory"},
}

func TestLoad(t *testing.T) {
	for _, uc := range loadUseCases {
		custom := &customFields{}
		err := custom.load(uc.filePath)
		if !reflect.DeepEqual(uc.expectedMap, custom) {
			t.Fatalf("Expecting load to generate the map %v. Got: %v\n", uc.expectedMap, custom)
		}
		if err != nil {
			if err.Error() != uc.expectedErr {
				t.Fatalf("Expecting load of %s to return the error \"%s\". Got: \"%s\"\n", uc.filePath, uc.expectedErr, err)
			}
		} else {
			if uc.expectedErr != "" {
				t.Fatalf("Expecting load of %s to return the error \"%s\". Got: \"%s\"\n", uc.filePath, uc.expectedErr, err)
			}
		}
	}
}
