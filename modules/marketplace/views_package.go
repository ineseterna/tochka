package marketplace

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocraft/web"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func (c *Context) ListPackages(w web.ResponseWriter, r *web.Request) {
	avs := FindPaidPackagesForBuyer(c.ViewUser.Uuid)
	c.ViewPackages = avs.ViewPackages()
	util.RenderTemplate(w, "package/list", c)
}

// Edit Package Wizard

func (c *Context) EditPackageStep1(w web.ResponseWriter, r *web.Request) {
	util.RenderTemplate(w, "package/edit-step-1", c)
}

func (c *Context) EditPackageStep2(w web.ResponseWriter, r *web.Request) {
	c.Countries = GetAllCountries()

	var template string
	switch c.Package.Type {
	case "drop", "drop preorder":
		template = "package/edit-step-2-drop"
		break

	case "mail":
		template = "package/edit-step-2-mail"
		break

	case "digital":
		template = "package/edit-step-2-digital"
		break

	default:
		http.NotFound(w, r.Request)
		return
	}

	c.ViewPackage = c.Package.ViewPackage()
	util.RenderTemplate(w, template, c)
}

func (c *Context) EditPackagesStep2DropAndDropPreorder(w web.ResponseWriter, r *web.Request) {
	c.GeoCities = FindCitiesByCountryNameEn(c.Package.CountryNameEnShippingTo)

	var template string
	switch c.Package.Type {
	case "drop":
		template = "package/edit-step-3-drop"
		break
	case "drop preorder":
		template = "package/edit-step-3-drop-preorder"
		break
	}

	c.ViewPackage = c.Package.ViewPackage()
	util.RenderTemplate(w, template, c)
}

func (c *Context) EditPackageStep3(w web.ResponseWriter, r *web.Request) {

	var template string
	switch c.Package.Type {
	case "drop":
		template = "package/edit-step-3-drop"
		break

	case "drop preorder":
		template = "package/edit-step-3-drop-preorder"
		break

	default:
		http.NotFound(w, r.Request)
		return
	}

	c.ViewPackage = c.Package.ViewPackage()
	util.RenderTemplate(w, template, c)
}

func (c *Context) EditPackageStep4(w web.ResponseWriter, r *web.Request) {

	var template string

	switch c.Package.Type {
	case "drop", "drop preorder":
		template = "package/edit-step-4-metro"
		break

	default:
		http.NotFound(w, r.Request)
		return
	}

	c.ViewPackage = c.Package.ViewPackage()
	util.RenderTemplate(w, template, c)
}

// Save Package Wizard
func (c *Context) SavePackage(w web.ResponseWriter, r *web.Request) {
	switch r.FormValue("step") {
	case "1":
		c.SavePackageStep1(w, r)
		break
	case "2":
		c.SavePackageStep2(w, r)
		break
	case "3":
		c.SavePackageStep3(w, r)
		break
	case "4":
		c.SavePackageStep4(w, r)
		break
	default:
		http.NotFound(w, r.Request)
		return
	}
}

func (c *Context) SavePackageStep1(w web.ResponseWriter, r *web.Request) {
	var err error

	c.Package, err = c.parsePackageFormStep1(w, r)
	if err != nil {
		c.Error = err.Error()
		c.ViewPackage = c.Package.ViewPackage()
		c.EditPackageStep1(w, r)
		return
	}

	switch c.Package.Type {
	case "drop", "digital", "mail", "drop preorder":
		c.EditPackageStep2(w, r)
		break
	default:
		http.NotFound(w, r.Request)
	}
}

func (c *Context) SavePackageStep2(w web.ResponseWriter, r *web.Request) {
	var err error
	c.Package, err = c.parsePackageFormStep2(w, r)
	if err != nil {
		c.Error = err.Error()
		c.ViewPackage = c.Package.ViewPackage()
		c.EditPackageStep2(w, r)
		return
	}

	switch c.Package.Type {

	case "drop", "drop preorder":
		c.SavePackageStep2DropAndDropPreorder(w, r)
		break

	case "mail":
		countryTo, err := FindCountryByNameEn(r.FormValue("country_name_en_to"))
		if err != nil {
			c.Error = "No such country (shipping to) in database"
			c.ViewPackage = c.Package.ViewPackage()
			c.EditPackageStep2(w, r)
			return
		}
		c.Package.CountryNameEnShippingTo = countryTo.NameEn

		countryFrom, err := FindCountryByNameEn(r.FormValue("country_name_en_from"))
		if err != nil {
			c.Error = "No such country (shipping from) in database"
			c.ViewPackage = c.Package.ViewPackage()
			c.EditPackageStep2(w, r)
			return
		}
		c.Package.CountryNameEnShippingFrom = countryFrom.NameEn
		c.SavePackageComplete(w, r)

		break

	case "digital":
		c.Package.CountryNameEnShippingTo = "Interwebs"
		c.Package.CountryNameEnShippingFrom = "Interwebs"
		c.SavePackageComplete(w, r)
		break

	default:
		http.NotFound(w, r.Request)
		return
	}

}

