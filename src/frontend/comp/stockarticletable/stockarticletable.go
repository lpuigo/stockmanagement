package stockarticletable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle/articleconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festock"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements"
	"github.com/lpuigo/hvue"
	"sort"
	"strconv"
	"strings"
)

const (
	template string = `
<el-container style="height: 100%">
	<el-header style="height: auto; padding: 0 15px; margin-bottom: 15px">
		<el-row :gutter="10" type="flex" align="middle">
			<el-col :span="5">
				<h2><i class="fa-solid fa-chart-line icon--left"></i>Etat du Stock</h2>
			</el-col>
			<el-col :span="5">
			</el-col>
			<el-col :span="10">
                <el-input v-model="filter" size="mini" style="width: 25vw; min-width: 130px"
                          @input="ApplyFilter">
                    <el-select v-model="filtertype"
                               @change="ApplyFilter"
                               slot="prepend" placeholder="Tous"
                               style="width: 10vw; min-width: 60px; max-width: 120px; margin-right: -10px">
                        <el-option
                                v-for="item in GetFilterType()"
                                :key="item.value"
                                :label="item.label"
                                :value="item.value"
                        ></el-option>
                    </el-select>
                    <el-button slot="append" icon="far fa-times-circle" @click="ClearFilter"></el-button>
                </el-input>
			</el-col>
		</el-row>
	</el-header>
	<el-main style="padding: 0 15px">
		<el-table ref="stockArticleTable"
				:border=true
				:data="filteredArticles"
				:default-sort = "{prop: 'Category', order: 'ascending'}"        
				:row-class-name="TableRowClassName" height="100%" size="mini"
		>
		<!--		:default-sort = "{prop: 'Stay.EndDate', order: 'descending'}"-->
			
			<!--	Index   -->
			<el-table-column
				label="N??" width="40px"
				type="index"
				index=1 
			></el-table-column>
		
			<!--	Availability   -->
			<el-table-column label="Disponibilit??" prop="Status" width="120px"
				:resizable="true" :show-overflow-tooltip=true
				sortable :sort-by="['Status', 'Category', 'SubCategory', 'Designation']"
				:filters="FilterList('Status')" :filter-method="FilterHandler" filter-placement="bottom-end"
			>
				<template slot-scope="scope">
					<span>{{FormatStatus(scope.row)}}</span>
				</template>
			</el-table-column>
		
			<!--	Category   -->
			<el-table-column label="Cat??gorie" prop="Category" width="200px"
				:resizable="true" :show-overflow-tooltip=true
				sortable :sort-by="['Category', 'SubCategory', 'Designation']"
				:filters="FilterList('Category')" :filter-method="FilterHandler" filter-placement="bottom-end"
			></el-table-column>
		
			<!--	Ss Category   -->
			<el-table-column label="Sous-Cat." prop="SubCategory" width="200px"
				:resizable="true" :show-overflow-tooltip=true
				sortable :sort-by="['SubCategory', 'Category', 'Designation']" 
				:filters="FilterList('SubCategory')" :filter-method="FilterHandler" filter-placement="bottom-end"
			></el-table-column>
		
			<!--	Designation   -->
			<el-table-column label="D??signation" prop="Designation"
				:resizable="true" :show-overflow-tooltip=true
				sortable 
			></el-table-column>
		
			<!--	Manufacturer   -->
			<el-table-column label="Constructeur" prop="Manufacturer" width="300px"
				:resizable="true" :show-overflow-tooltip=true
				sortable :sort-by="['Manufacturer', 'Category', 'SubCategory', 'Designation']" 
				:filters="FilterList('Manufacturer')" :filter-method="FilterHandler" filter-placement="bottom-end"
			></el-table-column>
		
			<!--	Ref   -->
			<el-table-column
				:resizable="true" :show-overflow-tooltip=true 
				prop="Ref" label="R??f??rence" width="150px"
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
				prop="RetailUnit" label="D??tail" width="150px" align="center"
				sortable :sort-by="['RetailUnit', 'Category', 'SubCategory', 'Designation']"  
				:filters="FilterList('RetailUnit')" :filter-method="FilterHandler" filter-placement="bottom-end"
			></el-table-column>
			
			<!--	StockUnit   -->
			<el-table-column
				:resizable="true" :show-overflow-tooltip=true 
				prop="StockUnit" label="Gros" width="150px" align="center"
				sortable :sort-by="['StockUnit', 'Category', 'SubCategory', 'Designation']"  
				:filters="FilterList('StockUnit')" :filter-method="FilterHandler" filter-placement="bottom-end"
			>
				<template slot-scope="scope">
					<p class="article-unit">{{scope.row.StockUnit}}</p>
					<p v-if="scope.row.StockUnit != scope.row.RetailUnit" class="article-unit light">{{GetStockRetailQty(scope.row)}}</p>
				</template>
			</el-table-column>
			
		</el-table>
	</el-main>
</el-container>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("stock-articles-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "articles", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewStockArticlesTableModel(vm)
		}),
		hvue.MethodsOf(&StockArticlesTableModel{}),
		hvue.Computed("filteredArticles", func(vm *hvue.VM) interface{} {
			satm := StockArticlesTableModelFromJS(vm.Object)
			return satm.GetFilteredArticles()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type StockArticlesTableModel struct {
	*js.Object

	Stock            *festock.Stock          `js:"value"`
	Articles         *fearticle.ArticleStore `js:"articles"`
	FilteredArticles []*fearticle.Article    `js:"FilteredArticles"`
	User             *feuser.User            `js:"user"`
	Filter           string                  `js:"filter"`
	FilterType       string                  `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewStockArticlesTableModel(vm *hvue.VM) *StockArticlesTableModel {
	satm := &StockArticlesTableModel{Object: tools.O()}
	satm.VM = vm
	satm.Articles = fearticle.NewArticleStore()
	satm.FilteredArticles = []*fearticle.Article{}
	satm.User = feuser.NewUser()
	satm.Filter = ""
	satm.FilterType = ""

	return satm
}

func StockArticlesTableModelFromJS(o *js.Object) *StockArticlesTableModel {
	return &StockArticlesTableModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Functions

// Filter related methods

func (satm *StockArticlesTableModel) ApplyFilter(vm *hvue.VM) {
	//satm = StockArticlesTableModelFromJS(vm.Object)
}

func (satm *StockArticlesTableModel) GetFilterType(vm *hvue.VM) []*elements.ValueLabel {
	return fearticle.GetFilterTypeValueLabel()
}

func (satm *StockArticlesTableModel) ClearFilter(vm *hvue.VM) {
	satm = StockArticlesTableModelFromJS(vm.Object)
	satm.FilterType = articleconst.FilterValueAll
	satm.Filter = ""
}

// Table related methods

func (satm *StockArticlesTableModel) TableRowClassName(vm *hvue.VM, rowInfo *js.Object) string {
	satm = StockArticlesTableModelFromJS(vm.Object)
	ar := fearticle.ArticleFromJS(rowInfo.Get("row"))
	return fearticle.GetStatusClass(ar.Status)
}

// Table column format related methods

func (satm *StockArticlesTableModel) FormatStatus(ar *fearticle.Article) string {
	return fearticle.GetStatusLabel(ar.Status)
}

func (satm *StockArticlesTableModel) GetStockRetailQty(vm *hvue.VM, ar *fearticle.Article) string {
	return "( " + strconv.FormatFloat(ar.RetailUnitStockQty, 'f', 1, 64) + " x " + ar.RetailUnit + " )"
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

// FilteredStatusValue returns pre filtered values for Status
func (satm *StockArticlesTableModel) FilteredStatusValue() []string {
	return []string{articleconst.StatusLabelAvailable, articleconst.StatusLabelOutOfStock}
}

func (satm *StockArticlesTableModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	switch prop {
	case "Status":
		return fearticle.GetStatusLabel(p.Get(prop).String()) == value
	default:
		return p.Get(prop).String() == value
	}
}

func (satm *StockArticlesTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	satm = StockArticlesTableModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	var getValue func(ar *fearticle.Article) string
	switch prop {
	case "Status":
		getValue = func(ar *fearticle.Article) string {
			return fearticle.GetStatusLabel(ar.Status)
		}
	default:
		getValue = func(ar *fearticle.Article) string {
			return ar.Get(prop).String()
		}
	}

	attrib := ""
	for _, ar := range satm.FilteredArticles {
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
// Table Functions

func (satm *StockArticlesTableModel) GetFilteredArticles() []*fearticle.Article {
	filter := func(ar *fearticle.Article) bool {
		return true
	}
	if !(satm.FilterType == articleconst.FilterValueAll && satm.Filter == "") {
		expected := strings.ToUpper(satm.Filter)
		filter = func(ar *fearticle.Article) bool {
			ss := ar.SearchString(satm.FilterType)
			if ss == "" {
				return false
			}
			return strings.Contains(strings.ToUpper(ss), expected)
		}
	}

	// filter articles in accs slice to prevent change on satm.Articles caused by sort
	var accs []*fearticle.Article
	for _, a := range satm.Articles.Articles {
		if a.Status != articleconst.StatusValueUnavailable && filter(a) {
			accs = append(accs, a)
		}
	}
	satm.FilteredArticles = accs
	return accs
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Table Functions
