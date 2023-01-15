package feactor

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

type ActorStore struct {
	*js.Object

	Actors []*Actor `js:"Actors"`

	Ref *ref.Ref `js:"Ref"`
}

func NewActorStore() *ActorStore {
	as := &ActorStore{Object: tools.O()}
	as.Actors = []*Actor{}
	as.Ref = ref.NewRef(func() string {
		return json.Stringify(as.Actors)
	})
	return as
}

func (ws *ActorStore) CallGetActors(vm *hvue.VM, onSuccess func()) {
	go ws.callGetActors(vm, onSuccess)
}

func (ws *ActorStore) callGetActors(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("GET", "/api/actors")
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
	loadedActors := []*Actor{}
	req.Response.Call("forEach", func(item *js.Object) {
		w := ActorFromJS(item)
		loadedActors = append(loadedActors, w)
	})
	ws.Actors = loadedActors
	ws.Ref.SetReference()
	onSuccess()
}

func (ws *ActorStore) CallUpdateActors(vm *hvue.VM, onSuccess func()) {
	go ws.callUpdateActors(vm, onSuccess)
}

func (ws *ActorStore) callUpdateActors(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("PUT", "/api/actors")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	toUpdates := ws.getUpdatedActors()
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
	msg := " Actor mis à jour"
	if nbToUpd > 1 {
		msg = " actors mis à jour"
	}
	message.NotifySuccess(vm, "Sauvegarde des actors", strconv.Itoa(nbToUpd)+msg)
	onSuccess()

}

func (ws *ActorStore) getUpdatedActors() []*Actor {
	refActors := []*Actor{}
	json.Parse(ws.Ref.Reference).Call("forEach", func(acc *Actor) {
		refActors = append(refActors, acc)
	})
	refDict := makeDictActors(refActors)

	updtActors := []*Actor{}
	for _, Actor := range ws.Actors {
		refAcc := refDict[Actor.Id]
		if !(refAcc != nil && json.Stringify(Actor) == json.Stringify(refAcc)) {
			// this Actor has been updated ...
			updtActors = append(updtActors, Actor)
		}
	}
	return updtActors
}

func makeDictActors(accs []*Actor) map[int]*Actor {
	res := make(map[int]*Actor)
	for _, acc := range accs {
		res[acc.Id] = acc
	}
	return res
}
