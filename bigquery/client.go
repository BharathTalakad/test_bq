package bigquery

import (
	"context"
	"log"

	"cloud.google.com/go/bigquery"
	"gcp_learn/global"
	"google.golang.org/api/option"
)

// getBQClient returns the BigQuery client using credentials in key folder
func getBQClient(ctx context.Context) *bigquery.Client {

	client, err := bigquery.NewClient(ctx, global.ProjectID, option.WithCredentialsFile(global.CredPath))
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}

	return client
}
