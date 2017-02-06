package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// customFields is a map containing the entire custom fields hierarchy using the
// following format: customFields[fieldClassifier{hostgroup: hostgroup, service: checkName}][fieldName]=value.
// checkName can be "all" for the defaults at hostgroup and common level.
// The content of the common.yaml will use the following format
// customFields[fieldClassifier{hostgroup: "##common##", service: "all"}][fieldName]=value
var customFields map[fieldClassifier]map[string]interface{}

// fieldClassifier is used as the key in the customFields
type fieldClassifier struct {
	hostgroup, service string
}

// loadCustomFields loads in memory the custom fields based on the yaml
// hierarchy on disk
func loadCustomFields(rootPath string) error {
	customFields = make(map[fieldClassifier]map[string]interface{})
	if _, err := os.Stat(rootPath); err != nil {
		return err
	}
	// 1st common.yaml. Just ignore the errors for now
	commonFields, _ := processYamlFile(filepath.Join(rootPath, "common.yaml"))
	customFields[fieldClassifier{"##common##", "all"}] = make(map[string]interface{})
	if commonFields != nil {
		customFields[fieldClassifier{"##common##", "all"}] = commonFields
	}
	// Then service/hostgroup/*.yaml. Just ignore the errors for now
	files, _ := filepath.Glob(filepath.Join(rootPath, "service", "*", "*.yaml"))
	for _, f := range files {
		key := fieldClassifier{filepath.Base(filepath.Dir(f)), filepath.Base(f[:len(f)-5])}
		fields, _ := processYamlFile(f)
		customFields[key] = make(map[string]interface{})
		if fields != nil {
			customFields[key] = fields
		}
	}
	return nil
}

// processYamlFile tries to process a yaml file return its content hashed and/or
// an error
func processYamlFile(path string) (map[string]interface{}, error) {
	returnValue := make(map[string]interface{})
	fc, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(fc, &returnValue)

	return returnValue, err
}

// mergeMapFromCustomFields merges the data from a given customFields key (aka
// yaml file loaded from the disk) into the map given as reference
func mergeMapFromCustomFields(reference map[string]interface{}, key *fieldClassifier) {
	if m, ok := customFields[*key]; ok {
		for f := range m {
			reference[f] = m[f]
		}
	}
}

// getCustomFields returns a hash containing the custom fields specific for this
// hostname and checkName
func getCustomFields(hostname, checkName string) map[string]interface{} {
	re, _ := regexp.Compile("[^A-Za-z0-9]$")
	hostgroup := re.ReplaceAllString(hostname, "")
	resultFields := make(map[string]interface{})
	mergeMapFromCustomFields(resultFields, &fieldClassifier{"##common##", "all"})
	mergeMapFromCustomFields(resultFields, &fieldClassifier{hostgroup, "all"})
	mergeMapFromCustomFields(resultFields, &fieldClassifier{hostgroup, checkName})
	return resultFields
}
