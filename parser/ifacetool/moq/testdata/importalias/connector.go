package importalias

import (
	dstclient "github.com/protogodev/protogo/parser/ifacetool/moq/testdata/importalias/dst/client"
	srcclient "github.com/protogodev/protogo/parser/ifacetool/moq/testdata/importalias/src/client"
)

type Connector interface {
	Connect(src srcclient.Client, dst dstclient.Client)
}
