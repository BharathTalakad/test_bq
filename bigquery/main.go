package bigquery

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/bigquery"

	"gcp_learn/bigquery/crud"
	"gcp_learn/bigquery/queries"
)

func Run() {
	ctx := context.Background()
	bqClient := getBQClient(ctx)

	year, month, day := time.Now().Date()

	date := strconv.Itoa(year) + month.String() + strconv.Itoa(day)

	table, err := crud.CreateTable(ctx, bqClient, date)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created table: " + date)

	fmt.Println("Executing query ")

	q := bqClient.Query(queries.SimpleGet)
	q.Parameters = []bigquery.QueryParameter{
		{
			Name:  "startDate",
			Value: time.Now().Format("20060102"),
		},
		{
			Name:  "endDate",
			Value: time.Now().AddDate(0, 0, 1).Format("20060102"),
		},
	}

	err = crud.ExecuteQuery(ctx, q, table)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finished executing: ")

	//err = crud.DeleteTable(ctx, bqClient, date)
	//if err != nil {
	//	panic(err)
	//}

	//fmt.Println("Deleted table: " + date)

}
