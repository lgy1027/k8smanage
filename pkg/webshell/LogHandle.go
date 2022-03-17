package webshell

import (
	"bufio"
	"encoding/json"
	log "github.com/cihub/seelog"
	"github.com/lgy1027/kubemanage/inital"
	"github.com/lgy1027/kubemanage/inital/client"
	k8s2 "github.com/lgy1027/kubemanage/k8s"
	"github.com/lgy1027/kubemanage/protocol"
	"github.com/lgy1027/kubemanage/utils"
	"io"
	v1 "k8s.io/api/core/v1"
	"net/http"
)

// @Tags cluster
// @Summary  获取pod日志
// @Produce  json
// @Accept  json
// @Param   namespace query string true "命名空间名 名字"
// @Param   podName query string true "Pod名字"
// @Param   container query string false "容器名"
// @Param   follow query boolean false "是否开启实时日志"
// @Success 200 {string} LogLines "level=error ts=2020-10-22T01:50:38.331Z ..."
// @Router /v1/pod/log [get]
func LogHandle(w http.ResponseWriter, req *http.Request) {
	nameSpace := req.FormValue("namespace")
	podName := req.FormValue("podName")
	container := req.FormValue("container")
	follow := req.FormValue("follow")
	pod, err := client.GetBaseClient().Pod.Get(nameSpace, podName)
	resp := protocol.NewResponse()
	resp.Data = nil
	resp.Errno = -1
	if err != nil {
		_ = log.Error("获取pod失败,err:", err.Error())
		resp.Errmsg = k8s2.ErrPodGet.Error()
		encoder := json.NewEncoder(w)
		_ = encoder.Encode(resp)
		return
	}
	if container == "" {
		container = pod.Spec.Containers[0].Name
	}
	follows, _ := utils.StringToBool(follow)
	lines := int64(200)

	opts := v1.PodLogOptions{
		Container:  container,
		Follow:     follows,
		Timestamps: true,
		TailLines:  &lines,
		Previous:   false,
	}

	reqs := inital.GetGlobal().GetClientSet().CoreV1().Pods(nameSpace).GetLogs(podName, &opts)
	r, err := reqs.Stream()
	if err != nil {
		_ = log.Error("获取容器日志错误,err:", err.Error())
		resp.Errmsg = k8s2.ErrLogGet.Error()
		encoder := json.NewEncoder(w)
		_ = encoder.Encode(resp)
		return
	}
	defer r.Close()
	var refLineNum = 0
	var offsetFrom = 2000000000
	var offsetTo = 2000000100

	refTimestamp := NewestTimestamp
	_ = DefaultSelection

	_ = &Selection{
		ReferencePoint: LogLineId{
			LogTimestamp: LogTimestamp(refTimestamp),
			LineNum:      refLineNum,
		},
		OffsetFrom:      offsetFrom,
		OffsetTo:        offsetTo,
		LogFilePosition: "end",
	}
	bufReader := bufio.NewReader(r)
	for {
		str := ""
		line, _, err := bufReader.ReadLine()
		// line = []byte(fmt.Sprintf("%s", string(line)))
		parsedLines := ToLogLines(string(line))
		for _, v := range parsedLines {
			str = str + v.Content + "\n"
		}
		//logLines, _, _, _, _ := parsedLines.SelectLogs(logSelector)
		//line = utils.ToValidUTF8(line, []byte(""))
		if err != nil {
			if err == io.EOF {
				//encoder := json.NewEncoder(w)
				//_ = encoder.Encode(parsedLines)
				_, err = w.Write([]byte(str))
			}
			return
		}
		// line = append(line, []byte("\r\n")...)
		// line = append(bytes.Trim(line, " "), []byte("\r\n")...)
		_, err = w.Write([]byte(str))
		//encoder := json.NewEncoder(w)
		//_ = encoder.Encode(parsedLines)
		//if err != nil {
		//	return
		//}
	}
}
