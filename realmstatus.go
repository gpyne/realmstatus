package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

// PvPZone status values
var ZoneStatus = map[int]string{
	-1: "Unknown",
	0:  "Idle",
	1:  "Populating",
	2:  "Active",
	3:  "Concluded",
}

// PvPZone Controlling faction
var ZoneFaction = map[int]string{
	0: "Horde",
	1: "Alliance",
}

var RealmStatus = map[bool]string{
	true:  "up",
	false: "down",
}

var RealmQueue = map[bool]string{
	true:  "a",
	false: "no",
}

type PvPZone struct {
	Area               int   // Internal ID of the zone
	ControllingFaction int   `json:"controlling-faction"` // Which faction is controlling the zone
	Status             int   // Current status of the zone
	Next               int64 // Timestamp of when the next battle starts
}

type Realm struct {
	RealmType   string `json:"type"`
	Queue       bool
	Wintergrasp PvPZone
	TolBarad    PvPZone `json:"tol-barad"`
	Status      bool
	Population  string
	Name        string
	Slug        string
	Battlegroup string
}

type Response struct {
	Realms []Realm
}

var (
	region  = flag.String("region", "us", "specify the region to query")
	realms  = flag.String("realms", "", "comma seperated list of realms to display")
	verbose = flag.Bool("verbose", false, "display more info.")
)

func main() {
	flag.Parse()
	addr := "http://" + *region + ".battle.net/api/wow/realm/status?realms=" + *realms
	resp, err := http.Get(addr)
	if err != nil {
		log.Fatal(err)
	}

	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		log.Fatal(err)
	}

	for _, Realm := range r.Realms {
		fmt.Printf("\n %s is %s, it has %s population and %s queue.\n", Realm.Name, RealmStatus[Realm.Status], Realm.Population, RealmQueue[Realm.Queue])
		if *verbose {
			fmt.Printf(" - It is a %s realm and a member of the %s battlegroup.\n", Realm.RealmType, Realm.Battlegroup)
			fmt.Printf(" - Tol Barad is under %s control, the next battle begins at %s.\n", ZoneFaction[Realm.TolBarad.ControllingFaction], time.Unix(Realm.TolBarad.Next/1000, 0))
			fmt.Printf(" - Wintergrasp is under %s control, the next battle begins at %s.\n", ZoneFaction[Realm.Wintergrasp.ControllingFaction], time.Unix(Realm.Wintergrasp.Next/1000, 0))
		}
	}
	resp.Body.Close()
}
