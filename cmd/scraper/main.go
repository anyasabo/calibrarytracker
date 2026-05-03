package main

import (
	"fmt"
	"os"

	"github.com/calibrarytracker/calibrarytracker/internal/scraper"
)

func main() {
	outputDir := "data"
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output dir: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Fetching administrative headquarters directory...")
	systems, err := scraper.ParseAdminHQ()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing admin HQ: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  Found %d library systems\n", len(systems))

	fmt.Println("Fetching branch/outlet directory...")
	branches, err := scraper.ParseBranches()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing branches: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  Found %d branches\n", len(branches))

	fmt.Println("Cross-referencing systems and branches...")
	scraper.CrossReference(systems, branches)

	scraper.SortSystems(systems)
	scraper.SortBranches(branches)

	fmt.Println("Writing output files...")
	if err := scraper.WriteJSON(outputDir, "systems.json", systems); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := scraper.WriteJSON(outputDir, "branches.json", branches); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done! Wrote %d systems and %d branches to %s/\n", len(systems), len(branches), outputDir)
}
