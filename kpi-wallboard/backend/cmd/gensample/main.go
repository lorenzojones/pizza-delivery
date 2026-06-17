// Command gensample writes an example sales spreadsheet (data/sales-sample.xlsx)
// matching the default column mapping, so the wallboard has real data to parse
// out of the box. Run with: go run ./cmd/gensample
package main

import (
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	const path = "data/sales-sample.xlsx"
	const sheet = "Sales"

	f := excelize.NewFile()
	defer f.Close()

	idx, err := f.NewSheet(sheet)
	if err != nil {
		log.Fatalf("new sheet: %v", err)
	}
	f.SetActiveSheet(idx)
	f.DeleteSheet("Sheet1")

	// Header row matches the default EXCEL_COL_* configuration.
	headers := []string{"Name", "Sales Count", "Sales Value", "Target"}
	for c, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(c+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// One row per salesperson. Names match config/staff-mapping.example.json.
	// Note: "Priya Patel" is intentionally absent to demonstrate the
	// missing-sales-data indicator on the wallboard.
	type row struct {
		name   string
		count  int
		value  int
		target int
	}
	rows := []row{
		{"Jane Doe", 7, 18250, 20000},
		{"Marcus Lee", 5, 14100, 20000},
		{"Aisha Khan", 9, 22400, 20000},
		{"Tom Brennan", 3, 8600, 20000},
		{"Sofia Rossi", 6, 16750, 20000},
		// Priya Patel deliberately omitted (unmatched / missing data demo).
		{"Unknown Person", 2, 4000, 20000}, // in sheet but not in staff mapping
	}

	for i, r := range rows {
		rowNum := i + 2
		set := func(col int, v any) {
			cell, _ := excelize.CoordinatesToCellName(col, rowNum)
			f.SetCellValue(sheet, cell, v)
		}
		set(1, r.name)
		set(2, r.count)
		set(3, r.value)
		set(4, r.target)
	}

	if err := f.SaveAs(path); err != nil {
		log.Fatalf("saving %s: %v", path, err)
	}
	log.Printf("wrote %s", path)
}
