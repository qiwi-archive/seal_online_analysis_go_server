package main
/*
import (
	"testing"
	"seal_online_go_server/src/innotrio"
	"net/http/httptest"
)

var (
	Tester *innotrio.Tester
)

func init() {
	httptest.NewServer(GetMainEngine)
	Tester = innotrio.NewTester("http://localhost:8080")
}

func TestReportsManagers(t *testing.T) {
	Tester.Get(t,"/reports/managers", nil)

	//t.Log(data, isOk)
}

func TestReportsGeneral(t *testing.T) {
	data := map[string][]string {
		"main_tag_code":{"PULL"},
		"inner_join_tag_codes":{"UNSEEN","ALL_BLOCKED","IN_WORK"},
		"left_join_tag_codes":{"FIRM","PSP","IK"},
	}

	// /reports/general?inner_join_tag_codes=UNSEEN&inner_join_tag_codes=ALL_BLOCKED&inner_join_tag_codes=IN_WORK&left_join_tag_codes=FIRM&left_join_tag_codes=PSP&left_join_tag_codes=IK&main_tag_code=PULL


	Tester.Get(t,"/reports/general", data)

	//t.Log(data, isOk)
}*/