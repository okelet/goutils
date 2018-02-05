package goutils

import "sort"
import "unicode"

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
	if !found {
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

func StringListsAreEqual(a1 []string, a2 []string) bool {
	sort.Strings(a1)
	sort.Strings(a2)
	if len(a1) == len(a2) {
		for i, v := range a1 {
			if v != a2[i] {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

func EnsureFirstSentenceLetterLowercase(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func EnsureFirstSentenceLetterUppercase(s string) string {
	a := []rune(s)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}
