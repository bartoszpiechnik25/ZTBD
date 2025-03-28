package models

import (
	"archive/zip"
	"context"
	"encoding/json"
	"io"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type ParseJob struct {
	Year       string
	File       *zip.File
	BaseDir    string
	PostgresDB *gorm.DB
	MySqlDB    *gorm.DB
	Mongo7     *mongo.Client
	Mongo8     *mongo.Client
}

func (j *ParseJob) ParseJson() (*Award, []byte, error) {

	rc, err := j.File.Open()
	if err != nil {
		return nil, nil, err
	}
	body, err := io.ReadAll(rc)
	if err != nil {
		return nil, body, err
	}

	var data Award
	err = json.Unmarshal(body, &data)
	if err != nil {
		slog.Error(err.Error())
		return nil, body, err
	}
	var po ProgramOfficer
	err = json.Unmarshal(body, &po)
	if err != nil {
		slog.Error(err.Error())
		return nil, body, err
	}
	var dir Directorate
	err = json.Unmarshal(body, &dir)
	if err != nil {
		slog.Error(err.Error())
		return nil, body, err
	}

	var div Division
	err = json.Unmarshal(body, &div)
	if err != nil {
		slog.Error(err.Error())
		return nil, body, err
	}

	var org Organization
	err = json.Unmarshal(body, &org)
	if err != nil {
		slog.Error(err.Error())
		return nil, body, err
	}
	org.Directorate = dir
	org.Division = div
	data.Organization = &org
	data.ProgramOfficer = po

	return &data, body, nil
}

func Insert(award *Award, db *gorm.DB) error {
	var directorate Directorate
	db.FirstOrCreate(&directorate, Directorate{Abbreviation: award.Organization.Directorate.Abbreviation, LongName: award.Organization.Directorate.LongName})

	var division Division
	db.FirstOrCreate(&division, Division{Abbreviation: award.Organization.Division.Abbreviation, LongName: award.Organization.Division.LongName})

	var organization Organization
	db.FirstOrCreate(&organization, Organization{Code: award.Organization.Code, DirectorateID: directorate.ID, DivisionID: division.ID, Directorate: directorate, Division: division})
	award.Organization = &organization
	award.OrganizationID = organization.ID

	var institution Institution
	db.FirstOrCreate(&institution, Institution{
		Name:          award.Institution.Name,
		CityName:      award.Institution.CityName,
		ZipCode:       award.Institution.ZipCode,
		PhoneNumber:   award.Institution.PhoneNumber,
		StreetAddress: award.Institution.StreetAddress,
		CountryName:   award.Institution.CountryName,
		StateCode:     award.Institution.StateCode,
	})
	award.InstitutionID = institution.ID
	award.Institution = institution

	var programOfficer ProgramOfficer
	db.FirstOrCreate(&programOfficer, ProgramOfficer{Email: award.ProgramOfficer.Email, SignBlockName: award.ProgramOfficer.SignBlockName, Phone: award.ProgramOfficer.Phone})
	award.ProgramOfficerID = programOfficer.ID
	award.ProgramOfficer = programOfficer

	el := []ProgramElement{}
	for _, elem := range award.ProgramElements {
		var pgmElem ProgramElement
		db.FirstOrCreate(&pgmElem, ProgramElement{Code: elem.Code, Text: elem.Text})
		el = append(el, pgmElem)
	}
	award.ProgramElements = el

	refs := []ProgramReference{}
	for _, ref := range award.ProgramReferences {
		var pgmRef ProgramReference
		db.FirstOrCreate(&pgmRef, ProgramReference{Code: ref.Code, Text: ref.Text})
		refs = append(refs, pgmRef)
	}
	award.ProgramReferences = refs

	inv := []Investigator{}
	for _, ref := range award.Investigators {
		var pgmRef Investigator
		db.FirstOrCreate(&pgmRef, Investigator{
			FirstName:    ref.FirstName,
			LastName:     ref.LastName,
			EmailAddress: ref.EmailAddress,
			StartDate:    ref.StartDate,
			EndDate:      ref.EndDate,
			Role:         ref.Role,
			NsfID:        ref.NsfID,
		})
		inv = append(inv, pgmRef)
	}
	award.Investigators = inv

	return db.Create(&award).Error
}

func InsertMongo(award *Award, client *mongo.Client) error {
	collection := client.Database("ztbd").Collection("awards")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, award)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
