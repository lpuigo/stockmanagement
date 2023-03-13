package articlepicktable

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
<el-table ref="articlePickTable"
        :border=true
        :data="filteredArticles"
        :default-sort = "{prop: 'Category', order: 'ascending'}" 
		highlight-current-row
        :row-class-name="TableRowClassName" height="100%" size="mini"
		@current-change="HandleSelectionChange"
>
<!--		:default-sort = "{prop: 'Stay.EndDate', order: 'descending'}"-->
	
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
	<el-table-column label="Constructeur" prop="Manufacturer" width="150px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="['Manufacturer', 'Category', 'SubCategory', 'Designation']" 
		:filters="FilterList('Manufacturer')" :filter-method="FilterHandler" filter-placement="bottom-end"
	></el-table-column>

	<!--	RetailUnit   -->
	<el-table-column
		:resizable="true" :show-overflow-tooltip=true 
		prop="RetailUnit" label="Détail" width="150px" align="center"
	></el-table-column>
	
	<!--	StockUnit   -->
	<el-table-column
		:resizable="true" :show-overflow-tooltip=true 
		prop="StockUnit" label="Gros" width="150px" align="center"
	>
		<template slot-scope="scope">
			<p class="article-unit">{{scope.row.StockUnit}}</p>
			<p v-if="scope.row.StockUnit != scope.row.RetailUnit" class="article-unit light">{{GetStockRetailQty(scope.row)}}</p>
		</template>
	</el-table-column>
	
	<!--	PhotoId  -->
<!--	<el-table-column-->
<!--			:resizable="true" :show-overflow-tooltip=true -->
<!--			label="Photo" width="180px"-->
<!--	>-->
<!--		<template slot-scope="scope">-->
<!--			<span>{{scope.row.PhotoId}}</span>-->
<!--		</template>-->
<!--	</el-table-column>-->
	
</el-table>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("articles-pick-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewArticlePickTableModel(vm)
		}),
		hvue.MethodsOf(&ArticlePickTableModel{}),
		hvue.Computed("filteredArticles", func(vm *hvue.VM) interface{} {
			atm := ArticlePickTableModelFromJS(vm.Object)
			return atm.GetFilteredArticles()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ArticlePickTableModel struct {
	*js.Object

	Articles   *fearticle.ArticleStore `js:"value"`
	User       *feuser.User            `js:"user"`
	Filter     string                  `js:"filter"`
	FilterType string                  `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewArticlePickTableModel(vm *hvue.VM) *ArticlePickTableModel {
	atm := &ArticlePickTableModel{Object: tools.O()}
	atm.VM = vm
	atm.Articles = fearticle.NewArticleStore()
	atm.User = feuser.NewUser()
	atm.Filter = ""
	atm.FilterType = ""

	return atm
}

func ArticlePickTableModelFromJS(o *js.Object) *ArticlePickTableModel {
	return &ArticlePickTableModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Functions

func (atm *ArticlePickTableModel) TableRowClassName(vm *hvue.VM, rowInfo *js.Object) string {
	//atm = ArticlePickTableModelFromJS(vm.Object)
	//as := fearticle.ArticleFromJS(rowInfo.Get("row"))
	//return as.GetAvailabilityRowClass()
	return ""
}

func (atm *ArticlePickTableModel) HandleDoubleClickedRow(vm *hvue.VM, ar *fearticle.Article) {
	atm = ArticlePickTableModelFromJS(vm.Object)
	message.NotifyWarning(vm, "Double Click Article", "front/comp/articletable/HandleDoubleClickedRow à implémenter")
}

func (atm *ArticlePickTableModel) HandleSelectionChange(vm *hvue.VM, selArt *fearticle.Article) {
	atm = ArticlePickTableModelFromJS(vm.Object)
	atm.VM.Emit("picked-article", selArt)
}

func (atm *ArticlePickTableModel) GetStockRetailQty(vm *hvue.VM, art *fearticle.Article) string {
	return "( " + strconv.FormatFloat(art.RetailUnitStockQty, 'f', 1, 64) + " x " + art.RetailUnit + " )"
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

// FilteredCategoryValue returns pre filtered values for Category
func (vtm *ArticlePickTableModel) FilteredCategoryValue() []string {
	return []string{}
}

func (vtm *ArticlePickTableModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (vtm *ArticlePickTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	vtm = ArticlePickTableModelFromJS(vm.Object)
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

func (atm *ArticlePickTableModel) GetFilteredArticles() []*fearticle.Article {
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
