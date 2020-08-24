package crud

import (
	"context"
	"fmt"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

func CreateDatabase(ctx context.Context, _ string) error {
	projectID := "kouzoh-p-bharath"
	instanceId := "spdbtest"
	dbID := "test"

	dbname := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectID, instanceId,dbID)

	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(dbname)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("Invalid database id %s", dbname)
	}

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
	fmt.Printf( "Created database [%s]\n", dbname)
	return nil
}
