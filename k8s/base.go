package k8s

type Base interface {
	Delete(namespace, name string) error
	Get(namespace, name string) (interface{}, error)
	List(namespace string) (interface{}, error)
}
