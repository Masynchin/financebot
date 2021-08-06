package main

import (
	"log"
	"os"

	"github.com/NicoNex/echotron/v3"
	_ "github.com/joho/godotenv/autoload"
)

var token = os.Getenv("BOT_TOKEN")
var expenceService *ExpenceService

func main() {
	es, err := NewExpenceService("expences.db")
	if err != nil {
		log.Fatal(err)
	}
	expenceService = es
	defer expenceService.Close()

	dsp := echotron.NewDispatcher(token, newBot)
	log.Println(dsp.Poll())
}
