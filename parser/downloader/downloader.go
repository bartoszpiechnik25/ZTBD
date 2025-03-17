package downloader

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/melbahja/got"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
	"ztbd/models"
)

const BASE_URL = "https://www.nsf.gov/awardsearch/download?DownloadFileName=%s&All=true&isJson=true"

type Downloader struct {
	BaseDir string
}

func New(baseDir string) *Downloader {
	return &Downloader{
		BaseDir: baseDir,
	}
}

func (d *Downloader) Cached(year string) (bool, error) {
	files, err := filepath.Glob(filepath.Join(d.BaseDir, year) + "/*")
	if err != nil {
		return false, fmt.Errorf("error occurred scanning directory")
	}

	if len(files) > 0 {
		slog.Info("directories are cached!")
		return true, nil
	}
	return false, nil

}

func (d *Downloader) Download(year string) (*bytes.Reader, int64, error) {
	dirPath := filepath.Join(d.BaseDir, year)
	err := os.MkdirAll(dirPath, os.ModePerm)
	cleaner := func(success bool) {
		if !success {
			os.RemoveAll(dirPath)
		}
	}
	if err != nil {
		slog.Error(fmt.Sprintf("could not creat directory for: %s due to %v", dirPath, err.Error()))
		cleaner(false)
		return nil, 0, err
	}
	tmpZip := filepath.Join(dirPath, year+"_tmp.zip")

	g := got.New()
	err = g.Download(fmt.Sprintf(BASE_URL, year), tmpZip)
	if err != nil {
		slog.Error(fmt.Sprintf("could not get data for year: %s due to : %v", year, err))
		cleaner(false)
		return nil, 0, err
	}
	file, err := os.Open(tmpZip)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not open temporary zip file: %s", tmpZip))
		cleaner(false)
		return nil, 0, err
	}
	defer os.RemoveAll(tmpZip)
	zipData, _ := io.ReadAll(file)
	byteReader := bytes.NewReader(zipData)
	return byteReader, int64(len(zipData)), nil
}

func (d *Downloader) Produce(year string, jobChan chan<- *models.ParseJob) {

	slog.Info(fmt.Sprintf("Starting pipeline for: %s", year))
	start := time.Now()
	bytes, length, err := d.Download(year)
	end := time.Now()
	slog.Info(fmt.Sprintf("Downloading %s took %f seconds", year, end.Sub(start).Seconds()))
	if err != nil {
		slog.Error(fmt.Sprintf("could not download zip due to: %v", err.Error()))
		return
	}

	zipReader, err := zip.NewReader(bytes, length)
	for _, zipFile := range zipReader.File {
		jobChan <- &models.ParseJob{
			Year:    year,
			File:    zipFile,
			BaseDir: d.BaseDir,
		}
	}
}

func Consume(job *models.ParseJob) {

	path := filepath.Join(job.BaseDir, job.Year, job.File.Name)
	outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, job.File.Mode())
	if err != nil {
		return
	}

	rc, err := job.File.Open()
	if err != nil {
		outFile.Close()
		return
	}

	_, err = io.Copy(outFile, rc)

	outFile.Close()
	rc.Close()

	if err != nil {
		return
	}
}
