package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"seal_online_go_server/src/innotrio/error"
	"seal_online_go_server/src/config"
)

   type ApiRouter struct {
	Router *gin.Engine
	Name   string
}

func NewApiRouter(router *gin.Engine, name string) *ApiRouter {
	return &ApiRouter{router, name}
}

type IContext interface {
	JSON(code int, obj interface{})
}


func (self *ApiRouter) Send(c IContext, obj interface{}, err safeerror.ISafeError) {
	if err == nil {
		self.SendResult(c, obj)
	} else {
		self.SendError(c, err)
	}
}

func (self *ApiRouter) SendResult(c IContext, obj interface{}) {
	c.JSON(200, gin.H{
		"status": "done",
		"result": obj,
	})
}

func (self *ApiRouter) SendError(c IContext, err safeerror.ISafeError) {
	var code string
	if err != nil {
		fmt.Println(err.Error())
		code = err.Code()
	}

	self.sendError(c, code)
}

func (self *ApiRouter) sendError(c IContext, code string) {
	if code != "" {
		code = "_" + code
	}
	//Добавляем название роутера
	code = self.Name + code

	c.JSON(200, gin.H{
		"status": "error",
		"code":   "ERROR_" + code,
	})
}

type LogicFunc func(c *Context)

func (self *ApiRouter) GET(path string, logicFunc LogicFunc) {
	self.Router.GET(config.BASE_API_PREFIX + path, func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(safeerror.ISafeError); ok == true {
					self.sendError(c, r.(safeerror.ISafeError).Code())
				} else {
					self.sendError(c, "UNKNOWN_GET_ERROR")
				}
			}
		}()
		logicFunc(NewContext(c));
	})
}

func (self *ApiRouter) POST(path string, logicFunc LogicFunc) {
	self.Router.POST(config.BASE_API_PREFIX + path, func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				if _, ok := r.(safeerror.ISafeError); ok == true {
					self.sendError(c, r.(safeerror.ISafeError).Code())
				} else {
					self.sendError(c, "UNKNOWN_POST_ERROR")
				}
			}
		}()
		logicFunc(NewContext(c));
	})
}


//func (self *ApiRouter)Validate(c *gin.Context, func func{})
