package articlescatalog

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/articletable"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle/articleconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuigo/hvue"
)

const (
	template string = `
<el-container style="height: 100%">

	<el-header style="height: auto; padding: 0 15px; margin-bottom: 15px">
		<el-row gutter=10 type="flex" align="middle">
			<el-col span=6>
				<h3>Catalogue d'Articles</h3>
			</el-col>
			<el-col :span="10">
                <el-input v-model="Filter" size="mini" style="width: 25vw; min-width: 130px"
                          @input="ApplyFilter">
                    <el-select v-model="FilterType"
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
			<el-col :span="3">
				<el-button-group>
					<el-tooltip content="Export vers un fichier XLSx" placement="bottom" effect="light" open-delay="500">
						<el-button type="warning" plain icon="fa-solid fa-file-import icon--big" @click="ExportToXLSx()" size="mini"></el-button>
					</el-tooltip>
					<el-tooltip content="Mise à jour depuis un fichier XLSx" placement="bottom" effect="light" open-delay="500">
						<el-button type="warning" plain icon="fa-solid fa-file-export icon--big" @click="ImportFromXLSx()" size="mini"></el-button>
					</el-tooltip>
				</el-button-group>
			</el-col>
		</el-row>
	</el-header>

	<el-main style="padding: 0 15px">
		<articles-table
				v-model="value"
				:user="user"
				:filter="Filter" :filtertype="FilterType"
		></articles-table>
	</el-main>
</el-container>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("articles-catalog", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		articletable.RegisterComponent(),
		hvue.Props("value", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewArticlesCatalogModel(vm)
		}),
		hvue.MethodsOf(&ArticlesCatalogModel{}),
		//hvue.Computed("filteredArticles", func(vm *hvue.VM) interface{} {
		//	atm := ArticlesCatalogModelFromJS(vm.Object)
		//	return atm.GetFilteredArticles()
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type ArticlesCatalogModel struct {
	*js.Object

	Articles   *fearticle.ArticleStore `js:"value"`
	User       *feuser.User            `js:"user"`
	Filter     string                  `js:"Filter"`
	FilterType string                  `js:"FilterType"`

	VM *hvue.VM `js:"VM"`
}

func NewArticlesCatalogModel(vm *hvue.VM) *ArticlesCatalogModel {
	mum := &ArticlesCatalogModel{Object: tools.O()}
	mum.Articles = fearticle.NewArticleStore()
	mum.User = feuser.NewUser()
	mum.Filter = ""
	mum.FilterType = ""

	mum.VM = vm

	return mum
}

func ArticlesCatalogModelFromJS(o *js.Object) *ArticlesCatalogModel {
	return &ArticlesCatalogModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Functions

// Filter related methods

func (mum *ArticlesCatalogModel) ApplyFilter(vm *hvue.VM) {
	//mum = ArticlesCatalogModelFromJS(vm.Object)
}

func (mum *ArticlesCatalogModel) GetFilterType(vm *hvue.VM) []*elements.ValueLabel {
	return fearticle.GetFilterTypeValueLabel()
}

func (mum *ArticlesCatalogModel) ClearFilter(vm *hvue.VM) {
	mum = ArticlesCatalogModelFromJS(vm.Object)
	mum.FilterType = articleconst.FilterValueAll
	mum.Filter = ""
}

// Action Methods

func (mum *ArticlesCatalogModel) ExportToXLSx(vm *hvue.VM) {
	tools.OpenUri("/api/articles/export")
}

func (mum *ArticlesCatalogModel) ImportFromXLSx(vm *hvue.VM) {
	message.NotifyError(vm, "ImportFromXLSx", "method to be implemented")
}
