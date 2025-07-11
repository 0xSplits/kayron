package runtime

import (
	"encoding/json"

	"github.com/xh3b4sd/tracer"
)

var (
	byt []byte = nil
)

var (
	dic = map[string]string{
		"arc": Arc(),
		"gos": Gos(),
		"sha": Sha(),
		"src": Src(),
		"tag": Tag(),
		"ver": Ver(),
	}
)

func init() {
	var err error

	{
		byt, err = json.Marshal(dic)
		if err != nil {
			tracer.Panic(tracer.Mask(err))
		}
	}
}

func JSON() []byte {
	return byt
}
