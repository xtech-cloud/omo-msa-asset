package tool

import "math"

func Page(items []interface{}, current int, size int) (int, []interface{}) {
	var list = make([]interface{}, 0, size)
	total := len(items)
	max := int(math.Ceil(float64(total) / float64(size)))
	if current > max {
		return total, nil
	}
	start := (current - 1) * size
	length := size
	for i := 0; i < length; i++ {
		index := start + i
		if index < total {
			list = append(list, items[index])
		}
	}
	return total, list
}
