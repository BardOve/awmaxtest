package internal

import (
	"fmt"

	"codeberg.org/go-pdf/fpdf"
)

func GeneratePDF(filepath string, name string, inputPNG string) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "AWMax Report - Measurement Data")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("This is a report generated from AWMAX measurement data for %s",name))
	pdf.Ln(10)

	pdf.Image(inputPNG, 10, 20, 140, 0, false, "", 0, "")


	err := pdf.OutputFileAndClose(filepath)
	if err != nil {
		return err
	}
	return nil
}

