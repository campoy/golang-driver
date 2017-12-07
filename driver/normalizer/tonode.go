package normalizer

import (
	"log"

	"gopkg.in/bblfsh/sdk.v1/uast"
)

// ToNode is an instance of `uast.ObjectToNode`, defining how to transform an
// into a UAST (`uast.Node`).
//
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast#ObjectToNode
var ToNode = &uast.ObjectToNode{
	InternalTypeKey:    "InternalType",
	OffsetKey:          "StartOffset",
	EndOffsetKey:       "EndOffset",
	TopLevelIsRootNode: true,

	Modifier: func(m map[string]interface{}) error {
		props, ok := m["Properties"].(map[string]interface{})
		if !ok {
			return nil
		}

		for k, v := range props {
			if _, ok := m[k]; ok {
				log.Printf("ignoring already defined property %s", k)
			} else {
				m[k] = v
				delete(props, k)
			}
		}
		delete(m, "Properties")
		return nil
	},
}
