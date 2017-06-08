SELECT id_case,
			id_resolution,
			case_status,
			id_tag,
			case_amount,
			case_stat_amount,
			case_lost_amount,
			case_count,
			case_stat_count,
			case_code,
			EXTRACT(epoch FROM(case_start_timestamp))::INTEGER AS case_start_timestamp,
			EXTRACT(epoch FROM(case_end_timestamp))::INTEGER AS case_end_timestamp,
			case_start_dtime,
			case_end_dtime,
			(1+EXTRACT(epoch FROM(case_end_timestamp-case_start_timestamp))/288)::INTEGER as case_duration
			FROM obj_case_tag
			WHERE case_end_timestamp >= date_trunc('day', now())::timestamp - interval '12 hours'
			AND case_parent_id IS NULL
			ORDER BY case_duration DESC, id_tag, case_end_timestamp ASC