func (c *Context) SavePackageStep2DropAndDropPreorder(w web.ResponseWriter, r *web.Request) {
	countryTo, err := FindCountryByNameEn(r.FormValue("country_name_en_to"))
	if err != nil {
		c.Error = "No such country in database"
		c.ViewPackage = c.Package.ViewPackage()
		c.EditPackagesStep2DropAndDropPreorder(w, r)
		return
	}
	c.Package.CountryNameEnShippingTo = countryTo.NameEn

	c.EditPackagesStep2DropAndDropPreorder(w, r)
}

func (c *Context) SavePackageStep3(w web.ResponseWriter, r *web.Request) {
	var err error
	c.Package, err = c.parsePackageFormStep3(w, r)
	c.ViewPackage = c.Package.ViewPackage()
	if err != nil {
		c.Error = err.Error()
		c.ViewPackage = c.Package.ViewPackage()
		c.EditPackageStep3(w, r)
		return
	}

	switch c.Package.Type {
	case "drop", "drop preorder":
		c.CityMetroStations = FindCityMetroStationsByCity(c.Package.DropCityId)
		if len(c.CityMetroStations) == 0 {
			c.SavePackageComplete(w, r)
		} else {
			c.EditPackageStep4(w, r)
			return
		}
		break
	default:
		http.NotFound(w, r.Request)
	}
}

func (c *Context) SavePackageStep4(w web.ResponseWriter, r *web.Request) {
	var err error
	c.Package, err = c.parsePackageFormStep4(w, r)
	if err != nil {
		c.Error = err.Error()
		c.ViewPackage = c.Package.ViewPackage()
		c.EditPackageStep4(w, r)
		return
	}

	switch c.Package.Type {
	case "drop", "drop preorder":
		c.SavePackageComplete(w, r)
	default:
		http.NotFound(w, r.Request)
	}
}

func (c *Context) SavePackageComplete(w web.ResponseWriter, r *web.Request) {
	if c.Package.Uuid == "" {
		c.Package.Uuid = util.GenerateUuid()
		c.Package.PackagePrice.Uuid = c.Package.Uuid
	}

	c.Package.Item = Item{}
	c.Package.CityMetroStation = CityMetroStation{}
	c.Package.GeoCity = City{}
	c.Package.GeoCountryFrom = Country{}
	c.Package.GeoCountryTo = Country{}

	err := c.Package.Save()
	if err != nil {
		c.Error = err.Error()
		c.ViewPackage = c.Package.ViewPackage()
		c.EditPackageStep2(w, r)
		return
	}
	err = c.Package.PackagePrice.Save()
	if err != nil {
		c.Error = err.Error()
		c.ViewPackage = c.Package.ViewPackage()
		c.EditPackageStep2(w, r)
		return
	}

	url := fmt.Sprintf("/user/%s/item/%s/", c.ViewSeller.Username, c.Item.Uuid)
	http.Redirect(w, r.Request, url, 302)
}

func (c *Context) parsePackageFormStep1(w web.ResponseWriter, r *web.Request) (Package, error) {
	currency := r.FormValue("currency")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		return c.Package, err
	}

	c.PackagePrice = r.FormValue("price")
	c.PackageCurrency = currency
	c.Package.PackagePrice.Currency = currency
	c.Package.PackagePrice.Price = price
	c.Package.ItemUuid = c.Item.Uuid
	c.Package.Item = c.Item
	c.Package.Type = r.FormValue("type")
	c.Package.Name = r.FormValue("name")

	err = c.Package.PreValidate(1)
	if err != nil {
		return c.Package, err
	}

	return c.Package, nil
}

