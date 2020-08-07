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
		{Name: "ID", Required: true, Type: bigquery.IntegerFieldType},
		{Name: "Events", Required: false, Repeated: true, Type: bigquery.StringFieldType},
	}

	if err := table.Create(ctx, &bigquery.TableMetadata{Schema: schema1}); err != nil {
		return nil, err
	}

	return table, nil
}
