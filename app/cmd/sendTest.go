package main

import (
	"asyncMessageSystem/app/middleware/mysql"
	"asyncMessageSystem/app/model"
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main(){
	mysql.InitDB()
	model.SendTest{}.Prepare()

	file,err := os.Open("./index.log")
	if err != nil {
		log.Panic(err.Error())
	}
	defer file.Close()

	forever := make(chan bool)

	limit := make(chan int ,100)

	rd := bufio.NewReader(file)
	for {

		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if io.EOF == err {
			log.Println("finish")
		}
		if err != nil || io.EOF == err {
			break
		}
        limit<- 1
		go func() {
			//fmt.Print(line)
			pattern:= "((([0-9]{3}[1-9]|[0-9]{2}[1-9][0-9]{1}|[0-9]{1}[1-9][0-9]{2}|[1-9][0-9]{3})\\/(((0[13578]|1[02])\\/(0[1-9]|[12][0-9]|3[01]))|"+
				"((0[469]|11)\\/(0[1-9]|[12][0-9]|30))|(02\\/(0[1-9]|[1][0-9]|2[0-8]))))|((([0-9]{2})(0[48]|[2468][048]|[13579][26])|"+
				"((0[48]|[2468][048]|[3579][26])00))\\/02\\/29))\\s([0-1][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])\\s([0-9]\\d*)"
			res,regerr := regexp.MatchString(pattern,line)
			if regerr != nil{
				log.Println("regexp match error",regerr.Error())
				<-limit
				return
			}
			//fmt.Println(res)
			if res {
				strs := strings.Split(strings.TrimRight(line,"\n")," ")
				//fmt.Printf("0:%s,1:%s,2:%s",strs[0],strs[1],strs[2])
				order,erratoi := strconv.Atoi(strs[2])
				if erratoi != nil{
					log.Println("Atoi error ",erratoi.Error(),strs)
				}
				insert,errInsert := model.SendTest{}.InsertSend(model.SendTest{Time:strs[0]+" "+strs[1],Order:int64(order)})
				if errInsert != nil{
					log.Println("InsertSend error",errInsert.Error(),strs)
				}
				if !insert {
					log.Println("insert failed",strs)
				}else{
					//fmt.Println("insert success",strs)
				}
			}
			<-limit
		}()

	}

	<- forever
	log.Println("finish")
}