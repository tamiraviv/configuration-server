# Configuration Server Using Viper & ETCD

Currently, viper integration with ETCD isn't fully implemented:
for example if you want to run function on every remote config changes you can't do it (on local changes viper offer the function "OnConfigChange" but it doesn't work on the remote config)

Therefore I created this repository to give an example of how to use viper with ETCD including hooks.

Enjoy :)