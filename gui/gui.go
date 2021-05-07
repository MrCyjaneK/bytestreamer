// +build gui

package gui

import (
	"log"

	"github.com/pkg/browser"
)

func Start() {
	log.Println("Starting gui...")
	browser.OpenURL("http://127.0.0.1:8081")
}
