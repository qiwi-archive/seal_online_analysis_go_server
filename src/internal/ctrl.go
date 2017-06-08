package internal

import (
	"seal_online_go_server/src/innotrio/router"
	"seal_online_go_server/src/config"
	"seal_online_go_server/src/innotrio/error"
	"seal_online_go_server/src/data"
)

type Ctrl struct {
	Model *Model
	*router.ApiRouter
}

func (self *Ctrl) Init() {
	self.GET("internal/refresh", self.getInternalHandler(self.refresh))
	self.GET("internal/refresh_providers", self.getInternalHandler(self.refreshProviders))
	self.GET("internal/refresh_tags", self.getInternalHandler(self.refreshTags))
}

/**
Нужно было быстро впилить проверку ключа для всех internal методов =(
 */
func (self *Ctrl) getInternalHandler(handler func(*router.Context)) func(*router.Context) {
	return func(c *router.Context) {
		key := c.GetQuery("key")
		if key == config.ADMIN_REQUEST_CODE {
			handler(c)
		} else {
			self.SendError(c, safeerror.NewByCode(`U_R_VERY_BAD_HACKER`))
		}
	}
}

func (self *Ctrl) refresh(c *router.Context) {
	if err := self.Model.RefreshProviders(); err != nil {
		self.SendError(c, err)
		return
	}

	err := self.Model.UpdateTags([]string{data.TAG_TYPE_PROVIDER_CODE, data.TAG_TYPE_FLOW_CODE})
	self.Send(c, true, err)
}

func (self *Ctrl) refreshProviders(c *router.Context) {
	if err := self.Model.RefreshProviders(); err != nil {
		self.SendError(c, err)
		return
	}
	self.Send(c, true, nil)
}

func (self *Ctrl) refreshTags(c *router.Context) {
	err := self.Model.UpdateTags([]string{data.TAG_TYPE_PROVIDER_CODE, data.TAG_TYPE_FLOW_CODE})
	self.Send(c, true, err)
}