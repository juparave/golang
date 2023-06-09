## Excelize

Library to read and write `xlsx` files

### Styling

Example file

```go
package report

import "github.com/xuri/excelize/v2"

var xlsBorders = []excelize.Border{
	{
		Type:  "left",
		Color: "#000000",
		Style: 1,
	},
	{
		Type:  "right",
		Color: "#000000",
		Style: 1,
	},
	{
		Type:  "top",
		Color: "#000000",
		Style: 1,
	},
	{
		Type:  "bottom",
		Color: "#000000",
		Style: 1,
	},
}

type xlsStyle struct {
	f *excelize.File
}

// setBorder bordered style cell
func (xs *xlsStyle) setBorder() int {
	var bor = excelize.Style{
		Border: xlsBorders,
	}
	s, _ := xs.f.NewStyle(&bor)
	return s
}

// setBorder and decimal places
func (xs *xlsStyle) setBorderDec(dec int) int {
	var bor = excelize.Style{
		NumFmt:        2,
		DecimalPlaces: dec,
		Border:        xlsBorders,
	}
	s, _ := xs.f.NewStyle(&bor)
	return s
}

// setPct percentage style cell
func (xs *xlsStyle) setPct() int {
	var pct = excelize.Style{
		NumFmt:        10, // percentage
		DecimalPlaces: 2,
		Border:        xlsBorders,
	}
	pctStyle, _ := xs.f.NewStyle(&pct)
	return pctStyle
}

// setPctRed red percentage style cell
func (xs *xlsStyle) setPctRed() int {
	var redPct = excelize.Style{
		NumFmt:        10, // percentage
		DecimalPlaces: 2,
		Font: &excelize.Font{
			Bold:   false,
			Family: "Arial",
			Size:   11,
			Color:  "#DC143C",
		},
		Border: xlsBorders,
	}
	var redPctStyle, _ = xs.f.NewStyle(&redPct)
	return redPctStyle
}

// serPctNoBorder percentage style cell without border
func (xs *xlsStyle) setPctNoBorder() int {
	var pct = excelize.Style{
		NumFmt:        10, // percentage
		DecimalPlaces: 2,
	}
	pctStyle, _ := xs.f.NewStyle(&pct)
	return pctStyle
}
```

Example usage

```go
func fillSpecieTable(f *excelize.File, sheet string, samples []models.QASample,
	startCol string, startRow int, specieCalibre string) {
	app := config.GetAppConfig()
	row := startRow + 1
	styles := xlsStyle{f: f}

...
    // add weights
    f.SetCellValue(sheet, col+fmt.Sprint(row), *sample.PesoNeto)
    f.SetCellStyle(sheet, col+fmt.Sprint(row), col+fmt.Sprint(row), styles.setBorder())
    col = string(col[0] + 1)
    // peso base
    f.SetCellValue(sheet, col+fmt.Sprint(row), *sample.PesoBase)
    f.SetCellStyle(sheet, col+fmt.Sprint(row), col+fmt.Sprint(row), styles.setBorder())
    col = string(col[0] + 1)
    // difference
    diff := *sample.PesoBase - *sample.PesoNeto
    f.SetCellValue(sheet, col+fmt.Sprint(row), diff)
    f.SetCellStyle(sheet, col+fmt.Sprint(row), col+fmt.Sprint(row), styles.setBorderDec(2))
    col = string(col[0] + 1)
...
```
