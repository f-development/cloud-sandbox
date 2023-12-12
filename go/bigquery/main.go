package main

import (
	"context"
	"encoding/csv"
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
)

type Row struct {
	Recipient     string `bigquery:"recipient"`
	RecipientHash string `bigquery:"recipient_hash"`
	CreatedAt     string `bigquery:"created_at"`
}

const (
	dataset           = "flicspy_stg_1_ms_content"
	gcpProject        = "f-development"
	partitionCount    = 4000
	recipientCount    = 4000 * 1000
	rowsPerRecipients = 1

	totalRows              = recipientCount * rowsPerRecipients
	partitionSize          = totalRows / partitionCount
	recipientsPerPartition = 1000

	threadCount = 100
	batchSize   = 10000
	iteration   = totalRows / batchSize / threadCount
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

func writeCsv() {
	records := [][]string{
		{"first_name", "last_name", "username"},
		{"Rob", "Pike", "rob"},
		{"Ken", "Thompson", "ken"},
		{"Robert", "Griesemer", "gri"},
	}

	w := csv.NewWriter(os.Stdout)
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}

}

func main() {
	cli, err := bigquery.NewClient(context.Background(), gcpProject)

	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	for t := 0; t < threadCount; t++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, t int) {
			defer wg.Done()
			for i := 0; i < iteration; i++ {
				arr := []Row{}
				p1 := time.Now()
				for j := 0; j < batchSize; j++ {
					r := rand.Int()
					recipient := r % recipientCount
					arr = append(arr, Row{
						Recipient:     strconv.Itoa(recipient),
						RecipientHash: strconv.Itoa(recipient % partitionCount),
						CreatedAt:     strconv.Itoa(rand.Int()),
					})
				}
				doInsert(cli, arr)
				p2 := time.Now()
				slog.Info("Batch finished", "thread", t, "iter", i, "throughput", batchSize/p2.Sub(p1).Seconds())
				if err != nil {
					slog.Error("err", "err", err)
				}
			}
		}(&wg, t)
	}
	wg.Wait()
}
