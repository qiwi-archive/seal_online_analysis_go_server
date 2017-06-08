DELETE FROM obj_case_tag
WHERE
  case_end_timestamp >= (now() - '1 day' :: INTERVAL)
  AND id_resolution IS NULL
  AND case_status = 'W'
  AND (case_parent_id IS NULL OR (SELECT case_status
                                  FROM obj_case_tag parent_case
                                  WHERE parent_case.id_case = obj_case_tag.case_parent_id) = 'W')
  AND (
    NOT EXISTS(SELECT 1
               FROM (
                      SELECT
                        aggr_bills_tags.id_tag,
                        bills_add_timestamp
                      FROM
                        aggr_bills_tags
                        INNER JOIN normalized_stat_bills_tag
                          ON normalized_stat_bills_tag.id_tag = aggr_bills_tags.id_tag
                             AND normalized_stat_bills_tag.sbills_add_timestamp = aggr_bills_tags.bills_add_timestamp
                        INNER JOIN v_stat_bills_accuracy_rates
                          ON v_stat_bills_accuracy_rates.id_tag = aggr_bills_tags.id_tag AND
                             v_stat_bills_accuracy_rates.add_section = aggr_bills_tags.bills_add_section
                             AND
                             (
                               (
                                 (bills_payments_count - sbills_payments_count) :: FLOAT /
                                 sbills_payments_count :: FLOAT * 100 < -4 * bills_payments_coeff
                                 OR
                                 (bills_count - sbills_all_count) :: FLOAT / sbills_all_count :: FLOAT
                                 * 100 < -4 * bills_count_coeff
                                 OR
                                 (bills_paid_count :: FLOAT / bills_count :: FLOAT * 100 -
                                  sbills_conversion) < -4 * bills_conversion_coeff
                               )
                               OR (
                                 sbills_add_timestamp < (SELECT max(bills_add_timestamp)
                                                         FROM aggr_bills_tags
                                                         WHERE
                                                           aggr_bills_tags.id_tag = v_stat_bills_accuracy_rates.id_tag)
                                 AND (
                                   (bills_payments_count - sbills_payments_count) :: FLOAT /
                                   sbills_payments_count :: FLOAT * 100 < -1 * bills_payments_coeff
                                   OR
                                   (bills_count - sbills_all_count) :: FLOAT / sbills_all_count :: FLOAT
                                   * 100 < -1 * bills_count_coeff
                                   OR
                                   (bills_paid_count :: FLOAT / bills_count :: FLOAT * 100 -
                                    sbills_conversion) < -1 * bills_conversion_coeff
                                 )
                               )
                             )
                      WHERE
                        (bills_add_timestamp >= (now() - '1 DAY ' :: INTERVAL))
                        AND bills_count > 0
                    ) AS subq
               WHERE subq.id_tag = obj_case_tag.id_tag AND
                     subq.bills_add_timestamp BETWEEN obj_case_tag.case_start_timestamp AND obj_case_tag.case_end_timestamp)
    OR (
      EXISTS(SELECT 1
             FROM obj_case_tag child_case
             WHERE case_parent_id = obj_case_tag.id_case
             LIMIT 1)
      AND (SELECT max(case_end_timestamp) - min(case_start_timestamp) AS case_duration
           FROM obj_case_tag child_case
           WHERE case_parent_id = obj_case_tag.id_case
           GROUP BY case_parent_id
          ) < INTERVAL '5 minutes')
  );