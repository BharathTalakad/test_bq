package crud

import (
	"context"

	"cloud.google.com/go/bigquery"

	"gcp_learn/global"
)

func DeleteTable(ctx context.Context, client *bigquery.Client, name string) error {

	dataset := client.DatasetInProject(global.ProjectID, global.DatasetID)
	table := dataset.Table(name)

	err := table.Delete(ctx)

	if err != nil {
		panic(err)
	}

	return nil
}
