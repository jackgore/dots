package yaml

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

type Yaml struct {
	data map[string]interface{}

	// For caching
	bools   map[string]bool
	strings map[string]string
	ints    map[string]int
}

// Constructs a Yaml config object from the file at path
func New(path string) (*Yaml, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Unable to find yaml file at path: %v\n%v\n", path, err)
		return nil, err
	}

	// Create our data map
	data := make(map[string]interface{})

	// Unmarshal the contents into Yaml object
	err = yaml.Unmarshal(contents, &data)
	if err != nil {
		log.Printf("Unable to umarshal yaml file at path: %v\n%v\n", path, err)
		return nil, err
	}

	y := &Yaml{data: data}
	y.strings = make(map[string]string)
	y.ints = make(map[string]int)
	y.bools = make(map[string]bool)

	return y, nil
}

// Helper function for accessing dot values in yaml file
func (y *Yaml) unwrap(key string, data map[string]interface{}) (interface{}, error) {
	// Split the key if it is dot seperated
	keys := strings.Split(key, ".")

	for i, val := range keys {
		if i == len(keys)-1 {
			return data[val], nil
		}
		d2 := make(map[string]interface{})
		// Make sure sub key exists
		if _, ok := data[val]; !ok {
			return nil, fmt.Errorf("Key %v could not be found\n", key)
		}

		// Convert sub key interface value to map of interfaces to be iterated
		m, ok := data[val].(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("Could not convert value into map for key %v", key)
		}

		for key, value := range m {
			d2[key.(string)] = value
		}
		data = d2
	}

	return nil, errors.New("something unexpected happened while unwrapping value")
}

func (y *Yaml) GetString(key string) (string, error) {
	// Before we try to unwrap in yaml file lets check our cache
	if val, ok := y.strings[key]; ok {
		return val, nil
	}

	// Otherwise its not in our cache so unwrap the value
	val, err := y.unwrap(key, y.data)
	if err != nil {
		return "", fmt.Errorf("Unable to unwrap nested value: %v\n", err)
	}

	// See if what we get is actually a string
	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("Expected string value for key %v but found %T", key, reflect.TypeOf(val))
	}

	// Insert into our cache
	y.strings[key] = s

	return s, nil
}

func (y *Yaml) GetInt(key string) (int, error) {
	// Before we try to unwrap in yaml file lets check our cache
	if val, ok := y.ints[key]; ok {
		return val, nil
	}

	// Otherwise its not in our cache so unwrap the value
	val, err := y.unwrap(key, y.data)
	if err != nil {
		return 0, fmt.Errorf("Unable to unwrap nested value: %v\n", err)
	}

	// See if what we get is actually a string
	i, ok := val.(int)
	if !ok {
		return 0, fmt.Errorf("Expected string value for key %v but found %T", key, reflect.TypeOf(val))
	}

	// Insert into our cache
	y.ints[key] = i

	return i, nil
}

func (y *Yaml) GetBool(key string) (bool, error) {
	// Before we try to unwrap in yaml file lets check our cache
	if val, ok := y.bools[key]; ok {
		return val, nil
	}

	// Otherwise its not in our cache so unwrap the value
	val, err := y.unwrap(key, y.data)
	if err != nil {
		return false, fmt.Errorf("Unable to unwrap nested value: %v\n", err)
	}

	// See if what we get is actually a string
	b, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("Expected string value for key %v but found %T", key, reflect.TypeOf(val))
	}

	// Insert into our cache
	y.bools[key] = b

	return b, nil
}
