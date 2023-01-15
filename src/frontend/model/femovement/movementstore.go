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
}

func NewMovementStore() *MovementStore {
	as := &MovementStore{Object: tools.O()}
	as.Movements = []*Movement{}
	as.Ref = ref.NewRef(func() string {
		return json.Stringify(as.Movements)
	})
	return as
}

func (as *MovementStore) CallGetMovements(vm *hvue.VM, onSuccess func()) {
	go as.callGetMovements(vm, onSuccess)
}

func (as *MovementStore) callGetMovements(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("GET", "/api/movements")
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
	as.Movements = loadedMovements
	as.Ref.SetReference()
	onSuccess()
}

func (as *MovementStore) CallUpdateMovements(vm *hvue.VM, onSuccess func()) {
	go as.callUpdateMovements(vm, onSuccess)
}

func (as *MovementStore) callUpdateMovements(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("PUT", "/api/movements")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	toUpdates := as.getUpdatedMovements()
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

	as.Ref.SetReference()
	msg := " movement mis à jour"
	if nbToUpd > 1 {
		msg = " movements mis à jour"
	}
	message.NotifySuccess(vm, "Sauvegarde des movements", strconv.Itoa(nbToUpd)+msg)
	onSuccess()

}

func (as *MovementStore) getUpdatedMovements() []*Movement {
	refMovements := []*Movement{}
	json.Parse(as.Ref.Reference).Call("forEach", func(acc *Movement) {
		refMovements = append(refMovements, acc)
	})
	refDict := makeDictMovements(refMovements)

	updtMovements := []*Movement{}
	for _, movement := range as.Movements {
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
