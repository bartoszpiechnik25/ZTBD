package models

type Investigator struct {
	ID           uint    `gorm:"primaryKey" bson:"-"`
	FirstName    *string `gorm:"size:64" json:"pi_first_name" bson:"first_name"`
	LastName     *string `gorm:"size:64" json:"pi_last_name" bson:"last_name"`
	EmailAddress *string `gorm:"size:128" json:"pi_email_addr" bson:"email_address"`
	StartDate    *string `gorm:"type:date" json:"pi_start_date" bson:"start_date"`
	EndDate      *string `gorm:"type:date" json:"pi_end_date" bson:"end_date"`
	Role         *string `gorm:"size128" json:"pi_role" bson:"role"`
	NsfID        *string `gorm:"size:255" json:"nsf_id" bson:"nsf_id"`
}

type ProgramElement struct {
	ID   uint    `gorm:"primaryKey" bson:"-"`
	Code *string `gorm:"size:32,uniqueIndex" json:"pgm_ele_code" bson:"code"`
	Text *string `gorm:"type:text" json:"pgm_ele_name" bson:"text"`
}

type ProgramReference struct {
	ID   uint    `gorm:"primaryKey" bson:"-"`
	Code *string `gorm:"size:32,uniqueIndex" json:"pgm_ref_code" bson:"code"`
	Text *string `gorm:"type:text" json:"pgm_ref_txt" bson:"text"`
}

type Institution struct {
	ID            uint    `gorm:"primaryKey" bson:"-"`
	Name          *string `gorm:"size:128" json:"inst_name" bson:"name"`
	CityName      *string `gorm:"size:64" json:"inst_city_name" bson:"city_name"`
	ZipCode       *string `gorm:"size:32" json:"inst_zip_codex" bson:"zip_code"`
	PhoneNumber   *string `gorm:"size:16" json:"inst_phone_num" bson:"phone_number"`
	StreetAddress *string `gorm:"size:128" json:"inst_street_address" bson:"street_address"`
	CountryName   *string `gorm:"size:64" json:"inst_country_name" bson:"country_name"`
	StateCode     *string `gorm:"size:16" json:"inst_state_code" bson:"state_code"`
}

type Directorate struct {
	ID           uint    `gorm:"primaryKey" bson:"-"`
	LongName     *string `gorm:"type:text" json:"org_dir_long_name" bson:"long_name"`
	Abbreviation *string `gorm:"size:64,uniqueIndex" json:"dir_abbr" bson:"abbreviation"`
}

type Division struct {
	ID           uint    `gorm:"primaryKey" bson:"-"`
	LongName     *string `gorm:"type:text" json:"org_div_long_name" bson:"long_name"`
	Abbreviation *string `gorm:"size:64,uniqueIndex" json:"div_abbr" bson:"abbreviation"`
}

type Organization struct {
	ID            uint        `gorm:"primaryKey" bson:"-"`
	Code          *string     `gorm:"size:64;uniqueIndex;not null" json:"org_code" bson:"code"`
	Directorate   Directorate `gorm:"foreignKey:DirectorateID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"directorate" bson:"directorate"`
	DirectorateID uint        `bson:"-"`
	Division      Division    `gorm:"foreignKey:DivisionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"division" bson:"division"`
	DivisionID    uint        `bson:"-"`
}

type ProgramOfficer struct {
	ID            uint    `gorm:"primaryKey" bson:"-"`
	SignBlockName *string `gorm:"size:128" json:"sign_blck_name" bson:"sign_block_name"`
	Email         *string `gorm:"size:128,uniqueIndex" json:"po_email" bson:"email"`
	Phone         *string `gorm:"size:16" json:"po_phone" bson:"phone"`
}

type Award struct {
	ID                  uint               `gorm:"primaryKey" bson:"-"`
	Title               string             `gorm:"type:text" json:"awd_titl_txt" bson:"title"`
	AwardEffectiveDate  *string            `gorm:"type:date" json:"awd_eff_date" bson:"award_effective_date"`
	AwardExpirationDate *string            `gorm:"type:date" json:"awd_exp_date" bson:"award_expiration_date"`
	AwardAmount         *float64           `gorm:"type:numeric" json:"tot_intn_awd_amt" bson:"award_amount"`
	AbstractText        *string            `gorm:"type:text" json:"abst_narr_txt" bson:"abstract_text"`
	ProgramOfficerID    uint               `json:"program_officer_id" bson:"-"`
	ProgramOfficer      ProgramOfficer     `gorm:"foreignKey:ProgramOfficerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"program_officer" bson:"program_officer"`
	OrganizationID      uint               `json:"organization_id" bson:"-"`
	Organization        *Organization      `gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"organization" bson:"organization"`
	InstitutionID       uint               `json:"institution_id" bson:"-"`
	Institution         Institution        `gorm:"foreignKey:InstitutionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"inst" bson:"institution"`
	Investigators       []Investigator     `gorm:"many2many:award_investigator;" json:"pi" bson:"investigators"`
	ProgramElements     []ProgramElement   `gorm:"many2many:award_program_element;" json:"pgm_ele" bson:"program_elements"`
	ProgramReferences   []ProgramReference `gorm:"many2many:award_program_reference;" json:"pgm_ref" bson:"program_references"`
}
