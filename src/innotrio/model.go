package innotrio

import (
	"gopkg.in/gorp.v1"
	"fmt"
	"seal_online_go_server/src/innotrio/error"
	"github.com/Sirupsen/logrus"
	"gopkg.in/guregu/null.v2/zero"
	"database/sql"
)

type Model struct {
	DB *gorp.DbMap
	//LogTime bool
}

func (self *Model)Run(query string, args ...interface{}) (safeerror.ISafeError) {
	_, err := self.DB.Exec(query, args...)
	if (err != nil) {
		logrus.Println(query, args, err)
	}
	return self.logError(err, query)
}

func (self *Model)Select(i interface{}, query string, args ...interface{}) (safeerror.ISafeError) {
	_, err := self.DB.Select(i, query, args...)
	return self.logError(err, query)
}

func (self *Model)SelectOne(holder interface{}, query string, args ...interface{}) (safeerror.ISafeError) {
	err := self.DB.SelectOne(holder, query, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	return self.logError(err, query)
}

func (self *Model)MustSelectOne(error string, holder interface{}, query string, args ...interface{}) (safeerror.ISafeError) {
	err := self.SelectOne(holder, query, args...)
	if err != nil {
		return err
	}
	if (holder == nil) {
		return safeerror.NewByCode("NO_SUCH_"+error)
	}

	return nil
}

func (self *Model)SelectStr(query string, args ...interface{}) (string, safeerror.ISafeError) {
	str, err := self.DB.SelectStr(query, args...)
	if err != nil {
		return "", self.logError(err, query)
	}
	return str, nil
}

func (self *Model)SelectInt(query string, args ...interface{}) (int64, safeerror.ISafeError) {
	value, err := self.DB.SelectInt(query, args...)
	if err != nil {
		return 0, self.logError(err, query)
	}
	return value, nil
}

func (self *Model)SelectRowsAndKeys(query string, args ...interface{}) (*sql.Rows, []string, safeerror.ISafeError){
	var rows *sql.Rows
	var keys []string
	var err error

	rows, err = self.DB.Db.Query(query, args...)
	if err != nil {
		return rows, keys, self.logError(err, query)
	}

	keys, err = rows.Columns()
	if err != nil {
		return rows, keys, self.logError(err, query)
	}
	return rows, keys, nil
}


func (self *Model)SelectTable(query string, args ...interface{}) ([][]string, safeerror.ISafeError) {
	var (
		result    [][]string
		container []zero.String
		pointers  []interface{}
	)

	rows, keys, err := self.SelectRowsAndKeys(query, args...)
	if err != nil {
		return result, err
	}

	length := len(keys)
	result = append(result, keys)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]zero.String, length)

		for i := range pointers {
			pointers[i] = &container[i]
		}

		err := rows.Scan(pointers...)
		if err != nil {
			return result, self.logError(err, query)
		}

		tempArr := make([]string, length)
		for i := range container {
			tempArr[i] = container[i].String
		}

		result = append(result, tempArr)
	}
	return result, nil
}

func (self *Model)SelectMap(query string, args ...interface{}) ([]map[string]*zero.String, safeerror.ISafeError) {
	var result []map[string]*zero.String

	rows, keys, err := self.SelectRowsAndKeys(query, args...)
	if err != nil {
		return result, err
	}

	length := len(keys)

	for rows.Next() {
		current := make(map[string]*zero.String)
		pointers := make([]interface{}, length)

		for i, key := range keys {
			handler := zero.String{}
			pointers[i] = &handler
			current[key] = &handler
		}

		err := rows.Scan(pointers...)
		if err != nil {
			return result, self.logError(err, query)
		}

		result = append(result, current)
	}
	return result, nil
}

func (self *Model)logError(err error, message string) (safeerror.ISafeError) {
	if (err != nil) {
		if len(message) > LOG_MESSAGE_LENGTH {
			message = message[:LOG_MESSAGE_LENGTH]
		}
		fmt.Println("DB_ERROR", err, message)
		return safeerror.New(err, "DB_ERROR")
	}
	return nil
}

