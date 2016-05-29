package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/osheroff/spotcontrol"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Device struct {
	Name   string `json:"name"`
	Ident  string `json:"ident"`
	Volume int    `json:"volume"`
}

var sController *spotcontrol.SpircController

func deviceList(w http.ResponseWriter, r *http.Request) {
	controllerDevices := sController.ListDevices()
	jsonDevices := make([]Device, len(controllerDevices))

	for i, d := range controllerDevices {
		jsonDevices[i].Name = d.Name
		jsonDevices[i].Ident = d.Ident
		jsonDevices[i].Volume = int(d.Volume)
	}

	json, _ := json.Marshal(jsonDevices)
	w.Write(json)
}

func setVolume(w http.ResponseWriter, r *http.Request) {
	ident := r.FormValue("ident")
	volume := r.FormValue("volume")

	if ident == "" || volume == "" {
		fmt.Fprintf(w, "{ \"status\": \"ERR\" }")
	} else {
		vol, _ := strconv.ParseUint(volume, 10, 32)
		sController.SendVolume(ident, uint32(vol))
		fmt.Fprintf(w, "{ \"status\": \"OK\" }")
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	index, _ := ioutil.ReadFile("index.html")
	w.Write(index)
}

func main() {
	username := flag.String("username", "", "spotify username")
	password := flag.String("password", "", "spotify password")
	port := flag.Int("port", 8080, "http port to listen on")
	appkey := flag.String("appkey", "./spotify_appkey.key", "spotify appkey file path")
	devicename := "spotcontrol-webvolume"
	flag.Parse()

	if *username != "" && *password != "" && *appkey != "" {
		sController = spotcontrol.Login(*username, *password, *appkey, devicename)

		http.HandleFunc("/", index)
		http.HandleFunc("/devices", deviceList)
		http.HandleFunc("/set_volume", setVolume)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
	} else {
		flag.PrintDefaults()
		fmt.Println("need to supply a username and password")
		fmt.Println("./spotcontrol-webvolume --username SPOTIFY_USERNAME --password SPOTIFY_PASSWORD")
		return
	}

}
