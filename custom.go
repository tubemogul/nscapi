package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type customFields struct {
	// customFields contains a map containing the entire custom fields hierarchy using the
	// following format: fields[fieldClassifier{hostgroup: hostgroup, service: checkName}][fieldName]=value.
	// checkName can be "all" for the defaults at hostgroup and common level.
	// The content of the common.yaml will use the following format
	// fields[fieldClassifier{hostgroup: "##common##", service: "all"}][fieldName]=value
	fields map[fieldClassifier]map[string]interface{}
}

// fieldClassifier is used as the key in the customFields
type fieldClassifier struct {
	hostgroup, service string
}

// loadCustomFields loads in memory the custom fields based on the yaml
// hierarchy on disk
func (f *customFields) load(rootPath string) error {
	f.fields = make(map[fieldClassifier]map[string]interface{})
	if _, err := os.Stat(rootPath); err != nil {
		return err
	}
	// 1st common.yaml. Just ignore the errors for now
	commonFields, _ := f.processYamlFile(filepath.Join(rootPath, "common.yaml"))
	f.fields[fieldClassifier{"##common##", "all"}] = make(map[string]interface{})
	if commonFields != nil {
		f.fields[fieldClassifier{"##common##", "all"}] = commonFields
	}
	// Then service/hostgroup/*.yaml. Just ignore the errors for now
	files, _ := filepath.Glob(filepath.Join(rootPath, "service", "*", "*.yaml"))
	for _, file := range files {
		key := fieldClassifier{filepath.Base(filepath.Dir(file)), filepath.Base(file[:len(file)-5])}
		fields, _ := f.processYamlFile(file)
		f.fields[key] = make(map[string]interface{})
		if fields != nil {
			f.fields[key] = fields
		}
	}
	return nil
}

// processYamlFile tries to process a yaml file return its content hashed and/or
// an error
func (f *customFields) processYamlFile(path string) (map[string]interface{}, error) {
	returnValue := make(map[string]interface{})
	fc, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(fc, &returnValue)

	return returnValue, err
}

// lookup gets all the fields available on the given customFields key and merges it with the reference map
func (f *customFields) lookup(reference map[string]interface{}, key *fieldClassifier) {
	if m, ok := f.fields[*key]; ok {
		for field := range m {
			reference[field] = m[field]
		}
	}
}

// get returns a hash containing the custom fields specific for this hostname and checkName
func (f *customFields) get(hostname, checkName string) map[string]interface{} {
	re, _ := regexp.Compile("[^A-Za-z0-9]$")
	hostgroup := re.ReplaceAllString(hostname, "")
	resultFields := make(map[string]interface{})
	f.lookup(resultFields, &fieldClassifier{"##common##", "all"})
	f.lookup(resultFields, &fieldClassifier{hostgroup, "all"})
	f.lookup(resultFields, &fieldClassifier{hostgroup, checkName})
	return resultFields
}
