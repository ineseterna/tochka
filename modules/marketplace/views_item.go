package marketplace

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gocraft/web"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) viewShowItem(w web.ResponseWriter, r *web.Request) {
	packages := c.ViewItem.Item.PackagesWithoutReservation()
	c.ViewPackages = packages.ViewPackages()
	table := packages.GroupsTable()
	c.GroupPackagesByTypeOriginDestination = table.GetGroupPackagesByTypeOriginDestination()
}

func (c *Context) ViewShowItem(w web.ResponseWriter, r *web.Request) {
	c.viewShowItem(w, r)
	util.RenderTemplate(w, "item/show", c)
}

func (c *Context) EditItem(w web.ResponseWriter, r *web.Request) {

	translateCat := func(ic ItemCategory) Category {
		cat := Category{ID: fmt.Sprintf("%d", ic.ID), Name: ic.NameEn}
		if c.ViewUser.Language == "ru" {
			cat.Name = ic.NameRu
		}
		pic := ic.ParentCategory()
		for {
			if pic == nil {
				break
			}
			if c.ViewUser.Language == "ru" {
				cat.Name = pic.NameRu + " - " + cat.Name
			} else {
				cat.Name = pic.NameEn + " - " + cat.Name
			}
			pic = pic.ParentCategory()
		}

		return cat
	}

	for _, cat1 := range FindAllCategories() {
		if len(cat1.Subcategories) == 0 {
			c.Categories = append(c.Categories, translateCat(cat1))
		} else {
			for _, cat2 := range cat1.Subcategories {
				if len(cat2.Subcategories) == 0 {
					c.Categories = append(c.Categories, translateCat(cat2))
				} else {
					for _, cat3 := range cat2.Subcategories {
						c.Categories = append(c.Categories, translateCat(cat3))
					}
				}
			}
		}
	}

	c.CategoryID = c.Item.ItemCategoryID
	util.RenderTemplate(w, "item/edit", c)
}

func (c *Context) SaveItem(w web.ResponseWriter, r *web.Request) {
	if r.PathParams["item"] == "new" {
		c.Item.Uuid = util.GenerateUuid()
	}

	categoryId, err := strconv.ParseInt(r.FormValue("category"), 10, 64)
	if err != nil {
		c.Error = err.Error()
		c.EditItem(w, r)
		return
	}

	category, err := FindCategoryByID(int(categoryId))
	if err != nil {
		c.Error = err.Error()
		c.EditItem(w, r)
		return
	}

	c.Item.ItemCategory = *category
	c.Item.Name = r.FormValue("name")
	c.Item.Description = r.FormValue("description")
	c.Item.UserUuid = c.ViewSeller.Uuid

	validationError := c.Item.Validate()
	if validationError != nil {
		c.Error = validationError.Error()
		c.EditItem(w, r)
		return
	}
	err = util.SaveImage(r, "image", 500, c.Item.Uuid)
	if err != nil && r.PathParams["item"] == "new" {
		c.Error = "Image: " + err.Error()
		c.EditItem(w, r)
		return
	}
	err = c.Item.Save()
	if err != nil {
		c.Error = err.Error()
		c.EditItem(w, r)
		return
	}

	if r.PathParams["item"] == "new" {
		CreateFeedItem(c.Item.UserUuid, "new_item", "added new item", c.Item.Uuid)
	}

	if c.Item.UserUuid != c.ViewUser.Uuid {
		now := time.Now()
		c.Item.ReviewedByUserUuid = c.ViewUser.Uuid
		c.Item.ReviewedAt = &now
		c.Item.Save()
	}

	if c.ViewUser.IsStaff {
		CreateFeedItem(c.ViewUser.Uuid, "staff_edit_item", "edited item", c.Item.Uuid)
	}

	EventNewItem(c.Item)
	http.Redirect(w, r.Request, "/user/"+c.ViewSeller.Username+"/item/"+c.Item.Uuid, 302)
}

func (c *Context) DeleteItem(w web.ResponseWriter, r *web.Request) {
	c.Item.Remove()
	if c.ViewUser.IsStaff {
		CreateFeedItem(c.ViewUser.Uuid, "staff_delete_item", "deleted item", c.Item.Uuid)
	}
	http.Redirect(w, r.Request, "/user/"+c.ViewUser.Username, 302)
}

func (c *Context) ItemImage(w web.ResponseWriter, r *web.Request) {
	itemUuid := r.PathParams["item"]
	size := "normal"
	if len(r.URL.Query()["size"]) > 0 {
		size = r.URL.Query()["size"][0]
	}
	util.ServeImage(itemUuid, size, w, r)
}
