package cases

import (
	"seal_online_go_server/src/innotrio"
	"seal_online_go_server/src/innotrio/error"
	"seal_online_go_server/src/data"
	"strings"
	"io/ioutil"
	"path/filepath"
)

type CalculationModel struct {
	*innotrio.Model
	CreatedCases map[string]bool
}

func (self *CalculationModel) RefreshCases() (map[string]bool, safeerror.ISafeError) {
	tracker := innotrio.NewTracker()
	defer func() {
		tracker.Log(`CASES_REFRESH_DONE`)
	}()
	self.CreatedCases = make(map[string]bool)

	//Обновляем расчетные данные
	err := self.generateStatBills()
	if err != nil {
		return self.CreatedCases, err
	}
	//Получаем данные для анализа
	statBills, err := self.getStatAndActualBills()
	if err != nil {
		return self.CreatedCases, err
	}

	//Анализируем статистику и создаем одиночные инциденты
	err = self.deleteOldCases()
	if err != nil {
		return self.CreatedCases, err
	}

	//Анализируем статистику и создаем одиночные инциденты
	err = self.generateCases(statBills)
	if err != nil {
		return self.CreatedCases, err
	}

	//Объединяем инциденты
	err = self.associateCases()
	if err != nil {
		return self.CreatedCases, err
	}

	//Гарантирует корректность сумм статов в родительских инцидентах
	err = self.refreshParentCasesStats()
	if err != nil {
		return self.CreatedCases, err
	}
	return self.CreatedCases, nil
}

func (self *CalculationModel) generateStatBills() (safeerror.ISafeError) {
	absPath, _ := filepath.Abs("cases/sql/generate_stat_bills.sql")
	sql, err := ioutil.ReadFile(absPath)
	if err != nil {
		return safeerror.New(err, `ERR_GENERATE_STAT_BILLS_READ_FILE`)
	}
	return self.Run(string(sql))
}

func (self *CalculationModel) getStatAndActualBills() (*[]data.StatBillsForCalculation, safeerror.ISafeError) {
	var items []data.StatBillsForCalculation
	absPath, _ := filepath.Abs("cases/sql/get_stat_and_actual_bills.sql")
	sql, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, safeerror.New(err, `ERR_GET_STAT_BILLS_READ_FILE`)
	}
	error := self.Select(&items, string(sql))
	return &items, error
}

func (self *CalculationModel) deleteOldCases() (safeerror.ISafeError) {
	absPath, _ := filepath.Abs("cases/sql/delete_old_cases.sql")
	sql, err := ioutil.ReadFile(absPath)
	if err != nil {
		return safeerror.New(err, `ERR_GENERATE_STAT_BILLS_READ_FILE`)
	}
	return self.Run(string(sql))
}

func (self *CalculationModel) generateCases(dayStat *[]data.StatBillsForCalculation) (safeerror.ISafeError) {
	if len(*dayStat) < 1 {
		return nil
	}

	status := ""
	description := ""

	for _, stat := range *dayStat {
		//По умолчанию считаем, что все хорошо
		status = data.CASE_STATUS_STAT
		description = data.PRV_CASE_DESCRIPTION_OK

		var caseLostAmount float64 = stat.PaidBillsAmount - stat.StatPaidBillsAmount

		if stat.StatPaidBillsCount < data.CASE_STAT_COUNT_LOWEST_VALUE {
			// Меньше 10 не мониторим
			continue
		}

		status = data.CASE_STATUS_WAITING
		description = data.PRV_CASE_DESCRIPTION_TOO_LOW

		var newCase data.Case = data.Case{
			description,
			data.TAG_CASE_CODE_MINUTES,
			status,
			stat.BillsAddDate,
			stat.BillsAddDate,
			stat.BillsAddTimestamp,
			stat.BillsAddTimestamp,
			int64(stat.PaidBillsAmount),
			int64(stat.StatPaidBillsAmount),
			int64(caseLostAmount),
			int64(stat.BillsCount),
			int64(stat.StatBillsCount),
			int64(stat.PaidBillsCount),
			int64(stat.StatPaidBillsCount),
			int64(stat.BillsPaymentsCount),
			int64(stat.StatPaymentsBillsCount),
		}
		var entityId = stat.IdTag.String
		if err := self.createCaseIfNotExist(&newCase, entityId); err != nil {
			return err
		}
	}
	return nil
}

