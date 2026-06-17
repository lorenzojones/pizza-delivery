// Package excel parses per-person sales metrics from an .xlsx spreadsheet.
// Column names are configurable rather than hard-coded so the same parser works
// across different teams' spreadsheets.
package excel

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"

	"wallboard/internal/config"
	"wallboard/internal/models"
)

// Result is the outcome of parsing a spreadsheet.
type Result struct {
	// BySpreadsheetName maps the normalised person name to their sales metrics.
	BySpreadsheetName map[string]models.SalesMetrics
	// Rows is the number of data rows successfully read.
	Rows int
	// SkippedRows counts rows that were ignored (e.g. missing a name).
	SkippedRows int
}

// ParseFile opens the .xlsx at path and extracts sales metrics using the
// configured column names.
func ParseFile(path string, cols config.ExcelColumns) (*Result, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()
	return parse(f, cols)
}

// ParseBytes parses an .xlsx provided as raw bytes (used by the upload handler).
func ParseBytes(b []byte, cols config.ExcelColumns) (*Result, error) {
	f, err := excelize.OpenReader(strings.NewReader(string(b)))
	if err != nil {
		return nil, fmt.Errorf("reading uploaded xlsx: %w", err)
	}
	defer f.Close()
	return parse(f, cols)
}

func parse(f *excelize.File, cols config.ExcelColumns) (*Result, error) {
	sheet := cols.Sheet
	if sheet == "" {
		sheet = f.GetSheetName(0)
	}
	if sheet == "" {
		return nil, fmt.Errorf("spreadsheet has no sheets")
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("reading sheet %q: %w", sheet, err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("sheet %q has no data rows", sheet)
	}

	header := rows[0]
	idx := indexHeaders(header)

	nameCol, ok := idx[norm(cols.Name)]
	if !ok {
		return nil, fmt.Errorf("name column %q not found; available columns: %s",
			cols.Name, strings.Join(header, ", "))
	}
	// Optional columns: -1 when absent.
	countCol := lookup(idx, cols.SalesCount)
	valueCol := lookup(idx, cols.SalesValue)
	targetCol := lookup(idx, cols.Target)

	res := &Result{BySpreadsheetName: make(map[string]models.SalesMetrics)}

	for _, row := range rows[1:] {
		name := strings.TrimSpace(cell(row, nameCol))
		if name == "" {
			res.SkippedRows++
			continue
		}
		var m models.SalesMetrics
		if countCol >= 0 {
			m.SalesCount = parseInt(cell(row, countCol))
		}
		if valueCol >= 0 {
			m.SalesValue = parseNumber(cell(row, valueCol))
		}
		if targetCol >= 0 {
			m.Target = parseNumber(cell(row, targetCol))
		}
		m.TargetProgress = computeProgress(m)
		res.BySpreadsheetName[normName(name)] = m
		res.Rows++
	}
	return res, nil
}

// computeProgress derives 0..1+ progress toward target, preferring sales value
// when present and falling back to sales count.
func computeProgress(m models.SalesMetrics) float64 {
	if m.Target <= 0 {
		return 0
	}
	if m.SalesValue > 0 {
		return m.SalesValue / m.Target
	}
	if m.SalesCount > 0 {
		return float64(m.SalesCount) / m.Target
	}
	return 0
}

// indexHeaders builds a normalised header->column-index map.
func indexHeaders(header []string) map[string]int {
	idx := make(map[string]int, len(header))
	for i, h := range header {
		key := norm(h)
		if key == "" {
			continue
		}
		if _, exists := idx[key]; !exists {
			idx[key] = i
		}
	}
	return idx
}

func lookup(idx map[string]int, name string) int {
	if name == "" {
		return -1
	}
	if i, ok := idx[norm(name)]; ok {
		return i
	}
	return -1
}

func cell(row []string, i int) string {
	if i < 0 || i >= len(row) {
		return ""
	}
	return row[i]
}

// norm normalises a header/key for case- and whitespace-insensitive matching.
func norm(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// NormName normalises a person's name for matching against the staff mapping.
// Exported so the aggregator uses the exact same rule.
func NormName(s string) string { return normName(s) }

func normName(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(s), " "))
}

// parseInt parses an integer, tolerating currency/grouping symbols and decimals.
func parseInt(s string) int {
	return int(parseNumber(s) + 0.0001)
}

// parseNumber parses a float from a messy spreadsheet cell, stripping currency
// symbols, thousands separators and surrounding whitespace.
func parseNumber(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var b strings.Builder
	for _, r := range s {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' {
			b.WriteRune(r)
		}
		// commas, £, $, €, %, spaces and other symbols are dropped
	}
	cleaned := b.String()
	if cleaned == "" || cleaned == "-" || cleaned == "." {
		return 0
	}
	var f float64
	if _, err := fmt.Sscanf(cleaned, "%g", &f); err != nil {
		return 0
	}
	return f
}
