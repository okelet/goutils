package goutils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type MapHelper struct {
	Data map[string]interface{}
}

func NewEmptyMapHelper() *MapHelper {
	return &MapHelper{Data: map[string]interface{}{}}
}

func NewMapHelperFromData(data map[string]interface{}) *MapHelper {
	return &MapHelper{Data: data}
}

func NewMapHelperFromJsonFile(configPath string, failIfNotFound bool) (*MapHelper, error) {

	// TODO: Use funvtion LoadJsonFileAsMap

	exists, err := FileExists(configPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		if failIfNotFound {
			return nil, errors.Errorf("File %v doesn't exist.", configPath)
		} else {
			return NewEmptyMapHelper(), nil
		}
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

func (h *MapHelper) Clear() {
	h.Data = map[string]interface{}{}
}

func (h *MapHelper) Count() int {
	return len(h.Data)
}

func (h *MapHelper) Keys() []string {
	keys := []string{}
	for key, _ := range h.Data {
		keys = append(keys, key)
	}
	return keys
}

func (h *MapHelper) GetBoolean(key string, value bool) bool {
	if val, ok := h.Data[key]; ok {
		if w, ok := val.(bool); ok {
			return w
		}
	}
	return value
}

func (h *MapHelper) SetBoolean(key string, value bool) {
	h.Data[key] = value
}

func (h *MapHelper) GetInt(key string, value int) int {
	if val, ok := h.Data[key]; ok {
		if w, ok := val.(int); ok {
			return w
		}
		if w, ok := val.(float64); ok {
			return int(w)
		}
	}
	return value
}

func (h *MapHelper) SetInt(key string, value int) {
	h.Data[key] = value
}

func (h *MapHelper) GetInt64(key string, value int64) int64 {
	if val, ok := h.Data[key]; ok {
		if w, ok := val.(int64); ok {
			return w
		}
	}
	return value
}

func (h *MapHelper) SetInt64(key string, value int64) {
	h.Data[key] = value
}

func (h *MapHelper) GetString(key string, value string) string {
	if val, ok := h.Data[key]; ok {
		if w, ok := val.(string); ok {
			return w
		}
	}
	return value
}

func (h *MapHelper) SetString(key string, value string) {
	h.Data[key] = value
}

func (h *MapHelper) GetListOfStrings(key string, value []string) []string {
	if val, ok := h.Data[key]; ok {
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

func (h *MapHelper) SetListOfStrings(key string, value []string) {
	h.Data[key] = value
}

func (h *MapHelper) GetHelper(key string) *MapHelper {
	if val, ok := h.Data[key]; ok {
		if w, ok := val.(map[string]interface{}); ok {
			helper := &MapHelper{Data: w}
			h.Data[key] = helper
			return helper
		} else if w, ok := val.(*MapHelper); ok {
			return w
		}
	}
	helper := NewEmptyMapHelper()
	h.Data[key] = helper
	return helper
}

func (h *MapHelper) SetHelper(key string, value *MapHelper) {
	h.Data[key] = value
}

func (h *MapHelper) GetListOfHelpers(key string) []*MapHelper {
	var list []*MapHelper
	if val, ok := h.Data[key]; ok {
		if w, ok := val.([]interface{}); ok {
			for _, curval := range w {
				list = append(list, NewMapHelperFromData(curval.(map[string]interface{})))
			}
		}
	}
	return list
}

func (h *MapHelper) SetListOfHelpers(key string, data []*MapHelper) {
	l := []interface{}{}
	for _, v := range data {
		l = append(l, v.Data)
	}
	h.Data[key] = l
}

func (h *MapHelper) GetList(key string, value []interface{}) []interface{} {
	if val, ok := h.Data[key]; ok {
		if w, ok := val.([]interface{}); ok {
			return w
		}
	}
	return value
}

func (h *MapHelper) SetList(key string, value []interface{}) {
	h.Data[key] = value
}

func (h *MapHelper) Delete(key string) {
	delete(h.Data, key)
}

func (h *MapHelper) GenerateMap() map[string]interface{} {
	d := map[string]interface{}{}
	for key, val := range h.Data {
		if val == nil {
			// Nil values
			d[key] = nil
		} else {
			t := convertValueToBasic(val)
			if t != nil {
				// If nil, it is not a basic type
				d[key] = t
			}
		}
	}
	return d
}

func convertValueToBasic(val interface{}) interface{} {
	if w, ok := val.(bool); ok {
		return w
	} else if w, ok := val.(int); ok {
		return w
	} else if w, ok := val.(int8); ok {
		return w
	} else if w, ok := val.(int16); ok {
		return w
	} else if w, ok := val.(int32); ok {
		return w
	} else if w, ok := val.(int64); ok {
		return w
	} else if w, ok := val.(float32); ok {
		return w
	} else if w, ok := val.(float64); ok {
		return w
	} else if w, ok := val.(string); ok {
		return w
	} else if w, ok := val.(map[string]interface{}); ok {
		return w
	} else if w, ok := val.(*MapHelper); ok {
		return w.GenerateMap()
	} else if w, ok := val.([]string); ok {
		return w
	} else if w, ok := val.([]bool); ok {
		return w
	} else if w, ok := val.([]int); ok {
		return w
	} else if w, ok := val.([]int8); ok {
		return w
	} else if w, ok := val.([]int16); ok {
		return w
	} else if w, ok := val.([]int32); ok {
		return w
	} else if w, ok := val.([]int64); ok {
		return w
	} else if w, ok := val.([]float32); ok {
		return w
	} else if w, ok := val.([]float64); ok {
		return w
	} else if w, ok := val.([]interface{}); ok {
		a := []interface{}{}
		for _, v := range w {
			t := convertValueToBasic(v)
			if t != nil {
				a = append(a, t)
			}
		}
		return a
	}
	return nil
}

func (h *MapHelper) GenerateJson() ([]byte, error) {
	return json.MarshalIndent(h.GenerateMap(), "", "  ")
}

func (h *MapHelper) SaveToJsonFile(path string) error {
	data, err := h.GenerateJson()
	if err != nil {
		return errors.Wrapf(err, "Error marshalling map")
	}
	err = ioutil.WriteFile(path, data, 0600)
	if err != nil {
		return errors.Wrapf(err, "Error saving JSON file")
	}
	return nil
}
