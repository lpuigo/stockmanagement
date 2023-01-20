package article

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/date"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"
	"github.com/lpuig/batec/stockmanagement/src/backend/tools/xlstools"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
)

func checkValue(xf *excelize.File, sheetname, axis, value string) error {
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

const (
	rowArticleHeader int = 1
	rowArticleInfo   int = 2
)

var colArticlesTitle = []string{
	"Id",
	"Category",
	"SubCategory",
	"Manufacturer",
	"Designation",
	"Ref",
	"PhotoId",
	"Location",
	"UnitStock",
	"UnitAccounting",
	"ConvStockAccounting",
	"Status",
}

func LoadArticlesFromXlsx(r io.Reader) ([]*Article, error) {
	xf, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	sheetName := xf.GetSheetName(0)

	rows, err := xf.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	// parse header
	colTitleDict := make(map[string]int)
	for i, _ := range colArticlesTitle {
		title, _ := xf.GetCellValue(sheetName, xlstools.RcToAxis(rowArticleHeader, i+1))
		colTitleDict[title] = i + 1
	}

	line := 0
	row := []string{}
	getCol := func(title string) string {
		col := colTitleDict[title]
		if col == 0 {
			return ""
		}
		if col-1 >= len(row) {
			return ""
		}
		return row[col-1]
	}

	newID := -1

	getNextNewId := func() int {
		id := newID
		newID--
		return id
	}

	// parse each row
	articles := []*Article{}
	for line, row = range rows {
		if line+1 < rowArticleInfo {
			continue
		}

		// Get ID : if not blank (not provided) assign a new Id (negative value)
		idStr := getCol("Id")
		id, err := strconv.Atoi(idStr)
		if err != nil || idStr == "" {
			id = getNextNewId()
		}

		convStockAccounting, err := strconv.ParseFloat(getCol("ConvStockAccounting"), 64)
		if err != nil {
			convStockAccounting = 0
		}
		art := &Article{
			Id:                  id,
			Category:            getCol("Category"),
			SubCategory:         getCol("SubCategory"),
			Designation:         getCol("Designation"),
			Ref:                 getCol("Ref"),
			Manufacturer:        getCol("Manufacturer"),
			PhotoId:             getCol("PhotoId"),
			Location:            getCol("Location"),
			UnitStock:           getCol("UnitStock"),
			UnitAccounting:      getCol("UnitAccounting"),
			ConvStockAccounting: convStockAccounting,
			Status:              getCol("Status"),
			TimeStamp:           timestamp.TimeStamp{},
		}

		articles = append(articles, art)
	}
	return articles, nil
}

func WriteArticlesToXlsx(w io.Writer, articles []*Article) error {
	xf := excelize.NewFile()
	sheetName := date.Now().TimeStampShort()
	xf.SetSheetName(xf.GetSheetName(0), sheetName)

	// write header
	for i, s := range colArticlesTitle {
		xf.SetCellStr(sheetName, xlstools.RcToAxis(rowArticleHeader, i+1), s)
	}

	for num, article := range articles {
		rowNum := num + 2
		for i, s := range colArticlesTitle {
			colNum := i + 1
			switch s {
			case "Id":
				xf.SetCellInt(sheetName, xlstools.RcToAxis(rowNum, colNum), article.Id)
			case "Category":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.Category)
			case "SubCategory":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.SubCategory)
			case "Manufacturer":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.Manufacturer)
			case "Designation":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.Designation)
			case "Ref":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.Ref)
			case "PhotoId":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.PhotoId)
			case "Location":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.Location)
			case "UnitStock":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.UnitStock)
			case "UnitAccounting":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.UnitAccounting)
			case "ConvStockAccounting":
				xf.SetCellFloat(sheetName, xlstools.RcToAxis(rowNum, colNum), article.ConvStockAccounting, 3, 64)
			case "Status":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.Status)
			}
		}
	}
	err := xf.Write(w)
	if err != nil {
		return fmt.Errorf("could not write XLS file:%s", err.Error())
	}
	return nil
}
