package marketplace

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gocraft/web"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) AdminMessagesShow(w web.ResponseWriter, r *web.Request) {
	if uuid := r.PathParams["uuid"]; uuid != "" {
		thread, err := FindThreadByUuid(uuid)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.Thread = *thread
		viewThread := thread.ViewThread(c.ViewUser.Language, c.ViewUser.User)
		c.ViewThread = &viewThread
	}
	util.RenderTemplate(w, "board/admin/message", c)
}

func (c *Context) AdminMessageboardSections(w web.ResponseWriter, r *web.Request) {
	c.MessageboardSections = FindParentMessageboardSections()
	util.RenderTemplate(w, "board/admin/sections", c)
}

func (c *Context) AdminMessageboardSectionsEdit(w web.ResponseWriter, r *web.Request) {
	if r.PathParams["id"] != "new" {
		msId, err := strconv.ParseInt(r.PathParams["id"], 10, 64)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		section, err := FindMessageboardSectionByID(int(msId))
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		c.MessageboardSection = *section
	}

	c.MessageboardSections = FindParentMessageboardSections()

	ms := MessageboardSection{
		NameEn: "",
		ID:     0,
	}
	c.MessageboardSections = append(c.MessageboardSections, ms)

	util.RenderTemplate(w, "board/admin/sections_edit", c)
}

func (c *Context) AdminMessageboardSectionsEditPOST(w web.ResponseWriter, r *web.Request) {
	var section MessageboardSection
	if r.PathParams["id"] != "new" {
		msId, err := strconv.ParseInt(r.PathParams["id"], 10, 64)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		sec, err := FindMessageboardSectionByID(int(msId))
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		section = *sec
	}
	err := r.ParseForm()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}

	section.NameEn = r.FormValue("name_en")
	section.NameRu = r.FormValue("name_ru")
	section.NameDe = r.FormValue("name_de")
	section.NameEs = r.FormValue("name_es")
	section.NameFr = r.FormValue("name_fr")

	section.DescriptionRu = r.FormValue("description_ru")
	section.DescriptionEn = r.FormValue("description_en")
	section.DescriptionDe = r.FormValue("description_de")
	section.DescriptionEs = r.FormValue("description_es")
	section.DescriptionFr = r.FormValue("description_fr")

	section.Icon = r.FormValue("icon")
	section.Flag = r.FormValue("flag")

	if r.FormValue("priority") != "" {
		priority, err := strconv.ParseInt(r.FormValue("priority"), 10, 64)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		section.Priority = int(priority)
	}

	if r.FormValue("parent_id") != "" {
		parentID, err := strconv.ParseInt(r.FormValue("parent_id"), 10, 64)
		if err != nil {
			http.NotFound(w, r.Request)
			return
		}
		if parentID != 0 {
			ps, _ := FindMessageboardSectionByID(int(parentID))
			if ps == nil {
				http.NotFound(w, r.Request)
				return
			}
		}
		section.ParentID = int(parentID)
	}

	if r.FormValue("heading_section") != "" {
		hs := r.FormValue("heading_section")
		if hs == "1" {
			section.HeadingSection = true
		} else {
			section.HeadingSection = false
		}
	}

	err = section.Save()
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r.Request, fmt.Sprintf("/messageboard_sections/admin/%d", section.ID), 302)
}

func (c *Context) AdminMessageboardSectionsDelete(w web.ResponseWriter, r *web.Request) {
	catId, err := strconv.ParseUint(r.PathParams["id"], 10, 64)
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	section, err := FindMessageboardSectionByID(int(catId))
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	err = section.Remove()
	if err != nil {
		http.NotFound(w, r.Request)
		return
	}
	http.Redirect(w, r.Request, "/messageboard_sections/admin", 302)
}
