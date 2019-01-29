package etcd

// EVERYTHING IN THIS FILE IS DEPRECATED.
// It is here only for compatibility with older code, and will be removed
// as soon as we can purge builder.

import (
	"time"

	"github.com/Masterminds/cookoo"
	"github.com/coreos/etcd/client"
)

// Getter describes the Get behavior of an Etcd client.
//
// Usually you will want to use go-etcd/etcd.Client to satisfy this.
//
// We use an interface because it is more testable.
type Getter interface {
	Get(string, bool, bool) (*client.Response, error)
}

// DirCreator describes etcd's CreateDir behavior.
//
// Usually you will want to use go-etcd/etcd.Client to satisfy this.
type DirCreator interface {
	CreateDir(string, uint64) (*client.Response, error)
}

// Setter sets a value in Etcd.
type Setter interface {
	Set(string, string, uint64) (*client.Response, error)
}

// GetterSetter performs get and set operations.
type GetterSetter interface {
	Getter
	Setter
}

// CreateSimpleClient creates a legacy simple client.
//
// DO NOT USE unless you must for backward compatibility.
//
// Params:
// 	- url (string): A server to connect to. This runs through os.ExpandEnv().
// 	- retries (int): Number of times to retry a connection to the server
// 	- retrySleep (time.Duration): How long to sleep between retries
//
// Returns:
// 	This puts a SimpleEtcdClient into context (implements Getter, Setter, etc.)
func CreateSimpleClient(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	r, e := CreateClient(c, p)
	if e != nil {
		return c, e
	}

	return &SimpleEtcdClient{
		realClient: r.(client.Client),
	}, nil
}

// NewSimpleClient Provides a simple wrapper around the old API.
//
// DO NOT USE for new code. Instead, use NewClient().
func NewSimpleClient(hosts []string) (*SimpleEtcdClient, error) {
	cfg := client.Config{
		Endpoints: hosts,
	}

	r, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	return &SimpleEtcdClient{
		realClient: r,
	}, nil
}

// SimpleEtcdClient provides an interface compatible with the old Etcd client.
type SimpleEtcdClient struct {
	realClient client.Client
}

// Get client.Response
func (c *SimpleEtcdClient) Get(key string, sort bool, rec bool) (*client.Response, error) {
	k := client.NewKeysAPI(c.realClient)
	return k.Get(dctx(), key, &client.GetOptions{Sort: sort, Recursive: rec})
}

// Set client.Response
func (c *SimpleEtcdClient) Set(key, val string, ttl uint64) (*client.Response, error) {
	k := client.NewKeysAPI(c.realClient)
	// We're banking on people not using really uge ttls. In the code base, the
	// highest is only a few hundred.
	return k.Set(dctx(), key, val, &client.SetOptions{TTL: time.Duration(ttl) * time.Second})
}

// CreateDir by name
func (c *SimpleEtcdClient) CreateDir(name string, ttl uint64) (*client.Response, error) {
	k := client.NewKeysAPI(c.realClient)
	return k.Set(dctx(), name, "", &client.SetOptions{TTL: time.Duration(ttl) * time.Second, Dir: true})
}
