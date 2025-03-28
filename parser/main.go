package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"
	"ztbd/downloader"
	"ztbd/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	dateRange := []string{}
	for i := beginVal; i <= endVal; i++ {
		if i == 2020 {
			continue
		}
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

	postgres_string := "user=postgres password=extensive_testing host=postgres dbname=postgres port=5432"
	postgres, err := gorm.Open(postgres.Open(postgres_string), &gorm.Config{})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = postgres.AutoMigrate(&models.Investigator{}, &models.Award{}, &models.ProgramReference{}, &models.ProgramElement{}, &models.Institution{}, &models.ProgramOfficer{}, &models.Organization{}, &models.Division{}, &models.Directorate{})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// mysql_string := "testing:extensive_testing@tcp(mysql:3306)/ztbd"
	// ms, err := gorm.Open(mysql.Open(mysql_string), &gorm.Config{})
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	os.Exit(1)
	// }
	//
	// err = postgres.AutoMigrate(&models.Investigator{}, &models.Award{}, &models.ProgramReference{}, &models.ProgramElement{}, &models.Institution{}, &models.ProgramOfficer{}, &models.Organization{}, &models.Division{}, &models.Directorate{})
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	os.Exit(1)
	// }

	clientOptions := options.Client().ApplyURI("mongodb://tester:passwd@mongo7:27017/ztbd")
	mongo7Client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	clientOptions8 := options.Client().ApplyURI("mongodb://tester:passwd@mongo8:27017/ztbd")
	mongo8client, err := mongo.Connect(context.TODO(), clientOptions8)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	for _, year := range dateRange {
		wg.Add(1)
		go func() {
			if cached, _ := d.Cached(year); !cached {
				d.Produce(year, jobChan, postgres, nil, mongo7Client, mongo8client)
			}
			defer wg.Done()
		}()
		time.Sleep(2 * time.Second)
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

func f() {

	postgres_string := "user=postgres password=extensive_testing host=localhost dbname=postgres port=5432"
	postgres, err := gorm.Open(postgres.Open(postgres_string), &gorm.Config{})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	clientOptions := options.Client().ApplyURI("mongodb://tester:passwd@localhost:27017/ztbd")
	mongo7Client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = postgres.AutoMigrate(&models.Investigator{}, &models.Award{}, &models.ProgramReference{}, &models.ProgramElement{}, &models.Institution{}, &models.ProgramOfficer{}, &models.Organization{}, &models.Division{}, &models.Directorate{})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	file, err := os.Open("/Users/bpiechnik/Downloads/2014(1)/1463722.json")
	if err != nil {
		slog.Error(err.Error())
		return
	}
	body, err := io.ReadAll(file)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	var data models.Award
	err = json.Unmarshal(body, &data)
	if err != nil {
		slog.Error(err.Error())
	}
	var po models.ProgramOfficer
	err = json.Unmarshal(body, &po)
	if err != nil {
		slog.Error(err.Error())
	}
	var dir models.Directorate
	err = json.Unmarshal(body, &dir)
	if err != nil {
		slog.Error(err.Error())
	}

	var div models.Division
	err = json.Unmarshal(body, &div)
	if err != nil {
		slog.Error(err.Error())
	}

	var org models.Organization
	err = json.Unmarshal(body, &org)
	if err != nil {
		slog.Error(err.Error())
	}
	org.Directorate = dir
	org.Division = div
	data.Organization = &org
	data.ProgramOfficer = po

	fmt.Printf("%+v\n", data)

	err = postgres.Create(&data).Error
	if err != nil {
		slog.Error(err.Error())
		return
	}

	err = models.InsertMongo(&data, mongo7Client)
	if err != nil {
		slog.Error(err.Error())
		return
	}

}

func main() {
	slog.Info("Parser starting")

	workToDo := make(chan *models.ParseJob, 100)
	wg := &sync.WaitGroup{}

	worker(workToDo, wg)
	producer(workToDo)

	wg.Wait()
	// f()
	slog.Info("App ending")
}
