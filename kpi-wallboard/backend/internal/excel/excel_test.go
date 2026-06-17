package excel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"

	"wallboard/internal/config"
)

// writeXLSX creates a temporary spreadsheet with the given header and rows.
func writeXLSX(t *testing.T, sheet string, rows [][]any) string {
	t.Helper()
	f := excelize.NewFile()
	defer f.Close()
	idx, err := f.NewSheet(sheet)
	if err != nil {
		t.Fatalf("new sheet: %v", err)
	}
	f.SetActiveSheet(idx)
	f.DeleteSheet("Sheet1")
	for r, row := range rows {
		for c, v := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+1)
			f.SetCellValue(sheet, cell, v)
		}
	}
	path := filepath.Join(t.TempDir(), "sales.xlsx")
	if err := f.SaveAs(path); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path
}

func TestParseFile_ConfigurableColumns(t *testing.T) {
	// Deliberately non-default, oddly-cased column names with messy values.
	path := writeXLSX(t, "Data", [][]any{
		{"Rep Name", "Deals", "Revenue", "Quota"},
		{"Jane Doe", 7, "£18,250", "20,000"},
		{"  marcus lee ", "5", "$14,100.50", 20000},
		{"", 99, 99, 99}, // missing name -> skipped
	})

	cols := config.ExcelColumns{
		Sheet:      "Data",
		Name:       "rep name", // case-insensitive match
		SalesCount: "Deals",
		SalesValue: "Revenue",
		Target:     "Quota",
	}
	res, err := ParseFile(path, cols)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if res.Rows != 2 {
		t.Errorf("Rows = %d, want 2", res.Rows)
	}
	if res.SkippedRows != 1 {
		t.Errorf("SkippedRows = %d, want 1", res.SkippedRows)
	}

	jane, ok := res.BySpreadsheetName[NormName("Jane Doe")]
	if !ok {
		t.Fatalf("Jane Doe missing from result")
	}
	if jane.SalesCount != 7 {
		t.Errorf("Jane count = %d, want 7", jane.SalesCount)
	}
	if jane.SalesValue != 18250 {
		t.Errorf("Jane value = %v, want 18250 (currency/comma stripped)", jane.SalesValue)
	}
	if jane.Target != 20000 {
		t.Errorf("Jane target = %v, want 20000", jane.Target)
	}
	if got := jane.TargetProgress; got < 0.91 || got > 0.92 {
		t.Errorf("Jane progress = %v, want ~0.9125", got)
	}

	// Name normalisation: "  marcus lee " collapses to "marcus lee".
	marcus, ok := res.BySpreadsheetName[NormName("Marcus Lee")]
	if !ok {
		t.Fatalf("Marcus Lee not matched after normalisation")
	}
	if marcus.SalesValue != 14100.5 {
		t.Errorf("Marcus value = %v, want 14100.5", marcus.SalesValue)
	}
}

func TestParseFile_MissingNameColumn(t *testing.T) {
	path := writeXLSX(t, "Sheet1", [][]any{
		{"Foo", "Bar"},
		{"x", "y"},
	})
	_, err := ParseFile(path, config.ExcelColumns{Name: "Name"})
	if err == nil {
		t.Fatal("expected error when name column is absent")
	}
}

func TestParseFile_OptionalColumnsAbsent(t *testing.T) {
	// Only a name and a count column; no value or target.
	path := writeXLSX(t, "Sheet1", [][]any{
		{"Name", "Sales Count"},
		{"Solo Rep", 4},
	})
	cols := config.ExcelColumns{
		Name:       "Name",
		SalesCount: "Sales Count",
		SalesValue: "Sales Value", // not present
		Target:     "Target",      // not present
	}
	res, err := ParseFile(path, cols)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	m := res.BySpreadsheetName[NormName("Solo Rep")]
	if m.SalesCount != 4 {
		t.Errorf("count = %d, want 4", m.SalesCount)
	}
	if m.SalesValue != 0 || m.Target != 0 || m.TargetProgress != 0 {
		t.Errorf("expected zero value/target/progress when columns absent, got %+v", m)
	}
}

func TestParseFile_FileNotFound(t *testing.T) {
	_, err := ParseFile(filepath.Join(os.TempDir(), "does-not-exist-12345.xlsx"), config.ExcelColumns{Name: "Name"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
