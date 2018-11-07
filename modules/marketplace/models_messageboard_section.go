package marketplace

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

/*
	Models
*/

type MessageboardSection struct {
	ID       int `json:"id" gorm:"primary_key"`
	Priority int `json:"priority"`
	ParentID int `json:"parent_id"`

	Icon string `json:"icon"`
	Flag string `json:"flag"`

	NameEn string `json:"name_en"`
	NameRu string `json:"name_ru"`
	NameDe string `json:"name_de"`
	NameEs string `json:"name_es"`
	NameFr string `json:"name_fr"`

	DescriptionEn string `json:"description_en"`
	DescriptionRu string `json:"description_ru"`
	DescriptionDe string `json:"description_de"`
	DescriptionEs string `json:"description_es"`
	DescriptionFr string `json:"description_fr"`

	// Subsections    []MessageboardSection `json:"subsections"`
	HeadingSection bool `json:"heading_section"`

	// View Model Extensions
	NumberOfThreads            int       `sql:"-"`
	NumberOfMessages           int       `sql:"-"`
	LastUpdated                time.Time `sql:"-"`
	LastUpdatedStr             string    `sql:"-"`
	LastThreadUuid             string    `sql:"-"`
	LastThreadTitle            string    `sql:"-"`
	LastThreadNumberOfMessages int       `sql:"-"`
	LastMessageByUsername      string    `sql:"-"`
}

type MessageboardSections []MessageboardSection

/*
	Database methods
*/

func (cat MessageboardSection) Validate() error {
	return nil
}

func (cat MessageboardSection) Remove() error {
	return database.Delete(&cat).Error
}

func (cat *MessageboardSection) Save() error {
	err := cat.Validate()
	if err != nil {
		return err
	}
	return cat.SaveToDatabase()
}

func (cat *MessageboardSection) SaveToDatabase() error {
	if existing, _ := FindMessageboardSectionByID(cat.ID); existing == nil {
		cat.ID = rand.Intn(100000)
		return database.Create(cat).Error
	}
	return database.Save(cat).Error
}

/*
	Queries
*/

func FindAllMessageboardSections() MessageboardSections {
	var sections MessageboardSections

	database.
		Table("v_messageboard_sections").
		Find(&sections)

	return sections
}

func FindParentMessageboardSections() MessageboardSections {
	var sections MessageboardSections

	database.
		Table("v_messageboard_sections").
		Where("parent_id=?", 0).
		Find(&sections)

	return sections
}

func FindMessageboardSectionsByParentID(parentID int) MessageboardSections {
	var sections MessageboardSections

	database.
		Model(MessageboardSection{}).
		Where("parent_id=?", parentID).
		Order("priority DESC").
		Find(&sections)

	return sections
}

func FindMessageboardSectionByID(id int) (*MessageboardSection, error) {
	var (
		section MessageboardSection
	)

	err := database.
		Where("id = ?", id).
		Find(&section).Error

	if err != nil {
		return nil, err
	}

	return &section, err
}

/*
	View Models
*/

type ViewMessageboardSection struct {
	*MessageboardSection

	Name            string
	Description     string
	SubViewsections []ViewMessageboardSection
}

func (section MessageboardSection) ViewMessageboardSection(lang string) ViewMessageboardSection {
	vs := ViewMessageboardSection{
		MessageboardSection: &section,
	}

	switch lang {
	case "ru":
		vs.Name = section.NameRu
		vs.Description = section.DescriptionRu
	case "fr":
		vs.Name = section.NameFr
		vs.Description = section.DescriptionFr
	case "es":
		vs.NameEs = section.NameEs
		vs.Description = section.DescriptionEs
	case "de":
		vs.Name = section.NameDe
		vs.Description = section.DescriptionDe
	default:
		vs.Name = section.NameEn
		vs.Description = section.DescriptionEn
	}

	vs.LastUpdatedStr = util.HumanizeTime(section.LastUpdated, lang)

	return vs
}

func (sections MessageboardSections) ViewMessageboardSections(lang string) []ViewMessageboardSection {

	vss := []ViewMessageboardSection{}

	for i, _ := range sections {
		vs := sections[i].ViewMessageboardSection(lang)
		vss = append(vss, vs)
	}

	sort.SliceStable(vss, func(i, j int) bool {
		if vss[i].Priority != vss[j].Priority {
			return vss[i].Priority > vss[j].Priority
		}
		return strings.Compare(vss[i].Name, vss[j].Name) < 0
	})

	return vss
}

func (sections MessageboardSections) ByParentID(parentId int) MessageboardSections {
	filteredSections := MessageboardSections{}
	for i, _ := range sections {
		if sections[i].ParentID == parentId {
			filteredSections = append(filteredSections, sections[i])
		}
	}

	return filteredSections
}

func (sections MessageboardSections) AsNestedViewMessageboardSections(lang string) []ViewMessageboardSection {
	parentSections := sections.ByParentID(0).ViewMessageboardSections(lang)
	for i, _ := range parentSections {
		parentSections[i].SubViewsections = sections.ByParentID(parentSections[i].ID).ViewMessageboardSections(lang)
	}

	return parentSections
}

/*
	Cache Optimizations
*/

func CacheGetAllMessageboardSections() MessageboardSections {
	key := "extended-messageboard-sections"
	cSections, _ := cache15m.Get(key)
	if cSections == nil {
		vs := FindAllMessageboardSections()
		cache15m.Set(key, vs)

		return vs
	}
	return cSections.(MessageboardSections)

}

/*
	DB Views
*/

func setupMessageboardCategoriesViews() {
	database.Exec("DROP VIEW IF EXISTS v_messageboard_sections CASCADE;")
	database.Exec(`
	CREATE VIEW v_messageboard_sections AS (
		WITH t_messageboard_sections as (
			select 
				messageboard_section_id as id, 
				count(*) as number_of_threads, 
				sum(number_of_messages) as number_of_messages
			from 
				v_messageboard_threads 
			group by 
				messageboard_section_id
		)
		
		SELECT 
			ms.*,
			tms.number_of_threads,
			tms.number_of_messages,
			vt.uuid as last_thread_uuid,
			vt.last_updated as last_updated,
			vt.title as last_thread_title,
			vt.number_of_messages as last_thread_number_of_messages,
			vt.last_message_by_username as last_message_by_username
			from 
				messageboard_sections ms 
			left join
				t_messageboard_sections tms on ms.id = tms.id
			left join
				v_threads vt on vt.uuid=(select uuid from v_threads where messageboard_section_id=ms.id order by last_updated DESC limit 1) 
	)`)
}
