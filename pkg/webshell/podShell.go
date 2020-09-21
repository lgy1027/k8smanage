package webshell

import (
	"fmt"
	"github.com/emicklei/go-restful/log"
	"github.com/gorilla/mux"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
	"relaper.com/kubemanage/inital"
	k8s2 "relaper.com/kubemanage/k8s"
)

var (
	cmd = []string{"/bin/sh"}
)

func ServeWsTerminal(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	namespace := pathParams["namespace"]
	podName := pathParams["pod"]
	containerName := pathParams["container"]
	log.Printf("exec pod: %s, container: %s, namespace: %s\n", podName, containerName, namespace)
	writer, err := NewTerminalSession(w, r, nil)

	if err != nil {
		log.Printf("获取socket客户端失败: %v\n", err)
		return
	}
	defer func() {
		log.Print("关闭会话.\n")
		writer.Close()
	}()
	pod, err := k8s2.NewPod().Get(namespace, podName)
	if err != nil {
		log.Printf("获取pod失败: %v\n", err)
		return
	}
	podDetail := pod.(*v1.Pod)
	ok, err := ValidatePod(podDetail, containerName)
	if !ok {
		msg := fmt.Sprintf("Validate pod error! err: %v", err)
		log.Print(msg + "\n")
		writer.Write([]byte(msg))
		writer.Close()
		return
	}
	err = Exec(cmd, writer, namespace, podName, containerName)
	if err != nil {
		msg := fmt.Sprintf("Exec to pod error! err: %v", err)
		log.Print(msg + "\n")
		writer.Write([]byte(msg))
		writer.Done()
	}
	return
}

// Exec exec into a pod
func Exec(cmd []string, ptyHandler PtyHandler, namespace, podName, containerName string) error {
	defer func() {
		ptyHandler.Done()
	}()

	req := inital.GetGlobal().GetClientSet().CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	req.VersionedParams(&v1.PodExecOptions{
		Container: containerName,
		Command:   cmd,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(inital.GetGlobal().GetK8sConfig(), "POST", req.URL())
	if err != nil {
		return err
	}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             ptyHandler,
		Stdout:            ptyHandler,
		Stderr:            ptyHandler,
		TerminalSizeQueue: ptyHandler,
		Tty:               true,
	})
	return err
}
