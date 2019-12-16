package user

import (
	"asyncMessageSystem/app/controller/producer"
	"asyncMessageSystem/app/model"
	"github.com/kataras/iris"
	"strconv"
)

type User struct {

}

func (u *User)UserInfo(ctx iris.Context){
	uid,_ := strconv.Atoi(ctx.FormValueDefault("uid","0"))

	userModel := new(model.User)
	user := userModel.GetUserInfoByUid(uint64(uid))
	_,_ = ctx.JSON(producer.ReturnJson{Code:10000,Msg:"success",Data: user})
	return
}