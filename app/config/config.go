package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"time"
)

type Config struct {
	Web Web
	Mysql Mysql
	Xorm Xorm
	RabbitMq RabbitMq
	Redis Redis
}

type Web struct {
	Debug bool
	ServerAddr string
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	IdleTimeout time.Duration
}

type Xorm struct {
	Debug bool
	DbType string
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLifetime time.Duration
	TimeZone string
}

type Mysql struct {
	Host     string
	Port     int
	UserName string
	Password string
	DbName   string
	Parameters string
	Dsn string
}

func (m Mysql) DSN()string{
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		m.UserName, m.Password, m.Host, m.Port, m.DbName, m.Parameters)
}

type RabbitMq struct {
	Host     string
	Port     int
	UserName string
	Password string
	ConnectionNum int
	ChannelNum int
	Dsn string
}

func (r RabbitMq) DSN()string{
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",r.UserName,r.Password,r.Host,r.Port)
}

type Redis struct {
	Host     string
	Port     int
	Password string
	Database int
}

var Conf = Config{}

// 加载配置
func init(){
	path,_:=filepath.Abs(".")
	v := viper.New()
	v.SetConfigFile(path+"/app/config/config.yaml")
	v.SetConfigType("yaml")
	if err1 := v.ReadInConfig(); err1 != nil {
		log.Panic(err1.Error())
		return
	}
	Conf.Web.Debug = v.GetBool("web.debug")
	Conf.Web.ServerAddr = v.GetString("web.server_addr")
	Conf.Web.ReadTimeout = v.GetDuration("web.read_timeout") * time.Second
	Conf.Web.WriteTimeout = v.GetDuration("web.write_timeout") * time.Second
	Conf.Web.IdleTimeout = v.GetDuration("web.idle_timeout") * time.Second

	Conf.Mysql.Host = v.GetString("mysql.host")
	Conf.Mysql.Port = v.GetInt("mysql.port")
	Conf.Mysql.UserName = v.GetString("mysql.user_name")
	Conf.Mysql.Password = v.GetString("mysql.password")
	Conf.Mysql.DbName = v.GetString("mysql.db_name")
	Conf.Mysql.Parameters = v.GetString("mysql.parameters")
	Conf.Mysql.Dsn = Conf.Mysql.DSN()

	Conf.Xorm.Debug = v.GetBool("xorm.debug")
	Conf.Xorm.DbType = v.GetString("xorm.db_type")
	Conf.Xorm.ConnMaxLifetime = v.GetDuration("xorm.max_lifetime") * time.Second
	Conf.Xorm.MaxOpenConns = v.GetInt("xorm.max_open_conns")
	Conf.Xorm.MaxIdleConns = v.GetInt("xorm.max_idle_conns")
	Conf.Xorm.TimeZone = v.GetString("Xorm.timezone")

	Conf.RabbitMq.Host = v.GetString("rabbitMq.host")
	Conf.RabbitMq.Port = v.GetInt("rabbitMq.port")
	Conf.RabbitMq.UserName = v.GetString("rabbitMq.user_name")
	Conf.RabbitMq.Password = v.GetString("rabbitMq.password")
	Conf.RabbitMq.ConnectionNum = v.GetInt("rabbitMq.connection_num")
	Conf.RabbitMq.ChannelNum = v.GetInt("rabbitMq.channel_num")
	Conf.RabbitMq.Dsn = Conf.RabbitMq.DSN()

	Conf.Redis.Host = v.GetString("redis.host")
	Conf.Redis.Port = v.GetInt("redis.port")
	Conf.Redis.Password = v.GetString("redis.password")
	Conf.Redis.Database = v.GetInt("redis.database")
	return
}