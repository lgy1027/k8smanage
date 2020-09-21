package file

import (
	"encoding/json"
	goyaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"net/http"
	k8s2 "relaper.com/kubemanage/k8s"
	"relaper.com/kubemanage/protocol"
	"relaper.com/kubemanage/utils"
)

var (
	deployment  *k8s2.Deployment
	statefulSet *k8s2.Sf
	service     *k8s2.Sv
)

func init() {
	deployment = k8s2.NewDeploy()
	statefulSet = k8s2.NewStateFulSet()
	service = k8s2.NewSv()
}

// @Summary  下载yaml文件
// @Accept  octet-stream
// @Param   namespace query string true "命名空间名 名字"
// @Param   kind query string true "资源类型"
// @Param   name query string true "资源名"
// @success 200 {object} string "success"
// @Router /cluster/v1/uploadYaml [get]
func HandleDownload(w http.ResponseWriter, req *http.Request) {
	kind := req.FormValue("kind")
	namespace := req.FormValue("namespace")
	name := req.FormValue("name")
	resp := protocol.NewResponse()
	if kind == "" || namespace == "" || name == "" {
		resp.Data = nil
		resp.Errno = -1
		resp.Errmsg = "版本号不能为空"
		encoder := json.NewEncoder(w)
		encoder.Encode(resp)
		return
	}
	var (
		err error
		obj *unstructured.Unstructured
	)
	switch kind {
	case utils.DEPLOY_DEPLOYMENT:
		obj, err = deployment.DynamicGet(namespace, name)
	case utils.DEPLOY_STATEFULSET:
		obj, err = statefulSet.DynamicGet(namespace, name)
	case utils.DEPLOY_Service:
		obj, err = service.DynamicGet(namespace, name)
	}
	if err != nil {
		return
	}

	delete(obj.Object, "status")
	delete(obj.Object["metadata"].(map[string]interface{}), "creationTimestamp")
	delete(obj.Object["metadata"].(map[string]interface{}), "generation")
	delete(obj.Object["metadata"].(map[string]interface{}), "resourceVersion")
	delete(obj.Object["metadata"].(map[string]interface{}), "selfLink")
	delete(obj.Object["metadata"].(map[string]interface{}), "uid")
	delete(obj.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{}), "terminationGracePeriodSeconds")
	delete(obj.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["metadata"].(map[string]interface{}), "creationTimestamp")
	cont := obj.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})
	delete(cont, "terminationMessagePath")
	delete(cont, "terminationMessagePolicy")

	data, err := goyaml.Marshal(obj.Object)
	if err != nil {
		return
	}
	//将文件写至responseBody
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("content-disposition", "attachment; filename=default.yaml")
	_, _ = w.Write(data)
}
