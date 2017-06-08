package data

import (
	"gopkg.in/guregu/null.v2/zero"
	"gopkg.in/guregu/null.v2"
)

const (
	//Ожидает решения пользователя
	CASE_STATUS_WAITING = "W"
	//Все хорошо, просто записали данные для статистики
	CASE_STATUS_STAT = "S"

	PRV_CASE_DESCRIPTION_ZERO     = "Нет транзакций"
	PRV_CASE_DESCRIPTION_OK       = "Показатели в норме"
	PRV_CASE_DESCRIPTION_TOO_HIGH = "Превышение"
	PRV_CASE_DESCRIPTION_TOO_LOW  = "Ниже ожиданий"

	CASE_STAT_COUNT_LOWEST_VALUE    = 10

	TAG_CASE_LOST_AMOUNT_LOWER_LIMIT = -60000
	TAG_CASE_CODE_MINUTES            = "TAG_CASE_CODE_MINUTES"

	REPLICATION_QW_MINUTES_MAIN_CODE              = "QW_MINUTES"
	REPLICATION_QW_MINUTES_EXTRAS_MAIN_CODE       = "QW_MINUTES_EXTRAS"

	TAG_TYPE_PROVIDER_CODE = "P"
	TAG_TYPE_FLOW_CODE     = "F"
)

type Case struct {
	Description    string `db:"case_description"`
	Code           string `db:"case_code"`
	Status         string `db:"case_status"`
	StartDate      string `db:"case_start_dtime"`
	EndDate        string `db:"case_end_dtime"`
	StartTimestamp int64 `db:"case_start_timestamp"`
	EndTimestamp   int64 `db:"case_end_timestamp"`
	Amount         int64 `db:"case_amount"`
	StatAmount     int64 `db:"case_stat_amount"`
	LostAmount     int64 `db:"case_lost_amount"`
	Count          int64 `db:"case_count"`
	StatCount      int64 `db:"case_stat_count"`
	PaidCount         int64 `db:"case_paid_count"`
	PaidStatCount     int64 `db:"case_paid_stat_count"`
	PaymentsCount     int64 `db:"case_payments_count"`
	PaymentsStatCount int64 `db:"case_payments_stat_count"`
}

type FullCase struct {
	Case
	Id             string `db:"id_case"`
	Message        zero.String `db:"case_message"`
	Duration       int `db:"case_duration"`
	LostAmount     int64 `db:"case_lost_amount"`
	IdResolution   null.String `db:"id_resolution"`
	Diff           float64 `db:"case_diff"`
	PrvName        zero.String `db:"prv_name"`
	PrvManager     zero.String `db:"prs_manager"`
	ResolutionName null.String `db:"resolution_name"`
}

type FullTagCase struct {
	FullCase
	IdTag      string `db:"id_tag"`
	StatByDays *[]StatBillsTag
}

type StatBills struct {
	AllCount      float32 `db:"sbills_all_count"`
	PaidCount     float32 `db:"sbills_paid_count"`
	AllAmount     float32 `db:"sbills_all_amount"`
	PaidAmount    float32 `db:"sbills_paid_amount"`
	Conversion    float32 `db:"sbills_conversion"`
	Date          string `db:"sbills_date"`
	DateTimestamp int64 `db:"sbills_add_timestamp"`
	AddSection    float32 `db:"sbills_add_section"`
	TimeCode      string `db:"sbills_time_code"`
}
type StatBillsTag struct {
	StatBills
	Id    string `db:"id_stat_bills_tag"`
	IdTag string `db:"id_tag"`
}