func (self *CalculationModel) createCaseIfNotExist(item *data.Case, entityId string) (safeerror.ISafeError) {
	var dbCase *data.DbCase
	err := self.SelectOne(&dbCase, `SELECT id_case FROM obj_case_tag WHERE id_tag = $1 AND case_start_dtime = $2 AND case_code = $3 AND case_start_dtime = case_end_dtime`, entityId, item.StartDate, item.Code)

	if err != nil {
		return err
	}

	if dbCase == nil {
		absPath, _ := filepath.Abs("cases/sql/insert_obj_case.sql")
		sql, error := ioutil.ReadFile(absPath)
		if error != nil {
			return safeerror.New(error, `ERR_GENERATE_STAT_BILLS_READ_FILE`)
		}
		err = self.SelectOne(&dbCase, string(sql),
			entityId, item.Description, item.Status, item.Code, item.StartDate, item.EndDate, item.StartTimestamp, item.EndTimestamp, item.LostAmount, item.Amount, item.StatAmount, item.Count, item.StatCount, item.PaidCount, item.PaidStatCount, item.PaymentsCount, item.PaymentsStatCount)
		if dbCase != nil && item.LostAmount < data.TAG_CASE_LOST_AMOUNT_LOWER_LIMIT {
			self.CreatedCases[dbCase.IdCase] = true
		}
		return err

	}
	return nil
}

func (self *CalculationModel) associateCases() (safeerror.ISafeError) {
	var items []data.CaseForCalculation
	//Выбираем только что сгенерированные инциденты и родительские инциденты, которые заканчиваются в последние 12 часов
	//При этом первыми строками будут родительские инциденты (длятся несколько временных промежутков)
	absPath, _ := filepath.Abs("cases/sql/select_cases_for_association.sql")
	sql, error := ioutil.ReadFile(absPath)
	if error != nil {
		return safeerror.New(error, `ERR_GENERATE_STAT_BILLS_READ_FILE`)
	}
	err := self.Select(&items, string(sql))
	if err != nil {
		return err
	}

	parentCases := []data.CaseForCalculation{}

	var findParentCase = func(dbCase *data.CaseForCalculation, parentCases []data.CaseForCalculation) (data.CaseForCalculation) {
		var result data.CaseForCalculation
		for _, parentCase := range parentCases {
			// Совпадает id тега и время конца \ начала близко (в пределах 15 минут считается что инцидент продолжается)
			if parentCase.IdTag.String == dbCase.IdTag.String && dbCase.StartTimestamp - parentCase.EndTimestamp <= 60 * 15 {
				result = parentCase
				break
			}
		}
		return result
	}

	var currentIdEntity string
	var currentDbCase *data.CaseForCalculation
	var currentTimestamp int64
	var startIndex int = -1
	for currentIndex, dbCase := range items {
		//Сначала идут все "родительские" - длительные
		if dbCase.Duration > 1 {
			parentCases = append(parentCases, dbCase)
		} else {
			// затем все одноинтервальные
			// Если пошел следующий тег или начало следующего отличается от конца предыдущего на более чем 5 минут
			if dbCase.IdTag.String != currentIdEntity || dbCase.StartTimestamp - currentTimestamp > 60 * 5 {
				if currentIdEntity != `` && startIndex != -1 {
					self.createOrUpdateParentCase(
						items[startIndex:currentIndex],
						findParentCase(currentDbCase, parentCases))
				}
				//Пошел следующий тег или серия
				startIndex = currentIndex
				currentIdEntity = dbCase.IdTag.String
				currentDbCase = &dbCase
				currentTimestamp = dbCase.EndTimestamp
			}
		}
	}

	if startIndex == -1 {
		return nil
	}

	return self.createOrUpdateParentCase(items[startIndex:], findParentCase(currentDbCase, parentCases))
}

