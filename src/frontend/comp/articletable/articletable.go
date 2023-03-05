package articletable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle/articleconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuigo/hvue"
	"sort"
	"strconv"
	"strings"
)

const (
	template string = `
<el-table ref="articleTable"
        :border=true
        :data="filteredArticles"
        :default-sort = "{prop: 'Category', order: 'ascending'}"        
        :row-class-name="TableRowClassName" height="100%" size="mini"
		@row-dblclick="HandleDoubleClickedRow"
>
<!--		:default-sort = "{prop: 'Stay.EndDate', order: 'descending'}"-->
	
	<!--	Index   -->
	<el-table-column
		label="N°" width="40px"
		type="index"
		index=1 
	></el-table-column>

	<!--	Category   -->
	<el-table-column label="Catégorie" prop="Category" width="140px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="['Category', 'SubCategory', 'Designation']"
		:filters="FilterList('Category')" :filter-method="FilterHandler" filter-placement="bottom-end" :filtered-value="FilteredCategoryValue()"
	></el-table-column>

	<!--	Ss Category   -->
	<el-table-column label="Sous-Cat." prop="SubCategory" width="140px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="['SubCategory', 'Category', 'Designation']" 
		:filters="FilterList('SubCategory')" :filter-method="FilterHandler" filter-placement="bottom-end"
	></el-table-column>

	<!--	Designation   -->
	<el-table-column label="Désignation" prop="Designation"
		:resizable="true" :show-overflow-tooltip=true
		sortable 
	></el-table-column>

	<!--	Manufacturer   -->
	<el-table-column label="Constructeur" prop="Manufacturer" width="250px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="['Manufacturer', 'Category', 'SubCategory', 'Designation']" 
		:filters="FilterList('Manufacturer')" :filter-method="FilterHandler" filter-placement="bottom-end"
	></el-table-column>

	<!--	Ref   -->
	<el-table-column
		:resizable="true" :show-overflow-tooltip=true 
		prop="Ref" label="Référence" width="150px"
		sortable :sort-by="['Ref', 'Category', 'SubCategory', 'Designation']" 
	></el-table-column>
	
	<!--	PhotoId  -->
<!--	<el-table-column-->
<!--			:resizable="true" :show-overflow-tooltip=true -->
<!--			label="Photo" width="180px"-->
<!--	>-->
<!--		<template slot-scope="scope">-->
<!--			<span>{{scope.row.PhotoId}}</span>-->
<!--		</template>-->
<!--	</el-table-column>-->
	
	<!--	RetailUnit   -->
	<el-table-column
		:resizable="true" :show-overflow-tooltip=true 
		prop="RetailUnit" label="Unité" width="200px"
		sortable :sort-by="['RetailUnit', 'Category', 'SubCategory', 'Designation']"  
		:filters="FilterList('RetailUnit')" :filter-method="FilterHandler" filter-placement="bottom-end"
	>
		<template slot-scope="scope">
			<span>{{FormatRetailUnit(scope.row)}}</span>
		</template>
	</el-table-column>
	
	<!--	InvoiceUnitPrice   -->
	<el-table-column
		:resizable="true" :show-overflow-tooltip=true 
		prop="InvoiceUnitPrice" label="Unité" width="200px"
	>
		<template slot-scope="scope">
			<span>{{FormatRetailPrice(scope.row)}}</span>
		</template>
	</el-table-column>
	
</el-table>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("articles-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewArticlesTableModel(vm)
		}),
		hvue.MethodsOf(&ArticlesTableModel{}),
		hvue.Computed("filteredArticles", func(vm *hvue.VM) interface{} {
			atm := ArticlesTableModelFromJS(vm.Object)
			return atm.GetFilteredArticles()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ArticlesTableModel struct {
	*js.Object

	Articles   *fearticle.ArticleStore `js:"value"`
	User       *feuser.User            `js:"user"`
	Filter     string                  `js:"filter"`
	FilterType string                  `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewArticlesTableModel(vm *hvue.VM) *ArticlesTableModel {
	atm := &ArticlesTableModel{Object: tools.O()}
	atm.VM = vm
	atm.Articles = fearticle.NewArticleStore()
	atm.User = feuser.NewUser()
	atm.Filter = ""
	atm.FilterType = ""

	return atm
}

func ArticlesTableModelFromJS(o *js.Object) *ArticlesTableModel {
	return &ArticlesTableModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Functions

func (atm *ArticlesTableModel) TableRowClassName(vm *hvue.VM, rowInfo *js.Object) string {
	//atm = ArticlesTableModelFromJS(vm.Object)
	//as := fearticle.ArticleFromJS(rowInfo.Get("row"))
	//return as.GetAvailabilityRowClass()
	return ""
}

func (atm *ArticlesTableModel) HandleDoubleClickedRow(vm *hvue.VM, ar *fearticle.Article) {
	atm = ArticlesTableModelFromJS(vm.Object)
	message.NotifyWarning(vm, "Double Click Article", "front/comp/articletable/HandleDoubleClickedRow à implémenter")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

// FilteredCategoryValue returns pre filtered values for Category
func (vtm *ArticlesTableModel) FilteredCategoryValue() []string {
	return []string{}
}

func (vtm *ArticlesTableModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (vtm *ArticlesTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	vtm = ArticlesTableModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	var getValue func(ar *fearticle.Article) string
	switch prop {
	//case "Status":
	//	getValue = func(ar *fearticle.Article) string {
	//		return ar.Status()
	//	}
	default:
		getValue = func(ar *fearticle.Article) string {
			return ar.Get(prop).String()
		}
	}

	attrib := ""
	for _, ar := range vtm.Articles.Articles {
		attrib = getValue(ar)
		if _, exist := count[attrib]; !exist {
			attribs = append(attribs, attrib)
		}
		count[attrib]++
	}
	sort.Strings(attribs)
	res := []*elements.ValText{}
	for _, a := range attribs {
		fa := a
		if fa == "" {
			fa = "Vide"
		}
		res = append(res, elements.NewValText(a, fa+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Actions Functions

func (atm *ArticlesTableModel) GetFilteredArticles() []*fearticle.Article {
	filter := func(ar *fearticle.Article) bool {
		return true
	}
	if !(atm.FilterType == articleconst.FilterValueAll && atm.Filter == "") {
		expected := strings.ToUpper(atm.Filter)
		filter = func(ar *fearticle.Article) bool {
			ss := ar.SearchString(atm.FilterType)
			if ss == "" {
				return false
			}
			return strings.Contains(strings.ToUpper(ss), expected)
		}
	}

	// filter articles in accs slice to prevent change on atm.Articles caused by sort
	var accs []*fearticle.Article
	for _, a := range atm.Articles.Articles {
		if filter(a) {
			accs = append(accs, a)
		}
	}
	return accs
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Format Functions

func (atm *ArticlesTableModel) FormatRetailUnit(vm *hvue.VM, art *fearticle.Article) string {
	suffix := ""
	if art.StockUnit != art.RetailUnit && art.RetailUnitStockQty != 1 {
		suffix = " (" + strconv.FormatFloat(art.RetailUnitStockQty, 'f', 1, 64) + " par " + art.StockUnit + ")"
	}
	return art.RetailUnit + suffix
}

func (atm *ArticlesTableModel) FormatRetailPrice(vm *hvue.VM, art *fearticle.Article) string {
	suffix := ""
	unitPrice := art.InvoiceUnitPrice
	if art.InvoiceUnit != art.RetailUnit || art.InvoiceUnitRetailQty != 1 {
		unitPrice *= art.InvoiceUnitRetailQty
		suffix = " (" + strconv.FormatFloat(art.InvoiceUnitPrice, 'f', 2, 64) + " €/" + art.InvoiceUnit + ")"
	}
	return strconv.FormatFloat(unitPrice, 'f', 2, 64) + " €/" + art.RetailUnit + suffix
}
