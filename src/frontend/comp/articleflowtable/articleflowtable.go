package articleflowtable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/articlepicktable"
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
        height="100%" size="mini"
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
				<el-popover
					v-model="PickArticleVisible"
					placement="right"
					title="Ajout d'un article"
					width="80vw"
					@show="ResetPickedArticle()"
					trigger="click">
				<el-container style="height: 60vh">
					<el-header>
						selectors
					</el-header>
					<el-main>
						<articles-pick-table
							v-model="pickableArticles"
							:user="user"
							:filter="Filter"
							:filtertype="FilterType"
							@picked-article="HandlePickedArticle"
						></articles-pick-table>
					</el-main>
					<div style="margin: 0 20px; text-align: right">
						<el-button size="mini" @click="PickArticleVisible = false">Fermer</el-button>
						<el-button type="success" @click="AddPickedArticle" plain size="mini":disabled="!isPickedArticle">Ajouter</el-button>
					</div>
				</el-container>
				<el-button slot="reference" type="success" plain icon="fa-solid fa-dolly fa-fw" size="mini"></el-button>
				</el-popover>
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
		articlepicktable.RegisterComponent(),
		hvue.Props("value", "articles", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewArticleFlowTableModel(vm)
		}),
		hvue.MethodsOf(&ArticleFlowTableModel{}),
		hvue.Computed("pickableArticles", func(vm *hvue.VM) interface{} {
			atm := ArticleFlowTableModelFromJS(vm.Object)
			return atm.getPickableArticleStore()
		}),
		hvue.Computed("isPickedArticle", func(vm *hvue.VM) interface{} {
			atm := ArticleFlowTableModelFromJS(vm.Object)
			return atm.PickedArticle.Object != nil && atm.PickedArticle.Id >= 0
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ArticleFlowTableModel struct {
	*js.Object

	ArticleFlows  []*femovement.ArticleFlow `js:"value"`
	StockArticles *fearticle.ArticleStore   `js:"articles"`
	User          *feuser.User              `js:"user"`

	PickArticleVisible bool   `js:"PickArticleVisible"`
	Filter             string `js:"Filter"`
	FilterType         string `js:"FilterType"`

	PickedArticle *fearticle.Article `js:"PickedArticle"`

	VM *hvue.VM `js:"VM"`
}

func NewArticleFlowTableModel(vm *hvue.VM) *ArticleFlowTableModel {
	aftm := &ArticleFlowTableModel{Object: tools.O()}
	aftm.VM = vm
	aftm.ArticleFlows = []*femovement.ArticleFlow{}
	aftm.StockArticles = fearticle.NewArticleStore()
	aftm.User = feuser.NewUser()

	aftm.PickArticleVisible = false
	aftm.Filter = ""
	aftm.FilterType = ""

	aftm.PickedArticle = fearticle.NewArticle()

	return aftm
}

func ArticleFlowTableModelFromJS(o *js.Object) *ArticleFlowTableModel {
	return &ArticleFlowTableModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Functions

//func (aftm *ArticleFlowTableModel) TableRowClassName(vm *hvue.VM, rowInfo *js.Object) string {
//	//aftm = ArticleFlowTableModelFromJS(vm.Object)
//	//as := femovement.MovementFromJS(rowInfo.Get("row"))
//	return ""
//}

func (aftm *ArticleFlowTableModel) HandlePickedArticle(vm *hvue.VM, pickedArt *fearticle.Article) {
	aftm = ArticleFlowTableModelFromJS(vm.Object)
	aftm.PickedArticle = pickedArt
}

func (aftm *ArticleFlowTableModel) ResetPickedArticle(vm *hvue.VM) {
	aftm = ArticleFlowTableModelFromJS(vm.Object)
	aftm.PickedArticle = fearticle.NewArticle()
}

func (aftm *ArticleFlowTableModel) AddPickedArticle(vm *hvue.VM) {
	aftm = ArticleFlowTableModelFromJS(vm.Object)
	af := femovement.NewArticleFlow()
	af.ArtId = aftm.PickedArticle.Id
	af.Qty = 1
	//aftm.ArticleFlows = append(aftm.ArticleFlows, af)
	aftm.Get("value").Call("unshift", af)
	aftm.PickArticleVisible = false
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
	art, found := aftm.StockArticles.ArticleIndex[id]
	if !found {
		return "article " + strconv.Itoa(id) + " inconnu"
	}
	return art.Category
}

func (aftm *ArticleFlowTableModel) GetArticleSubCat(vm *hvue.VM, id int) string {
	aftm = ArticleFlowTableModelFromJS(vm.Object)
	art, found := aftm.StockArticles.ArticleIndex[id]
	if !found {
		return "article " + strconv.Itoa(id) + " inconnu"
	}
	return art.SubCategory
}

func (aftm *ArticleFlowTableModel) GetArticleDesignation(vm *hvue.VM, id int) string {
	aftm = ArticleFlowTableModelFromJS(vm.Object)
	art, found := aftm.StockArticles.ArticleIndex[id]
	if !found {
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

func (aftm *ArticleFlowTableModel) getPickableArticleStore() *fearticle.ArticleStore {
	pArts := []*fearticle.Article{}
	artFlows := make(map[int]bool)
	for _, flow := range aftm.ArticleFlows {
		artFlows[flow.ArtId] = true
	}
	for _, article := range aftm.StockArticles.Articles {
		if artFlows[article.Id] {
			continue
		}
		pArts = append(pArts, article)
	}
	pas := fearticle.NewArticleStore()
	pas.SetArticles(pArts)
	return pas
}
