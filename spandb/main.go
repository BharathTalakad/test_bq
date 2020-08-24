package spandb

import (
	"context"
	"fmt"
)

func Run() {
	ctx := context.Background()

	projectID := "kouzoh-p-bharath"
	instanceId := "spdbtest"
	dbID := "test"
	//
	dbname := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceId,dbID)

	//err := crud.CreateDatabase(ctx, "dbname")
	//if err != nil {
	//	panic(err)
	//}

	_,_,err := CreateClients(ctx,dbname)
	if err != nil {
		panic(err)
	}
}
