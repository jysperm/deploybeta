package utils

func StringInSlice(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}

func AddStringToUniqueSlice(list []string, str string) []string {
	if StringInSlice(list, str) {
		return list
	} else {
		return append(list, str)
	}
}

func PullStringFromSlice(list []string, str string) []string {
	result := make([]string, 0)

	for _, item := range list {
		if item != str {
			result = append(result, item)
		}
	}

	return result
}
