// database/etcd/etcd.go
package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"General_Framework_Gin/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	Client *clientv3.Client
	Kv     clientv3.KV
	Lease  clientv3.Lease
	once   sync.Once
)

// Init 初始化 ETCD 数据库连接
func Init() {
	once.Do(func() {
		etcdCfg := config.AppConfig.Database.ETCD

		// 加载 CA 证书
		caCert, err := ioutil.ReadFile(etcdCfg.CACert)
		if err != nil {
			log.Fatalf("加载 CA 证书失败: %v", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			log.Fatalf("解析 CA 证书失败")
		}

		// 加载客户端证书和私钥
		cert, err := tls.LoadX509KeyPair(etcdCfg.CertFile, etcdCfg.KeyFile)
		if err != nil {
			log.Fatalf("加载客户端证书和私钥失败: %v", err)
		}

		// 创建 TLS 配置
		tlsConfig := &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{cert},
			ServerName:   etcdCfg.ServerName,
		}

		// 配置 ETCD 客户端
		config := clientv3.Config{
			Endpoints:   etcdCfg.Endpoints,
			DialTimeout: time.Duration(etcdCfg.DialTimeout) * time.Millisecond,
			TLS:         tlsConfig,
			Username:    etcdCfg.EtcdName,
			Password:    etcdCfg.Password,
		}

		// 创建 ETCD 客户端
		Client, err = clientv3.New(config)
		if err != nil {
			log.Fatalf("连接 ETCD 失败: %v", err)
		}

		// 初始化 Kv 和 Lease
		Kv = clientv3.NewKV(Client)
		Lease = clientv3.NewLease(Client)

		log.Println("成功连接到 ETCD 数据库")
	})
}

// Close 关闭 ETCD 数据库连接
func Close() {
	if Client != nil {
		log.Println("正在关闭 ETCD 数据库连接...")
		err := Client.Close()
		if err != nil {
			log.Printf("关闭 ETCD 数据库失败: %v", err)
		} else {
			log.Println("ETCD 数据库连接已关闭")
		}
	}
}
