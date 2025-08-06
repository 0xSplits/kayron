package registry

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

func isImageNotFound(err error) bool {
	var inf *types.ImageNotFoundException
	var rnf *types.RepositoryNotFoundException

	return errors.As(err, &inf) || errors.As(err, &rnf)
}
