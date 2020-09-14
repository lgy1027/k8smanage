package k8s

import (
	"github.com/pkg/errors"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
)

func PrintErr(err error) error {
	if err != nil {
		if k8serror.IsAlreadyExists(err) {
			return errors.New("资源已存在")
		}
		if k8serror.IsNotFound(err) {
			return errors.New("资源不存在")
		}
		statusError, isStatus := err.(*k8serror.StatusError)
		if isStatus {
			return errors.New(statusError.ErrStatus.Message)
		}
		return err
	}
	return nil
}
