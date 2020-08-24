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
)

func SetupAndRun(ctx context.Context) (err error ) {
	// Check and verify firestore env is set
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		_ = os.Setenv("SPANNER_EMULATOR_HOST", "0.0.0.0:9010")
	}

	projectID := "emu"
	instanceId := "test-instance"
	dbID := "dummy"

	dbname := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceId,dbID)

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


	return nil

}

func createDatabase( db string) error {
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
			`CREATE TABLE Albums (
                                SingerId     INT64 NOT NULL,
                                AlbumId      INT64 NOT NULL,
                                AlbumTitle   STRING(MAX)
                        ) PRIMARY KEY (SingerId, AlbumId),
                        INTERLEAVE IN PARENT Singers ON DELETE CASCADE`,
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

	stmt := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
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
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Printf("%d %d %s\n", singerID, albumID, albumTitle)
	}
}