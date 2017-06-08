package innotrio

import (
	"testing"
	"github.com/parnurzeal/gorequest"
	"encoding/json"
)

type Tester struct {
	Request *gorequest.SuperAgent
	Host    string
}

func NewTester(host string) *Tester {
	return &Tester{gorequest.New(),host}
}

func (self *Tester)Get(t *testing.T, url string, data map[string][]string) (interface{}, bool) {
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
		t.Fatal(errs);
		return nil, false;
	}
	var x map[string]interface{}
	json.Unmarshal([]byte(body), &x)
	return x["result"], true;
}
