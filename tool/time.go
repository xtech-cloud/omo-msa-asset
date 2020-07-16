package tool

import (
	"math/rand"
	"time"
)

func ParseDate(date string) (time.Time, error) {
	if date == "" {
		return time.Now(), nil
	}
	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		return dt, err
	}
	return dt, nil
}

func ParseTime(date string) (time.Time, error) {
	if date == "" {
		return time.Now(), nil
	}
	dt, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		return dt, err
	}
	return dt, nil
}

func unixToTime(unix int64) time.Time {
	return time.Unix(unix, 0) //获取时间对象
}

func RandomRange(num int) []uint {
	random := rand.New(rand.NewSource(time.Now().Unix()))
	array := make([]uint, 0, num)
	for i := 0; i < num; i++ {
		val := random.Intn(num)
		array = append(array, uint(val))
	}
	return array
}

func RandomInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func SwitchWeekday(week string) uint8 {
	weeks := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	for i := 0; i < len(weeks); i += 1 {
		if weeks[i] == week {
			return uint8(i + 1)
		}
	}
	return 0
}
