package front

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func apiTorrentInfo(w http.ResponseWriter, req *http.Request) {
	log.Println("Downloading torrent info...")
	log.Println(os.Getwd())
	name := "pte-" + req.URL.Query().Get("torrentid") + ".torrent"
	if _, err := os.Stat(name); err != nil {
		r, err := http.NewRequest("GET", "https://api.pte.nu/torrents/download/"+req.URL.Query().Get("torrentid"), nil)
		if err != nil {
			log.Println(err)
			return
		}
		r.Header.Set("API-Key", req.URL.Query().Get("apikey"))
		client := &http.Client{}
		resp, err := client.Do(r)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "appplication/json")
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		ioutil.WriteFile(name, content, 0750)
	}
	T, err := cl.AddTorrentFromFile(name)
	if err != nil {
		log.Println(err)
		return
	}
	body, err := json.Marshal(T.Info())
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
}
