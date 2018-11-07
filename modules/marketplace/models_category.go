package marketplace

import (
	"errors"
	"fmt"
	"sort"
)

/*
	Models
*/

type ItemCategory struct {
	ID            int            `json:"id" gorm:"primary_key"`
	ParentID      int            `json:"parent_id"`
	Icon          string         `json:"icon"`
	NameEn        string         `json:"price_en"`
	NameRu        string         `json:"name_ru"`
	NameDe        string         `json:"name_de"`
	NameEs        string         `json:"name_es"`
	NameFr        string         `json:"name_fr"`
	NameRs        string         `json:"name_rs"`
	NameTr        string         `json:"name_tr"`
	ItemCount     int            `gorm:"-" sql:"-" json:"item_count"`
	UserCount     int            `gorm:"-" sql:"-" json:"user_count"`
	Subcategories []ItemCategory `json:"subcategories"`
}

func (ic ItemCategory) String() string {
	return fmt.Sprintf("/staff/item_categories/%d/", ic.ID)
}

/*
	Sort
*/

type ItemCaterogiesByCount []ItemCategory

func (s ItemCaterogiesByCount) Len() int {
	return len(s)
}

func (s ItemCaterogiesByCount) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ItemCaterogiesByCount) Less(i, j int) bool {
	return s[i].ItemCount > s[j].ItemCount
}

/*
	Model Methods
*/

func (ic ItemCategory) ParentCategory() *ItemCategory {
	if ic.ParentID == 0 {
		return nil
	}

	pc, _ := FindCategoryByID(ic.ParentID)
	return pc
}

func (ic ItemCategory) FindSubcategories() []ItemCategory {
	var cat []ItemCategory
	database.
		Where("parent_id = ?", ic.ID).
		Find(&cat)

	return cat
}

/*
	Model Interface Implementation
*/

func (cat ItemCategory) Validate() error {
	numberOfLevels := 0
	var cc *ItemCategory
	cc = &cat

	for {
		cc = cc.ParentCategory()
		if cc == nil {
			break
		}
		numberOfLevels += 1
	}

	if numberOfLevels >= 3 {
		return errors.New("Too many levels")
	}

	return nil
}

func (cat ItemCategory) Remove() error {
	return database.Delete(&cat).Error
}

func (cat ItemCategory) Save() error {
	err := cat.Validate()
	if err != nil {
		return err
	}
	return cat.SaveToDatabase()
}

func (cat *ItemCategory) SaveToDatabase() error {
	if existing, _ := FindCategoryByID(cat.ID); existing == nil {
		return database.Create(cat).Error
	}
	return database.Save(cat).Error
}

/*
	Queries
*/

func GetAllCategories() []Category {
	var categories []Category
	database.Unscoped().Find(&categories)
	return categories
}

func FindCategoryByID(id int) (*ItemCategory, error) {
	var cat ItemCategory
	err := database.
		First(&cat, "id = ?", id).
		Error
	if err != nil {
		return nil, err
	}
	return &cat, err
}

func FindCategoryByNameEn(name string) (*ItemCategory, error) {
	var cat ItemCategory
	err := database.
		First(&cat, "name_en = ?", name).
		Error
	if err != nil {
		return nil, err
	}
	return &cat, err
}

func FindCategoriesByParentID(id int) ([]ItemCategory, error) {
	var cat []ItemCategory
	err := database.
		Where("parent_id = ?", id).
		Find(&cat).
		Error
	if err != nil {
		return nil, err
	}
	return cat, err
}

func FindAllCategories() []ItemCategory {
	cats, err := FindCategoriesByParentID(0)
	if err != nil {
		return cats
	}

	for i, _ := range cats {
		subcats, _ := FindCategoriesByParentID(cats[i].ID)
		cats[i].Subcategories = subcats

		for j, _ := range cats[i].Subcategories {
			subcats, _ := FindCategoriesByParentID(cats[i].Subcategories[j].ID)
			cats[i].Subcategories[j].Subcategories = subcats
		}
	}

	return cats
}

