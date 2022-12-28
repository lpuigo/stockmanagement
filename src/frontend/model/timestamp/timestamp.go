package timestamp

type FeTimeStamp struct {
	CTime string `js:"CTime"`
	UTime string `js:"UTime"`
	DTime string `js:"DTime"`
}

func (ts *FeTimeStamp) Init() {
	ts.CTime = ""
	ts.UTime = ""
	ts.DTime = ""
}
