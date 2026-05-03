package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/calibrarytracker/calibrarytracker/internal/scraper"
)

func main() {
	outputDir := flag.String("output", "data", "output directory for JSON files")
	overdriveOnly := flag.Bool("overdrive", false, "only update digital-access.json from OverDrive API")
	flag.Parse()

	if err := os.MkdirAll(*outputDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating output dir: %v\n", err)
		os.Exit(1)
	}

	digitalAccessPath := filepath.Join(*outputDir, "digital-access.json")

	if *overdriveOnly {
		runOverdriveUpdate(*outputDir, digitalAccessPath)
		return
	}

	fmt.Println("Fetching administrative headquarters (Excel XML)...")
	systems, err := scraper.ParseAdminHQExcel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Excel failed (%v), falling back to HTML...\n", err)
		systems, err = scraper.ParseAdminHQ()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing admin HQ: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("  Found %d library systems\n", len(systems))

	fmt.Println("Fetching branch/outlet directory (Excel XML)...")
	branches, err := scraper.ParseBranchesExcel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Excel failed (%v), falling back to HTML...\n", err)
		branches, err = scraper.ParseBranches()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing branches: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("  Found %d branches\n", len(branches))

	fmt.Println("Cross-referencing systems and branches...")
	scraper.CrossReference(systems, branches)

	fmt.Println("Loading cooperatives for CLSA membership...")
	cooperativesPath := filepath.Join(*outputDir, "cooperatives.json")
	if err := scraper.PopulateCLSA(cooperativesPath, systems); err != nil {
		fmt.Fprintf(os.Stderr, "  warning: %v (skipping CLSA population)\n", err)
	}

	fmt.Println("Loading digital access data...")
	if err := scraper.LoadDigitalAccess(digitalAccessPath, systems); err != nil {
		fmt.Fprintf(os.Stderr, "  warning: %v (skipping digital access enrichment)\n", err)
	}

	scraper.SortSystems(systems)
	scraper.SortBranches(branches)

	fmt.Println("Writing output files...")
	if err := scraper.WriteJSON(*outputDir, "systems.json", systems); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if err := scraper.WriteJSON(*outputDir, "branches.json", branches); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Done! Wrote %d systems and %d branches to %s/\n", len(systems), len(branches), *outputDir)
}

func runOverdriveUpdate(outputDir, digitalAccessPath string) {
	fmt.Println("Fetching systems for matching...")
	systems, err := scraper.ParseAdminHQExcel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Excel failed (%v), falling back to HTML...\n", err)
		systems, err = scraper.ParseAdminHQ()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing admin HQ: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Fetching OverDrive libraries...")
	odLibs, err := scraper.FetchOverdriveLibraries()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching OverDrive libraries: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  Found %d OverDrive libraries\n", len(odLibs))

	fmt.Println("Matching to systems...")
	if err := scraper.MatchOverdriveToSystems(odLibs, systems, digitalAccessPath); err != nil {
		fmt.Fprintf(os.Stderr, "error updating digital access: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done! Updated", digitalAccessPath)
}
