package common

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"reflect"
)


const (
	ExchangeNameNotice = "exchange_wxforum_notice"
	RouteKeyNotice     = "route_wxforum_notice"
	QueueNameNotice    = "queue_wxforum_notice"

	ExchangeNameRead = "exchange_wxforum_read"
	RouteKeyRead     = "route_wxforum_read"
	QueueNameRead    = "queue_wxforum_read"

	TypeTopicLiked   = 1
	TypeReplyLiked   = 2
	TypeTopicReplied = 3
	TypeUserFocused  = 4
	TypeDeleteAvatar = 5
	TypeForbidUser   = 6
	TypeTopicDeleted = 7
	TypeAdminMsg     = 8
	TypeBeanPrize    = 9
	TypeReplyDeleted = 10
	TypeTopicLabel   = 11

    TYPE_LIKE      = 1   //赞、桃
    TYPE_REPLY     = 2   //回复
    TYPE_FOCUS     = 3   //关注
    TYPE_SYSTEM    = 4   //系统
)

func IsEmpty(a interface{}) bool {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v=v.Elem()
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// 生成32位MD5
func MD5(text string) string {
	md := md5.New()
	md.Write([]byte(text))
	return hex.EncodeToString(md.Sum(nil))
}

func GetHaseValue(uid int)int {
	return (uid % 16)+1
}

/**
记录失败任务
 */
func LogErrorJobs(db *sql.DB,orderSn string,data string,typ string) bool{
	var failed_count int
	_ = db.QueryRow("select failed_count from failed_queues where order_sn = ? and type = ?", orderSn,typ).Scan(&failed_count)
	//if err != nil {
	//	log.Println(err)
	//}
	if failed_count >= 5 {
		return true
	}else if failed_count >= 1 && failed_count < 5 {
		exec,err := db.Prepare("update failed_queues set failed_count=failed_count+1 where order_sn=? and type = ?")
		if err != nil {
			log.Println(err)
		}
		result,err := exec.Exec(orderSn,typ)
		if err != nil {
			log.Println(err)
		}
		result.LastInsertId()
		exec.Close()
		return false
	}else {
		exec,err := db.Prepare("insert into failed_queues (order_sn,data,failed_count,type) values (?,?,?,?)")
		if err != nil {
			log.Println(err)
		}
		result,err := exec.Exec(orderSn,data,1,typ)
		if err != nil {
			log.Println(err)
		}
		result.LastInsertId()
		exec.Close()
		return false
	}
}

func NoticeType(typ int)string{
	switch typ {
	case TYPE_LIKE:
		return "1,2"
	case TYPE_REPLY:
		return "3"
	case TYPE_FOCUS:
		return "4"
	case TYPE_SYSTEM:
		return "5,6,7,8,9,10"
	default:
		return "0"
	}
}


func Log(file string,content string){
	//写入文件
	file6, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766);

	loger := log.New(file6, "前缀", log.Ldate|log.Ltime|log.Lshortfile)

	//SetFlags设置输出选项
	loger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//设置输出前缀
	loger.SetPrefix("test_")

	loger.Output(2, content)

	file6.Close()
}