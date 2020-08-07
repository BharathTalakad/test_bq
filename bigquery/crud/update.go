package crud

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
)

func ExecuteQuery(ctx context.Context, q *bigquery.Query, t *bigquery.Table) error {

	q.Location = "US" // Location must match the dataset(s) referenced in query.
	q.QueryConfig.Dst = t
	// Run the query and print results when the query job is completed.
	fmt.Println("Executing....")

	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
		return err
	}

	return nil
}
