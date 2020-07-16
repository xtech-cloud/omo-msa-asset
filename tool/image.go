package tool

import (
	"image"
	"io/ioutil"
	"net/http"
	"os"
)

func ImageLoad(url string, path string) error {
	//fmt.Println(url+"---"+path)
	response, error1 := http.Get(url)
	if error1 != nil {
		return error1
	}

	buffer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		CheckError(err)
		return err
	}
	err2 := ioutil.WriteFile(path, buffer, 0666)
	if err2 != nil {
		return err2
	}
	return nil
}

func ImageOpen(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}
