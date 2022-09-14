package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/FujikuraYohei/masterdata-convertor/ssaccessor"
	"github.com/FujikuraYohei/masterdata-convertor/ssinterpretor"
)

var (
	ignorablePrefix = "_"
	csvpath         = "dist/csv/"
	sqlpath         = "dist/sql/"
)

func handler(ctx context.Context) (err error) {
	spreadsheetIDs := flag.Args()

	for _, spreadsheetID := range spreadsheetIDs {
		// get all sheet titles of spreadsheet
		titles, err := ssaccessor.GetSheetTitles(spreadsheetID)
		if err != nil {
			return err
		}

		// filter out untargeted sheet
		titles = filterSheet(titles)

		// process by sheet
		for _, title := range titles {

			// output csv file
			c, err := ssaccessor.GetCSVFromSpreadsheet(spreadsheetID, title)
			if err != nil {
				return err
			}

			f, err := os.Create(csvpath + title + ".csv")
			if err != nil {
				return err
			}
			defer f.Close()

			f.WriteString(c)

			// construct and output sql file
			griddata, err := ssaccessor.GetSliceFromSpreadsheet(spreadsheetID, title)
			if err != nil {
				return err
			}

			// カラム
			columnList, err := ssinterpretor.GetInsertColumns(griddata)
			columns := strings.Join(columnList, ", ")

			// 値リスト
			dataValues, err := ssinterpretor.GetDataValues(griddata)
			valueList := []string{}
			for _, dataValue := range dataValues {
				var strVal []string

				for i := 0; i < len(dataValue); i++ {
					strVal = append(strVal, fmt.Sprintf("%v", dataValue[i]))
				}
				strValList := fmt.Sprintf("(%s)", strings.Join(strVal, ","))
				valueList = append(valueList, strValList)
			}
			values := strings.Join(valueList, ", ")

			tableName := title + "_master"
			sql := sqlTemplate(tableName, columns, values)

			sf, err := os.Create(sqlpath + title + ".sql")
			if err != nil {
				return err
			}
			defer sf.Close()

			sf.WriteString(sql)
		}
	}

	return nil
}

// filter out sheet which name start with "_"
func filterSheet(titles []string) []string {
	var validTitles []string

	for _, title := range titles {
		if strings.HasPrefix(title, ignorablePrefix) == false {
			validTitles = append(validTitles, title)
		}
	}

	return validTitles
}

func sqlTemplate(tableName, columns, values string) string {
	sql := ""
	sql = sql + fmt.Sprintf("begin;")
	sql = sql + fmt.Sprintf("delete from %s;\n", tableName)
	sql = sql + fmt.Sprintf("insert into %s (%s) values %s;\n", tableName, columns, values)
	sql = sql + fmt.Sprintf("commit;")

	return sql
}

func main() {
	flag.Parse()
	ctx := context.TODO()
	err := handler(ctx)
	if err != nil {
		panic(err)
	}
}
