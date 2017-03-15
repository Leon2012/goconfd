#golang 配置中心

```
    1，ETCD启动
    etcd -name etc0 \
    -initial-advertise-peer-urls http://192.168.174.114:2380 \
    -listen-peer-urls http://192.168.174.114:2380 \
    -listen-client-urls http://192.168.174.114:2379,http://127.0.0.1:2379 \
    -advertise-client-urls http://192.168.174.114:2379 

    2，编译
    make

    3，运行
    修改agent.toml，指定hosts为etcd监听的地址，修改需监控的key前缀
    ./run.sh

    4，SDK用法

    php-sdk:
    https://github.com/Leon2012/goconfd-php-sdk
    include_once("../vendor/autoload.php");
    use goconfd\phpsdk\Goconfd;
    $config = [
        'save_path' => '/tmp/goconfd',
        'key_prefix' => 'key.',
    ];
    $sdk = new Goconfd($config);
    $kv = $sdk->get("key.k1");
    echo $kv->getValue();

    go-sdk:
    gconfd, err := NewGoconfd()
	if err != nil {
		t.Error(err)
	}
	k, err := gconfd.Get("key.k1")

```