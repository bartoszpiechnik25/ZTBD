package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"ztbd/downloader"
	"ztbd/models"
)

const N_WORKERS = 64

func parseEnv() ([]string, error) {

	begin := os.Getenv("BEGIN_YEAR")
	end := os.Getenv("END_YEAR")
	beginVal, err := strconv.Atoi(begin)
	if err != nil {
		slog.Error("could not convert 'BEGIN_YEAR' to int")
		return nil, err
	}
	endVal, err := strconv.Atoi(end)
	if err != nil {
		slog.Error("could not convert 'END_VAL' to int")
		return nil, err
	}
	dateRange := make([]string, endVal-beginVal)
	for i := beginVal; i <= endVal; i++ {
		dateRange = append(dateRange, strconv.Itoa(i))
	}
	return dateRange, nil
}

func producer(jobChan chan<- *models.ParseJob) {
	slog.Info("Producer starting")

	dateRange, err := parseEnv()
	if err != nil {
		slog.Error(fmt.Sprintf("could not parse env due to %v", err.Error()))
		os.Exit(1)
	}

	volumeDirectory := os.Getenv("VOLUME_DIR")
	d := downloader.New(volumeDirectory)

	wg := sync.WaitGroup{}

	for _, year := range dateRange {
		wg.Add(1)
		go func() {
			if cached, _ := d.Cached(year); !cached {
				d.Produce(year, jobChan)
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
	close(jobChan)
}

func worker(jobChan <-chan *models.ParseJob, waitGroup *sync.WaitGroup) {
	for i := range N_WORKERS {
		slog.Debug(fmt.Sprintf("Worker %d starting", i))
		waitGroup.Add(1)
		go func() {
			for job := range jobChan {
				downloader.Consume(job)
			}
			defer waitGroup.Done()
		}()
	}

}

func main() {
	slog.Info("Parser starting")

	workToDo := make(chan *models.ParseJob, 100)
	wg := &sync.WaitGroup{}

	worker(workToDo, wg)
	producer(workToDo)

	wg.Wait()

	slog.Info("App ending")
}
