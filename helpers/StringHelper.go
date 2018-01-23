package helpers

import (
	"sort"
	"strconv"
)

func IsContain(arrs []string,el string) bool{
	for _,v := range arrs{
		if el == v{
			return true
		}
	}
	return false
}


func RemoveDuplicatesAndEmpty(a []string) (ret []string){
	sort.Strings(a)
	a_len := len(a)
	for i:=0; i < a_len; i++{
		if (i > 0 && a[i-1] == a[i]) || len(a[i])==0{
			continue;
		}
		ret = append(ret, a[i])
	}
	return
}

func GetString(v interface{}) string{
	var result string
	switch v := v.(type) {
	case int:
		result = string(v)
	case int8:
		result = strconv.Itoa(int(v))
	case int32:
		result = strconv.Itoa(int(v))
	case int64:
		result = strconv.Itoa(int(v))
	}
	return result
}

