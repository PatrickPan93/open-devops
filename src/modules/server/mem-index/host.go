package mem_index

import (
	"github.com/go-kit/log"
	ii "github.com/ning1875/inverted-index"
)

type HostIndex struct {
	Ir      *ii.HeadIndexReader
	Logger  log.Logger
	Modulus int // 静态分片的模
	Num     int
}

func (hi *HostIndex) FlushIndex() {

}

func (hi *HostIndex) GetIndexReader() *ii.HeadIndexReader {
	return hi.Ir
}
