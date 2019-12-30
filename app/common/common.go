package common

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"reflect"
	"sort"
)


const (
	ExchangeNameNotice = "exchange_msg_notice"
	RouteKeyNotice     = "route_msg_notice"
	QueueNameNotice    = "queue_msg_notice"

	ExchangeNameRead = "exchange_msg_read"
	RouteKeyRead     = "route_msg_read"
	QueueNameRead    = "queue_msg_read"

	ExchangeNameQuizzesEvent = "exchange_msg_quizzes_event"
	RouteKeyQuizzesEvent     = "route_msg_quizzes_event"
	QueueNameQuizzesEvent    = "queue_msg_quizzes_event"

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

    LAYOUT_STYLE     = "2006-01-02 15:04:05"
    LAYOUT_STYLE_ONE = "01/02 15:04:05"
    LAYOUT_DATE      = "2006-01-02"
	LAYOUT_DATE_ONE  = "2006/01/02"
	LAYOUT_TIME      = "15:04:05"

	SUCCESS = 10000
	FAILED  = 10001
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

func GetHaseValue(uid uint64)uint64 {
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
		return "5,6,7,8,9,10,11"
	default:
		return "0"
	}
}


func Log(file string,content string){
	//写入文件
	file6, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766);

	loger := log.New(file6, "", log.Ldate|log.Ltime|log.Llongfile)

	//SetFlags设置输出选项
	loger.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	//设置输出前缀
	//loger.SetPrefix("test_")

	loger.Output(2, content)

	file6.Close()
}

func SortMap(mp map[string][]interface{})(list map[string][]interface{}) {
	var newMp = make([]string, 0)
	for k, _ := range mp {
		newMp = append(newMp, k)
	}
	//sort.Strings(newMp)
	sort.Sort(sort.Reverse(sort.StringSlice(newMp)))
	list = make(map[string][]interface{})
	for _, v := range newMp {
		list[v] = mp[v]
		fmt.Println("根据key排序后的新集合》》   key:", v, "    value:", mp[v])
	}
	return
}