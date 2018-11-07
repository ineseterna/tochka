package marketplace

import (
	"errors"
	"github.com/russross/blackfriday"
	"html/template"
)

/*
	Models
*/

type Invitation struct {
	Uuid           string `form:"uuid" json:"uuid" gorm:"primary_key" sql:"size:1024"`
	Username       string `form:"username" json:"username" sql:"size:1024"`
	InvitationText string `form:"invitation_text" json:"invitation_text" sql:"size:1024"`
	InviterUuid    string `form:"inviter_uuid" json:"inviter_uuid" sql:"size:1024"`
	IsActivated    bool
}

type ViewInvitation struct {
	Invitation
	InvitationTextHtml template.HTML
}

/*
	Model Methods
*/

func (i Invitation) Validate() error {
	if i.Username == "" {
		return errors.New("Name is not valid")
	}
	u, _ := FindUserByUsername(i.Username)
	if u != nil && i.Uuid == "" {
		return errors.New("Username is taken")
	}
	if i.InvitationText == "" {
		return errors.New("Invitation text is not valid")
	}
	if i.InviterUuid == "" {
		return errors.New("Inviter is not valid")
	}
	return nil
}

func (r Invitation) Save() error {
	err := r.Validate()
	if err != nil {
		return err
	}
	return r.SaveToDatabase()
}

func (itm Invitation) SaveToDatabase() error {
	if existing, _ := FindInvitationByUuid(itm.Uuid); existing == nil {
		return database.Create(&itm).Error
	}
	return database.Save(&itm).Error
}

func (r Invitation) Remove() error {
	return database.Delete(r).Error
}

func (i Invitation) ViewInvitation() ViewInvitation {
	return ViewInvitation{
		Invitation:         i,
		InvitationTextHtml: template.HTML(htmlPolicy.Sanitize(string(blackfriday.MarkdownCommon([]byte(i.InvitationText))))),
	}
}

/*
	Queries
*/

func FindInvitationByUuid(uuid string) (*Invitation, error) {
	var item Invitation
	err := database.
		Where(&Invitation{Uuid: uuid}).
		First(&item).
		Error
	if err != nil {
		return nil, err
	}
	return &item, err
}

func FindInvitationByUsername(username string) (*Invitation, error) {
	var item Invitation
	err := database.
		Where(&Invitation{Username: username}).
		First(&item).
		Error
	if err != nil {
		return nil, err
	}
	return &item, err
}

func FindInvitationsByInviterUuid(uuid string) []Invitation {
	var invitations []Invitation
	err := database.
		Where(&Invitation{InviterUuid: uuid}).
		Find(&invitations).
		Error
	if err != nil {
		return invitations
	}
	return invitations
}
