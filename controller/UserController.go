package controller

import (
	"net/http"
	"goadmin/model"
	"goadmin/lolgo"

	"github.com/gin-gonic/gin"

	"goadmin/entity"
	"goadmin/service"
	"strconv"

	"github.com/gin-gonic/gin/binding"


)

//用户管理控制器
type UserController struct {
	lolgo.Controller
}

//用户服务层
var userService service.UserService

//路由注册
func (ctrl *UserController) Router(router *gin.Engine) {

	r := router.Group("user")
	r.POST("search", ctrl.query)
	r.POST("findOne", ctrl.findOne)
	r.POST("register", ctrl.register)
	r.POST("login", ctrl.login)
	r.Any("quit", ctrl.quit)
	r.POST("updatestat", ctrl.updateStat)

}

//基于全部的搜索
func (ctrl *UserController) query(ctx *gin.Context) {
	var userArg model.UserArg

	ctx.ShouldBindWith(&userArg, binding.FormPost)
	ret := userService.Query(userArg)
	num := userService.Count(userArg)
	//最后响应数据列表到前端
	lolgo.ResultList(ctx, ret, num)
}

func (ctrl *UserController) findOne(ctx *gin.Context) {
	userId, _ := strconv.ParseInt(ctx.PostForm("userId"), 10, 64)
	ret := userService.FindOne(userId)
	lolgo.ResultOk(ctx, ret)
}

func (ctrl *UserController) updateStat(ctx *gin.Context) {
	userId, _ := strconv.ParseInt(ctx.PostForm("id"), 10, 64)
	stat, _ := strconv.Atoi(ctx.PostForm("stat"))
	_, err := userService.UpdateStat(userId, stat)
	if err != nil {
		lolgo.ResultFail(ctx, "修改失败,请稍后再试"+err.Error())
	} else {
		lolgo.ResultOkMsg(ctx, nil, "修改成功,请稍后再试")
	}

}

//注册用户信息
func (ctrl *UserController) register(ctx *gin.Context) {

	code := ctx.PostForm("verify")
	r := lolgo.CheckVerify(ctx, code)
	if !r {
		lolgo.ResultFail(ctx, "验证码错误请重试")
		return
	}
	var user entity.User
	ctx.ShouldBindWith(&user, binding.FormPost)

	ret, err := userService.Register(ctx, &user)
	if err != nil {
		lolgo.ResultFail(ctx, err.Error())
	} else {
		lolgo.ResultOkMsg(ctx, ret, "恭喜你注册成功")
	}

}

//注册用户信息
func (ctrl *UserController) login(ctx *gin.Context) {

	code := ctx.PostForm("verify")
	r := lolgo.CheckVerify(ctx, code)
	if !r {
		lolgo.ResultFail(ctx, "验证码错误请重试")
		return
	}

	ret, err := userService.Login(ctx, ctx.PostForm("kword"), ctx.PostForm("passwd"))
	if err != nil {
		lolgo.ResultFail(ctx, err.Error())
	} else {
		lolgo.ResultOkMsg(ctx, ret, "恭喜你登录成功")
	}

}


//退出系统
func (ctrl *UserController) quit(ctx *gin.Context) {
	lolgo.ClearAllSession(ctx)
	ctx.Redirect(http.StatusMovedPermanently, "/")
}
