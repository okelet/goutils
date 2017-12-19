package goutils

type JsonHelper struct {
	data map[string]interface{}
}

func NewJsonHelper(data map[string]interface{}) *JsonHelper {
	return &JsonHelper{data: data}
}

func (json *JsonHelper) GetBoolean(key string, value bool) bool {
	if val, ok := json.data[key]; ok {
		if w, ok := val.(bool); ok {
			return w
		}
	}
	return value
}

func (json *JsonHelper) SetBoolean(key string, value bool) {
	json.data[key] = value
}

func (json *JsonHelper) GetInt(key string, value int) int {
	if val, ok := json.data[key]; ok {
		if w, ok := val.(int); ok {
			return w
		}
	}
	return value
}

func (json *JsonHelper) SetInt(key string, value int) {
	json.data[key] = value
}

func (json *JsonHelper) GetInt64(key string, value int64) int64 {
	if val, ok := json.data[key]; ok {
		if w, ok := val.(int64); ok {
			return w
		}
	}
	return value
}

func (json *JsonHelper) SetInt64(key string, value int64) {
	json.data[key] = value
}

func (json *JsonHelper) GetString(key string, value string) string {
	if val, ok := json.data[key]; ok {
		if w, ok := val.(string); ok {
			return w
		}
	}
	return value
}

func (json *JsonHelper) SetString(key string, value string) {
	json.data[key] = value
}

func (json *JsonHelper) GetStringArray(key string, value []string) []string {
	if val, ok := json.data[key]; ok {
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

func (json *JsonHelper) SetStringArray(key string, value []string) {
	json.data[key] = value
}

func (json *JsonHelper) GetMap(key string) *JsonHelper {
	if val, ok := json.data[key]; ok {
		if w, ok := val.(map[string]interface{}); ok {
			return &JsonHelper{data: w}
		}
	}
	return &JsonHelper{data: make(map[string]interface{})}
}

func (json *JsonHelper) GetArray(key string) []*JsonHelper {
	var list []*JsonHelper
	if val, ok := json.data[key]; ok {
		if w, ok := val.([]interface{}); ok {
			for _, curval := range w {
				list = append(list, NewJsonHelper(curval.(map[string]interface{})))
			}
		}
	}
	return list
}
