package grpclb

import (
	"google.golang.org/grpc/naming"
)

type Watcher struct{}

func (w Watcher) Next() ([]*naming.Update, error) {
	return nil, nil
}

func (w Watcher) Close() {

}