/*
	Cache
*/

func CacheFillCategories(packageType, countryNameEnTo, countryNameEnFrom, vendorUuid string, cityId int) []ItemCategory {
	activeCategories := FindNonEmptyCategories(packageType, countryNameEnTo, countryNameEnFrom, vendorUuid, cityId)
	key := fmt.Sprintf(
		"active-categories-%s-%s-%s-%s-%d",
		packageType, countryNameEnTo, countryNameEnFrom, vendorUuid, cityId,
	)
	cache15m.Set(key, activeCategories)
	return activeCategories
}

func CacheGetCategories(packageType, countryNameEnTo, countryNameEnFrom, vendorUuid string, cityId int) []ItemCategory {
	key := fmt.Sprintf(
		"active-categories-%s-%s-%s-%s-%d",
		packageType, countryNameEnTo, countryNameEnFrom, vendorUuid, cityId,
	)
	cCats, _ := cache15m.Get(key)
	if cCats == nil {
		return CacheFillCategories(packageType, countryNameEnTo, countryNameEnFrom, vendorUuid, cityId)
	}
	categories := cCats.([]ItemCategory)
	return categories
}

/*
	Database Views
*/

func findNonEmptyCategories() []ItemCategory {
	var cats []ItemCategory
	database.
		Table("v_categories").
		Model(&ItemCategory{}).
		Find(&cats)
	return cats
}

func findNonEmptyCategoriesByUserUuid(userUuid string) []ItemCategory {
	query := `SELECT * from (
			SELECT 
				*,
				(
					select 
						count(distinct(items.uuid))
					from items 
					join users 
					on items.user_uuid=users.uuid 
					join packages 
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						items.user_uuid=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as item_count ,
				(
					select 
						count(distinct(items.user_uuid)) 
					from items 
					join users 
					on items.user_uuid=users.uuid
					join packages
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						items.user_uuid=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as user_count
			FROM 
				item_categories ic 
			) ic
			WHERE 
				ic.item_count > 0`

	var cats []ItemCategory
	database.Raw(query, userUuid, userUuid).
		Model(&ItemCategory{}).
		Find(&cats)
	return cats
}

func findNonEmptyCategoriesByPackageType(packageType string) []ItemCategory {
	query := `SELECT * from (
			SELECT 
				*,
				(
					select 
						count(distinct(items.uuid))
					from items 
					join users 
					on items.user_uuid=users.uuid 
					join packages 
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false
						and packages.type=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as item_count ,
				(
					select 
						count(distinct(items.user_uuid)) 
					from items 
					join users 
					on items.user_uuid=users.uuid
					join packages
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false
						and packages.type=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as user_count
			FROM 
				item_categories ic 
			) ic
			WHERE 
				ic.item_count > 0`

	var cats []ItemCategory
	database.Raw(query, packageType, packageType).
		Model(&ItemCategory{}).
		Find(&cats)
	return cats
}

func findNonEmptyCategoriesByPackageTypeAndUserUuid(packageType, userUuid string) []ItemCategory {
	query := `SELECT * from (
			SELECT 
				*,
				(
					select 
						count(distinct(items.uuid))
					from items 
					join users 
					on items.user_uuid=users.uuid 
					join packages 
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						packages.type=? and
						items.user_uuid=? and
						items.delted_at IS null and
						packages.deleted_at IS null and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as item_count ,
				(
					select 
						count(distinct(items.user_uuid)) 
					from items 
					join users 
					on items.user_uuid=users.uuid
					join packages
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						packages.type=? and
						items.user_uuid=? and
						items.delted_at IS null and
						packages.deleted_at IS null and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as user_count
			FROM 
				item_categories ic 
			) ic
			WHERE 
				ic.item_count > 0`

	var cats []ItemCategory
	database.Raw(query, packageType, userUuid, packageType, userUuid).
		Model(&ItemCategory{}).
		Find(&cats)
	return cats
}

