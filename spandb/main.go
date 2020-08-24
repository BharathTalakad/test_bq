package spandb

import (
	"context"
	"gcp_learn/spandb/crud"
)

func Run() {
	ctx := context.Background()

	dbname := "dummy"

	err := crud.CreateDatabase(ctx, dbname)
	if err != nil {
		panic(err)
	}
}
