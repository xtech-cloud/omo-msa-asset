package tool

func CheckError(err error) bool {
	if err != nil {
		panic(err.Error())
		return true
	}
	return false
}
