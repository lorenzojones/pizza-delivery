// Package config loads all runtime configuration from environment variables
// (optionally seeded from a .env file) plus the staff-mapping JSON file.
package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"wallboard/internal/models"
)

// Config is the fully-resolved application configuration.
type Config struct {
	// HTTP
	Port string

	// Aircall
	AircallAPIID    string
	AircallAPIToken string
	AircallBaseURL  string

	// Data sources
	SalesXLSXPath    string
	StaffMappingPath string

	// Behaviour
	RefreshInterval time.Duration
	StaleAfter      time.Duration
	DashboardTitle  string
	Timezone        string
	Location        *time.Location

	// Scoring
	Weights models.ScoreWeights

	// Excel column mapping (configurable, not hard-coded)
	Excel ExcelColumns

	// Loaded staff list
	Staff []models.StaffMember
}

// ExcelColumns names the spreadsheet columns to read. Header matching is
// case-insensitive and whitespace-trimmed (see internal/excel).
type ExcelColumns struct {
	Sheet       string // optional; empty = first sheet
	Name        string
	SalesCount  string
	SalesValue  string
	Target      string
}

// DemoMode reports whether the Aircall client should synthesise data because no
// API credentials were supplied.
func (c *Config) DemoMode() bool {
	return c.AircallAPIID == "" || c.AircallAPIToken == ""
}

// Load resolves configuration from the environment. It first attempts to read a
// .env file (path from ENV_FILE, default ".env") if present, then reads env
// vars, applying sensible defaults. It also loads the staff mapping file.
func Load() (*Config, error) {
	loadDotEnv(getenv("ENV_FILE", ".env"))

	cfg := &Config{
		Port:             getenv("PORT", "8080"),
		AircallAPIID:     os.Getenv("AIRCALL_API_ID"),
		AircallAPIToken:  os.Getenv("AIRCALL_API_TOKEN"),
		AircallBaseURL:   getenv("AIRCALL_BASE_URL", "https://api.aircall.io/v1"),
		SalesXLSXPath:    getenv("SALES_XLSX_PATH", "data/sales-sample.xlsx"),
		StaffMappingPath: getenv("STAFF_MAPPING_PATH", "config/staff-mapping.json"),
		DashboardTitle:   getenv("DASHBOARD_TITLE", "Sales Wallboard"),
		Timezone:         getenv("DASHBOARD_TIMEZONE", "Europe/London"),
		Weights: models.ScoreWeights{
			Sales:          getenvFloat("WEIGHT_SALES", 0.50),
			CallsOut:       getenvFloat("WEIGHT_CALLS_OUT", 0.25),
			CallsIn:        getenvFloat("WEIGHT_CALLS_IN", 0.10),
			MinutesOnPhone: getenvFloat("WEIGHT_MINUTES", 0.15),
		},
		Excel: ExcelColumns{
			Sheet:      os.Getenv("EXCEL_SHEET"),
			Name:       getenv("EXCEL_COL_NAME", "Name"),
			SalesCount: getenv("EXCEL_COL_SALES_COUNT", "Sales Count"),
			SalesValue: getenv("EXCEL_COL_SALES_VALUE", "Sales Value"),
			Target:     getenv("EXCEL_COL_TARGET", "Target"),
		},
	}

	refreshSecs := getenvInt("REFRESH_INTERVAL_SECONDS", 300)
	if refreshSecs < 10 {
		refreshSecs = 10 // guard against hammering upstreams
	}
	cfg.RefreshInterval = time.Duration(refreshSecs) * time.Second

	// Data is considered stale once it is older than ~2.5 refresh cycles.
	staleSecs := getenvInt("STALE_AFTER_SECONDS", refreshSecs*5/2)
	cfg.StaleAfter = time.Duration(staleSecs) * time.Second

	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		// Fall back to UTC rather than failing to boot.
		loc = time.UTC
		cfg.Timezone = "UTC"
	}
	cfg.Location = loc

	staff, err := LoadStaff(cfg.StaffMappingPath)
	if err != nil {
		return nil, fmt.Errorf("loading staff mapping: %w", err)
	}
	cfg.Staff = staff

	return cfg, nil
}

// staffFile is the on-disk shape of the staff mapping JSON.
type staffFile struct {
	Staff []models.StaffMember `json:"staff"`
}

// LoadStaff reads and validates the staff mapping JSON file.
func LoadStaff(path string) ([]models.StaffMember, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var sf staffFile
	if err := json.Unmarshal(b, &sf); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if len(sf.Staff) == 0 {
		return nil, fmt.Errorf("staff mapping %s contains no staff", path)
	}
	for i := range sf.Staff {
		if strings.TrimSpace(sf.Staff[i].ID) == "" {
			return nil, fmt.Errorf("staff entry %d is missing an id", i)
		}
		if strings.TrimSpace(sf.Staff[i].SpreadsheetName) == "" {
			// Default the spreadsheet name to the display name so a minimal
			// mapping still works.
			sf.Staff[i].SpreadsheetName = sf.Staff[i].Name
		}
	}
	return sf.Staff, nil
}

// --- env helpers ---

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return n
		}
	}
	return def
}

func getenvFloat(key string, def float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			return f
		}
	}
	return def
}

// loadDotEnv reads simple KEY=VALUE lines from a .env file if it exists. Values
// already present in the real environment win, so the file only fills gaps.
// This intentionally avoids an external dependency.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // no .env file is fine
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		val = strings.Trim(val, `"'`)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, val)
		}
	}
}
