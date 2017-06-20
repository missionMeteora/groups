package groups

import (
	"encoding/json"

	"github.com/itsmontoya/turtle"
)

func marshal(val turtle.Value) (b []byte, err error) {
	var (
		gm groupMap
		ok bool
	)

	if gm, ok = val.(groupMap); !ok {
		err = ErrInvalidType
		return
	}

	return json.Marshal(gm.Slice())
}

func unmarshal(b []byte) (val turtle.Value, err error) {
	var (
		gs []string
		gm groupMap
	)

	if err = json.Unmarshal(b, &gs); err != nil {
		return
	}

	gm = make(groupMap, len(gs))
	for _, group := range gs {
		gm.Set(group)
	}

	val = gm
	return
}
