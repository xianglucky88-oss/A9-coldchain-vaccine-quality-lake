package main

import (
	"flag"
	"fmt"
	"os"

	"example.com/coldchain-lake/internal/demo"
	"example.com/coldchain-lake/internal/engine"
	"example.com/coldchain-lake/internal/model"
	"example.com/coldchain-lake/internal/service"
	"example.com/coldchain-lake/internal/store"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		usage()
		return nil
	}
	switch args[0] {
	case "ingest":
		return runIngest(args[1:])
	case "compute":
		return runCompute(args[1:])
	case "demo":
		return runDemo(args[1:])
	case "help", "-h", "--help":
		usage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func usage() {
	fmt.Print(`coldchain-lake - vaccine cold-chain quality lake

Commands:
  ingest  --samples FILE [--events FILE] --data-dir DIR
  compute --data-dir DIR [--min-temp C] [--max-temp C]
  demo    [--data-dir DIR]
`)
}

func runIngest(args []string) error {
	fs := flag.NewFlagSet("ingest", flag.ContinueOnError)
	samplesPath := fs.String("samples", "", "sample JSONL file")
	eventsPath := fs.String("events", "", "door event JSONL file")
	dataDir := fs.String("data-dir", "data", "lake root")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *samplesPath == "" && *eventsPath == "" {
		return fmt.Errorf("at least one of --samples or --events is required")
	}
	lake := service.New(*dataDir)
	if *samplesPath != "" {
		rows, err := store.ReadJSONL[model.Sample](*samplesPath)
		if err != nil {
			return err
		}
		if err := lake.IngestSamples(rows); err != nil {
			return err
		}
	}
	if *eventsPath != "" {
		rows, err := store.ReadJSONL[model.DoorEvent](*eventsPath)
		if err != nil {
			return err
		}
		if err := lake.IngestEvents(rows); err != nil {
			return err
		}
	}
	return writeJSON(map[string]string{"status": "ok"})
}

func runCompute(args []string) error {
	fs := flag.NewFlagSet("compute", flag.ContinueOnError)
	dataDir := fs.String("data-dir", "data", "lake root")
	minTemp := fs.Float64("min-temp", 2, "minimum safe temperature")
	maxTemp := fs.Float64("max-temp", 8, "maximum safe temperature")
	if err := fs.Parse(args); err != nil {
		return err
	}
	result, err := service.New(*dataDir).Compute(engine.Config{
		MinTempC:            *minTemp,
		MaxTempC:            *maxTemp,
		MinExcursionMs:      60_000,
		AttributionWindowMs: 5 * 60_000,
		LoggerOffsetsMs:     map[string]int64{"LOGGER-A": 30_000},
	})
	if err != nil {
		return err
	}
	return writeJSON(result)
}

func runDemo(args []string) error {
	fs := flag.NewFlagSet("demo", flag.ContinueOnError)
	dataDir := fs.String("data-dir", ".demo-data", "demo lake root")
	if err := fs.Parse(args); err != nil {
		return err
	}
	lake := service.New(*dataDir)
	if err := lake.IngestSamples(demo.Samples()); err != nil {
		return err
	}
	if err := lake.IngestEvents(demo.Events()); err != nil {
		return err
	}
	result, err := lake.Compute(engine.Config{
		MinTempC:            2,
		MaxTempC:            8,
		MinExcursionMs:      60_000,
		AttributionWindowMs: 5 * 60_000,
		LoggerOffsetsMs:     map[string]int64{"LOGGER-A": 30_000},
	})
	if err != nil {
		return err
	}
	demo.PrintSummary(len(result.Quality), len(result.Excursions))
	return writeJSON(result)
}
