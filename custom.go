package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// customFields is a map containing the entire custom fields hierarchy using the
// following format: customFields[FieldClassifier{hostgroup: hostgroup, service: checkName}][fieldName]=value.
// checkName can be "all" for the defaults at hostgroup and common level.
// The content of the common.yaml will use the following format
// customFields[FieldClassifier{hostgroup: "##common##", service: "all"}][fieldName]=value
var customFields map[FieldClassifier]map[string]interface{}

// FieldClassifier is used as the key in the customFields
type FieldClassifier struct {
	Hostgroup, Service string
}

// loadCustomFields loads in memory the custom fields based on the yaml
// hierarchy on disk
func loadCustomFields(rootPath string) error {
	customFields = make(map[FieldClassifier]map[string]interface{})
	if _, err := os.Stat(rootPath); err != nil {
		return err
	}
	// 1st common.yaml. Just ignore the errors for now
	commonFields, _ := processYamlFile(filepath.Join(rootPath, "common.yaml"))
	customFields[FieldClassifier{"##common##", "all"}] = make(map[string]interface{})
	if commonFields != nil {
		customFields[FieldClassifier{"##common##", "all"}] = commonFields
	}
	// Then service/hostgroup/*.yaml. Just ignore the errors for now
	files, _ := filepath.Glob(filepath.Join(rootPath, "service", "*", "*.yaml"))
	for _, f := range files {
		key := FieldClassifier{filepath.Base(filepath.Dir(f)), filepath.Base(f[:len(f)-5])}
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

func mergeMapFromCustomFields(reference map[string]interface{}, key *FieldClassifier) {
	if m, ok := customFields[*key]; ok {
		for f := range m {
			reference[f] = m[f]
		}
	}
}

// getCustomFields returns a hash containing the custom fields
func getCustomFields(hostname, checkName string) map[string]interface{} {
	re, _ := regexp.Compile("[^A-Za-z0-9]$")
	hostgroup := re.ReplaceAllString(hostname, "")
	resultFields := make(map[string]interface{})
	mergeMapFromCustomFields(resultFields, &FieldClassifier{"##common##", "all"})
	mergeMapFromCustomFields(resultFields, &FieldClassifier{hostgroup, "all"})
	mergeMapFromCustomFields(resultFields, &FieldClassifier{hostgroup, checkName})
	return resultFields
}
