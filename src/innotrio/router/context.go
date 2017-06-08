package router
import "github.com/gin-gonic/gin"

type Context struct {
	*gin.Context
}

func NewContext(context *gin.Context) *Context {
	return &Context{context}
}

func (self *Context) GetQuery(code string) string {
	return self.Request.URL.Query().Get(code)
}

func (self *Context) GetQueryArr(code string) []string {
	return self.Request.URL.Query()[code]
}

func (self *Context) GetBody(code string) string {
	return self.PostForm(code)
}

/*
func (self *Context) GetBodyArr(code string) []string {
	return self.Request.Body.Query()[code]
}*/