func findNonEmptyCategoriesByPackageTypeAndCountry(packageType, countryNameEn string) []ItemCategory {
	query := `SELECT * from (
			SELECT 
				*,
				(
					select 
						count(distinct(items.uuid))
					from items 
					join users 
					on items.user_uuid=users.uuid 
					join packages 
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						packages.type=? and
						packages.country_name_en_shipping_to=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as item_count ,
				(
					select 
						count(distinct(items.user_uuid)) 
					from items 
					join users 
					on items.user_uuid=users.uuid
					join packages
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						packages.type=? and
						packages.country_name_en_shipping_to=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as user_count
			FROM 
				item_categories ic 
			) ic
			WHERE 
				ic.item_count > 0`

	var cats []ItemCategory
	database.Raw(query, packageType, countryNameEn, packageType, countryNameEn).
		Model(&ItemCategory{}).
		Find(&cats)
	return cats
}

func findNonEmptyCategoriesByPackageTypeAndCountryAndUserUuid(packageType, countryNameEn, userUuid string) []ItemCategory {
	query := `SELECT * from (
			SELECT 
				*,
				(
					select 
						count(distinct(items.uuid))
					from items 
					join users 
					on items.user_uuid=users.uuid 
					join packages 
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						packages.type=? and
						packages.country_name_en_shipping_to=? and
						items.user_uuid=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as item_count ,
				(
					select 
						count(distinct(items.user_uuid)) 
					from items 
					join users 
					on items.user_uuid=users.uuid
					join packages
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						packages.type=? and
						packages.country_name_en_shipping_to=? and
						items.user_uuid=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as user_count
			FROM 
				item_categories ic 
			) ic
			WHERE 
				ic.item_count > 0`

	var cats []ItemCategory
	database.Raw(query, packageType, countryNameEn, userUuid, packageType, countryNameEn, userUuid).
		Model(&ItemCategory{}).
		Find(&cats)
	return cats
}

func findNonEmptyCategoriesByPackageTypeAndCityId(packageType string, cityId int) []ItemCategory {
	query := `SELECT * from (
			SELECT 
				*,
				(
					select 
						count(distinct(items.uuid))
					from items 
					join users 
					on items.user_uuid=users.uuid 
					join packages 
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						packages.type=? and
						packages.drop_city_id=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as item_count ,
				(
					select 
						count(distinct(items.user_uuid)) 
					from items 
					join users 
					on items.user_uuid=users.uuid
					join packages
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						packages.type=? and
						packages.drop_city_id=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as user_count
			FROM 
				item_categories ic 
			) ic
			WHERE 
				ic.item_count > 0`

	var cats []ItemCategory
	database.Raw(query, packageType, cityId, packageType, cityId).
		Model(&ItemCategory{}).
		Find(&cats)
	return cats
}

func findNonEmptyCategoriesByPackageTypeAndCityIdAndUserUuid(packageType, userUuid string, cityId int) []ItemCategory {
	query := `SELECT * from (
			SELECT 
				*,
				(
					select 
						count(distinct(items.uuid))
					from items 
					join users 
					on items.user_uuid=users.uuid 
					join packages 
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						packages.type=? and
						packages.drop_city_id=? and
						items.user_uuid=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as item_count ,
				(
					select 
						count(distinct(items.user_uuid)) 
					from items 
					join users 
					on items.user_uuid=users.uuid
					join packages
					on packages.item_uuid = items.uuid
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						packages.type=? and
						packages.drop_city_id=? and
						items.user_uuid=? and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as user_count
			FROM 
				item_categories ic 
			) ic
			WHERE 
				ic.item_count > 0`

	var cats []ItemCategory
	database.Raw(query, packageType, cityId, userUuid, packageType, cityId, userUuid).
		Model(&ItemCategory{}).
		Find(&cats)
	return cats
}

