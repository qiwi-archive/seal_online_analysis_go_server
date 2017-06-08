package cases

import (
	"seal_online_go_server/src/innotrio/router"
	"seal_online_go_server/src/data"
	"net/http"
	"seal_online_go_server/src/config"
	"net/url"
	"github.com/Sirupsen/logrus"
	"seal_online_go_server/src/internal"
	"time"
	"seal_online_go_server/src/innotrio/error"
)

type Ctrl struct {
	CalculationModel *CalculationModel
	InternalModel    *internal.Model
	*router.ApiRouter
}

func (self *Ctrl) Init() {
	self.GET("cases/refresh5m", self.refresh5mCases)
}

func (self *Ctrl) Refresh5mCases() safeerror.ISafeError {
	hasReplicationTroubles, err := self.checkAndInformReplicationTroubles()
	if hasReplicationTroubles {
		return err
	}
	err = self.InternalModel.UpdateTags([]string{data.TAG_TYPE_FLOW_CODE});
	if err != nil {
		logrus.Warn("FAIL Refresh Flows Tags", err)
		return err
	}
	logrus.Info("Done Refresh Flows Tags")

	createdCases, err := self.CalculationModel.RefreshCases()
	if err != nil {
		logrus.Warn("FAIL Refresh Cases", err)
		return err
	}
	var casesIds []string
	for caseId := range createdCases {
		casesIds = append(casesIds, caseId)
	}
	err = self.informNewCases(casesIds)
	return err
}

/**
Проверка наличия проблем с репликацией. Есть ли смысл в онлайн мониторинге без актуальных данных?
 */
func (self *Ctrl) checkAndInformReplicationTroubles() (bool, safeerror.ISafeError) {
	var hasMainTroubles bool = false

	var lastReplicationLogs []data.ReplicationLog
	var err safeerror.ISafeError
	lastReplicationLogs, err = self.InternalModel.GetLastReplicationLogs([]string{
		data.REPLICATION_QW_MINUTES_MAIN_CODE,
		data.REPLICATION_QW_MINUTES_EXTRAS_MAIN_CODE,
	})
	if err != nil {
		logrus.Warn("Replication Troubles Check Error!", err)
		return hasMainTroubles, err
	}
	var replicationTroubles []data.ReplicationLog
	if len(lastReplicationLogs) < 2 {
		replicationTroubles = lastReplicationLogs
	} else {
		for _, replicationLog := range lastReplicationLogs {
			if time.Since(*replicationLog.MaxDate).Minutes() > 6 {
				replicationTroubles = append(replicationTroubles, replicationLog)
			}
		}
	}
	if len(lastReplicationLogs) < 2 || len(replicationTroubles) > 0 {
		var replicationTroublesIds []string
		for _, replicationTrouble := range replicationTroubles {
			if replicationTrouble.Code == data.REPLICATION_QW_MINUTES_MAIN_CODE || replicationTrouble.Code == data.REPLICATION_QW_MINUTES_EXTRAS_MAIN_CODE {
				hasMainTroubles = true
			}
			replicationTroublesIds = append(replicationTroublesIds, replicationTrouble.Id)
		}
		_, error := http.PostForm(config.API_HOST+config.API_BOT_PREFIX+`/users/inform_adm_replication_troubles`,
			url.Values{"replication_log_ids": replicationTroublesIds, "key": {config.ORACLE_API_REQUEST_CODE}})
		if error != nil {
			logrus.Warn("FAIL Inform", error)
			return hasMainTroubles, err
		}
		logrus.Info("Done inform replication troubles")
	}
	return hasMainTroubles, err
}

func (self *Ctrl) informNewCases(casesIds []string) safeerror.ISafeError {
	if len(casesIds) > 0 {
		response, error := http.PostForm(config.API_HOST+config.API_BOT_PREFIX+`/users/inform_dadm_cases`,
			url.Values{"cases_ids": casesIds, "key": {config.ORACLE_API_REQUEST_CODE}})
		if error != nil || response.StatusCode != 200 {
			logrus.Warn("FAIL Inform", error, response.StatusCode)
			return safeerror.New(error, `INFORM_NEW_CASES`)
		}
		logrus.Info("Done inform new cases")
	}
	return nil
}

func (self *Ctrl) refresh5mCases(c *router.Context) {
	err := self.Refresh5mCases()
	self.Send(c, true, err)
}
