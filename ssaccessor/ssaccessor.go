package ssaccessor

import (
	"context"
	"errors"
	"os"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// spreadsheet accesskey location
var filename = os.Getenv("HOME") + "/.config/cred.json"

// get all spreadsheet titles by spreadsheet-id
func GetSheetTitles(spreadsheetID string) ([]string, error) {
	service, err := sheets.NewService(context.TODO(), option.WithCredentialsFile(filename))
	if err != nil {
		panic(err)
	}
	ssservice := sheets.NewSpreadsheetsService(service)

	sss, err := ssservice.Get(spreadsheetID).Do()
	if err != nil {
		return []string{}, errors.New("Given spreadsheet-id not found")
	}

	titles := make([]string, len(sss.Sheets))
	for key, sheet := range sss.Sheets {
		var title = Escape(sheet.Properties.Title)
		titles[key] = title
	}

	return titles, nil
}

// get all input by spreadsheet-id and sheet-title and then pack them into multi-dimention slice object
func GetSliceFromSpreadsheet(spreadsheetID, sheet string) (csv [][]string, err error) {
	service, err := sheets.NewService(context.TODO(), option.WithCredentialsFile(filename))
	if err != nil {
		panic(err)
	}
	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, sheet).Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Values) == 0 {
		return nil, errors.New("No data found")
	}

	rows := make([][]string, len(resp.Values))
	for i, row := range resp.Values {
		colomns := make([]string, len(row))
		for j, c := range row {
			colomns[j] = c.(string)
		}
		rows[i] = colomns
	}

	return rows, nil
}

// get all input by spreadsheet-id and sheet-title and then pack them into csv file
func GetCSVFromSpreadsheet(spreadsheetID, sheet string) (csv string, err error) {
	service, err := sheets.NewService(context.TODO(), option.WithCredentialsFile(filename))
	if err != nil {
		return "", nil
	}

	resp, err := service.Spreadsheets.Values.Get(spreadsheetID, sheet).Do()
	if err != nil {
		return "", nil
	}

	if len(resp.Values) == 0 {
		return "", errors.New("No data found")
	}

	rows := make([]string, len(resp.Values))
	for i, row := range resp.Values {
		colomns := make([]string, len(row))
		for j, c := range row {
			s := Escape(c.(string))
			colomns[j] = s
		}
		rows[i] = strings.Join(colomns, ",")
	}

	csv = strings.Join(rows, "\n")

	return csv, nil
}

func Escape(s string) string {
	x := strings.Replace(s, ",", "_+_", -1)
	x = strings.Replace(x, "\n", "\\n", -1)

	return x
}

func Unescape(s string) string {
	x := strings.Replace(s, "_+_", ",", -1)
	x = strings.Replace(x, "\\n", "\n", -1)

	return x
}