func FindNonEmptyCategories(packageType, countryNameEnTo, countryNameEnFrom, userUuid string, cityId int) []ItemCategory {
	var categories []ItemCategory

	if userUuid == "" {
		if packageType == "" {
			categories = findNonEmptyCategories()
		} else if countryNameEnTo == "" && cityId == 0 {
			categories = findNonEmptyCategoriesByPackageType(packageType)
		} else if countryNameEnTo != "" && cityId == 0 {
			categories = findNonEmptyCategoriesByPackageTypeAndCountry(packageType, countryNameEnTo)
		} else if cityId != 0 {
			categories = findNonEmptyCategoriesByPackageTypeAndCityId(packageType, cityId)
		}

	} else {
		if packageType == "" {
			categories = findNonEmptyCategoriesByUserUuid(userUuid)
		} else if countryNameEnTo == "" && cityId == 0 {
			categories = findNonEmptyCategoriesByPackageTypeAndUserUuid(packageType, userUuid)
		} else if countryNameEnTo != "" && cityId == 0 {
			categories = findNonEmptyCategoriesByPackageTypeAndCountryAndUserUuid(packageType, countryNameEnTo, userUuid)
		} else if cityId != 0 {
			categories = findNonEmptyCategoriesByPackageTypeAndCityIdAndUserUuid(packageType, userUuid, cityId)
		}
	}

	clearCategories := FindAllCategories()

	for i1, _ := range clearCategories {

		for _, nec := range categories {
			if clearCategories[i1].ID == nec.ID {
				clearCategories[i1].ItemCount += nec.ItemCount
				clearCategories[i1].UserCount += nec.UserCount
			}
		}

		for i2, _ := range clearCategories[i1].Subcategories {

			for _, nec := range categories {
				if clearCategories[i1].Subcategories[i2].ID == nec.ID {
					clearCategories[i1].Subcategories[i2].ItemCount += nec.ItemCount
					clearCategories[i1].Subcategories[i2].UserCount += nec.UserCount

					clearCategories[i1].ItemCount += nec.ItemCount
					clearCategories[i1].UserCount += nec.UserCount
				}
			}
			for i3, _ := range clearCategories[i1].Subcategories[i2].Subcategories {
				for _, nec := range categories {
					if clearCategories[i1].Subcategories[i2].Subcategories[i3].ID == nec.ID {
						clearCategories[i1].Subcategories[i2].Subcategories[i3].ItemCount += nec.ItemCount
						clearCategories[i1].Subcategories[i2].Subcategories[i3].UserCount += nec.UserCount

						clearCategories[i1].Subcategories[i2].ItemCount += nec.ItemCount
						clearCategories[i1].Subcategories[i2].UserCount += nec.UserCount

						clearCategories[i1].ItemCount += nec.ItemCount
						clearCategories[i1].UserCount += nec.UserCount
					}
				}
			}

			sort.Sort(ItemCaterogiesByCount(clearCategories[i1].Subcategories[i2].Subcategories))
		}
		sort.Sort(ItemCaterogiesByCount(clearCategories[i1].Subcategories))
	}
	sort.Sort(ItemCaterogiesByCount(clearCategories))

	return clearCategories
}

// Create views and other representatives
func setupCategoriesViews() {
	database.Exec("DROP VIEW IF EXISTS v_categories")
	database.Exec(`
		CREATE VIEW v_categories AS
			SELECT * from (
			SELECT 
				*,
				(
					select 
						count(*) 
					from items 
					join users 
					on items.user_uuid=users.uuid 
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						users.pgp <> '' and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as item_count ,
				(
					select 
						count(distinct(items.user_uuid)) 
					from items 
					join users 
					on items.user_uuid=users.uuid 
					where
						items.item_category_id=ic.id and
						users.last_login_date >= (now() - interval '7 days') and
						users.banned = false and
						users.pgp <> '' and
						users.is_signed = true and
						users.is_trusted_seller = true
				) as user_count
			FROM 
				item_categories ic 
			) ic
			WHERE 
				ic.item_count > 0
	`)
}
