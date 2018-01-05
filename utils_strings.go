package goutils

func ListContainsString(list []string, element string) bool {
	for _, v := range list {
		if v == element {
			return true
		}
	}
	return false
}

func FilterStrings(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func FilterEmptyStrings(vs []string) []string {
	return FilterStrings(vs, func(s string) bool {
		return s != ""
	})
}


func AddStringToList(list []string, elem string) []string {
	found := false
	for _, v := range list {
		if v == elem {
			found = true
			break
		}
	}
	if ! found {
		list = append(list, elem)
	}
	return list
}

func RemoveStringFromList(list []string, elem string) []string {
	newList := []string{}
	for _, v := range list {
		if v != elem {
			newList = append(newList, v)
		}
	}
	return newList
}
