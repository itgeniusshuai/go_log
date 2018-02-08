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

func DistinctByMap(arr []string ) []string{
	tmp_arr := make([]string, 0)
	if arr == nil || len(arr) == 0{
		return tmp_arr
	}
	tm := make(map[string]int,0)
	for _,e := range arr{
		tm[e] = 0
	}
	for k,_ := range tm{
		tmp_arr = append(tmp_arr, k)
	}
	return tmp_arr
}

func DistinctByCircle(arr []string ) []string{
	if arr == nil{
		return nil
	}
	size := len(arr)
	tmpArr := make([]string,0)
	if size == 0 {
		return tmpArr
	}
	for _,e := range arr{
		if !IsContain(tmpArr,e){
			tmpArr = append(tmpArr, e)
		}
	}
	return tmpArr
}

