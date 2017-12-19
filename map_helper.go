package goutils

import (
	"os"
	"github.com/pkg/errors"
	"encoding/json"
)

type MapHelper struct {
	Data map[string]interface{}
}

func NewMapHelper() *MapHelper {
	return &MapHelper{Data: map[string]interface{}{}}
}

func NewMapHelperFromData(data map[string]interface{}) *MapHelper {
	return &MapHelper{Data: data}
}

func NewMapHelperFromFile(configPath string) (*MapHelper, error) {

	stat, err := os.Stat(configPath)
	if err != nil && ! os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "Error checking if file %v exists.", configPath)
	}
	if ! stat.Mode().IsRegular() {
		return nil, errors.Errorf("File %v, exists, but is not a file.", configPath)
	}

	reader, err := os.Open(configPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Error opening file %v.", configPath)
	}
	decoder := json.NewDecoder(reader)

	var x map[string]interface{}
	err = decoder.Decode(&x)
	if err != nil {
		return nil, errors.Wrapf(err, "Error loading JSON from file %v.", configPath)
	}

	return &MapHelper{Data: x}, nil

}

func (json *MapHelper) GetBoolean(key string, value bool) bool {
	if val, ok := json.Data[key]; ok {
		if w, ok := val.(bool); ok {
			return w
		}
	}
	return value
}

func (json *MapHelper) SetBoolean(key string, value bool) {
	json.Data[key] = value
}

func (json *MapHelper) GetInt(key string, value int) int {
	if val, ok := json.Data[key]; ok {
		if w, ok := val.(int); ok {
			return w
		}
	}
	return value
}

func (json *MapHelper) SetInt(key string, value int) {
	json.Data[key] = value
}

func (json *MapHelper) GetInt64(key string, value int64) int64 {
	if val, ok := json.Data[key]; ok {
		if w, ok := val.(int64); ok {
			return w
		}
	}
	return value
}

func (json *MapHelper) SetInt64(key string, value int64) {
	json.Data[key] = value
}

func (json *MapHelper) GetString(key string, value string) string {
	if val, ok := json.Data[key]; ok {
		if w, ok := val.(string); ok {
			return w
		}
	}
	return value
}

func (json *MapHelper) SetString(key string, value string) {
	json.Data[key] = value
}

func (json *MapHelper) GetStringArray(key string, value []string) []string {
	if val, ok := json.Data[key]; ok {
		if w, ok := val.([]interface{}); ok {
			strList := []string{}
			for _, val2 := range w {
				if w2, ok := val2.(string); ok {
					strList = append(strList, w2)
				}
			}
			return strList
		}
	}
	return value
}

func (json *MapHelper) SetStringArray(key string, value []string) {
	json.Data[key] = value
}

func (json *MapHelper) GetMap(key string) *MapHelper {
	if val, ok := json.Data[key]; ok {
		if w, ok := val.(map[string]interface{}); ok {
			return &MapHelper{Data: w}
		}
	}
	return &MapHelper{Data: make(map[string]interface{})}
}

func (json *MapHelper) GetArray(key string) []*MapHelper {
	var list []*MapHelper
	if val, ok := json.Data[key]; ok {
		if w, ok := val.([]interface{}); ok {
			for _, curval := range w {
				list = append(list, NewMapHelperFromData(curval.(map[string]interface{})))
			}
		}
	}
	return list
}
