package config

import "time"

const (
	Debug = true

	ServerAddr = "localhost:3333" //开发机、测试机
	ReadTimeout = 3*time.Second
	WriteTimeout = 5*time.Second

	//开发机
	MysqlDataSource = "root:hello123@tcp(127.0.0.1:3306)/message" +
		"?timeout=3s&collation=utf8mb4_unicode_ci&readTimeout=5s&writeTimeout=5s"
	MaxOpenConns    = 2000
	MaxIdleConns    = 1000
    ConnMaxLifetime = 9*time.Second
    TimeZone        = "Asia/Shanghai"

	//测试机
	//MysqlDataSource = "root:xC6ch6I4u0X3h@tcp(127.0.0.1:3306)/wx_forum" +
	//		"?timeout=3s&collation=utf8mb4_unicode_ci&readTimeout=5s&writeTimeout=5s"


	AmqpUrl       = "amqp://guest:guest@localhost:5672/" //开发机、测试机
	ConnectionNum = 10
	ChannelNum    = 100

	OssDomain = "http://image.com"
	)