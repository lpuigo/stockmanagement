package movementtable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement/movementconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festatus"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feworksite"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/fedate"
	"github.com/lpuigo/hvue"
	"sort"
	"strconv"
	"strings"
)

const (
	template string = `
<el-table ref="movementsTable"
        :border=true
        :data="filteredMovements"
        :default-sort = "{prop: 'Date', order: 'descending'}"        
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

	<!--	Date   -->
	<el-table-column label="Date" prop="Date" width="120px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="['Date']"
	>
		<template slot-scope="scope">
			<span>{{FormatDate(scope.row.Date)}}</span>
		</template>
	</el-table-column>

	<!--	Type   -->
	<el-table-column label="Type" prop="Type" width="160px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="['Type', 'Date']"
		:filters="FilterList('Type')" :filter-method="FilterHandler" filter-placement="bottom-end"
	>
		<template slot-scope="scope">
			<div class="header-menu-container on-hover">
				<span>{{FormatMovementType(scope.row)}}</span>
				<i class="show link fa-solid fa-pen-to-square icon--left" @click="EditMovement(scope.row)"></i>
			</div>
		</template>
	</el-table-column>

	<!--	Actor   -->
	<el-table-column label="Acteur" prop="Actor" width="200px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="['Actor', 'Date']"
		:filters="FilterList('Actor')" :filter-method="FilterHandler" filter-placement="bottom-end"
	></el-table-column>

	<!--	Responsible   -->
	<el-table-column label="Responsable" prop="Responsible" width="200px"
		:resizable="true" :show-overflow-tooltip=true
		sortable :sort-by="['Responsible', 'Date']"
		:filters="FilterList('Responsible')" :filter-method="FilterHandler" filter-placement="bottom-end"
	></el-table-column>

	<!--	Worksite   -->
	<el-table-column label="Client" prop="WorksiteId" width="200px"
		:resizable="true" :show-overflow-tooltip=true
		sortable
		:filters="FilterList('WorksiteClient')" :filter-method="FilterHandler" filter-placement="bottom-end"
	>
		<template slot-scope="scope">
			<span>{{FormatWorksiteClient(scope.row.WorksiteId)}}</span>
		</template>
	</el-table-column>

	<!--	Worksite   -->
	<el-table-column label="Ville" prop="WorksiteId" width="200px"
		:resizable="true" :show-overflow-tooltip=true
		sortable
		:filters="FilterList('WorksiteCity')" :filter-method="FilterHandler" filter-placement="bottom-end"
	>
		<template slot-scope="scope">
			<span>{{FormatWorksiteCity(scope.row.WorksiteId)}}</span>
		</template>
	</el-table-column>

	<!--	Worksite   -->
	<el-table-column label="Chantier" prop="WorksiteId"
		:resizable="true" :show-overflow-tooltip=true
		sortable
	>
		<template slot-scope="scope">
			<span>{{FormatWorksiteRef(scope.row.WorksiteId)}}</span>
		</template>
	</el-table-column>

	<!--	Status   -->
	<el-table-column label="Status" prop="StatusHistory" width="150px"
		:resizable="true" :show-overflow-tooltip=true
		:filters="FilterList('StatusHistory')" :filter-method="FilterHandler" filter-placement="bottom-end"
	>
		<template slot-scope="scope">
			<span>{{FormatStatus(scope.row)}}</span>
		</template>
	</el-table-column>

	<!--	ArticleFlows   -->
	<el-table-column label="Articles" width="150px"
		:resizable="true" :show-overflow-tooltip=true
	>
		<template slot-scope="scope">
			<div class="header-menu-container on-hover">
				<span>{{scope.row.ArticleFlows.length}} article(s)</span>

				<el-popover placement="left" width="400" trigger="hover" :open-delay=250
					title="Liste des Articles"
				>
					<div v-for="(articleInfo, index) in GetArticlesInfoFor(scope.row)" :key="articleInfo" style="font-size: 0.85em">
						{{index+1}}) <span>{{articleInfo}}</span>
					</div>                        
					<i slot="reference" class="fas fa-info-circle icon--right show"></i>
				</el-popover>				
			</div>
		</template>
	</el-table-column>

	<!--	Price   -->
	<el-table-column label="Montant €" prop="StatusHistory" width="150px"
		:resizable="true" :show-overflow-tooltip=true align="center"
		sortable
	>
		<template slot-scope="scope">
			<div class="header-menu-container">
				<i v-if="ExternalMove(scope.row)" class="fa-solid fa-right-from-bracket"></i>
				<i v-else class="fa-solid fa-right-to-bracket"></i>
				<span>{{FormatPrice(scope.row)}}</span>
			</div>
		</template>
	</el-table-column>

</el-table>

`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("movements-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "articles", "worksites", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewMovementsTableModel(vm)
		}),
		hvue.MethodsOf(&MovementsTableModel{}),
		hvue.Computed("filteredMovements", func(vm *hvue.VM) interface{} {
			atm := MovementsTableModelFromJS(vm.Object)
			return atm.GetFilteredMovements()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Comp Model

type MovementsTableModel struct {
	*js.Object

	Movements  *femovement.MovementStore `js:"value"`
	User       *feuser.User              `js:"user"`
	Articles   *fearticle.ArticleStore   `js:"articles"`
	Worksites  *feworksite.WorksiteStore `js:"worksites"`
	Filter     string                    `js:"filter"`
	FilterType string                    `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewMovementsTableModel(vm *hvue.VM) *MovementsTableModel {
	atm := &MovementsTableModel{Object: tools.O()}
	atm.VM = vm
	atm.Movements = femovement.NewMovementStore()
	atm.User = feuser.NewUser()
	atm.Articles = fearticle.NewArticleStore()
	atm.Worksites = feworksite.NewWorksiteStore()
	atm.Filter = ""
	atm.FilterType = ""

	return atm
}

func MovementsTableModelFromJS(o *js.Object) *MovementsTableModel {
	return &MovementsTableModel{Object: o}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Functions

func (mtm *MovementsTableModel) TableRowClassName(vm *hvue.VM, rowInfo *js.Object) string {
	//mtm = MovementsTableModelFromJS(vm.Object)
	//as := femovement.MovementFromJS(rowInfo.Get("row"))
	//return as.GetAvailabilityRowClass()
	return ""
}

func (mtm *MovementsTableModel) HandleDoubleClickedRow(vm *hvue.VM, mvt *femovement.Movement) {
	mtm = MovementsTableModelFromJS(vm.Object)
	mtm.EditMovement(vm, mvt)
}

func (mtm *MovementsTableModel) EditMovement(vm *hvue.VM, mvt *femovement.Movement) {
	//mtm = MovementsTableModelFromJS(vm.Object)
	vm.Emit("edit-movement", mvt)
}

func (mtm *MovementsTableModel) FormatDate(d string) string {
	return fedate.DateString(d)
}

func (mtm *MovementsTableModel) FormatWorksiteClient(vm *hvue.VM, wsId int) string {
	mtm = MovementsTableModelFromJS(vm.Object)
	return mtm.Worksites.GetWorksiteById(wsId).GetClient()
}

func (mtm *MovementsTableModel) FormatWorksiteCity(vm *hvue.VM, wsId int) string {
	mtm = MovementsTableModelFromJS(vm.Object)
	return mtm.Worksites.GetWorksiteById(wsId).GetCity()
}

func (mtm *MovementsTableModel) FormatWorksiteRef(vm *hvue.VM, wsId int) string {
	mtm = MovementsTableModelFromJS(vm.Object)
	return mtm.Worksites.GetWorksiteById(wsId).GetRef()
}

func (mtm *MovementsTableModel) FormatMovementType(m *femovement.Movement) string {
	return m.GetTypeLabel()
}

func (mtm *MovementsTableModel) FormatStatus(m *femovement.Movement) string {
	return m.GetCurrentStatus().GetLabel()
}

func (mtm *MovementsTableModel) ExternalMove(m *femovement.Movement) bool {
	return m.IsExternalMove()
}

func (mtm *MovementsTableModel) FormatPrice(vm *hvue.VM, m *femovement.Movement) string {
	mtm = MovementsTableModelFromJS(vm.Object)
	var price float64
	for _, flow := range m.ArticleFlows {
		art, found := mtm.Articles.ArticleIndex[flow.ArtId]
		if !found {
			continue
		}
		price += art.GetRetailPrice(flow.Qty)
	}
	return strconv.FormatFloat(price, 'f', 2, 64) + " €"
}

func (mtm *MovementsTableModel) GetArticlesInfoFor(vm *hvue.VM, m *femovement.Movement) []string {
	mtm = MovementsTableModelFromJS(vm.Object)
	res := []string{}
	for _, articleFlow := range m.ArticleFlows {
		var artDesc string
		art, found := mtm.Articles.ArticleIndex[articleFlow.ArtId]
		if found {
			artDesc = art.RetailUnit + " de " + art.Designation
		} else {
			artDesc = "Article id " + strconv.Itoa(articleFlow.ArtId) + " inconnu"
		}
		res = append(res, strconv.Itoa(articleFlow.Qty)+" x "+artDesc)
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

// FilteredCategoryValue returns pre filtered values for Category
func (mtm *MovementsTableModel) FilteredCategoryValue() []string {
	return []string{}
}

func (mtm *MovementsTableModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	m := femovement.MovementFromJS(p)
	switch prop {
	case "StatusHistory":
		return m.GetCurrentStatus().Status == value
	default:
		return p.Get(prop).String() == value
	}
}

func (mtm *MovementsTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	mtm = MovementsTableModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	var getValue func(m *femovement.Movement) string
	var getLabel func(string) string
	switch prop {
	case "Type":
		getValue = func(m *femovement.Movement) string {
			return m.Type
		}
		getLabel = func(v string) string { return femovement.GetTypeLabel(v) }
	case "StatusHistory":
		getValue = func(m *femovement.Movement) string {
			return m.GetCurrentStatus().Status
		}
		getLabel = func(v string) string { return festatus.GetLabel(v) }
	case "WorksiteClient":
		getValue = func(m *femovement.Movement) string {
			return mtm.Worksites.GetWorksiteById(m.WorksiteId).GetClient()
		}
		getLabel = func(v string) string {
			return v
		}
	case "WorksiteCity":
		getValue = func(m *femovement.Movement) string {
			return mtm.Worksites.GetWorksiteById(m.WorksiteId).GetCity()
		}
		getLabel = func(v string) string {
			return v
		}
	default:
		getValue = func(m *femovement.Movement) string {
			return m.Get(prop).String()
		}
		getLabel = func(v string) string { return v }
	}

	attrib := ""
	for _, m := range mtm.Movements.Movements {
		attrib = getValue(m)
		if _, exist := count[attrib]; !exist {
			attribs = append(attribs, attrib)
		}
		count[attrib]++
	}
	sort.Strings(attribs)
	res := []*elements.ValText{}
	for _, a := range attribs {
		fa := getLabel(a)
		if fa == "" {
			fa = "Vide"
		}
		res = append(res, elements.NewValText(a, fa+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Actions Functions

func (mtm *MovementsTableModel) GetFilteredMovements() []*femovement.Movement {
	filter := func(*femovement.Movement) bool {
		return true
	}
	if !(mtm.FilterType == movementconst.FilterValueAll && mtm.Filter == "") {
		expected := strings.ToUpper(mtm.Filter)
		filter = func(m *femovement.Movement) bool {
			ss := m.SearchString(mtm.FilterType)
			if ss == "" {
				return false
			}
			return strings.Contains(strings.ToUpper(ss), expected)
		}
	}

	// filter movements in mvts slice to prevent change on mtm.Movements caused by sort
	var mvts []*femovement.Movement
	for _, m := range mtm.Movements.Movements {
		if filter(m) {
			mvts = append(mvts, m)
		}
	}
	return mvts
}