func (c *Context) parsePackageFormStep2(w web.ResponseWriter, r *web.Request) (Package, error) {

	var (
		err error
	)

	c.Package, err = c.parsePackageFormStep1(w, r)
	if err != nil {
		return c.Package, err
	}

	if r.FormValue("country_name_en_to") != "" {
		countryTo, err := FindCountryByNameEn(r.FormValue("country_name_en_to"))
		if err != nil {
			return c.Package, errors.New("No such source country in database")
		}
		c.Package.CountryNameEnShippingTo = countryTo.NameEn
	} else if c.Package.Type == "drop" || c.Package.Type == "drop preorer" || c.Package.Type == "mail" {
		return c.Package, errors.New("Empty country_name_en_to")
	}

	if r.FormValue("country_name_en_from") != "" {
		countryTo, err := FindCountryByNameEn(r.FormValue("country_name_en_from"))
		if err != nil {
			return c.Package, errors.New("No such destination country in database")
		}
		c.Package.CountryNameEnShippingFrom = countryTo.NameEn
	} else if c.Package.Type == "mail" {
		return c.Package, errors.New("Empty country_name_en_from")
	}

	c.Package.Description = r.FormValue("description")

	err = c.Package.PreValidate(2)
	if err != nil {
		return c.Package, err
	}

	return c.Package, nil
}

func (c *Context) parsePackageFormStep3(w web.ResponseWriter, r *web.Request) (Package, error) {

	var (
		err error
	)

	c.Package, err = c.parsePackageFormStep2(w, r)
	if err != nil {
		return c.Package, err
	}

	if r.FormValue("city_id") != "" {
		cityId, err := strconv.ParseInt(r.FormValue("city_id"), 10, 64)
		if err != nil {
			return c.Package, err
		}

		city, err := FindCityByID(int(cityId))
		if err != nil {
			return c.Package, err
		}
		c.Package.DropCityId = city.ID
	} else {
		return c.Package, errors.New("Empty country_name_en_to")
	}
	c.Package.Description = r.FormValue("description")

	if coordinates := r.FormValue("coordinates"); coordinates != "" {
		parts := strings.Split(coordinates, ",")
		if len(parts) != 2 {
			return c.Package, errors.New("Wrong coordinates format")
		}

		longitude, err := strconv.ParseFloat(strings.Trim(parts[0], " "), 64)
		if err != nil {
			return c.Package, err
		}

		latitude, err := strconv.ParseFloat(strings.Trim(parts[1], " "), 64)
		if err != nil {
			return c.Package, err
		}

		c.Package.Longitude = longitude
		c.Package.Latitude = latitude
	}

	err = c.Package.PreValidate(3)
	if err != nil {
		return c.Package, err
	}

	return c.Package, nil
}

func (c *Context) parsePackageFormStep4(w web.ResponseWriter, r *web.Request) (Package, error) {

	var (
		err error
	)

	c.Package, err = c.parsePackageFormStep3(w, r)
	if err != nil {
		return c.Package, err
	}

	if r.FormValue("city_metro_station_uuid") != "" {
		cityMentroStationUuid := r.FormValue("city_metro_station_uuid")
		cityMetroStation, err := FindCityMetroStationByUuid(cityMentroStationUuid)
		if err != nil {
			return c.Package, err
		}
		c.Package.CityMetroStationUuid = cityMetroStation.Uuid
	} else {
		return c.Package, errors.New("Empty city_metro_station_uuid")
	}

	err = c.Package.PreValidate(3)
	if err != nil {
		return c.Package, err
	}

	return c.Package, nil
}

func (c *Context) DeletePackage(w web.ResponseWriter, r *web.Request) {
	c.Package.Remove()
	url := fmt.Sprintf("/user/%s/item/%s/", c.ViewSeller.Username, c.Item.Uuid)
	http.Redirect(w, r.Request, url, 302)
}

func (c *Context) PackageImage(w web.ResponseWriter, r *web.Request) {
	file, _ := os.Open("./data/images/" + c.Package.Uuid + ".jpeg")
	w.Header().Set("Content-type", "image/jpeg")
	w.Header().Set("Cache-control", "public, max-age=259200")
	io.Copy(w, file)
}

func (c *Context) DuplicatePackage(w web.ResponseWriter, r *web.Request) {
	c.Package.Uuid = util.GenerateUuid()
	c.Package.Save()
	c.Package.PackagePrice.Uuid = c.Package.Uuid
	c.Package.PackagePrice.Save()
	url := fmt.Sprintf("/user/%s/item/%s/", c.ViewSeller.Username, c.Item.Uuid)
	http.Redirect(w, r.Request, url, 302)
}
