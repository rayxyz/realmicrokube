package grpclb

import (
	"google.golang.org/grpc/naming"
)

type KubeResolver struct{}

func NewResolver() *KubeResolver {
	return &KubeResolver{}
}

func (r *KubeResolver) Resolve(target string) (naming.Watcher, error) {
	return Watcher{}, nil
}
