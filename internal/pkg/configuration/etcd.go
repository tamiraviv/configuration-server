package configuration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/client"
)

type Etcd struct {
	client client.KeysAPI
	dir    string
}

func NewEtcd(etcdHost string, etcdDir string) (*Etcd, error) {
	c, err := client.New(client.Config{
		Endpoints:               []string{etcdHost},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 10 * time.Second,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Error while trying to create new etcd client")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if _, err := c.GetVersion(ctx); err != nil {
		return nil, errors.Wrap(err, "Failed to ping to etcd client")
	}

	kapi := client.NewKeysAPI(c)

	return &Etcd{
		client: kapi,
		dir:    etcdDir,
	}, nil
}

func (e *Etcd) Get(rp viper.RemoteProvider) (io.Reader, error) {
	res, err := e.client.Get(context.Background(), e.dir, &client.GetOptions{
		Recursive: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get configuration from etcd")
	}

	m := e.etcdNodeTreeToMap(res.Node)

	b, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal configuration")
	}

	return bytes.NewReader(b), nil
}

func (e *Etcd) etcdNodeTreeToMap(node *client.Node) map[string]interface{} {
	if !node.Dir {
		return map[string]interface{}{node.Key: node.Value}
	}

	m := make(map[string]interface{})

	for _, n := range node.Nodes {
		if n.Dir {
			m[keyLastChild(n.Key)] = e.etcdNodeTreeToMap(n)
		} else {
			m[node.Key] = node.Value
		}
	}

	return m
}

func keyLastChild(key string) string {
	kk := strings.Split(key, "/")

	return kk[len(kk)-1]
}

func (e *Etcd) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	return nil, nil
}

func (e *Etcd) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	watcher := e.client.Watcher(e.dir, &client.WatcherOptions{
		Recursive: true,
	})

	done := make(chan bool, 1)
	resChan := make(chan *viper.RemoteResponse)

	go func(watcher client.Watcher) {
		for {
			res, err := watcher.Next(context.Background())
			if err != nil {
				remoteRes := viper.RemoteResponse{
					Error: err,
				}
				resChan <- &remoteRes
				done <- true
				return
			}

			splitKeys := strings.Split(res.Node.Key, "/")

			m := map[string]interface{}{splitKeys[len(splitKeys)-1]: res.Node.Value}
			for i := len(splitKeys) - 2; i >= 0; i-- {
				if splitKeys[i] == "" || splitKeys[i] == e.dir {
					continue
				}

				m = map[string]interface{}{splitKeys[i]: m}
			}

			b, err := json.Marshal(m)
			if err != nil {
				remoteRes := viper.RemoteResponse{
					Error: err,
				}
				resChan <- &remoteRes
				done <- true
				return
			}

			remoteRes := viper.RemoteResponse{
				Value: b,
			}

			resChan <- &remoteRes
		}
	}(watcher)

	return resChan, nil
}

func (e *Etcd) Set(key string, value string) error {
	keyInDirFormat := strings.Replace(key, ".", "/", -1)
	remoteKey := strings.Join([]string{e.dir, keyInDirFormat}, "/")
	if _, err := e.client.Set(context.Background(), remoteKey, value, nil); err != nil {
		return errors.Wrapf(err, "Failed to set key (%s) value (%s) to etcd server", remoteKey, value)
	}
	return nil
}
