package validator
import (
	"strconv"
	"seal_online_go_server/src/innotrio/error"
	"golang.org/x/net/html"
	"regexp"
)

type StringValidator func(string) string

func IsInt(item string) (int) {
	num, err := strconv.Atoi(item)
	checkError(err, "VALIDATION_NOT_INT")
	return num
}

func IsAllowed(item string, allowedValues []string) (string) {
	item = html.EscapeString(item)
	for _, value := range allowedValues {
		if value == item {
			return item
		}
	}
	sendErrorByCode(`VALIDATION_NOT_ALLOWED_VALUE: ` + item)
	return ``
}
func Match(item string, pattern string) (string) {
	item = html.EscapeString(item)
	isOk, _ := regexp.MatchString(pattern, item)
	if (isOk) {
		return item
	}
	sendErrorByCode(`VALIDATION_NOT_MATCHES: ` + item)
	return ``
}

func IsCode(item string) (string) {
	isOk, err := regexp.MatchString("^[a-zA-Z0-9_]+$", item)
	checkError(err, "VALIDATION_CODE_REGEXP")

	if (isOk) {
		return item
	}
	sendErrorByCode("VALIDATION_NOT_CODE")
	return ""
}

func IsValidStrArr(items []string, validateFunc StringValidator, required bool) []string {
	if (items == nil) {
		if (required){
			sendErrorByCode("NO_VALID_STR_ARR")
			return nil
		}else{
			return make([]string, 0)
		}

	}
	results := make([]string, len(items))
	for i, item := range items {
		results[i] = validateFunc(item)
	}
	return results
}

func Escape(item string) (string) {
	return html.EscapeString(item)
}

func checkError(err error, errCode string) {
	if err != nil {
		panic(safeerror.New(err, errCode))
	}
}

func sendErrorByCode(errCode string) {
	panic(safeerror.NewByCode(errCode))
}