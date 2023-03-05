package article

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/date"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"
	"github.com/lpuig/batec/stockmanagement/src/backend/tools/xlstools"
	"github.com/xuri/excelize/v2"
	"io"
	"strconv"
	"strings"
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
	"InvoiceUnit",
	"InvoiceUnitPrice",
	"InvoiceUnitRetailQty",
	"RetailUnit",
	"RetailUnitStockQty",
	"StockUnit",
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
	for i := range colArticlesTitle {
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
	getFloat := func(title string, def float64) float64 {
		val := getCol(title)
		if val == "" {
			return def
		}
		res, err := strconv.ParseFloat(strings.Replace(val, ",", ".", -1), 64)
		if err != nil {
			return 0
		}
		return res
	}
	getColWithDefault := func(title, def string) string {
		val := getCol(title)
		if val == "" {
			return def
		}
		return val
	}

	newID := 0
	getNextNewId := func() int {
		newID--
		return newID
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
		invoiceUnit := getCol("InvoiceUnit")
		retailUnit := getColWithDefault("RetailUnit", invoiceUnit)
		art := &Article{
			Id:                   id,
			Category:             getCol("Category"),
			SubCategory:          getCol("SubCategory"),
			Designation:          getCol("Designation"),
			Ref:                  getCol("Ref"),
			Manufacturer:         getCol("Manufacturer"),
			PhotoId:              getCol("PhotoId"),
			Location:             getCol("Location"),
			InvoiceUnit:          invoiceUnit,
			InvoiceUnitPrice:     getFloat("InvoiceUnitPrice", 0),
			InvoiceUnitRetailQty: getFloat("InvoiceUnitRetailQty", 1),
			RetailUnit:           retailUnit,
			RetailUnitStockQty:   getFloat("RetailUnitStockQty", 1),
			StockUnit:            getColWithDefault("StockUnit", retailUnit),

			Status:    getCol("Status"),
			TimeStamp: timestamp.TimeStamp{},
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

			case "InvoiceUnit":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.InvoiceUnit)
			case "InvoiceUnitPrice":
				xf.SetCellFloat(sheetName, xlstools.RcToAxis(rowNum, colNum), article.InvoiceUnitPrice, 3, 64)
			case "InvoiceUnitRetailQty":
				xf.SetCellFloat(sheetName, xlstools.RcToAxis(rowNum, colNum), article.InvoiceUnitRetailQty, 3, 64)
			case "RetailUnit":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.RetailUnit)
			case "RetailUnitStockQty":
				xf.SetCellFloat(sheetName, xlstools.RcToAxis(rowNum, colNum), article.RetailUnitStockQty, 3, 64)
			case "StockUnit":
				xf.SetCellStr(sheetName, xlstools.RcToAxis(rowNum, colNum), article.StockUnit)

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
