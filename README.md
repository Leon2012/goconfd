#golang 配置中心

```
    ETCD启动
    etcd -name etc0 \
-initial-advertise-peer-urls http://192.168.174.114:2380 \
-listen-peer-urls http://192.168.174.114:2380 \
-listen-client-urls http://192.168.174.114:2379,http://127.0.0.1:2379 \
-advertise-client-urls http://192.168.174.114:2379 

```