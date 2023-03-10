package festock

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/ref"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/json"
	"github.com/lpuigo/hvue"
	"honnef.co/go/js/xhr"
	"sort"
	"strconv"
)

type StockStore struct {
	*js.Object

	Stocks []*Stock `js:"Stocks"`

	Ref *ref.Ref `js:"Ref"`
}

func NewStockStore() *StockStore {
	ss := &StockStore{Object: tools.O()}
	ss.Stocks = []*Stock{}
	ss.Ref = ref.NewRef(func() string {
		return json.Stringify(ss.Stocks)
	})
	return ss
}

func (ss *StockStore) CallGetStockById(vm *hvue.VM, sid int, onSuccess func()) {
	go ss.callGetStockById(vm, sid, onSuccess)
}

func (ss *StockStore) callGetStockById(vm *hvue.VM, sid int, onSuccess func()) {
	req := xhr.NewRequest("GET", "/api/stocks/"+strconv.Itoa(sid))
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
	loadedStock := StockFromJS(req.Response)
	//ss.Get("Stocks").Call("push", loadedStock)
	ss.Stocks = append(ss.Stocks, loadedStock)
	ss.Ref.SetReference()
	onSuccess()
}

func (ss *StockStore) CallGetStocks(vm *hvue.VM, onSuccess func()) {
	go ss.callGetStocks(vm, onSuccess)
}

func (ss *StockStore) callGetStocks(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("GET", "/api/stocks")
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
	loadedStocks := []*Stock{}
	req.Response.Call("forEach", func(item *js.Object) {
		w := StockFromJS(item)
		loadedStocks = append(loadedStocks, w)
	})

	sort.Slice(loadedStocks, func(i, j int) bool {
		return loadedStocks[i].Ref < loadedStocks[j].Ref
	})

	ss.Stocks = loadedStocks
	ss.Ref.SetReference()
	onSuccess()
}

func (ss *StockStore) CallUpdateStocks(vm *hvue.VM, onSuccess func()) {
	go ss.callUpdateStocks(vm, onSuccess)
}

func (ss *StockStore) callUpdateStocks(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("PUT", "/api/stocks")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	toUpdates := ss.getUpdatedStocks()
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

	ss.Ref.SetReference()
	msg := " Stock mis ?? jour"
	if nbToUpd > 1 {
		msg = " stocks mis ?? jour"
	}
	message.NotifySuccess(vm, "Sauvegarde des stocks", strconv.Itoa(nbToUpd)+msg)
	onSuccess()

}

func (ss *StockStore) getUpdatedStocks() []*Stock {
	refStocks := []*Stock{}
	json.Parse(ss.Ref.Reference).Call("forEach", func(acc *Stock) {
		refStocks = append(refStocks, acc)
	})
	refDict := makeDictStocks(refStocks)

	updtStocks := []*Stock{}
	for _, Stock := range ss.Stocks {
		refAcc := refDict[Stock.Id]
		if !(refAcc != nil && json.Stringify(Stock) == json.Stringify(refAcc)) {
			// this Stock has been updated ...
			updtStocks = append(updtStocks, Stock)
		}
	}
	return updtStocks
}

func makeDictStocks(accs []*Stock) map[int]*Stock {
	res := make(map[int]*Stock)
	for _, acc := range accs {
		res[acc.Id] = acc
	}
	return res
}
