package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"zgluco/internal/formatters"
	"zgluco/internal/sources/nightscout"
	"zgluco/internal/types"

	"github.com/atotto/clipboard"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	var (
		all      bool
		sections []string
		days     int
		nsURL    string
		nsKey    string
	)

	formatCmd := &cobra.Command{
		Use:   "format",
		Short: "Export Nightscout data as formatted text",
		Example: `  zgluco format --all --days 7
  zgluco format --all --days 7 --nightscout-url https://fqdn.tld --nightscout-api-key "apikey"
  zgluco format --sections sgv,treatments -n 5 --nightscout-url https://fqdn.tld`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = godotenv.Load()

			if nsURL == "" {
				nsURL = os.Getenv("NIGHTSCOUT_URL")
			}
			if nsKey == "" {
				nsKey = os.Getenv("NIGHTSCOUT_API_KEY")
			}
			if nsURL == "" {
				return fmt.Errorf("nightscout URL required: use --nightscout-url or set NIGHTSCOUT_URL in .env")
			}

			if all && len(sections) > 0 {
				return fmt.Errorf("--all and --sections are mutually exclusive")
			}

			var targets []string
			switch {
			case all:
				targets = []string{"profile", "sgv", "treatments"}
			case len(sections) > 0:
				targets = sections
			default:
				return fmt.Errorf("specify --all or --sections sgv,treatments,profile")
			}

			ns, err := nightscout.New(nsURL, nsKey)
			if err != nil {
				return fmt.Errorf("error initializing Nightscout client: %w", err)
			}

			tr := types.NewTimeRange(days)
			loc := time.Now().Location()

			p, err := ns.FetchProfile(tr)
			if err != nil {
				return fmt.Errorf("could not fetch profile: %w", err)
			}

			var sb strings.Builder

			for i, target := range targets {
				if i > 0 {
					sb.WriteString("\n\n")
				}
				switch target {
				case "sgv":
					sgvs, err := ns.FetchSensorGlucoseValues(tr)
					if err != nil {
						return fmt.Errorf("could not fetch glucose values: %w", err)
					}
					if err := formatters.FormatSGVs(&sb, p, sgvs, loc); err != nil {
						return fmt.Errorf("could not format glucose data: %w", err)
					}
				case "treatments":
					t, err := ns.FetchTreatments(tr)
					if err != nil {
						return fmt.Errorf("could not fetch treatments: %w", err)
					}
					formatters.FormatTreatments(&sb, p, t, tr)
				case "profile":
					formatters.FormatProfile(&sb, p, tr, loc)
				default:
					return fmt.Errorf("unknown section: %s", target)
				}
			}

			fmt.Println(sb.String())

			if err = clipboard.WriteAll(sb.String()); err != nil {
				return fmt.Errorf("could not copy to clipboard: %w", err)
			}

			return nil
		},
	}

	f := formatCmd.Flags()
	f.BoolVar(&all, "all", false, "export all sections (profile, sgv, treatments)")
	f.StringSliceVar(&sections, "sections", nil, "comma-separated sections to export: sgv, treatments, profile")
	f.IntVar(&days, "days", 7, "number of days to look back")
	f.IntVarP(&days, "n", "n", 7, "alias for --days")
	f.StringVar(&nsURL, "nightscout-url", "", "Nightscout base URL (overrides NIGHTSCOUT_URL env)")
	f.StringVar(&nsKey, "nightscout-api-key", "", "Nightscout API key (overrides NIGHTSCOUT_API_KEY env; omit for public instances)")

	root := &cobra.Command{Use: "zgluco"}
	root.AddCommand(formatCmd)

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
