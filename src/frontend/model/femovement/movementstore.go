package femovement

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/ref"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/json"
	"github.com/lpuigo/hvue"
	"honnef.co/go/js/xhr"
	"strconv"
)

type MovementStore struct {
	*js.Object

	Movements []*Movement `js:"Movements"`

	Ref *ref.Ref `js:"Ref"`

	NextNewId int `js:"NextNewId"`
}

func NewMovementStore() *MovementStore {
	as := &MovementStore{Object: tools.O()}
	as.Movements = []*Movement{}
	as.Ref = ref.NewRef(func() string {
		return json.Stringify(as.Movements)
	})
	as.NextNewId = 0
	return as
}

func (ms *MovementStore) GetNextNewId() int {
	ms.NextNewId--
	return ms.NextNewId
}

func (ms *MovementStore) AddNewMovement(mvt *Movement) {
	mvt.Id = ms.GetNextNewId()
	ms.Movements = append(ms.Movements, mvt)
}

func (ms *MovementStore) DeleteMovement(mvt *Movement) {
	for i, movement := range ms.Movements {
		if movement.Id == mvt.Id {
			ms.Get("Movements").Call("splice", i, 1)
			break
		}
	}
}

func (ms *MovementStore) CallGetMovementsForStockId(vm *hvue.VM, stockId int, onSuccess func()) {
	url := "/api/movements/stock/" + strconv.Itoa(stockId)
	go ms.callGetMovements(vm, url, onSuccess)
}

func (ms *MovementStore) CallGetMovements(vm *hvue.VM, onSuccess func()) {
	url := "/api/movements"
	go ms.callGetMovements(vm, url, onSuccess)
}

func (ms *MovementStore) callGetMovements(vm *hvue.VM, url string, onSuccess func()) {
	req := xhr.NewRequest("GET", url)
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(vm, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(vm, req)
		return
	}
	loadedMovements := []*Movement{}
	req.Response.Call("forEach", func(item *js.Object) {
		a := MovementFromJS(item)
		loadedMovements = append(loadedMovements, a)
	})
	ms.Movements = loadedMovements
	ms.Ref.SetReference()
	onSuccess()
}

func (ms *MovementStore) CallUpdateMovements(vm *hvue.VM, onSuccess func()) {
	go ms.callUpdateMovements(vm, onSuccess)
}

func (ms *MovementStore) callUpdateMovements(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("PUT", "/api/movements")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	toUpdates := ms.getUpdatedMovements()
	nbToUpd := len(toUpdates)
	if nbToUpd == 0 {
		onSuccess()
		return
	}

	err := req.Send(json.Stringify(toUpdates))
	if err != nil {
		message.ErrorStr(vm, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(vm, req)
		return
	}

	ms.Ref.SetReference()
	msg := " mouvement de stock mis à jour"
	if nbToUpd > 1 {
		msg = " mouvements de stock mis à jour"
	}
	message.NotifySuccess(vm, "Sauvegarde des movements", strconv.Itoa(nbToUpd)+msg)
	onSuccess()

}

func (ms *MovementStore) getUpdatedMovements() []*Movement {
	refMovements := []*Movement{}
	json.Parse(ms.Ref.Reference).Call("forEach", func(acc *Movement) {
		refMovements = append(refMovements, acc)
	})
	refDict := makeDictMovements(refMovements)

	updtMovements := []*Movement{}
	for _, movement := range ms.Movements {
		refAcc := refDict[movement.Id]
		if !(refAcc != nil && json.Stringify(movement) == json.Stringify(refAcc)) {
			// this movement has been updated ...
			updtMovements = append(updtMovements, movement)
		}
	}
	return updtMovements
}

func makeDictMovements(accs []*Movement) map[int]*Movement {
	res := make(map[int]*Movement)
	for _, acc := range accs {
		res[acc.Id] = acc
	}
	return res
}
