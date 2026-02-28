package infrastructure

import (
	"bytes"

	"github.com/xuri/excelize/v2"
)

type ExcelProvider interface {
	GenerateSimpleExcel(sheetName string, headers []string, rows [][]interface{}) ([]byte, error)
	NewFile() *excelize.File
	WriteToBuffer(file *excelize.File) ([]byte, error)
}

type excelProvider struct{}

func NewExcelProvider() ExcelProvider {
	return &excelProvider{}
}

func (p *excelProvider) GenerateSimpleExcel(sheetName string, headers []string, rows [][]interface{}) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	if sheetName == "" {
		sheetName = "Sheet1"
	}

	// Rename the default sheet
	f.SetSheetName("Sheet1", sheetName)

	// Add headers
	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return nil, err
		}
		f.SetCellValue(sheetName, cell, header)
	}

	// Make header bold
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err == nil {
		f.SetRowStyle(sheetName, 1, 1, style)
	}

	// Add rows
	for r, row := range rows {
		for c, val := range row {
			cell, err := excelize.CoordinatesToCellName(c+1, r+2)
			if err != nil {
				return nil, err
			}
			f.SetCellValue(sheetName, cell, val)
		}
	}

	var buffer bytes.Buffer
	if err := f.Write(&buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (p *excelProvider) NewFile() *excelize.File {
	return excelize.NewFile()
}

func (p *excelProvider) WriteToBuffer(file *excelize.File) ([]byte, error) {
	defer file.Close()
	var buffer bytes.Buffer
	if err := file.Write(&buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
