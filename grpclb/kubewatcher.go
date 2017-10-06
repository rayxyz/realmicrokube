package grpclb

import (
	"sync"

	"google.golang.org/grpc/naming"
)

type watchResult struct {
	ep  *Event
	err error
}

// A Watcher provides name resolution updates from Kubernetes endpoints
// identified by name.
type watcher struct {
	target    string
	endpoints map[string]interface{}
	stopCh    chan struct{}
	result    chan watchResult
	sync.Mutex
	stopped bool
}

// Close closes the watcher, cleaning up any open connections.
func (w *watcher) Close() {
	close(w.stopCh)
}

// Next updates the endpoints for the name being watched.
func (w *watcher) Next() ([]*naming.Update, error) {
	updates := make([]*naming.Update, 0)
	updatedEndpoints := make(map[string]interface{})
	var ep Event

	select {
	case <-w.stopCh:
		w.Lock()
		if !w.stopped {
			w.stopped = true
		}
		w.Unlock()
		return updates, nil
	case r := <-w.result:
		if r.err == nil {
			ep = *r.ep
		} else {
			return updates, r.err
		}
	}

	for _, subset := range ep.Object.Subsets {
		for _, address := range subset.Addresses {
			updatedEndpoints[address.IP] = nil
		}
	}

	// Create updates to add new endpoints.
	for addr, md := range updatedEndpoints {
		if _, ok := w.endpoints[addr]; !ok {
			updates = append(updates, &naming.Update{Op: naming.Add, Addr: addr, Metadata: md})
		}
	}

	// Create updates to delete old endpoints.
	for addr := range w.endpoints {
		if _, ok := updatedEndpoints[addr]; !ok {
			updates = append(updates, &naming.Update{Op: naming.Delete, Addr: addr, Metadata: nil})
		}
	}
	w.endpoints = updatedEndpoints
	return updates, nil
}
