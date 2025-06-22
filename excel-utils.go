package main

import (
	"bytes"
	"fmt"

	"github.com/tealeg/xlsx"
)

func createExcelFile(record AggregatedRecord) ([]byte, error) {

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")

	if err != nil {
		return nil, fmt.Errorf("failed to add sheet: %w", err)
	}

	// New Header row
	headers := []string{
		"ItemID", "ItemName", "Quantity", "Color", "Size", "FetchedAPIData", "APISource", "ProcessingDate",
	}

	headerRow := sheet.AddRow()

	for _, header := range headers {
		cell := headerRow.AddCell()
		cell.Value = header
	}

	// Data Rows
	dataRow := sheet.AddRow()
	dataRow.AddCell().Value = record.ItemID
	dataRow.AddCell().Value = record.ItemName
	dataRow.AddCell().SetInt(record.Quantity)
	dataRow.AddCell().Value = record.Color
	dataRow.AddCell().Value = record.Size
	dataRow.AddCell().Value = record.FetchedAPIData
	dataRow.AddCell().Value = record.APISource
	dataRow.AddCell().Value = record.ProcessingDate

	var buffer bytes.Buffer

	if err := file.Write(&buffer); err != nil {
		return nil, fmt.Errorf("failed to write excel file to buffer %w", err)
	}

	return buffer.Bytes(), nil

}
