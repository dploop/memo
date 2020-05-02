package constraints

import (
	"github.com/dploop/memo/stl/types"
)

type Readable interface {
	Read() types.Data
}
