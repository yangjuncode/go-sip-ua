package ua

import (
	"sync"
)

const (
	defaultInitialCSeq uint32 = 1
)

type CSeqManager struct {
	mu      sync.Mutex
	nextSeq uint32
}

func NewCSeqManager(seed uint32) *CSeqManager {
	if seed == 0 {
		seed = defaultInitialCSeq
	}

	return &CSeqManager{
		nextSeq: seed,
	}
}

func (m *CSeqManager) Next() uint32 {
	m.mu.Lock()
	defer m.mu.Unlock()

	seq := m.nextSeq
	m.nextSeq++
	return seq
}