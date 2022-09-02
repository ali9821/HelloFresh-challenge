package tools

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func ExtractNumbersFromString(str string) []string {
	re := regexp.MustCompile("[0-9]+")
	allStrings := re.FindAllString(str, -1)
	return allStrings
}

func ConvertStringToInt(str string) int {
	i, _ := strconv.Atoi(str)
	return i
}

func StringContains(str string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			return true
		}
	}
	return false
}

func FindMaxInSyncMap(syncMap *sync.Map) (key string, value int) {
	syncMap.Range(func(k, v interface{}) bool {
		if v.(int) > value {
			key = k.(string)
			value = v.(int)
		}
		return true
	})
	return
}
