package feworksite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/ref"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/fedate"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/json"
	"github.com/lpuigo/hvue"
	"honnef.co/go/js/xhr"
	"strconv"
)

type WorksiteStore struct {
	*js.Object

	Worksites []*Worksite `js:"Worksites"`

	Ref *ref.Ref `js:"Ref"`
}

func NewWorksiteStore() *WorksiteStore {
	as := &WorksiteStore{Object: tools.O()}
	as.Worksites = []*Worksite{}
	as.Ref = ref.NewRef(func() string {
		return json.Stringify(as.Worksites)
	})
	return as
}

func (ws *WorksiteStore) GetWorksiteById(id int) *Worksite {
	for _, worksite := range ws.Worksites {
		if worksite.Id == id {
			return worksite
		}
	}
	return NewWorksite()
}

func (ws *WorksiteStore) CallGetWorksites(vm *hvue.VM, onSuccess func()) {
	go ws.callGetWorksites(vm, onSuccess)
}

func (ws *WorksiteStore) callGetWorksites(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("GET", "/api/worksites")
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
	loadedWorksites := []*Worksite{}
	req.Response.Call("forEach", func(item *js.Object) {
		w := WorksiteFromJS(item)
		loadedWorksites = append(loadedWorksites, w)
	})
	ws.Worksites = loadedWorksites
	ws.Ref.SetReference()
	onSuccess()
}

func (ws *WorksiteStore) CallUpdateWorksites(vm *hvue.VM, onSuccess func()) {
	go ws.callUpdateWorksites(vm, onSuccess)
}

func (ws *WorksiteStore) callUpdateWorksites(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("PUT", "/api/worksites")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	toUpdates := ws.getUpdatedWorksites()
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

	ws.Ref.SetReference()
	msg := " Worksite mis à jour"
	if nbToUpd > 1 {
		msg = " worksites mis à jour"
	}
	message.NotifySuccess(vm, "Sauvegarde des worksites", strconv.Itoa(nbToUpd)+msg)
	onSuccess()

}

func (ws *WorksiteStore) getUpdatedWorksites() []*Worksite {
	refWorksites := []*Worksite{}
	json.Parse(ws.Ref.Reference).Call("forEach", func(acc *Worksite) {
		refWorksites = append(refWorksites, acc)
	})
	refDict := makeDictWorksites(refWorksites)

	updtWorksites := []*Worksite{}
	for _, Worksite := range ws.Worksites {
		refAcc := refDict[Worksite.Id]
		if !(refAcc != nil && json.Stringify(Worksite) == json.Stringify(refAcc)) {
			// this Worksite has been updated ...
			updtWorksites = append(updtWorksites, Worksite)
		}
	}
	return updtWorksites
}

func makeDictWorksites(accs []*Worksite) map[int]*Worksite {
	res := make(map[int]*Worksite)
	for _, acc := range accs {
		res[acc.Id] = acc
	}
	return res
}

func (ws *WorksiteStore) GetActiveWorksites() []*Worksite {
	today := fedate.TodayAfter(0)
	res := []*Worksite{}
	for _, worksite := range ws.Worksites {
		if !worksite.IsActiveOn(today) {
			continue
		}
		res = append(res, worksite)
	}
	return res
}
