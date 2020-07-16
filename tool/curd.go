package tool

import "fmt"

func Remove(slice []interface{}, i int) []interface{} {
	return append(slice[:i], slice[i+1:]...)
}

func Add(slice []interface{}, value interface{}) []interface{} {
	return append(slice, value)
}

func Insert(slice *[]interface{}, index int, value interface{}) {
	rear := append([]interface{}{}, (*slice)[index:]...)
	*slice = append(append((*slice)[:index], value), rear...)
}

func HasItem(array []string, value string) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == value {
			return true
		}
	}
	return false
}

func HasItemByUint(array []uint, value uint) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == value {
			return true
		}
	}
	return false
}

func Update(slice []interface{}, index int, value interface{}) {
	slice[index] = value
}

func Find(slice []interface{}, index int) interface{} {
	return slice[index]
}

func Clear(slice *[]interface{}) {
	//    *slice = nil
	*slice = append([]interface{}{})
}

func List(slice []interface{}) {
	for _, v := range slice {
		fmt.Printf("%d ", v)
	}
}
