package stockarticletable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle/articleconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festock"
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
<el-container style="height: 100%">
	<el-header style="height: auto; padding: 0 15px; margin-bottom: 15px">
		<el-row :gutter="10" type="flex" align="middle">
			<el-col :span="3">
				<el-button type="primary" size="mini" @click="ToggleSelection" :disabled="SelectedArticles.length == 0">Ajout/Retrait</el-button>
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
				@row-dblclick="HandleDoubleClickedRow"
				@selection-change="HandleSelectionChange"
		>
		<!--		:default-sort = "{prop: 'Stay.EndDate', order: 'descending'}"-->
			
			<!--	Index   -->
			<el-table-column
			  type="selection"
			  width="55">
			</el-table-column>
				
			<!--	Index   -->
			<el-table-column
				label="N°" width="40px"
				type="index"
				index=1 
			></el-table-column>
		
			<!--	Availability   -->
			<el-table-column label="Disponibilité" prop="Status" width="120px"
				:resizable="true" :show-overflow-tooltip=true
				sortable :sort-by="['Status', 'Category', 'SubCategory', 'Designation']"
				:filters="FilterList('Status')" :filter-method="FilterHandler" filter-placement="bottom-end" :filtered-value="FilteredStatusValue()"
			>
				<template slot-scope="scope">
					<span>{{FormatStatus(scope.row)}}</span>
				</template>
			</el-table-column>
		
			<!--	Category   -->
			<el-table-column label="Catégorie" prop="Category" width="150px"
				:resizable="true" :show-overflow-tooltip=true
				sortable :sort-by="['Category', 'SubCategory', 'Designation']"
				:filters="FilterList('Category')" :filter-method="FilterHandler" filter-placement="bottom-end"
			></el-table-column>
		
			<!--	Ss Category   -->
			<el-table-column label="Sous-Cat." prop="SubCategory" width="150px"
				:resizable="true" :show-overflow-tooltip=true
				sortable :sort-by="['SubCategory', 'Category', 'Designation']" 
				:filters="FilterList('SubCategory')" :filter-method="FilterHandler" filter-placement="bottom-end"
			></el-table-column>
		
			<!--	Designation   -->
			<el-table-column label="Désignation" prop="Designation" width="300px"
				:resizable="true" :show-overflow-tooltip=true
				sortable 
			></el-table-column>
		
			<!--	Manufacturer   -->
			<el-table-column label="Constructeur" prop="Manufacturer" width="150px"
				:resizable="true" :show-overflow-tooltip=true
				sortable :sort-by="['Manufacturer', 'Category', 'SubCategory', 'Designation']" 
				:filters="FilterList('Manufacturer')" :filter-method="FilterHandler" filter-placement="bottom-end"
			></el-table-column>
		
			<!--	Ref   -->
			<el-table-column
				:resizable="true" :show-overflow-tooltip=true 
				prop="Ref" label="Référence" width="100px"
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
			
			<!--	UnitStock   -->
			<el-table-column
				:resizable="true" :show-overflow-tooltip=true 
				prop="UnitStock" label="Unité" width="100px" align="center"
				sortable :sort-by="['UnitStock', 'Category', 'SubCategory', 'Designation']"  
				:filters="FilterList('UnitStock')" :filter-method="FilterHandler" filter-placement="bottom-end"
			></el-table-column>
			
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
			atm := StockArticlesTableModelFromJS(vm.Object)
			return atm.GetFilteredArticles()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type StockArticlesTableModel struct {
	*js.Object

	Stock            *festock.Stock          `js:"value"`
	Articles         *fearticle.ArticleStore `js:"articles"`
	SelectedArticles []*fearticle.Article    `js:"SelectedArticles"`
	User             *feuser.User            `js:"user"`
	Filter           string                  `js:"filter"`
	FilterType       string                  `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewStockArticlesTableModel(vm *hvue.VM) *StockArticlesTableModel {
	atm := &StockArticlesTableModel{Object: tools.O()}
	atm.VM = vm
	atm.Articles = fearticle.NewArticleStore()
	atm.SelectedArticles = []*fearticle.Article{}
	atm.User = feuser.NewUser()
	atm.Filter = ""
	atm.FilterType = ""

	return atm
}

func StockArticlesTableModelFromJS(o *js.Object) *StockArticlesTableModel {
	return &StockArticlesTableModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Functions

// Filter related methods

func (atm *StockArticlesTableModel) ApplyFilter(vm *hvue.VM) {
	//atm = StockArticlesTableModelFromJS(vm.Object)
}

func (atm *StockArticlesTableModel) GetFilterType(vm *hvue.VM) []*elements.ValueLabel {
	return fearticle.GetFilterTypeValueLabel()
}

func (atm *StockArticlesTableModel) ClearFilter(vm *hvue.VM) {
	atm = StockArticlesTableModelFromJS(vm.Object)
	atm.FilterType = articleconst.FilterValueAll
	atm.Filter = ""
}

// Table related methods

func (atm *StockArticlesTableModel) TableRowClassName(vm *hvue.VM, rowInfo *js.Object) string {
	//atm = StockArticlesTableModelFromJS(vm.Object)
	//as := fearticle.ArticleFromJS(rowInfo.Get("row"))
	//return as.GetAvailabilityRowClass()
	return ""
}

func (atm *StockArticlesTableModel) HandleDoubleClickedRow(vm *hvue.VM, ar *fearticle.Article) {
	atm = StockArticlesTableModelFromJS(vm.Object)
	message.NotifyWarning(vm, "Double Click Article", "front/comp/articletable/HandleDoubleClickedRow à implémenter")
}

func (atm *StockArticlesTableModel) HandleSelectionChange(vm *hvue.VM, selArts *js.Object) {
	atm = StockArticlesTableModelFromJS(vm.Object)
	selectedArticles := []*fearticle.Article{}
	selArts.Call("forEach", func(art *fearticle.Article) {
		selectedArticles = append(selectedArticles, art)
	})
	atm.SelectedArticles = selectedArticles
}

func (atm *StockArticlesTableModel) ToggleSelection(vm *hvue.VM) {
	atm = StockArticlesTableModelFromJS(vm.Object)
	for _, article := range atm.SelectedArticles {
		article.ToggleInStock()
	}
}

// Table column format related methods

func (atm *StockArticlesTableModel) FormatStatus(ar *fearticle.Article) string {
	return fearticle.GetStatusLabel(ar.Status)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

// FilteredStatusValue returns pre filtered values for Status
func (vtm *StockArticlesTableModel) FilteredStatusValue() []string {
	return []string{articleconst.StatusValueAvailable}
}

func (vtm *StockArticlesTableModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	switch prop {
	case "Status":
		return fearticle.GetStatusLabel(p.Get(prop).String()) == value
	default:
		return p.Get(prop).String() == value
	}
}

func (vtm *StockArticlesTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	vtm = StockArticlesTableModelFromJS(vm.Object)
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

func (atm *StockArticlesTableModel) GetFilteredArticles() []*fearticle.Article {
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

	// Set Article status : available in stock or not
	isArticleInStockById := atm.Stock.GenHasArticleById()
	for _, art := range atm.Articles.Articles {
		if isArticleInStockById(art.Id) {
			art.Status = articleconst.StatusValueAvailable
			continue
		}
		art.Status = articleconst.StatusValueUnavailable
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
