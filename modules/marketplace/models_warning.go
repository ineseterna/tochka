package marketplace

import (
	"errors"
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Models
*/

type UserWarning struct {
	Uuid         string `json:"uuid" gorm:"primary_key"`
	UserUuid     string `json:"user_uuid" gorm:"index"`
	Severety     string `json:"severety" gorm:"index"`
	Text         string `json:"text"`
	IsApproved   bool   `json:"is_approved" gorm:"index"`
	ReporterUuid string `json:"reporter_uuid" gorm:"index"`

	User     User `json:"-"`
	Reporter User `json:"-"`

	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

/*
	Database methods
*/

func (w UserWarning) Validate() error {
	if w.Uuid == "" {
		return errors.New("Empty Public Key")
	}
	if w.UserUuid == "" {
		return errors.New("Empty User Uuid")
	}
	if w.Severety == "" {
		return errors.New("Empty Severety")
	}
	if w.Text == "" {
		return errors.New("Empty Text")
	}
	if w.ReporterUuid == "" {
		return errors.New("Empty Reporter User Uuid")
	}
	return nil
}

func (w UserWarning) Save() error {
	err := w.Validate()
	if err != nil {
		return err
	}
	return w.SaveToDatabase()
}

func (w UserWarning) Remove() error {
	return database.Delete(w).Error
}

func (w UserWarning) SaveToDatabase() error {
	if existing, _ := FindUserWarningByUuid(w.Uuid); existing == nil {
		return database.Create(&w).Error
	}
	return database.Save(&w).Error
}

/*
	Model Methods
*/

func (w UserWarning) HasExpired() bool {

	switch w.Severety {
	case "RED":
		d, _ := time.ParseDuration("672h")
		return w.CreatedAt.Before(time.Now().Add(-d))
	case "YELLOW":
		d, _ := time.ParseDuration("336h")
		return w.CreatedAt.Before(time.Now().Add(-d))
	case "GREEN":
		d, _ := time.ParseDuration("168h")
		return w.CreatedAt.Before(time.Now().Add(-d))
	}

	return false
}

func (w UserWarning) ExpiresIn(lang string) string {

	switch w.Severety {
	case "RED":
		d, _ := time.ParseDuration("672h")
		return util.HumanizeTime(w.CreatedAt.Add(d), lang)
	case "YELLOW":
		d, _ := time.ParseDuration("336h")
		return util.HumanizeTime(w.CreatedAt.Add(d), lang)
	case "GREEN":
		d, _ := time.ParseDuration("168h")
		return util.HumanizeTime(w.CreatedAt.Add(d), lang)
	}

	return ""
}

/*
	Model Methods
*/

func CreateUserWarning(userUuid, reporterUuid, text string) (UserWarning, error) {
	now := time.Now()
	uw := UserWarning{
		Uuid:         util.GenerateUuid(),
		UserUuid:     userUuid,
		ReporterUuid: reporterUuid,
		Text:         text,
		Severety:     "NEW",
		CreatedAt:    &now,
	}
	return uw, uw.Save()
}

func (uw *UserWarning) UpdateSeverety(severety string) error {
	uw.Severety = severety
	uw.IsApproved = true
	return uw.Save()
}

/*
	View Models
*/

type ViewUserWarning struct {
	*UserWarning
	CreatedAtStr string
	ExpiresIn    string
	HasExpired   bool
}

type UserWarnings []UserWarning

func (uw UserWarning) ViewUserWarning(lang string) ViewUserWarning {
	vw := ViewUserWarning{
		UserWarning:  &uw,
		CreatedAtStr: util.HumanizeTime(*uw.CreatedAt, lang),
		ExpiresIn:    uw.ExpiresIn(lang),
		HasExpired:   uw.HasExpired(),
	}
	return vw
}

func (uws UserWarnings) ViewUserWarnings(lang string) []ViewUserWarning {
	vuws := []ViewUserWarning{}
	for _, uw := range uws {
		if !uw.HasExpired() {
			vuw := uw.ViewUserWarning(lang)
			vuws = append(vuws, vuw)
		}
	}
	return vuws
}

/*
	Queries
*/

func FindUserWarningByUuid(uuid string) (*UserWarning, error) {
	var userWarning UserWarning
	err := database.
		Where("uuid = ?", uuid).
		Preload("User").
		First(&userWarning).
		Error
	if err != nil {
		return nil, err
	}
	return &userWarning, err
}

func FindActiveWarningsForUser(userUuid string) UserWarnings {
	var warnings UserWarnings
	database.
		Table("user_warnings").
		Where("user_uuid = ?", userUuid).
		Preload("User").
		Preload("Reporter").
		Find(&warnings)
	return warnings
}

func CountWarningsForUser(userUuid string) int {
	var count int
	database.
		Table("user_warnings").
		Where("user_uuid = ?", userUuid).
		Count(&count)
	return count
}

func FindAllActiveWarnings() UserWarnings {
	var warnings UserWarnings
	database.
		Table("user_warnings").
		Order("created_at DESC").
		Preload("User").
		Preload("Reporter").
		Find(&warnings)
	return warnings
}

/*
	Utility Methods
*/

func CanUserReportUser(u1, u2 User) bool {
	if u1.IsAdmin || u1.IsStaff {
		return true
	}
	if u2.IsSeller {
		return CountPayedTransactionsForBuyerAndSeller(u1.Uuid, u2.Uuid) > 0
	}
	return false
}
