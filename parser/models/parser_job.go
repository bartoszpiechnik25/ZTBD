package models

import "archive/zip"

type ParseJob struct {
	Year    string
	File    *zip.File
	BaseDir string
}
