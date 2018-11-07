package marketplace

import (
	"errors"
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

type Advertising struct {
	Uuid        string    `json:"uuid" gorm:"primary_key"`
	DateCreated time.Time `json:"created_at"`
	Comment     string    `json:"comment"`
	DateStart   time.Time `json:"date_start"`
	DateEnd     time.Time `json:"date_end"`
	Status      bool      `json:"status"`

	CountImpressions        int `json:"count_impressions"`
	CurrentCountImpressions int `json:"Current_count_impressions"`

	VendorUuid string `json:"vendor_uuid" sql:"index"`
	ItemUuid   string `json:"item_uuid" sql:"index"`
	Item       Item   `json:"-"`
}

type Advertisings []Advertising

/*
	Model Methods
*/

func (ad Advertising) Validate() error {
	if len(ad.Comment) > 50 || len(ad.Comment) == 0 {
		return errors.New("Text is limited to 50 symbols.")
	}
	return nil
}

func (ad Advertising) Save() error {
	err := ad.Validate()
	if err != nil {
		return err
	}
	return ad.SaveToDatabase()
}

func (ad *Advertising) SaveToDatabase() error {
	if existing, _ := FindAdvertisingByUuid(ad.Uuid); existing == nil {
		return database.Create(ad).Error
	}
	return database.Save(ad).Error
}

func (ad *Advertising) AddImpressions() error {
	ad.CurrentCountImpressions++
	if ad.CurrentCountImpressions >= ad.CountImpressions {
		ad.Status = false
		ad.DateEnd = time.Now()
		return database.Save(&ad).Error
	}

	return database.
		Model(&ad).
		UpdateColumn("CurrentCountImpressions", ad.CurrentCountImpressions).
		Error
}

func CreateAdvertising(comment string, count int, vendorUuid string, itemUuid string) error {

	ad := Advertising{
		Uuid:                    util.GenerateUuid(),
		DateCreated:             time.Now(),
		Comment:                 comment,
		DateStart:               time.Now(),
		Status:                  true,
		CountImpressions:        count,
		CurrentCountImpressions: 0,
		VendorUuid:              vendorUuid,
		ItemUuid:                itemUuid,
	}

	err := ad.Validate()
	if err != nil {
		return err
	}

	return ad.Save()
}

/*
	Queries
*/

func FindAllAdvertising() (Advertisings, error) {
	var ads Advertisings
	err := database.
		Preload("Item").
		Preload("Item.User").
		Order("date_created ASC").
		Find(&ads).Error
	if err != nil {
		return nil, err
	}
	return ads, err
}

func FindAllActiveAdvertisings() (Advertisings, error) {
	var ads Advertisings
	err := database.
		Table("v_advertisings").
		Preload("Item").
		Preload("Item.User").
		Find(&ads).Error
	return ads, err
}

func FindRandomActiveAdvertisings(limit int) (Advertisings, error) {
	var ads []Advertising
	err := database.
		Table("v_advertisings").
		Order("random()").
		Limit(limit).
		Preload("Item").
		Preload("Item.User").
		Find(&ads).
		Error
	return ads, err
}

func FindAdvertisingByUuid(uuid string) (*Advertising, error) {
	var ad Advertising
	err := database.
		Where("advertisings.item_uuid = ? and advertisings.status = true", uuid).
		Preload("Item").
		Preload("Item.User").
		Last(&ad).Error
	if err != nil {
		return nil, err
	}
	return &ad, err
}

func FindAdvertisingByVendor(uuid string) (Advertisings, error) {
	var ads Advertisings
	err := database.
		Where(&Advertising{VendorUuid: uuid}).
		Preload("Item").
		Preload("Item.User").
		Find(&ads).Error
	return ads, err
}

func FindAdvertisingByItem(uuid string) (*Advertising, error) {
	var ad Advertising
	err := database.
		Where(&Advertising{ItemUuid: uuid}).
		Preload("Item").
		Preload("Item.User").
		First(&ad).Error
	if err != nil {
		return nil, err
	}
	return &ad, err
}

func GetAdvertisings(limit int) (Advertisings, error) {
	var ads Advertisings
	ads, err := FindRandomActiveAdvertisings(limit)
	if err != nil {
		return ads, err
	}
	for i, _ := range ads {
		ads[i].AddImpressions()
	}
	return ads, nil
}

/*
	Database Views
*/

func setupAdvertisingViews() {
	database.Exec("DROP VIEW IF EXISTS v_advertisings CASCADE;")
	database.Exec(`
		CREATE VIEW v_advertisings AS (
			select 
				a.* 
			from 
				advertisings a
			join
				users u on u.uuid=a.vendor_uuid
			join
				items i on i.uuid=a.item_uuid
			where
				a.status = true and
				u.banned=false and
				i.reviewed_at is not null
	);`)
}
