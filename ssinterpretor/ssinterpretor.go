// Convert intputin spreadSheet into struct and json
package ssinterpretor

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	headerSize      = 4
	ignorablePrefix = "_"
	columnNameRow   = 0
	typeRow         = 1
	keyColRow       = 2
	jsonColRow      = 3
	dataStartRow    = 4
)

// interpret structured spreadsheet data, and then convert to json string
func Do(ssData [][]string) (string, error) {

	return "", nil
}

// check if spread sheet have enough input to process
func IsValid(ssInput [][]string) bool {
	if !(len(ssInput) > headerSize) {
		return false
	}

	return true
}

// get insert target columns
func GetInsertColumns(ssRawData [][]string) ([]string, error) {
	if !IsValid(ssRawData) {
		return []string{}, fmt.Errorf("input data is invalud: %v", ssRawData)
	}

	var colSize int = len(ssRawData[columnNameRow])
	var insertColumns []string
	hasJSON := false

	for colNum := 1; colNum < colSize; colNum++ {
		// skip the columns having their name with "_"
		if strings.HasPrefix(ssRawData[columnNameRow][colNum], ignorablePrefix) {
			continue
		}

		isKey, err := strconv.ParseBool(ssRawData[keyColRow][colNum])
		if err != nil {
			return []string{}, fmt.Errorf("%s cannot to be parsed to bool", ssRawData[keyColRow][colNum])
		}
		if isKey {
			insertColumns = append(insertColumns, ssRawData[columnNameRow][colNum])
		}

		isJSON, err := strconv.ParseBool(ssRawData[jsonColRow][colNum])
		if err != nil {
			return []string{}, fmt.Errorf("%s cannot to be parsed to bool", ssRawData[keyColRow][colNum])
		}
		if isJSON {
			hasJSON = true
		}
	}

	if hasJSON {
		insertColumns = append(insertColumns, "json")
	}

	return insertColumns, nil
}

func GetDataValues(ssRawData [][]string) ([][]interface{}, error) {
	if !IsValid(ssRawData) {
		return nil, fmt.Errorf("input data is invalud: %v", ssRawData)
	}

	var valueList [][]interface{}

	for _, row := range ssRawData[dataStartRow:] {
		ignoreRow, err := strconv.ParseBool(row[0])
		if err != nil {
			return nil, fmt.Errorf("%s cannot to be parsed to bool", row[0])
		}
		if ignoreRow {
			continue
		}

		var colSize int = len(row)
		//var values []string
		var values []interface{}
		jsonValues := map[string]interface{}{}

		for colNum := 1; colNum < colSize; colNum++ {
			// skip the columns having their name with "_"
			if strings.HasPrefix(ssRawData[columnNameRow][colNum], ignorablePrefix) {
				continue
			}

			isKey, err := strconv.ParseBool(ssRawData[keyColRow][colNum])
			if err != nil {
				return nil, fmt.Errorf("%s cannot to be parsed to bool", ssRawData[keyColRow][colNum])
			}
			if isKey {
				keyVal, err := castColumnType(ssRawData[typeRow][colNum], row[colNum])
				if err != nil {
					return nil, err
				}
				values = append(values, keyVal)

			}

			isJSON, err := strconv.ParseBool(ssRawData[jsonColRow][colNum])
			if err != nil {
				return nil, fmt.Errorf("%s cannot to be parsed to bool", ssRawData[jsonColRow][colNum])
			}
			if isJSON {
				jsonVal, err := castColumnType(ssRawData[typeRow][colNum], row[colNum])
				if err != nil {
					return nil, err
				}
				jsonValues[ssRawData[columnNameRow][colNum]] = jsonVal
			}
		}
		if len(jsonValues) != 0 {
			json, err := json.Marshal(jsonValues)
			if err != nil {
				return nil, fmt.Errorf("")
			}
			stringJSON := fmt.Sprintf("'%s'", json)
			values = append(values, stringJSON)
		}

		valueList = append(valueList, values)

	}

	return valueList, nil
}

func castColumnType(typename string, value string) (interface{}, error) {
	switch typename {
	case "string":
		return value, nil
	case "bool":
		castedVal, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("%s cannot to be parsed to bool", value)
		}
		return castedVal, nil
	case "int":
		castedVal, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("%s cannot to be parsed to int", value)
		}
		return castedVal, nil
	}

	return nil, fmt.Errorf("Invalid type is given: %s", typename)
}
