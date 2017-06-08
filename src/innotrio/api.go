package innotrio

import (
	"github.com/parnurzeal/gorequest"
	"encoding/json"
	"seal_online_go_server/src/innotrio/error"
	"fmt"
)

type Api struct {
	Request *gorequest.SuperAgent
	Host    string
}

func NewApi(host string) *Api {
	return &Api{gorequest.New(),host}
}

func (self *Api)Get(url string, data map[string][]string) (interface{}, safeerror.ISafeError) {
	request := self.Request.Get(self.Host+url)
	if (data!=nil){
		//inner_join_tag_codes=UNSEEN&inner_join_tag_codes=ALL_BLOCKED&inner_join_tag_codes=IN_WORK&left_join_tag_codes=FIRM&left_join_tag_codes=PSP&left_join_tag_codes=IK&main_tag_code=PULL
		for key, values := range data {
			for _, value := range values {
				request.Query(key+"="+value)
			}
		}
	}
	_, body, errs := request.End()
	if (errs != nil) {
		return nil, self.logErrors(errs,"");
	}
	var x map[string]interface{}
	json.Unmarshal([]byte(body), &x)
	return x["result"], nil;
}

func (self *Api)logErrors(errs []error, message string) (safeerror.ISafeError) {
	if (errs != nil) {
		errText := "";
		for _, err := range errs {
			errText += err.Error();
		}
		if len(message) > LOG_MESSAGE_LENGTH {
			message = message[:LOG_MESSAGE_LENGTH]
		}
		fmt.Println("API_ERROR", errText, message)
		return safeerror.NewByCode("API_ERROR")
	}
	return nil
}