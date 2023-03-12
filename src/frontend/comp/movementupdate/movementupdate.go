package movementupdate

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/movementeditmodal"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/movementtable"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement/movementconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festock"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feworksite"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements"
	"github.com/lpuigo/hvue"
)

const (
	template string = `
<el-container style="height: 100%">

	<movement-edit-modal
			ref="MovementEditModal"
			:stock="stock"
			:articles="articles"
			:worksites="worksites"
			@new-movement="StoreNewMovement"
	></movement-edit-modal>

	<el-header style="height: auto; padding: 0 15px; margin-bottom: 15px">
		<el-row :gutter="10" type="flex" align="middle">
			<el-col span=4>
				<h2><i class="fa-solid fa-truck-ramp-box icon--left"></i>Mouvements de stock</h2>
			</el-col>
			<el-col span=5>
				<el-button-group>
					<el-tooltip content="Déclarer un Retrait" placement="bottom" effect="light" open-delay="500">
						<el-button type="warning" plain icon="fa-solid fa-person-walking-arrow-right icon--big" @click="AddMovement('withdrawal')"></el-button>
					</el-tooltip>
					<el-tooltip content="Déclarer un Retour" placement="bottom" effect="light" open-delay="500">
						<el-button type="warning" plain icon="fa-solid fa-person-walking-arrow-loop-left icon--big" @click="AddMovement('return')"></el-button>
					</el-tooltip>
					<el-tooltip content="Déclarer un Approvisionnement" placement="bottom" effect="light" open-delay="500">
						<el-button type="warning" plain icon="fa-solid fa-truck-arrow-right icon--big" @click="AddMovement('supply')"></el-button>
					</el-tooltip>
					<el-tooltip content="Déclarer un Inventaire" placement="bottom" effect="light" open-delay="500">
						<el-button type="warning" plain icon="fa-solid fa-list-check icon--big" @click="AddMovement('inventory')"></el-button>
					</el-tooltip>
				</el-button-group>
			</el-col>
			<el-col span=10>
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
		<movements-table
				v-model="value"
				:articles="articles"
				:worksites="worksites"
				:user="user"
				:filter="filter" :filtertype="filtertype"
				@edit-movement="EditMovement"
		></movements-table>
	</el-main>
</el-container>
`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("movements-update", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		movementeditmodal.RegisterComponent(),
		movementtable.RegisterComponent(),
		hvue.Props("value", "stock", "articles", "worksites", "user"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewMovementsUpdateModel(vm)
		}),
		hvue.MethodsOf(&MovementsUpdateModel{}),
		//hvue.Computed("filteredArticles", func(vm *hvue.VM) interface{} {
		//	atm := MovementsUpdateModelFromJS(vm.Object)
		//	return atm.GetFilteredArticles()
		//}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type MovementsUpdateModel struct {
	*js.Object

	StockMovements *femovement.MovementStore `js:"value"`
	Stock          *festock.Stock            `js:"stock"`
	Articles       *fearticle.ArticleStore   `js:"articles"`
	Worksites      *feworksite.WorksiteStore `js:"worksites"`
	User           *feuser.User              `js:"user"`
	Filter         string                    `js:"filter"`
	FilterType     string                    `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewMovementsUpdateModel(vm *hvue.VM) *MovementsUpdateModel {
	mum := &MovementsUpdateModel{Object: tools.O()}
	mum.StockMovements = femovement.NewMovementStore()
	mum.Stock = festock.NewStock()
	mum.Articles = fearticle.NewArticleStore()
	mum.Worksites = feworksite.NewWorksiteStore()
	mum.User = feuser.NewUser()
	mum.Filter = ""
	mum.FilterType = ""

	mum.VM = vm

	return mum
}

func MovementsUpdateModelFromJS(o *js.Object) *MovementsUpdateModel {
	return &MovementsUpdateModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Functions

// Filter related methods

func (mum *MovementsUpdateModel) ApplyFilter(vm *hvue.VM) {
	//mum = MovementsUpdateModelFromJS(vm.Object)
}

func (mum *MovementsUpdateModel) GetFilterType(vm *hvue.VM) []*elements.ValueLabel {
	return femovement.GetFilterTypeValueLabel()
}

func (mum *MovementsUpdateModel) ClearFilter(vm *hvue.VM) {
	mum = MovementsUpdateModelFromJS(vm.Object)
	mum.FilterType = movementconst.FilterValueAll
	mum.Filter = ""
}

func (mum *MovementsUpdateModel) EditMovement(vm *hvue.VM, mvt *femovement.Movement) {
	mum = MovementsUpdateModelFromJS(vm.Object)
	mum.ShowEditMovement(mvt, false)
}

func (mum *MovementsUpdateModel) AddMovement(vm *hvue.VM, mvtType string) {
	mum = MovementsUpdateModelFromJS(vm.Object)

	switch mvtType {
	case movementconst.TypeValueWithdrawal:
		nm := femovement.CreateNewMovement(mum.Stock.Id, mvtType)
		mum.ShowEditMovement(nm, true)
	case movementconst.TypeValueSupply:
		nm := femovement.CreateNewMovement(mum.Stock.Id, mvtType)
		mum.ShowEditMovement(nm, true)
	case movementconst.TypeValueInventory:
	}
}

func (mum *MovementsUpdateModel) StoreNewMovement(vm *hvue.VM, mvt *femovement.Movement) {
	mum = MovementsUpdateModelFromJS(vm.Object)
	mum.StockMovements.AddNewMovement(mvt)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Inner Functions

func (mum *MovementsUpdateModel) ShowEditMovement(mvt *femovement.Movement, isNewMovement bool) {
	memm := movementeditmodal.MovementEditModalModelFromJS(mum.VM.Refs("MovementEditModal"))
	memm.Show(mvt, mum.User, isNewMovement)
}
