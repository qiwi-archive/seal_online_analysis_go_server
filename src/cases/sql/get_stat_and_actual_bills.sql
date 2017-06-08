SELECT
  aggr_bills_tags.id_tag,
  (bills_count - sbills_all_count) :: FLOAT / sbills_all_count :: FLOAT * 100 AS bills_count_deviation,
  bills_count_coeff,
  (bills_payments_count - sbills_payments_count) :: FLOAT / sbills_payments_count :: FLOAT *
  100                                                                         AS bills_payments_count_deviation,
  bills_payments_coeff,
  CASE WHEN bills_count = 0 THEN 100 ELSE
  (bills_paid_count :: FLOAT / bills_count :: FLOAT * 100 - sbills_conversion)
  END AS bills_conversion_deviation,
  bills_conversion_coeff,
  ROUND(sbills_paid_amount) :: INT8                                              sbills_paid_amount,
  ROUND(sbills_paid_count) :: INT8                                               sbills_paid_count,
  ROUND(sbills_all_amount) :: INT8                                               sbills_all_amount,
  ROUND(sbills_all_count) :: INT8                                                sbills_all_count,
  ROUND(sbills_conversion) :: INT8                                               sbills_conversion,
  ROUND(
      sbills_payments_count) :: INT8                                             sbills_payments_count,
  ROUND(
      sbills_payments_amount) :: INT8                                            sbills_payments_amount,
  ROUND(sbills_health) :: INT8                                                   sbills_health,
  ROUND(sbills_pay_seconds) :: INT8                                              sbills_pay_seconds,
  ROUND(bills_paid_amount) :: INT8                                               paid_bills_amount,
  ROUND(bills_paid_count) :: INT8                                                paid_bills_count,
  ROUND(bills_amount) :: INT8                                                    bills_amount,
  ROUND(bills_count) :: INT8                                                     bills_count,
  ROUND(bills_paid_count) :: INT8                                                bills_paid_count,
  ROUND(bills_paid_count) :: INT8                                                bills_payments_count,
  CASE WHEN bills_count = 0 THEN 0 ELSE
  ROUND(bills_paid_count :: FLOAT / bills_count :: FLOAT * 100)
  END AS bills_conversion,
  TO_CHAR(bills_add_timestamp, 'DD-MM-YYYY HH24:MI:SS')                         bills_add_date,
  EXTRACT(EPOCH FROM bills_add_timestamp) :: INTEGER                          AS bills_add_timestamp
FROM
  aggr_bills_tags
  INNER JOIN normalized_stat_bills_tag
    ON normalized_stat_bills_tag.id_tag = aggr_bills_tags.id_tag
       AND normalized_stat_bills_tag.sbills_add_timestamp = aggr_bills_tags.bills_add_timestamp
  INNER JOIN v_stat_bills_accuracy_rates ON v_stat_bills_accuracy_rates.id_tag = aggr_bills_tags.id_tag
                                            AND v_stat_bills_accuracy_rates.add_section = aggr_bills_tags.bills_add_section
                                            AND
                                            (
                                            (bills_count = 0 OR bills_paid_count = 0 OR bills_payments_count = 0)
                                            OR
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
                                                sbills_add_timestamp < (select max(bills_add_timestamp) from aggr_bills_tags where aggr_bills_tags.id_tag = v_stat_bills_accuracy_rates.id_tag)
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
  AND aggr_bills_tags.id_tag <> 90
ORDER BY aggr_bills_tags.id_tag, aggr_bills_tags.bills_add_timestamp