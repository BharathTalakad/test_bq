package spandb

import (
	"context"
	"fmt"

	"gcp_learn/spandb/emulator"
)

func Run() {
	runEmu()
}

func runReal() {
	ctx := context.Background()

	projectID := "kouzoh-p-bharath"
	instanceId := "spdbtest"
	dbID := "test"
	//
	dbname := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceId, dbID)

	//err := crud.CreateDatabase(ctx, "dbname")
	//if err != nil {
	//	panic(err)
	//}

	_, _, err := CreateClients(ctx, dbname)
	if err != nil {
		panic(err)
	}
}

func runEmu() {
	ctx := context.Background()
	err := emulator.SetupAndRun(ctx)
	if err != nil {
		panic(err)
	}
}
