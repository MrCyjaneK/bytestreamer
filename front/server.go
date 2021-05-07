package front

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/anacrolix/torrent"
)

//go:embed player.html index.html scripts/* info/* styles.css
var content embed.FS

var cl *torrent.Client

func Start() {
	var err error
	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	//os.RemoveAll(dir + "/ByteStreamer")
	config := torrent.NewDefaultClientConfig()
	config.NoDHT = true
	config.DataDir = dir + "/ByteStreamer"
	config.Seed = true
	config.HTTPUserAgent = "ByteStreamer v0.0.0"
	config.ExtendedHandshakeClientVersion = "ByteStreamer v0.0.0"
	config.Bep20 = "-BY0000-"
	cl, err = torrent.NewClient(config)
	if err != nil {
		log.Panic(err)
	}
	http.Handle("/", http.FileServer(http.FS(content)))
	http.HandleFunc("/play", apiPlay)
	http.HandleFunc("/api/info/pte", apiInfoPte)
	http.HandleFunc("/api/torrent_info", apiTorrentInfo)
	http.HandleFunc("/streamfile", streamfile)
	http.ListenAndServe(":8081", nil)
}
