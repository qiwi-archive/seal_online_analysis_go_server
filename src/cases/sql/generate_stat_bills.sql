INSERT INTO stat_bills_tag
(
  sbills_all_count,
  sbills_paid_count,
  sbills_all_amount,
  sbills_paid_amount,
  sbills_conversion,
  sbills_payments_count,
  sbills_payments_amount,
  sbills_health,
  sbills_pay_seconds,
  sbills_add_timestamp,
  sbills_add_section,
  sbills_time_code,
  id_tag
)
  SELECT
    round(avg(sums.all_bills_count))         all_bills_count,
    round(avg(sums.paid_bills_count))        paid_bills_count,
    round(avg(sums.all_bills_amount))        all_bills_amount,
    round(avg(sums.paid_bills_amount))       paid_bills_amount,
    round(avg(sums.paid_bills_count)) / round(avg(sums.all_bills_count)) * 100 AS sbills_conversion,
    round(avg(sums.payments_count))          paid_bills_count,
    round(avg(sums.payments_amount))         all_bills_amount,
    round(avg(sums.payments_count)) / round(avg(sums.all_bills_count)) * 100 AS sbills_health,
    round(avg(sums.bills_pay_seconds))       sbills_pay_seconds,
    dates.bills_add_timestamp,
    dates.bills_add_section,
    to_char(dates.bills_add_timestamp, 'dy') sbills_time_code,
    sums.id_tag
  FROM
    (
      SELECT
        date_timestamp :: TIMESTAMP                                                         AS bills_add_timestamp,
        EXTRACT(EPOCH FROM date_timestamp - date_trunc('day', date_timestamp)) / 5 / 60 + 1 AS bills_add_section
      FROM
            generate_series(
                date_trunc('day', now()) - INTERVAL '6 days' - INTERVAL '5 minutes',
                NOW() + INTERVAL '60 minutes',
                INTERVAL '5 minutes'
            ) date_timestamp
    ) dates
    INNER JOIN (
                 SELECT
                   id_tag,
                   bills_add_timestamp,
                   bills_add_section,
                   avg(bills_count)
                   OVER sbills_window AS all_bills_count,
                   avg(bills_paid_count)
                   OVER sbills_window AS paid_bills_count,
                   avg(bills_amount)
                   OVER sbills_window AS all_bills_amount,
                   avg(bills_paid_amount)
                   OVER sbills_window AS paid_bills_amount,
                   avg(bills_payments_count)
                   OVER sbills_window AS payments_count,
                   avg(bills_payments_amount)
                   OVER sbills_window AS payments_amount,
                   avg(bills_pay_seconds)
                   OVER sbills_window AS bills_pay_seconds
                 FROM
                   aggr_bills_tags
                 WHERE bills_add_timestamp < date_trunc('day', now()) - INTERVAL '6 days' + INTERVAL '65 minutes'
                       AND bills_count > 0
                 WINDOW sbills_window AS (
                   PARTITION BY id_tag
                   ORDER BY bills_add_timestamp
                   ROWS BETWEEN 3 PRECEDING AND 3 FOLLOWING )
               ) sums ON
                        to_char(dates.bills_add_timestamp, 'dy') = to_char(sums.bills_add_timestamp, 'dy')
                        AND dates.bills_add_section = sums.bills_add_section
    WHERE sums.all_bills_count > 0 AND sums.paid_bills_count > 0 AND sums.payments_count > 0
  GROUP BY sums.id_tag,
    dates.bills_add_section,
    dates.bills_add_timestamp
ON CONFLICT (id_tag, sbills_add_timestamp)
  DO NOTHING;
INSERT INTO normalized_stat_bills_tag (
  sbills_paid_amount, sbills_paid_count, sbills_all_amount, sbills_all_count, sbills_conversion, sbills_payments_count, sbills_payments_amount, sbills_health, sbills_pay_seconds, sbills_add_timestamp, sbills_add_section, sbills_time_code, id_stat_bills_tag, id_tag
) SELECT
    sbills_paid_amount,
    normalized_paid_count,
    sbills_all_amount,
    normalized_all_count,
    sbills_conversion,
    normalized_payments_count,
    sbills_payments_amount,
    sbills_health,
    sbills_pay_seconds,
    sbills_add_timestamp,
    sbills_add_section,
    sbills_time_code,
    id_stat_bills_tag,
    id_tag
  FROM
    (
      SELECT
        stat_bills_tag.*,
        ROUND(sbills_all_count + AVG(
            COALESCE(
                bills_count,
                sbills_all_count
            ) - sbills_all_count
        )
        OVER sbills_window) AS normalized_all_count,
        ROUND(sbills_paid_count + AVG(
            COALESCE(
                bills_paid_count,
                sbills_paid_count
            ) - sbills_paid_count
        )
        OVER sbills_window) AS normalized_paid_count,
        ROUND(sbills_payments_count + AVG(
            COALESCE(
                bills_payments_count,
                sbills_payments_count
            ) - sbills_payments_count
        )
        OVER sbills_window) AS normalized_payments_count
      FROM
        stat_bills_tag
        LEFT JOIN aggr_bills_tags
          ON aggr_bills_tags.id_tag = stat_bills_tag.id_tag
             AND aggr_bills_tags.bills_add_timestamp = sbills_add_timestamp
             AND bills_add_timestamp < (
          SELECT MAX(bills_add_timestamp)
          FROM
            aggr_bills_tags aggr_current_tag
          WHERE
            aggr_current_tag.id_tag = aggr_bills_tags.id_tag
        )
          AND bills_count > 1 AND bills_paid_count > 1 AND bills_payments_count > 1
             AND NOT EXISTS(
            SELECT 1
            FROM
              obj_case_tag
            WHERE
              obj_case_tag.id_tag = aggr_bills_Tags.id_tag
              AND bills_add_timestamp = case_start_timestamp
        )
      WHERE
        sbills_add_timestamp >= date_trunc('day', now()) - INTERVAL '6 days' - INTERVAL '5 minutes'
      WINDOW sbills_window AS (
        PARTITION BY stat_bills_tag.id_tag
        ORDER BY
          sbills_add_timestamp
        ROWS BETWEEN 3 PRECEDING AND 0 FOLLOWING
      )
    ) AS normalizing_subquery
ON CONFLICT (id_tag, sbills_add_timestamp)
  DO UPDATE SET
    sbills_paid_amount     = EXCLUDED.sbills_paid_amount,
    sbills_paid_count      = EXCLUDED.sbills_paid_count,
    sbills_all_amount      = EXCLUDED.sbills_all_amount,
    sbills_all_count       = EXCLUDED.sbills_all_count,
    sbills_payments_count  = EXCLUDED.sbills_payments_count,
    sbills_payments_amount = EXCLUDED.sbills_payments_amount,
    sbills_conversion      = EXCLUDED.sbills_conversion,
    sbills_health          = EXCLUDED.sbills_health,
    sbills_pay_seconds     = EXCLUDED.sbills_pay_seconds;