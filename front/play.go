package front

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"git.mrcyjanek.net/mrcyjanek/goprod/goprod"
	"github.com/anacrolix/torrent"
)

// SeekableContent describes an io.ReadSeeker that can be closed as well.
type SeekableContent interface {
	io.ReadSeeker
	io.Closer
}

// FileEntry helps reading a torrent file.
type FileEntry struct {
	*torrent.File
	torrent.Reader
}

func apiPlay(w http.ResponseWriter, req *http.Request) {
	if runtime.GOOS == "android" {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(`<script>
		
<script>
	window.location.href = "/player.html"
</script>`))
		var intent goprod.Intent
		intent.Type = "android-intent"
		intent.Data.Package = "org.videolan.vlc"
		intent.Data.URI = `http://127.0.0.1:8081/streamfile?torrentid=` + url.QueryEscape(req.URL.Query().Get("torrentid")) + `&file=` + url.QueryEscape(req.URL.Query().Get("file"))
		intent.Data.CustomComponent = true
		intent.Data.Component.PKG = intent.Data.Package
		intent.Data.Component.CLS = "org.videolan.vlc.gui.video.VideoPlayerActivity"
		goprod.CallIntent(intent)
		return
	}
	w.Header().Add("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Add("Content-Disposition", `filename=`+url.QueryEscape(req.URL.Query().Get("file"))+`.m3u8`)
	w.Write([]byte(`#EXTM3U
#EXTINF:-1,LiveStream
http://127.0.0.1:8081/streamfile?torrentid=` + url.QueryEscape(req.URL.Query().Get("torrentid")) + `&file=` + url.QueryEscape(req.URL.Query().Get("file"))))
}

func streamfile(w http.ResponseWriter, req *http.Request) {
	//log.Println(req.Header)
	filepath := req.URL.Query().Get("file")
	_, err := cl.AddTorrentFromFile("pte-" + req.URL.Query().Get("torrentid") + ".torrent")
	if err != nil {
		log.Println(err)
		return
	}
	torrents := cl.Torrents()
	var file *torrent.File
	for i := range torrents {
		files := torrents[i].Files()
		for j := range files {
			if files[j].DisplayPath() == filepath {
				file = files[j]
			}
		}
	}
	//file.Download()
	//firstPieceIndex := file.Offset() * int64(T.NumPieces()) / T.Length()
	//endPieceIndex := (file.Offset() + file.Length()) * int64(T.NumPieces()) / T.Length()
	//for idx := firstPieceIndex; idx <= endPieceIndex*5/100; idx++ {
	//	T.Piece(int(idx)).SetPriority(torrent.PiecePriorityNow)
	//}
	//entry, err := newFileReader(file)
	//if err != nil {
	//	log.Println(err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//defer func() {
	//	if err := entry.Close(); err != nil {
	//		log.Println(err)
	//		log.Printf("Error closing file reader: %s\n", err)
	//	}
	//}()
	w.Header().Set("Content-Disposition", "attachment; filename=\""+file.FileInfo().Path[len(file.FileInfo().Path)-1]+"\"")
	//log.Println("checking")
	//r := strings.Split(strings.Split(req.Header.Get("Range"), "bytes=")[1], "-")
	//rstart, err := strconv.ParseInt(r[0], 10, 64)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//rstart = rstart / file.State()[0].Bytes
	//rend, err := strconv.ParseInt(r[1], 10, 64)
	//if err != nil {
	//	rend = int64(len(file.State())) - 1
	//} else {
	//	rend = rend / file.State()[0].Bytes
	//}
	//T.DownloadPieces(int(rstart), int(rend))
	//log.Println(r, rstart, file.State()[0].Bytes, file.State()[rstart].PieceState.Complete)
	log.Println("serving content")

	// Don't wait for pieces to complete and be verified.
	//reader.SetResponsive()

	size := file.Length()
	reader := file.NewReader()
	if size > 0 {
		// Read ahead 3% of file.
		reader.SetReadahead((size * 3) / 100)
	}

	w.Header().Set("Content-Disposition", `filename="`+file.FileInfo().Path[len(file.FileInfo().Path)-1]+`"`)

	_, err = reader.Seek(0, 0)
	if err != nil {
		log.Println(err)
		return
	}

	http.ServeContent(w, req, "", time.Unix(0, 0), reader)
	log.Println("serving content.end")
}

func newFileReader(f *torrent.File) (SeekableContent, error) {
	torrent := f.Torrent()
	reader := torrent.NewReader()

	// We read ahead 1% of the file continuously.
	reader.SetReadahead(f.Length() / 100)
	reader.SetResponsive()
	_, err := reader.Seek(f.Offset(), io.SeekStart)
	if err != nil {
		log.Println(err)
	}
	return &FileEntry{
		File:   f,
		Reader: reader,
	}, err
}
