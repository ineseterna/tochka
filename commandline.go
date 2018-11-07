package main

import (
	"fmt"

	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/marketplace"
	"qxklmrhx7qkzais6.onion/Tochka/tochka-free-market/modules/util"
)

func manageRole(username, action, role string) {
	user, _ := marketplace.FindUserByUsername(username)
	if user == nil {
		fmt.Println("No such user")
		return
	}
	if action == "grant" && role == "seller" {
		user.IsSeller = !user.IsSeller
	} else if action == "grant" && role == "admin" {
		user.IsAdmin = !user.IsAdmin
	} else {
		fmt.Println("Wrong action")
		return
	}
	user.Save()
}

func indexItems() {
	util.Log.Debug("[Index] Indexing items...")
	for _, item := range marketplace.GetAllItems() {
		util.Log.Debug("[Index] %s", item.Name)
		err := item.Index()
		if err != nil {
			util.Log.Error("%s", err)
		}
	}
}

func syncModels() {
	marketplace.SyncModels()
}

func staffStats() {
	interval := "2018-09-07 13:32"
	sTickets, err := marketplace.StaffSupportTicketsResolutionStats(interval)
	if err != nil {
		return
	}
	sDisputes, err := marketplace.StaffSupportDisputesResolutionStats(interval)
	if err != nil {
		return
	}
	sItems, err := marketplace.StaffItemApprovalStats(interval)
	if err != nil {
		return
	}

	var (
		text = fmt.Sprintf(
			`
Support Agent | Ticket Status | Number Of Tickets
--- | --- | ---
`)
	)
	for _, si := range sTickets {
		text += fmt.Sprintf("%s | TICKET %s | %d\n", si.ResolverUsername, si.CurrentStatus, si.TicketCount)
	}
	for _, si := range sDisputes {
		text += fmt.Sprintf("%s | DISPUTE %s | %d\n", si.ResolverUsername, si.CurrentStatus, si.TicketCount)
	}
	for _, si := range sItems {
		text += fmt.Sprintf("%s | ITEM %s | %d\n", si.ResolverUsername, si.CurrentStatus, si.TicketCount)
	}

	println(text)
}

func importMetroStations() {
	marketplace.ImportCityMetroStations(524901, "./dumps/moscow-metro.json")
	marketplace.ImportCityMetroStations(498817, "./dumps/spb-metro.json")
}
