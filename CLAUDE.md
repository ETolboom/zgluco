# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the unified entrypoint
go run ./cmd/zgluco format --all
go run ./cmd/zgluco format sgv
go run ./cmd/zgluco format treatments
go run ./cmd/zgluco format profile

# Or the individual binaries (thin wrappers around the same format package)
go run ./cmd/sgv
go run ./cmd/treatments
go run ./cmd/profile

# Build all
go build ./...

# Run tests
go test ./...
```

## Configuration

All commands require a `.env` file in the project root:

```
NIGHTSCOUT_URL=https://your-nightscout-instance.example.com
NIGHTSCOUT_API_KEY=your_api_key
```

## Architecture

This is a CLI toolset that fetches diabetes management data from a Nightscout instance and formats it for display + clipboard copy.

**Entry points** — `cmd/zgluco` is the primary binary (`zgluco format [--all|sgv|treatments|profile]`). `cmd/sgv`, `cmd/treatments`, `cmd/profile` are thin standalone wrappers kept for convenience. All share the formatting logic from `internal/format/`.

**Source abstraction** — `internal/sources/source.go` defines the `Source` interface (`FetchProfile`, `FetchTreatments`, `FetchSensorGlucoseValues`). The only current implementation is `internal/sources/nightscout/`.

**Nightscout package layout:**
- `struct.go` — raw JSON API structs (`Profile`, `Treatment`, `Sgv`)
- `client.go` — public-facing `Nightscout` struct and its methods
- `profile.go`, `treatments.go`, `sgv.go` — private fetch + parse logic
- `helpers.go` — shared HTTP utilities (`doApiCall`)

**Domain models** — `internal/models/` holds the canonical domain types that the Nightscout package maps into:
- `SensorGlucoseValue`, `GlucoseDirection` in `sgv.go`
- `Treatment` interface + concrete treatment types (SMB, TempBasal, Bolus, etc.) under `treatments/`
- `profile.Profile` (with `BasalRates`, `CarbRatios`, `InsulinSensitivityFactors`, `GlucoseTargets`, `Changes`) under `profile/`

**Unit handling** — SGV values are stored internally in mg/dL. Display code in `internal/format/sgv.go` and `internal/format/profile.go` divides by 18 when `p.PreferredUnits == models.Mmol`. Known issue: `TemporaryTarget` still displays in mg/dL regardless of preferred unit.

**Profile changelog** — `sources/nightscout/profile.go` compares consecutive Nightscout profile documents to build a `[]profile.Change` showing what changed between profile versions.