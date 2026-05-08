package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"zgluco/internal/formatters"
	"zgluco/internal/sources/nightscout"
	"zgluco/internal/types"

	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ns, err := nightscout.New(os.Getenv("NIGHTSCOUT_URL"), os.Getenv("NIGHTSCOUT_API_KEY"))
	if err != nil {
		log.Fatalf("error initializing Nightscout client: %v", err)
	}

	tr := types.NewTimeRange(30)

	p, err := ns.FetchProfile(tr)
	if err != nil {
		log.Fatalf("could not fetch profile: %v", err)
	}

	t, err := ns.FetchTreatments(tr)
	if err != nil {
		log.Fatalf("could not fetch treatments: %v", err)
	}

	var sb strings.Builder

	formatters.FormatTreatments(&sb, p, t, tr)

	fmt.Println(sb.String())

	err = clipboard.WriteAll(sb.String())
	if err != nil {
		log.Fatalf("could not copy to clipboard: %v", err)
	}

}
