package xlstools

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/date"
	"github.com/xuri/excelize/v2"
	"strconv"
	"time"
)

func RcToAxis(row, col int) string {
	//res, err := excelize.CoordinatesToCellName(col, row)
	//if err != nil {
	//	res = "A1"
	//}
	//return res

	colname, err := excelize.ColumnNumberToName(col)
	if err != nil {
		colname = "A"
	}
	return colname + strconv.Itoa(row)
}

func CheckValue(xf *excelize.File, sheetname, axis, value string) error {
	foundValue, err := xf.GetCellValue(sheetname, axis)
	if err != nil {
		return err
	}
	if foundValue != value {
		return fmt.Errorf("misformated XLS file: cell %s!%s should contain '%s' (found '%s' instead)",
			sheetname, axis,
			foundValue, value,
		)
	}
	return nil
}

func GetCellInt(xf *excelize.File, sheetname, axis string) (int, error) {
	foundValue, err := xf.GetCellValue(sheetname, axis)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(foundValue)
	if err != nil {
		return 0, fmt.Errorf("misformated XLS file: cell %s!%s should contain int value", sheetname, axis)
	}
	return val, nil
}

func GetCellFloat(xf *excelize.File, sheetname, axis string) (float64, error) {
	foundValue, err := xf.GetCellValue(sheetname, axis)
	if err != nil {
		return 0, err
	}
	val, err := strconv.ParseFloat(foundValue, 64)
	if err != nil {
		return 0, fmt.Errorf("misformated XLS file: cell %s!%s should contain float value", sheetname, axis)
	}
	return val, nil
}

func GetCellDate(xf *excelize.File, sheetname, axis string) (string, error) {
	foundValue, _ := xf.GetCellValue(sheetname, axis)
	foundDate, err := time.Parse("01-02-06", foundValue)
	if err != nil {
		foundDate, err = time.Parse("1/2/06 15:04", foundValue)
		if err != nil {
			return "", fmt.Errorf("misformated XLS file: cell %s!%s should contain date value ('%s' found instead): %s", sheetname, axis, foundValue, err.Error())
		}
	}
	return date.Date(foundDate).String(), nil
}

func ParseDate(source string) (string, error) {
	pdate := ""
	if source != "" {
		tdate, err := time.Parse("2006-01-02", source)
		if err != nil {
			tdate, err = time.Parse("01-02-06", source)
			if err != nil {
				tdate, err = time.Parse("1/2/06 15:04", source)
				if err != nil {
					return "", fmt.Errorf("could not parse date '%s': %s", source, err.Error())
				}
			}
		}
		pdate = date.Date(tdate).String()
	}
	return pdate, nil
}

func GetColName(col int) string {
	colName, _ := excelize.ColumnNumberToName(col)
	return colName
}
