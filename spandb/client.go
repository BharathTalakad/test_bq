package spandb

import (
	"context"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
)

func CreateClients(ctx context.Context, db string) (admin *database.DatabaseAdminClient, cl *spanner.Client, err error) {
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	dataClient, err := spanner.NewClient(ctx, db)
	if err != nil {
		return nil, nil, err
	}

	defer adminClient.Close()
	defer dataClient.Close()

	return
}
