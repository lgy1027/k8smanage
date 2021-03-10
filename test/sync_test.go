package sync_test

import (
	"fmt"
	"github.com/spf13/viper"
	goyaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"relaper.com/kubemanage/pkg/deploy"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_Once(t *testing.T) {
	var once sync.Once
	for i := 0; i < 10; i++ {
		once.Do(func() {
			fmt.Println("Once success")
		})
		fmt.Println("this int is: ", i)
	}
}

func Test_Atomic(t *testing.T) {
	var a int32 = 100
	var b int32 = 200
	atomic.CompareAndSwapInt32(&a, a, b)
	fmt.Println(a)
}

func Test_Unmap(t *testing.T) {
	//k := "node-role.kubernetes.io/master"
	index := strings.Index("node-role.kubernetes.io/master", "node-role.kubernetes.io/")

	fmt.Println(index)
}

func Test_Upload(t *testing.T) {
	//读取yaml文件
	v := viper.New()
	//设置读取的配置文件名
	v.SetConfigName("server")
	//windows环境下为%GOPATH，linux环境下为$GOPATH
	v.AddConfigPath("E:\\workspace\\go\\src\\relaper.com\\kubemanage\\test\\")
	//设置配置文件类型
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}
	var userInfo interface{}
	if err := v.Unmarshal(&userInfo); err != nil {
		fmt.Printf("err:%s", err)
	}
	fmt.Println(userInfo)
}

func Test_FloatToString(t *testing.T) {
	num := float64(100.2568)
	string := strconv.FormatFloat(num, 'f', -1, 64)
	fmt.Println(string)
}

func Test_Parse(t *testing.T) {
	file, err := os.Open("E:\\workspace\\go\\src\\relaper.com\\kubemanage\\test\\server.yaml")
	fmt.Println(err)
	info, err := deploy.ExpandMultiYamlFileToObject(file)
	fmt.Println(err)
	fmt.Println(info[0].Object)
	data, err := goyaml.Marshal(info[0].Object)
	fmt.Println(err)
	src := "E:\\workspace\\go\\src\\relaper.com\\kubemanage\\config\\c.yaml"
	err = ioutil.WriteFile(src, data, 0777)
	fmt.Println(err)
}
