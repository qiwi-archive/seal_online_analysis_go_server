package internal

import (
	"seal_online_go_server/src/innotrio"
	"seal_online_go_server/src/innotrio/error"
	"seal_online_go_server/src/data"
	"strings"
	"strconv"
	"seal_online_go_server/src/config"
	"github.com/gin-gonic/gin"
)

type Model struct {
	*innotrio.Model
}

func (self *Model) RefreshProviders() (safeerror.ISafeError) {
	err := self.Run("UPDATE obj_provider SET prv_add_timestamp = to_timestamp( prv_add_date, 'YYYY-MM-DD\"T\"HH24:MI:SSZ' ), prv_first_txn_timestamp = to_timestamp( prv_first_txn_date, 'YYYY-MM-DD\"T\"HH24:MI:SSZ' ), prv_last_txn_timestamp = to_timestamp( prv_last_txn_date, 'YYYY-MM-DD\"T\"HH24:MI:SSZ' ), prv_contract_timestamp = to_timestamp( contract_date_from, 'YYYY-MM-DD\"T\"HH24:MI:SSZ' )")
	return err
}

func (self *Model) GetTags(tagsType string) ([]data.HiddenTag, safeerror.ISafeError) {
	var items []data.HiddenTag
	err := self.Select(&items, "SELECT id_tag, tag_code, tag_sql, tag_db, tag_type_code FROM spr_tags WHERE tag_type_code = $1 ORDER BY tag_run_order, tag_order, id_tag", tagsType)
	return items, err
}
func (self *Model) GetLastReplicationLogs(replicationCodes []string) ([]data.ReplicationLog, safeerror.ISafeError) {
	var items []data.ReplicationLog
	var replicationCodesFilter string
	if len(replicationCodes) == 0 {
		return items, nil
	}
	for _, value := range replicationCodes {
		if replicationCodesFilter == `` {
			replicationCodesFilter = `replication_code IN (`
		} else {
			replicationCodesFilter += `,`
		}
		replicationCodesFilter += `'` + value + `'`
	}
	replicationCodesFilter += `)`
	err := self.Select(&items, `
	SELECT
	id_replication_log,
	replication_code,
	replication_start_date,
	replication_end_date,
	replication_max_date,
	EXTRACT(epoch from replication_start_date - replication_max_date)::integer as replication_lag,
	replication_days_interval
FROM
    sys_replication_log
WHERE
    `+replicationCodesFilter+`
ORDER BY replication_start_date DESC LIMIT `+strconv.Itoa(len(replicationCodes)))
	return items, err
}

func (self *Model) UpdateTags(tagsTypes []string) (safeerror.ISafeError) {
	tracker := innotrio.NewTracker()
	defer tracker.Log(`DONE UPDATE TAGS`)
	for _, tagsType := range tagsTypes {
		tags, err := self.GetTags(tagsType)
		if err != nil {
			return err
		}
		if err = self.clearTags(tagsType); err != nil {
			return err
		}

		for _, tag := range tags {
			if err := self.updateTag(tag); err != nil {
				return err
			}
		}

		if err := self.refreshTagStat(tagsType); err != nil {
			return err
		}
	}
	return nil
}

func (self *Model) updateTag(tag data.HiddenTag) (safeerror.ISafeError) {
	if config.GIN_RELEASE_MODE == gin.DebugMode {
		println(tag.Code.String)
		tracker := innotrio.NewTracker()
		defer func() {
			tracker.Log(tag.Code.String)
		}()
	}

	if tag.Sql.Valid && tag.Sql.String != "" {
		switch tag.TypeCode.String {
		case `P`:
			//Не заполнена база, значит PG
			if !tag.Db.Valid || tag.Db.String == "" {
				tagSql := strings.Replace(tag.Sql.String, "$1", strconv.Itoa(tag.Id), 1)
				sql := "INSERT INTO rel_prv_tags (id_prv, id_tag) " + tagSql
				return self.Run(sql)
			} else {
				println("UNKNOWN TAG", tag.Id, tag.Code.String, tag.Db.String)
			}
			break;
		case `F`:
			tagSql := strings.Replace(tag.Sql.String, "$1", strconv.Itoa(tag.Id), 1)
			sql := "INSERT INTO rel_flow_tags (id_pay_flow, pflow_code, id_tag) " + tagSql
			if err := self.Run(sql); err != nil {
				return err
			}
			break;
		}
	}

	return nil
}

func (self *Model) clearTags(tagsType string) (safeerror.ISafeError) {
	var err safeerror.ISafeError
	switch tagsType {
	case data.TAG_TYPE_PROVIDER_CODE:
		err = self.Run("DELETE FROM rel_prv_tags WHERE id_tag NOT IN ( SELECT id_tag FROM spr_tags WHERE tag_code = 'P' AND (tag_sql IS NULL OR tag_sql = '' ))")
		break
	case data.TAG_TYPE_FLOW_CODE:
		err = self.Run("DELETE FROM rel_flow_tags WHERE id_tag NOT IN ( SELECT id_tag FROM spr_tags WHERE tag_code = 'F' AND (tag_sql IS NULL OR tag_sql = '' ))")
	default:
		err = safeerror.NewByCode("WRONG_TAG_TYPE_CODE_ERROR")
	}
	if err == nil {
		err = self.Run("UPDATE spr_tags SET tag_count = 0 WHERE tag_type_code = $1", tagsType);
	}
	return err
}

func (self *Model) refreshTagStat(tagsType string) (safeerror.ISafeError) {
	var err safeerror.ISafeError
	switch tagsType {
	case data.TAG_TYPE_PROVIDER_CODE:
		err = self.Run("UPDATE spr_tags SET tag_count = aggr.tag_count FROM  (SELECT id_tag, COUNT(id_prv) as tag_count FROM rel_prv_tags GROUP BY id_tag) as aggr WHERE spr_tags.id_tag = aggr.id_tag")
		if err != nil {
			return err
		}
		err = self.Run("REFRESH MATERIALIZED VIEW v_providers")
		break
	case data.TAG_TYPE_FLOW_CODE:
		err = self.Run("UPDATE spr_tags SET tag_count = aggr.tag_count FROM  (SELECT id_tag, COUNT(id_pay_flow) as tag_count FROM rel_flow_tags GROUP BY id_tag) as aggr	WHERE spr_tags.id_tag = aggr.id_tag")
		if err != nil {
			return err
		}
		err = self.Run("REFRESH MATERIALIZED VIEW v_pay_flow_tag")
	default:
		err = safeerror.NewByCode("WRONG_TAG_TYPE_CODE_ERROR")
	}
	return err
}
