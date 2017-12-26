package mmu

import (
	"github.com/lbarrios/yesSGMB/types"
)

type Peripheral interface {
	MapByte(logical_address types.Address, physical_address *byte)
}
