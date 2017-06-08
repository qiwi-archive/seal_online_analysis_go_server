WITH sums AS (
    SELECT
      min(case_start_timestamp)     AS case_start_timestamp,
      max(case_end_timestamp)       AS case_end_timestamp,
      SUM(case_lost_amount)         AS lost_amount_sum,
      SUM(case_amount)              AS amount_sum,
      SUM(case_stat_amount)         AS stat_amount_sum,
      SUM(case_count)               AS count_sum,
      SUM(case_stat_count)          AS stat_count_sum,
      SUM(case_paid_count)          AS paid_count_sum,
      SUM(case_stat_paid_count)     AS stat_paid_count_sum,
      SUM(case_payments_count)      AS payments_count_sum,
      SUM(case_stat_payments_count) AS stat_payments_count_sum,
      case_parent_id
    FROM
      obj_case_tag
    WHERE
      case_parent_id IS NOT NULL
    GROUP BY
      case_parent_id
) UPDATE
  obj_case_tag
SET
  case_start_timestamp     = sums.case_start_timestamp,
  case_start_dtime         = TO_CHAR(sums.case_start_timestamp, 'DD-MM-YYYY HH24:MI:SS'),
  case_end_timestamp       = sums.case_end_timestamp,
  case_end_dtime           = TO_CHAR(sums.case_end_timestamp, 'DD-MM-YYYY HH24:MI:SS'),
  case_lost_amount         = sums.lost_amount_sum,
  case_amount              = sums.amount_sum,
  case_stat_amount         = sums.stat_amount_sum,
  case_count               = sums.count_sum,
  case_stat_count          = sums.stat_count_sum,
  case_paid_count          = sums.paid_count_sum,
  case_stat_paid_count     = sums.stat_paid_count_sum,
  case_payments_count      = sums.payments_count_sum,
  case_stat_payments_count = sums.stat_payments_count_sum
FROM
  sums
WHERE
  sums.case_parent_id = obj_case_tag.id_case
  AND obj_case_tag.case_parent_id IS NULL
  AND obj_case_tag.case_end_timestamp > obj_case_tag.case_start_timestamp
RETURNING *;