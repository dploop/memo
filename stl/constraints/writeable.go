package constraints

import (
	"github.com/dploop/memo/stl/types"
)

type Writeable interface {
	Write(data types.Data)
}