func (self *CalculationModel) createOrUpdateParentCase(cases []data.CaseForCalculation, parentCase data.CaseForCalculation) (safeerror.ISafeError) {
	var casesLength int = len(cases)
	if casesLength <= 1 {
		return nil
	}

	var caseTable string = "obj_case_tag"

	var dbParentCase *data.DbCase
	var paidBillsAmount int64 = parentCase.Amount
	var statPaidBillsAmount int64 = parentCase.StatAmount
	var billsCount int64 = parentCase.Count
	var statBillsCount int64 = parentCase.StatCount
	var paidBillsCount int64 = parentCase.PaidCount
	var statPaidBillsCount int64 = parentCase.PaidStatCount
	var paymentsCount int64 = parentCase.PaymentsCount
	var statPaymentsCount int64 = parentCase.PaymentsStatCount
	var caseLostAmount int64 = parentCase.LostAmount
	var description string
	var status string


	var casesIds []string = []string{}
	var firstCase data.CaseForCalculation = cases[0]
	var lastCase data.CaseForCalculation = cases[casesLength - 1];
	for _, currentCase := range cases {
		casesIds = append(casesIds, currentCase.Id)
		paidBillsAmount += currentCase.Amount
		statPaidBillsAmount += currentCase.StatAmount
		billsCount += currentCase.Count
		billsCount += currentCase.StatCount
		paidBillsCount += currentCase.PaidCount
		statPaidBillsCount += currentCase.PaidStatCount
		paymentsCount += currentCase.PaymentsCount
		statPaymentsCount += currentCase.PaymentsStatCount
		caseLostAmount += currentCase.Amount - currentCase.StatAmount
	}

	if &parentCase == nil || parentCase.Id == `` {
		if casesLength == 1 {
			return nil
		}

		if statPaidBillsAmount == 0 {
			return nil
		}

		status = data.CASE_STATUS_WAITING
		if paidBillsAmount == 0 {
			description = data.PRV_CASE_DESCRIPTION_ZERO
		} else if paidBillsAmount < statPaidBillsAmount {
			description = data.PRV_CASE_DESCRIPTION_TOO_LOW
		} else if paidBillsAmount > statPaidBillsAmount {
			description = data.PRV_CASE_DESCRIPTION_TOO_HIGH
		}
		absPath, _ := filepath.Abs("cases/sql/insert_obj_case.sql")
		sql, error := ioutil.ReadFile(absPath)
		if error != nil {
			return safeerror.New(error, `ERR_GENERATE_STAT_BILLS_READ_FILE`)
		}
		err := self.SelectOne(&dbParentCase, string(sql),
			lastCase.IdTag.String, description, status, lastCase.Code, firstCase.StartDate, lastCase.EndDate,
			firstCase.StartTimestamp, lastCase.EndTimestamp, int64(caseLostAmount), paidBillsAmount, statPaidBillsAmount, billsCount, statBillsCount, paidBillsCount, statPaidBillsCount, paymentsCount, statPaymentsCount)
		if err != nil {
			return err
		}
		if caseLostAmount < data.TAG_CASE_LOST_AMOUNT_LOWER_LIMIT {
			self.CreatedCases[dbParentCase.IdCase] = true
		}
	} else {
		if parentCase.LostAmount > data.TAG_CASE_LOST_AMOUNT_LOWER_LIMIT && caseLostAmount < data.TAG_CASE_LOST_AMOUNT_LOWER_LIMIT {
			self.CreatedCases[parentCase.Id] = true
		}
		dbParentCase = &data.DbCase{parentCase.Id}
	}

	if &dbParentCase == nil {
		return safeerror.NewByCode(`COULD NOT CREATE OR UPDATE PARENT CASE`)
	}

	err := self.Run(`UPDATE `+caseTable+` SET
				case_parent_id = $1
				WHERE
				id_case IN (`+strings.Join(casesIds, `,`)+`)`, dbParentCase.IdCase)
	for _, caseId := range casesIds {
		delete(self.CreatedCases, caseId)
	}
	if err != nil {
		return err
	}

	return nil
}

func (self *CalculationModel) refreshParentCasesStats() (safeerror.ISafeError) {
	absPath, _ := filepath.Abs("cases/sql/refresh_parent_cases_stats.sql")
	sql, err := ioutil.ReadFile(absPath)
	if err != nil {
		return safeerror.New(err, `ERR_GENERATE_STAT_BILLS_READ_FILE`)
	}
	return self.Run(string(sql))
}
