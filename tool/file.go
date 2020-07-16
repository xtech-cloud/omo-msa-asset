package tool

import (
	"io/ioutil"
	"net"
	"os"
)

func PathIsExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func ByteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

func internalIP() string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops:" + err.Error())
		os.Exit(1)
	}
	var tmp = ""
	for _, a := range address {
		if ip, ok := a.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				os.Stdout.WriteString(ip.IP.String() + "\n")
				tmp += ip.IP.String() + ":"
			}
		}
	}
	os.Exit(0)
	return tmp
}

func ReadFile(path string) ([]byte, error) {
	f, err1 := os.OpenFile(path, os.O_RDWR, 0666)
	if err1 != nil {
		return nil, err1
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
