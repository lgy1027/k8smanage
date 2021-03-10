package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"io"
	"log"
	"net/http"
	"relaper.com/kubemanage/example/go-restful/api"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello

type Base struct {
	Optional bool     `json:"optional,omitempty"`
	Value    string   `json:"value,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Timestamp string `json:"timestamp,omitempty"`
	MType     string `json:"type,omitempty"`
}

type PropertyVisitorsBase struct {
	DataMax          Base `json:"data_max,omitempty"`
	DataMin          Base `json:"data_min,omitempty"`
	ExpectedDatatype Base `json:"expected_datatype,omitempty"`
}

type ModBusPropertyVisitors struct {
	AccessMode       Base `json:"access_mode,omitempty"`
	IsRegisterSwap   Base `json:"is_registerswap,omitempty"`
	IsSwap           Base `json:"is_swap,omitempty"`
	OriginalDatatype Base `json:"original_datatype,omitempty"`
	RegisterIndex    Base `json:"register_index,omitempty"`
	RegisterNum      Base `json:"register_num,omitempty"`
	RegisterType     Base `json:"register_type,omitempty"`
	SampleInterval   Base `json:"sample_interval,omitempty"`
	ScaleIndex       Base `json:"scale_index,omitempty"`
}

type OPCUAPropertyVisitors struct {
	BrowSeName Base `json:"browse_name,omitempty"`
	NodeId     Base `json:"node_id,omitempty"`
}

type PropertyVisitors struct {
	ModBusPropertyVisitors
	OPCUAPropertyVisitors
	PropertyVisitorsBase
}

type Expected struct {
	Value string `json:"value,omitempty"`
}

type Twin struct {
	Expected Expected `json:"expected,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
	Optional bool     `json:"optional,omitempty"`
}

type AccessConfigBase struct {
	ProtocolType Base `json:"protocol_type,omitempty"` // 传输模式
	ProtocolName Base `json:"protocol_name,omitempty"` // 访问名称
}

type RTUAccessConfig struct {
	SerialPort Base `json:"serial_port,omitempty"` // 串口
	BaudRate   Base `json:"baud_rate,omitempty"`   // 波特率
	DataBits   Base `json:"data_bits,omitempty"`   // 数据位
	StopBits   Base `json:"stop_bits,omitempty"`   // 停止位
	ParityBits Base `json:"parity_bits,omitempty"` // 校验位
}

type TCPAccessConfig struct {
	Ip   Base `json:"ip,omitempty"`
	Port Base `json:"port,omitempty"`
}

type ModBusAccessConfig struct {
	SlaveId Base `json:"slave_id,omitempty"`
	RTUAccessConfig
	TCPAccessConfig
}

type OPCUAAccessConfig struct {
	AuthType    Base  `json:"auth_type,omitempty"`
	Certificate Base  `json:"certificate,omitempty"`
	PassWord    Base  `json:"password,omitempty"`
	PrivateKey  Base  `json:"private_key,omitempty"`
	SecMode     Base  `json:"sec_mode,omitempty"`
	SecPolicy   Base  `json:"sec_policy,omitempty"`
	Timeout     int64 `json:"timeout,omitempty"`
	Url         Base  `json:"url,omitempty"`
	Username    Base  `json:"username,omitempty"`
}

type AccessConfig struct {
	AccessConfigBase
	ModBusAccessConfig
	OPCUAAccessConfig
}

type DeviceInfo struct {
	DeviceModelName    string                      `json:"deviceModel,omitempty"`
	DeviceInstanceName string                      `json:"deviceName,omitempty"`
	Description        string                      `json:"description,omitempty"`
	AccessProtocol     string                      `json:"access_protocol,omitempty"`
	Attributes         map[string]Base             `json:"attributes,omitempty"`
	PropertyVisitors   map[string]PropertyVisitors `json:"property_visitors,omitempty"`
	Tags               map[string]string           `json:"tags,omitempty"`
	Twin               map[string]Twin             `json:"twin,omitempty"`
	AccessConfig       AccessConfig                `json:"access_config,omitempty"`
	NodeName           string                      `json:"nodeName,omitempty"`
	ConnectionType     string                      `json:"connection_type,omitempty"`
}

func main() {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/devops/{devops}/credentials").
		To(hello).
		Param(ws.PathParameter("devops", "devops name")).
		Doc("list the credentials of the specified devops for the current user").
		Returns(http.StatusOK, api.StatusOK, api.ListResult{Items: []interface{}{}}))
	ws.Route(ws.POST("/post").
		To(post).
		Reads(DeviceInfo{}).
		Doc("list the credentials of the specified devops for the current user").
		Returns(http.StatusOK, api.StatusOK, api.ListResult{Items: []interface{}{}}))

	ws.Route(ws.POST("/test").
		To(test).
		Reads(map[string]interface{}{}).
		Doc("list the credentials of the specified devops for the current user").
		Returns(http.StatusOK, api.StatusOK, api.ListResult{Items: []interface{}{}}))
	restful.Add(ws)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func test(req *restful.Request, resp *restful.Response) {
	params := make(map[string]interface{})

	_ = req.ReadEntity(&params)
	resp.WriteAsJson(params)
}

func hello(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, req.PathParameter("devops"))
}

func post(req *restful.Request, resp *restful.Response) {
	var tokenReview DeviceInfo

	err := req.ReadEntity(&tokenReview)
	fmt.Println(err)
	resp.WriteAsJson(tokenReview)
}
