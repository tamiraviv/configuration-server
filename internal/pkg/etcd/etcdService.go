package etcd

import (
	"configuration-server/internal/pkg/viper"
	"context"
	"time"

	"github.com/pkg/errors"
	"go.etcd.io/etcd/client"
)

func NewEtcd(v *viper.ViperConf) (*Etcd, error){
	c, err := client.New(client.Config{
		Endpoints:               []string{v.GetString(etcdHost)},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 10*time.Second,
	})

	if err != nil {
		return nil, errors.Wrap(err,"Error while trying to create new client")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if _, err := c.GetVersion(ctx); err != nil {
		return nil, errors.Wrap(err,"Failed to ping to etcd client")
	}

	kapi := client.NewKeysAPI(c)

	return &Etcd {
		kapi,
	}, nil
}

func (e *Etcd) Get(key string) (string, error){
	res, err := e.client.Get(context.Background(), key, nil)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to get key (%s) from etcd", key)
	}

	return res.Node.Value, nil
}

