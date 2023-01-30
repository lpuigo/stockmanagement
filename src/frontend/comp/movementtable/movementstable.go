package movementtable

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle/articleconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
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
	<el-table-column label="Désignation" prop="Designation" width="200px"
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

`
)

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("movements-table", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.Props("value", "user", "filter", "filtertype"),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewMovementsTableModel(vm)
		}),
		hvue.MethodsOf(&MovementsTableModel{}),
		hvue.Computed("filteredArticles", func(vm *hvue.VM) interface{} {
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
	Filter     string                    `js:"filter"`
	FilterType string                    `js:"filtertype"`

	VM *hvue.VM `js:"VM"`
}

func NewMovementsTableModel(vm *hvue.VM) *MovementsTableModel {
	atm := &MovementsTableModel{Object: tools.O()}
	atm.VM = vm
	atm.Movements = femovement.NewMovementStore()
	atm.User = feuser.NewUser()
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
	message.NotifyWarning(vm, "Double Click Movement", "front/comp/movementtable/HandleDoubleClickedRow à implémenter")
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Column Filtering Related Methods

// FilteredCategoryValue returns pre filtered values for Category
func (mtm *MovementsTableModel) FilteredCategoryValue() []string {
	return []string{}
}

func (mtm *MovementsTableModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (mtm *MovementsTableModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	mtm = MovementsTableModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	var getValue func(m *femovement.Movement) string
	switch prop {
	//case "Status":
	//	getValue = func(m *femovement.Movement) string {
	//		return m.Status()
	//	}
	default:
		getValue = func(m *femovement.Movement) string {
			return m.Get(prop).String()
		}
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

func (mtm *MovementsTableModel) GetFilteredMovements() []*femovement.Movement {
	filter := func(ar *femovement.Movement) bool {
		return true
	}
	if !(mtm.FilterType == articleconst.FilterValueAll && mtm.Filter == "") {
		expected := strings.ToUpper(mtm.Filter)
		filter = func(ar *femovement.Movement) bool {
			ss := ar.SearchString(mtm.FilterType)
			if ss == "" {
				return false
			}
			return strings.Contains(strings.ToUpper(ss), expected)
		}
	}

	// filter movements in mvts slice to prevent change on mtm.Movements caused by sort
	var mvts []*femovement.Movement
	for _, a := range mtm.Movements.Movements {
		if filter(a) {
			mvts = append(mvts, a)
		}
	}
	return mvts
}
