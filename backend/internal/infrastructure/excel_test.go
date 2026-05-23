package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExcelProvider_GenerateSimpleExcel_WithHeadersAndRows(t *testing.T) {
	p := NewExcelProvider()

	headers := []string{"Name", "Age", "City"}
	rows := [][]interface{}{
		{"Alice", 30, "NYC"},
		{"Bob", 25, "LA"},
	}

	data, err := p.GenerateSimpleExcel("Users", headers, rows)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestExcelProvider_GenerateSimpleExcel_EmptySheetName(t *testing.T) {
	p := NewExcelProvider()

	headers := []string{"Col1"}
	rows := [][]interface{}{{"val1"}}

	data, err := p.GenerateSimpleExcel("", headers, rows)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestExcelProvider_NewFile(t *testing.T) {
	p := NewExcelProvider()

	f := p.NewFile()
	assert.NotNil(t, f)
	f.Close()
}

func TestExcelProvider_WriteToBuffer(t *testing.T) {
	p := NewExcelProvider()

	f := p.NewFile()
	f.SetCellValue("Sheet1", "A1", "Hello")

	data, err := p.WriteToBuffer(f)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}
