package emulator

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
)

func SetupAndRun(ctx context.Context) (err error) {
	// Check and verify firestore env is set
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		_ = os.Setenv("SPANNER_EMULATOR_HOST", "0.0.0.0:9010")
	}

	projectID := "emu"
	instanceId := "testIns"
	dbID := "dummy"

	dbname := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceId, dbID)

	err = createInstance(projectID, instanceId)
	if err != nil {
		panic(err)
	}

	err = createDatabase(dbname)
	if err != nil {
		panic(err)
	}

	err = writeUsingDML(dbname)
	if err != nil {
		panic(err)
	}

	err = query(dbname)
	if err != nil {
		panic(err)
	}

	err = deleteInstance(projectID, instanceId)
	if err != nil {
		return err
	}

	return nil

}

func createInstance(projectID, instanceID string) error {
	ctx := context.Background()

	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer instanceAdmin.Close()

	op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instancepb.Instance{
			Config:      fmt.Sprintf("projects/%s/instanceConfigs/%s", projectID, "regional-us-central1"),
			DisplayName: instanceID,
			NodeCount:   1,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create instance %s: %v", fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID), err)
	}
	// Wait for the instance creation to finish.
	i, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance creation to finish failed: %v", err)
	}
	// The instance may not be ready to serve yet.
	if i.State != instancepb.Instance_READY {
		fmt.Printf("instance state is not READY yet. Got state %v\n", i.State)
	}
	fmt.Printf("Created instance [%s]\n", instanceID)
	return nil
}

func createDatabase(db string) error {
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", db)
	}

	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          matches[1],
		CreateStatement: "CREATE DATABASE `" + matches[2] + "`",
		ExtraStatements: []string{
			`CREATE TABLE Singers (
                                SingerId   INT64 NOT NULL,
                                FirstName  STRING(1024),
                                LastName   STRING(1024),
                                SingerInfo BYTES(MAX)
                        ) PRIMARY KEY (SingerId)`,
		},
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Printf("Created database [%s]\n", db)
	return nil
}

func createClients(w io.Writer, db string) error {
	ctx := context.Background()

	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}

	dataClient, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}

	_ = adminClient
	_ = dataClient

	return nil
}

func writeUsingDML(db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT Singers (SingerId, FirstName, LastName) VALUES
                                (12, 'Melissa', 'Garcia'),
                                (13, 'Russell', 'Morales'),
                                (14, 'Jacqueline', 'Long'),
                                (15, 'Dylan', 'Shaw')`,
		}
		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Printf("%d record(s) inserted.\n", rowCount)
		return err
	})
	return err
}

func query(db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Printf("Reading now....\n")

	stmt := spanner.Statement{SQL: `SELECT SingerId, FirstName FROM Singers`}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID int64
		var firstName string
		if err := row.Columns(&singerID, &firstName); err != nil {
			return err
		}
		fmt.Printf("%d %s\n", singerID, firstName)
	}
}

func deleteInstance(projectID, instanceID string) error {
	ctx := context.Background()

	adminClient, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}

	instanceName := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)
	err = adminClient.DeleteInstance(ctx, &instancepb.DeleteInstanceRequest{Name: instanceName})
	if err != nil {
		return err
	}

	fmt.Printf("sucessfully deleted instance : %s ", instanceName)

	return nil
}
