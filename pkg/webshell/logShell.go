package webshell

import (
	"bufio"
	"fmt"
	"github.com/emicklei/go-restful/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io"
	v1 "k8s.io/api/core/v1"
	"net/http"
	"relaper.com/kubemanage/inital"
	"relaper.com/kubemanage/inital/client"
	"relaper.com/kubemanage/utils"
)

type Logger interface {
	io.WriteCloser
}

func (l *WsLogger) Write(p []byte) (n int, err error) {
	if err := l.wsConn.WriteMessage(websocket.TextMessage, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (l *WsLogger) Close() error {
	return l.wsConn.Close()
}

type WsLogger struct {
	wsConn *websocket.Conn
}

func newWsLogger(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*WsLogger, error) {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	session := &WsLogger{
		wsConn: conn,
	}
	return session, nil
}

// @Tags xshell
// @Summary  日志
// @Produce  json
// @Accept  json
// @Param namespace path string true "namespace"
// @Param pod path string true "pod"
// @Param container path string true "container"
// @Success 200 {string} json ""
// @Router /ws/{namespace}/{pod}/{container}/log [get]
func ServeWsLogs(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	namespace := pathParams["namespace"]
	podName := pathParams["pod"]
	containerName := pathParams["container"]
	tailLine, _ := utils.StringToInt64(r.URL.Query().Get("tail"))
	follow, _ := utils.StringToBool(r.URL.Query().Get("follow"))
	log.Printf("log pod: %s, container: %s, namespace: %s, tailLine: %d, follow: %v\n", podName, containerName, namespace, tailLine, follow)

	writer, err := newWsLogger(w, r, nil)

	if err != nil {
		log.Printf("获取输出流失败: %v\n", err)
		return
	}
	defer func() {
		log.Print("关闭会话.\n")
		_ = writer.Close()
	}()

	podDetail, err := client.GetBaseClient().Pod.Get(namespace, podName)
	if err != nil {
		log.Printf("获取pod失败: %v\n", err)
		return
	}
	ok, err := ValidatePod(podDetail, containerName)
	if !ok {
		msg := fmt.Sprintf("Validate pod error! err: %v", err)
		log.Print(msg + "\n")
		_, _ = writer.Write([]byte(msg))
		_ = writer.Close()
		return
	}
	opt := v1.PodLogOptions{
		Container: containerName,
		Follow:    follow,
		TailLines: &tailLine,
	}
	err = LogStreamLine(podName, namespace, &opt, writer)
	if err != nil {
		msg := fmt.Sprintf("log err: %v", err)
		log.Print(msg + "\n")
		writer.Write([]byte(msg))
		writer.Close()
	}
	return
}

func LogStreamLine(name, namespace string, opts *v1.PodLogOptions, writer io.Writer) error {
	req := inital.GetGlobal().GetClientSet().CoreV1().Pods(namespace).GetLogs(name, opts)
	r, err := req.Stream()
	if err != nil {
		return err
	}
	defer r.Close()
	bufReader := bufio.NewReaderSize(r, 256)
	// bufReader := bufio.NewReader(r)
	for {
		line, _, err := bufReader.ReadLine()
		// line = []byte(fmt.Sprintf("%s", string(line)))
		line = utils.ToValidUTF8(line, []byte(""))
		if err != nil {
			if err == io.EOF {
				_, err = writer.Write(line)
			}
			return err
		}
		// line = append(line, []byte("\r\n")...)
		// line = append(bytes.Trim(line, " "), []byte("\r\n")...)
		_, err = writer.Write(line)
		if err != nil {
			return err
		}
	}
}
