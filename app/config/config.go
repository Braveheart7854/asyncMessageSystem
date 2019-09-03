package config

const (
	Debug = true

	ServerAddr = "localhost" //开发机、测试机

	//开发机
	MysqlDataSource = "root:hello123@tcp(127.0.0.1:3306)/wx_forum" +
		"?timeout=3s&collation=utf8mb4_unicode_ci&readTimeout=5s&writeTimeout=5s"

	//测试机
	//MysqlDataSource = "root:xC6ch6I4u0X3h@tcp(127.0.0.1:3306)/wx_forum" +
	//		"?timeout=3s&collation=utf8mb4_unicode_ci&readTimeout=5s&writeTimeout=5s"


	AmqpUrl = "amqp://guest:guest@localhost:5672/" //开发机、测试机

	OssDomain = "http://image.com"
	)