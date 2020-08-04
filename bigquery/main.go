package bigquery

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"gcp_learn/bigquery/crud"
)

func Run() {
	ctx := context.Background()
	bqClient := getBQClient(ctx)

	year, month, day := time.Now().Date()

	date := strconv.Itoa(year) + month.String() + strconv.Itoa(day)

	err := crud.CreateTable(ctx, bqClient, date)
	if err != nil {
		panic(err)
	}

	fmt.Println("Created table: " + date)

	err = crud.DeleteTable(ctx, bqClient, date)
	if err != nil {
		panic(err)
	}

	fmt.Println("Deleted table: " + date)

}
