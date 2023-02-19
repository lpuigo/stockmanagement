package articleflowtable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuigo/hvue"
	"strconv"
)

const (
	template string = `
<el-table ref="articleFlowTable"
        :border=true
        :data="value"
        :default-sort = "{prop: 'ArtId', order: 'ascending'}"        
        :row-class-name="TableRowClassName" height="100%" size="mini"
		@row-dblclick=""
>
	<!--	Index   -->
	<el-table-column
		label="N°" width="40px"
		type="index"
		index=1 
	></el-table-column>

	<!--	Actions   -->
	<el-table-column label="" width="70px">
		<template slot="header" slot-scope="scope">
			<el-tooltip content="Ajouter un article" placement="bottom" effect="light" open-delay=400>
				<el-button type="success" plain icon="fa-solid fa-dolly fa-fw" size="mini" @click="AddNewArticleFlow()"></el-button>
			</el-tooltip>
		</template>
		<template slot-scope="scope">
			<el-tooltip content="Retirer cet article" placement="bottom" effect="light" open-delay=400>
				<el-button type="danger" plain icon="fa-solid fa-ban fa-fw" size="mini" @click="RemoveArticleFlow(scope.row)"></el-button>
			</el-tooltip>
		</template>
	</el-table-column>

	<!--	Article Cat   -->
	<el-table-column label="Catégorie" width="240px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="SortByCat"
	>
		<template slot-scope="scope">
			<span>{{GetArticleCat(scope.row.ArtId)}}</span>
		</template>
	</el-table-column>

	<!--	Article SubCat   -->
	<el-table-column label="Sous-Catégorie" width="240px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="SortBySubCat"
	>
		<template slot-scope="scope">
			<span>{{GetArticleSubCat(scope.row.ArtId)}}</span>
		</template>
	</el-table-column>

	<!--	Article Designation   -->
	<el-table-column label="Désignation" 
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="SortByDesignation"
	>
		<template slot-scope="scope">
			<span>{{GetArticleDesignation(scope.row.ArtId)}}</span>
		</template>
	</el-table-column>

	<!--	Quantities   -->
	<el-table-column label="Quantité" prop="Qty" width="210px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="['Qty']"
	>
		<template slot-scope="scope">
			<el-input-number v-model="scope.row.Qty" :min="1"></el-input-number>
		</template>
	</el-table-column>
</el-table>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("articleflow-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "articles", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewArticleFlowTableModel(vm)
		}),
		hvue.MethodsOf(&ArticleFlowTableModel{}),
		//hvue.Computed("filteredMovements", func(vm *hvue.VM) interface{} {
		//	atm := ArticleFlowTableModelFromJS(vm.Object)
		//	return atm.GetFilteredMovements()
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ArticleFlowTableModel struct {
	*js.Object

	ArticleFlows      []*femovement.ArticleFlow `js:"value"`
	StockArticleStore *fearticle.ArticleStore   `js:"articles"`
	User              *feuser.User              `js:"user"`

	VM *hvue.VM `js:"VM"`
}

func NewArticleFlowTableModel(vm *hvue.VM) *ArticleFlowTableModel {
	aftm := &ArticleFlowTableModel{Object: tools.O()}
	aftm.VM = vm
	aftm.ArticleFlows = []*femovement.ArticleFlow{}
	aftm.StockArticleStore = fearticle.NewArticleStore()
	aftm.User = feuser.NewUser()

	return aftm
}

func ArticleFlowTableModelFromJS(o *js.Object) *ArticleFlowTableModel {
	return &ArticleFlowTableModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Functions

func (aftm *ArticleFlowTableModel) TableRowClassName(vm *hvue.VM, rowInfo *js.Object) string {
	//aftm = ArticleFlowTableModelFromJS(vm.Object)
	//as := femovement.MovementFromJS(rowInfo.Get("row"))
	return ""
}

// Sort Methods --------------------------------------------------------------------------------------------------------

func (aftm *ArticleFlowTableModel) SortByCat(vm *hvue.VM, row *js.Object, index int) string {
	af := fearticle.ArticleFromJS(row)
	return aftm.GetArticleCat(vm, af.Id)
}

func (aftm *ArticleFlowTableModel) SortBySubCat(vm *hvue.VM, row *js.Object, index int) string {
	af := fearticle.ArticleFromJS(row)
	return aftm.GetArticleSubCat(vm, af.Id)
}

func (aftm *ArticleFlowTableModel) SortByDesignation(vm *hvue.VM, row *js.Object, index int) string {
	af := fearticle.ArticleFromJS(row)
	return aftm.GetArticleDesignation(vm, af.Id)
}

// Get Article aatribute Methods ---------------------------------------------------------------------------------------

func (aftm *ArticleFlowTableModel) GetArticleCat(vm *hvue.VM, id int) string {
	aftm = ArticleFlowTableModelFromJS(vm.Object)
	art := aftm.StockArticleStore.GetById(id)
	if art == nil {
		return "article " + strconv.Itoa(id) + " inconnu"
	}
	return art.Category
}

func (aftm *ArticleFlowTableModel) GetArticleSubCat(vm *hvue.VM, id int) string {
	aftm = ArticleFlowTableModelFromJS(vm.Object)
	art := aftm.StockArticleStore.GetById(id)
	if art == nil {
		return "article " + strconv.Itoa(id) + " inconnu"
	}
	return art.SubCategory
}

func (aftm *ArticleFlowTableModel) GetArticleDesignation(vm *hvue.VM, id int) string {
	aftm = ArticleFlowTableModelFromJS(vm.Object)
	art := aftm.StockArticleStore.GetById(id)
	if art == nil {
		return "article " + strconv.Itoa(id) + " inconnu"
	}
	return art.Designation
}

// ---------------------------------------------------------------------------------------------------------------------

func (aftm *ArticleFlowTableModel) AddNewArticleFlow(vm *hvue.VM) {
	message.NotifyError(vm, "ArticleFlowTableModel.AddNewArticleFlow", "method to be implemented")
}

func (aftm *ArticleFlowTableModel) RemoveArticleFlow(vm *hvue.VM, af *femovement.ArticleFlow) {
	aftm = ArticleFlowTableModelFromJS(vm.Object)

	for i, articleFlow := range aftm.ArticleFlows {
		if articleFlow.ArtId == af.ArtId {
			aftm.Get("value").Call("splice", i, 1)
			break
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Internal Functions
