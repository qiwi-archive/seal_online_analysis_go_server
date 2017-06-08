package data

import (
	"database/sql"
	"gopkg.in/guregu/null.v2/zero"
	"gopkg.in/guregu/null.v2"
	"time"
)

type HiddenTag struct {
	Id       int `db:"id_tag"`
	Code     sql.NullString `db:"tag_code"`
	Sql      sql.NullString `db:"tag_sql"`
	Db       sql.NullString `db:"tag_db"`
	TypeCode sql.NullString `db:"tag_type_code"`
}

type ReplicationLog struct {
	Id           string `db:"id_replication_log"`
	Code         string `db:"replication_code"`
	StartDate    *time.Time `db:"replication_start_date"`
	EndDate      *time.Time `db:"replication_end_date"`
	MaxDate      *time.Time `db:"replication_max_date"`
	Lag          int `db:"replication_lag"`
	DaysInterval float32 `db:"replication_days_interval"`
}

type StatBillsForCalculation struct {
	IdTag                       null.String `db:"id_tag"`
	BillsAddDate                string `db:"bills_add_date"`
	BillsCountDeviation         float64 `db:"bills_count_deviation"`
	BillsCountCoeff             float64 `db:"bills_count_coeff"`
	BillsPaymentsCountDeviation float64 `db:"bills_payments_count_deviation"`
	BillsPaymentsCoeff          float64 `db:"bills_payments_coeff"`
	BillsConversionDeviation    float64 `db:"bills_conversion_deviation"`
	BillsConversionCoeff        float64 `db:"bills_conversion_coeff"`
	PaidBillsAmount             float64 `db:"paid_bills_amount"`
	PaidBillsCount              float64 `db:"paid_bills_count"`
	BillsAmount                 float64 `db:"bills_amount"`
	BillsCount                  float64 `db:"bills_count"`
	BillsPaidCount              float64 `db:"bills_paid_count"`
	BillsPaymentsCount          float64 `db:"bills_payments_count"`
	BillsConversion             float64 `db:"bills_conversion"`
	BillsAddTimestamp           int64 `db:"bills_add_timestamp"`
	StatPaidBillsAmount         float64 `db:"sbills_paid_amount"`
	StatPaidBillsCount          float64 `db:"sbills_paid_count"`
	StatBillsAmount             float64 `db:"sbills_all_amount"`
	StatBillsCount              float64 `db:"sbills_all_count"`
	StatBillsConversion         float64 `db:"sbills_conversion"`
	StatPaymentsBillsCount      float64 `db:"sbills_payments_count"`
	StatPaymentsAmount          float64 `db:"sbills_payments_amount"`
	StatHealth                  float64 `db:"sbills_health"`
	StatPaySeconds              float64 `db:"sbills_pay_seconds"`
}

type CaseForCalculation struct {
	Code              string `db:"case_code"`
	Status            string `db:"case_status"`
	StartTimestamp    int64 `db:"case_start_timestamp"`
	EndTimestamp      int64 `db:"case_end_timestamp"`
	StartDate         string `db:"case_start_dtime"`
	EndDate           string `db:"case_end_dtime"`
	Amount            int64 `db:"case_amount"`
	StatAmount        int64 `db:"case_stat_amount"`
	LostAmount        int64 `db:"case_lost_amount"`
	Count             int64 `db:"case_count"`
	StatCount         int64 `db:"case_stat_count"`
	PaidCount         int64 `db:"case_paid_count"`
	PaidStatCount     int64 `db:"case_paid_stat_count"`
	PaymentsCount     int64 `db:"case_payments_count"`
	PaymentsStatCount int64 `db:"case_payments_stat_count"`
	Duration          int `db:"case_duration"`
	Id                string `db:"id_case"`
	IdResolution      null.String `db:"id_resolution"`
	IdTag             zero.String `db:"id_tag"`
}
type DbCase struct {
	IdCase string `db:"id_case"`
}

type DbPrvInfo struct {
	IdPrv         string `db:"id_prv"`
	MaxNoTxns     int `db:"prv_max_no_txns"`
	AvgPaidAmount int64
}

type DbReportSql struct {
	Id           string `db:"id_report"`
	Sql          string `db:"report_sql"`
	OrderSql     null.String `db:"report_order_sql"`
	HasComments  bool `db:"report_has_comments"`
	SqlUniqueCol zero.String `db:"report_sql_uniq_col"`
}
