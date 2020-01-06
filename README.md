# asyncMessageSystem
Async message system , is based of iris (golang) and rabbitmq

异步消息处理系统，支持高并发。

集成了http、grpc、mvc、中间件模块、yaml配置文件加载、系统日志等功能。

## 实现的功能
- 通知消息的接收、处理、展示

## 依赖存储中间件
- Rabbitmq
- Mysql

## 实践一下
- 下载代码
-     git clone https://github.com/Braveheart7854/asyncMessageSystem.git

- 编译&运行
-     cd asyncMessageSystem/
-     主程序 
      编译 go build app/main/http/main.go
      运行 ./main
-     消息通知消费者
      编译 go build app/consumer/consumer_notify.go
      运行 ./consumer_notify
-     消息读取消费者
      编译 go build app/consumer/consumer_read.go
      运行 ./consumer_read
      
## 接口 
（详见 https://github.com/Braveheart7854/asyncMessageSystem/wiki ）
- 通知消息接收api
-     /api/msg/product
- 通知消息标记为已读api
-     /api/msg/read
- 通知消息列表api
-     /api/msg/list
- 用户信息接口
-     /api/user/info
  
## 模块组成
- 主文件 : app/main/

  /grpc
  
  /http(iris)
  
- 消费者 : app/consumer
  
- 配置文件 : config/

  采用github.com/spf13/viper 加载配置文件config.yaml，如下
  ```yaml
  web:
     debug: true
     server_addr: 0.0.0.0:3335
     read_timeout: 3
     write_timeout: 5
     idle_timeout: 8
   xorm:
     debug: true
     db_type: mysql
     max_lifetime: 7200
     max_open_conns: 2000
     max_idle_conns: 1000
     timezone: Asia/Shanghai
   mysql:
     host: 127.0.0.1
     port: 3306
     user_name: root
     password: hello123
     db_name: message
     parameters: timeout=3s&collation=utf8mb4_unicode_ci&readTimeout=5s&writeTimeout=5s
   rabbitMq:
     host: localhost
     port: 5672
     user_name: guest
     password: guest
     connection_num: 10
     channel_num: 100
   redis:
     host: localhost
     port: 6379
     password:
     database: 0
     ```
  
- 路由   : app/router/url.go

- 中间件 : middleware/
  添加删除某个中间件如下代码
  
  ```go
  var RequireMiddleware = map[string]func(){
  	"mysql"    : mysql.Init,
  	"rabbitmq" : rabbitmq.Init,
  	//"redis"    : redis.Init,
  }
  ```
- 日志模块

  采用go.uber.org/zap + lumberjack 搭建的日志模块

  记录在目录logs/
   
  main.log：记载系统日志，包括info、error、panic等
  
  notice_retry.log：记载需要重试入库的通知消息
  
  read_retry.log：记载需要重新标记为已读的通知消息
  
- mysql 数据表

  主程序启动时自动创建如下数据表 failed_queues、th_notice_1～16、user;
  
  failed_queues：记录失败重试的消息，failed_count表示失败次数，failed_count >= 5的消息需要手动重试；
  
  th_notice_1～16：水平分表，根据uid将用户消息均分到16张表中，分表数量可以根据需要自己调整；
  
  user：用户表，记录消息通知数量；