package queries

const (
	pascalEventLog = "`mercari-bigquery-jp-prod.pascal_event_log.event_log_*`"
	localTable     = "`kouzoh-p-bharath.achievements.buy_complete_events`"

	Get = `
SELECT
  user_id,
  ARRAY_AGG(event_id) AS agg
FROM (
  SELECT
    user_id,
    event_id,
    COUNT(event_id) AS count
  FROM ` + pascalEventLog +
		`WHERE _table_suffix between @startDate and @endDate ` +
		` GROUP BY
    1,
    2 )
GROUP BY
  user_id`

	SimpleGet = `SELECT * FROM ` + localTable + `LIMIT 100`
)
