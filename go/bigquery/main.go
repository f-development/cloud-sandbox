package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
)

var (
	slogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	tables = []string{"perf_test", "perf_test_partition", "perf_test_clustering", "perf_test_partition_clustering"}
	// tables = []string{"perf_test_partition_clustering"}
)

type Row struct {
	Recipient     string `bigquery:"recipient"`
	RecipientHash string `bigquery:"recipient_hash"`
	CreatedAt     string `bigquery:"created_at"`
}

const (
	dataset               = "flicspy_stg_1_ms_content"
	gcpProject            = "f-development"
	partitionCount        = 10
	recipientCount        = 1000 * 100
	rowsPerRecipients     = 1000

	totalRows = recipientCount * rowsPerRecipients

	batchSize   = 100000
	iteration   = totalRows / batchSize
)

func doInsert(cli *bigquery.Client, rows []Row) error {
	for _, table := range tables {
		ins := cli.Dataset(dataset).Table(table).Inserter()

		err := ins.Put(context.Background(), rows)
		if err != nil {
			return err
		}

	}
	return nil
}

func importCSVFromFile(projectID, datasetID, tableID, filename string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	// f, err := os.Open(filename)
	// if err != nil {
	// 	return err
	// }
	// source := bigquery.NewReaderSource(f)
	source := bigquery.NewGCSReference("gs://test-tokyo/temp.csv")
	// source.AutoDetect = true   // Allow BigQuery to determine schema.
	source.SkipLeadingRows = 1 // CSV has a single header line.

	loader := client.Dataset(datasetID).Table(tableID).LoaderFrom(source)

	slog.Info("Start job")
	job, err := loader.Run(ctx)
	if err != nil {
		return err
	}
	slog.Info("Wait job")
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	slog.Info("Wait job finished")
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

func writeCsv(file string) {
	f, err := os.Create(file)
	f.Truncate(0)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	w.Write([]string{"recipient", "recipient_hash", "created_at"})
	for i := 0; i < iteration; i++ {
		arr := [][]string{}
		p1 := time.Now()
		for j := 0; j < batchSize; j++ {
			r := rand.Int()
			recipient := r % recipientCount
			arr = append(arr, []string{strconv.Itoa(recipient), strconv.Itoa(recipient % partitionCount), strconv.Itoa(rand.Int())})

		}

		w.WriteAll(arr)

		p2 := time.Now()
		slog.Info("Batch finished", "iter", i, "throughput", batchSize/p2.Sub(p1).Seconds())
		if err != nil {
			slog.Error("err", "err", err)
		}
	}

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

}

func loadTable(wg *sync.WaitGroup, table string) error {
	defer wg.Done()
	slog.Info("Load job", "table", table)
	err := importCSVFromFile(gcpProject, dataset, table, "temp.csv")
	if err != nil {
		slog.Error("Load job error", "err", err)
		return err
	}
	return nil
}

func loadJob() {
	var wg sync.WaitGroup
	for _, table := range tables {
		wg.Add(1)
		go loadTable(&wg, table)
	}
	wg.Wait()
}

func main() {
	// writeCsv("temp.csv")
	loadJob()
}

