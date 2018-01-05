package goutils

func ListContainsInt(list []int, element int) bool {
	for _, v := range list {
		if v == element {
			return true
		}
	}
	return false
}


