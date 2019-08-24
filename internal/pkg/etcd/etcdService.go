package etcd

import (
	"context"
	"strings"
	"time"

	"configuration-server/internal/pkg/viper"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/client"
)

func NewEtcd(v *viper.ViperConf) (*Etcd, error) {
	c, err := client.New(client.Config{
		Endpoints:               []string{v.GetString(etcdHost)},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 10 * time.Second,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Error while trying to create new client")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if _, err := c.GetVersion(ctx); err != nil {
		return nil, errors.Wrap(err, "Failed to ping to etcd client")
	}

	kapi := client.NewKeysAPI(c)

	return &Etcd{
		client: kapi,
		v:      v,
		dir:    v.GetString(etcdDir),
	}, nil
}

func (e *Etcd) GetConfig() (map[string]interface{}, error) {
	res, err := e.client.Get(context.Background(), e.dir, &client.GetOptions{
		Recursive: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get configuration from etcd")
	}

	m := e.etcdNodeTreeToMap(res.Node)
	return m, nil
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

func (e *Etcd) Set(key string, value string) error {
	keyInDirFormat := strings.Replace(key, ".", "/", -1)
	remoteKey := strings.Join([]string{e.dir, keyInDirFormat}, "/")
	if _, err := e.client.Set(context.Background(), remoteKey, value, nil); err != nil {
		return errors.Wrapf(err, "Failed to set key (%s) value (%s) to etcd server", remoteKey, value)
	}
	return nil
}

func (e *Etcd) WatchConfig() (<-chan *WatcherResponse, chan error) {
	watcher := e.client.Watcher(e.dir, &client.WatcherOptions{
		Recursive: true,
	})

	done := make(chan error)
	resChan := make(chan *WatcherResponse)

	go func(watcher client.Watcher) {
		for {
			res, err := watcher.Next(context.Background())
			if err != nil {
				done <- err
				return
			}

			splitKeys := strings.Split(res.Node.Key, "/")

			formatedKey := strings.Join(splitKeys[2:], ".")

			resChan <- &WatcherResponse{
				Key:   formatedKey,
				Value: res.Node.Value,
			}
		}
	}(watcher)

	return resChan, done
}
