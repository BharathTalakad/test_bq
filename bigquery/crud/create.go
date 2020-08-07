package crud

import (
	"context"
	"gcp_learn/global"

	"cloud.google.com/go/bigquery"
)

func CreateTable(ctx context.Context, client *bigquery.Client, name string) (*bigquery.Table, error) {

	dataset := client.DatasetInProject(global.ProjectID, global.DatasetID)
	table := dataset.Table(name)

	schema1 := bigquery.Schema{
		{Name: "user_id", Required: true, Type: bigquery.IntegerFieldType},
		{Name: "created", Required: true, Type: bigquery.TimestampFieldType},
		{Name: "buy_count", Required: true, Type: bigquery.IntegerFieldType},
		{Name: "light_user", Required: true, Type: bigquery.BooleanFieldType},
		{Name: "first_event", Required: true, Type: bigquery.TimestampFieldType},
	}

	if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema1}); err != nil {
		return nil, err
	}

	return table, nil
}
