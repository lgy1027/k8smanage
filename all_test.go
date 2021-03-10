package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_Metric(t *testing.T) {
	resp, err := http.Get("http://192.168.5.17:9100/metrics")
	fmt.Println(err)
	if err != nil {
		return
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(buf))
}
