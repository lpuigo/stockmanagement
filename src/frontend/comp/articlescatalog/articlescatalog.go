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
				<h2><i class="fa-solid fa-boxes-stacked icon--left"></i>Catalogue d'Articles</h2>
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
				<div>	
					<el-tooltip content="Export vers un fichier XLSx" placement="bottom" effect="light" open-delay="500">
						<el-button type="warning" plain icon="fa-solid fa-file-import icon--big" @click="ExportToXLSx()" size="mini"></el-button>
					</el-tooltip>
	
					<el-popover v-model="VisibleImportArticlesCat" 
							title="Import d'un catalogue d'articles:"
							placement="bottom" width="360" 
					>
						<el-upload 
								   :action="ImportFromXLSxURL()"
								   drag
								   style="width: 300px"
								   :before-upload="ImportFromXLSxBeforeUpload"
								   :on-success="ImportFromXLSxUploadSuccess"
								   :on-error="ImportFromXLSxUploadError"
						>
							<i class="el-icon-upload"></i>
							<div class="el-upload__text">Déposez un fichier XLSx ici ou <em>cliquez</em></div>
						</el-upload>
	
						<el-tooltip slot="reference" content="Mise à jour depuis un fichier XLSx" placement="bottom" effect="light" open-delay="500">
							<el-button type="warning" plain icon="fa-solid fa-file-export icon--big" @click="VisibleImportArticlesCat = !VisibleImportArticlesCat" size="mini"></el-button>
						</el-tooltip>
					</el-popover>
				</div>
			</el-col>
			<el-col span=3>
				<el-button-group>
                    <el-tooltip v-if="user.Permissions['Validate']" content="Enregistrer les modifications"
                                placement="bottom" effect="light" open-delay=500>
                        <el-button type="warning" class="icon" icon="fas fa-cloud-upload-alt icon--big" @click="SaveArticlesCatalog"
                                   :disabled="!IsDirty" size="mini"></el-button>
                    </el-tooltip>
                    <el-tooltip content="Raffraichir / Annuler les modifications" placement="bottom" effect="light"
                                open-delay="500">
                        <el-button type="warning" class="icon" icon="fas fa-undo-alt icon--big" @click="LoadArticlesCatalog"
                                   :disabled="!value.Loaded" size="mini"></el-button>
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
		hvue.Computed("IsDirty", func(vm *hvue.VM) interface{} {
			acm := ArticlesCatalogModelFromJS(vm.Object)
			return acm.Articles.IsDirty()
		}),
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

	VisibleImportArticlesCat bool `js:"VisibleImportArticlesCat"`

	VM *hvue.VM `js:"VM"`
}

func NewArticlesCatalogModel(vm *hvue.VM) *ArticlesCatalogModel {
	mum := &ArticlesCatalogModel{Object: tools.O()}
	mum.Articles = fearticle.NewArticleStore()
	mum.User = feuser.NewUser()
	mum.Filter = ""
	mum.FilterType = ""
	mum.VisibleImportArticlesCat = false

	mum.VM = vm

	return mum
}

func ArticlesCatalogModelFromJS(o *js.Object) *ArticlesCatalogModel {
	return &ArticlesCatalogModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Functions

// Filter related methods

func (acm *ArticlesCatalogModel) ApplyFilter(vm *hvue.VM) {
	//acm = ArticlesCatalogModelFromJS(vm.Object)
}

func (acm *ArticlesCatalogModel) GetFilterType(vm *hvue.VM) []*elements.ValueLabel {
	return fearticle.GetFilterTypeValueLabel()
}

func (acm *ArticlesCatalogModel) ClearFilter(vm *hvue.VM) {
	acm = ArticlesCatalogModelFromJS(vm.Object)
	acm.FilterType = articleconst.FilterValueAll
	acm.Filter = ""
}

// Action Methods

func (acm *ArticlesCatalogModel) ExportToXLSx(vm *hvue.VM) {
	tools.OpenUri(acm.Articles.GetExportArticlestoXlsxURL())
}

// Import Catalog from XLSx method
func (acm *ArticlesCatalogModel) ImportFromXLSxURL(vm *hvue.VM) string {
	return acm.Articles.GetImportArticlesFromXlsxURL()
}

func (acm *ArticlesCatalogModel) ImportFromXLSxBeforeUpload(vm *hvue.VM, file *js.Object) bool {
	if file.Get("type").String() != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		message.NotifyError(vm, "Import d'un catalogue d'articles", "Le fichier '"+file.Get("name").String()+"' n'est pas un document Xlsx")
		return false
	}
	return true
}

func (acm *ArticlesCatalogModel) ImportFromXLSxUploadError(vm *hvue.VM, err, file *js.Object) {
	acm = ArticlesCatalogModelFromJS(vm.Object)
	acm.VisibleImportArticlesCat = false
	message.NotifyError(vm, "Import d'un catalogue d'articles", err.String())
}

func (acm *ArticlesCatalogModel) ImportFromXLSxUploadSuccess(vm *hvue.VM, response, file *js.Object) {
	acm = ArticlesCatalogModelFromJS(vm.Object)
	importedArticles := fearticle.ArticleSliceFromJS(response)
	acm.Articles.UpdateWith(importedArticles)
	acm.VisibleImportArticlesCat = false
}

// Save & Reload methods
func (acm *ArticlesCatalogModel) SaveArticlesCatalog(vm *hvue.VM) {
	acm = ArticlesCatalogModelFromJS(vm.Object)
	onSavedArticle := func() {
		acm.Articles.CallGetArticles(vm, func() {})
	}
	acm.Articles.CallUpdateArticles(vm, onSavedArticle)
}

func (acm *ArticlesCatalogModel) LoadArticlesCatalog(vm *hvue.VM) {
	acm = ArticlesCatalogModelFromJS(vm.Object)
	acm.Articles.CallGetArticles(vm, func() {})
}
