package consensus

import (
	"doozer/store"
	"sort"
)


type Run struct {
	Seqn int64
	Cals []string

	sink sink
}


func (r *Run) Deliver(p Packet) {
	r.sink.done = true
}


func GenerateRuns(alpha int64, w <-chan store.Event, runs chan<- Run) {
	for e := range w {
		runs <- Run{Seqn: e.Seqn + alpha, Cals: getCals(e)}
	}
}

func getCals(g store.Getter) []string {
	slots := store.Getdir(g, "/doozer/slot")
	cals  := make([]string, len(slots))

	for i, slot := range slots {
		cals[i] = store.GetString(g, "/doozer/slot/"+slot)
	}

	sort.SortStrings(cals)

	return cals
}
