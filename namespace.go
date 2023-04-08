package socketio

import "sync"

// Namespace is the name space of socket.io handler.
type Namespace interface {
	// Name returns the name of namespace.
	Name() string

	// Of returns the namespace with given name.
	Of(name string) Namespace

	// On registers the function f to handle message.
	On(message string, f interface{}) error
}

type namespace struct {
	*baseHandler
	root map[string]Namespace
	lock sync.RWMutex
}

func newNamespace(broadcast BroadcastAdaptor) *namespace {
	ret := &namespace{
		baseHandler: newBaseHandler("", broadcast),
		root:        make(map[string]Namespace),
	}
	ret.root[ret.Name()] = ret
	return ret
}

func (n *namespace) Name() string {
	return n.name
}

func (n *namespace) Of(name string) Namespace {
	if name == "/" {
		name = ""
	}
	n.lock.Lock()
	defer n.lock.Unlock()
	if ret, ok := n.root[name]; ok {
		return ret
	}
	ret := &namespace{
		baseHandler: newBaseHandler(name, n.broadcast),
		root:        n.root,
		lock:        sync.RWMutex{},
	}
	n.root[name] = ret
	return ret
}
