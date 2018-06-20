/***************************************************************************
 *
 * Copyright (c) 2017 Baidu.com, Inc. All Rights Reserved
 * @author duanbing(duanbing@baidu.com)
 *
 **************************************************************************/

/**
 * @filename s.go
 * @desc
 * @create time 2018-05-11 16:56:54
**/
package core

import (
	"math/big"

	"github.com/duanbing/go-evm/state"
	"github.com/duanbing/go-evm/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

type StateObject struct {
	db           state.Database
	mdb          *ethdb.LDBDatabase
	stateStorage *state.StateDB
	xlog         log.Logger
}

func NewStateObject(root common.Hash, leveldbPath string) (*StateObject, error) {
	xlog := log.New("evm", "s")
	mdb, err := ethdb.NewLDBDatabase(leveldbPath, 100, 100)
	if err != nil {
		return nil, err
	}
	db := state.NewDatabase(mdb)
	s, err := state.New(root, db)
	return &StateObject{stateStorage: s, xlog: xlog, db: db, mdb: mdb}, nil
}

func (s *StateObject) Close() {
	s.mdb.Close()
}

func (s *StateObject) CreateAccount(addr common.Address) {
	s.stateStorage.CreateAccount(addr)
}

//BEGIN user-defined part combined with your own blockchain
func (s *StateObject) SubBalance(addr common.Address, b *big.Int) {
	s.stateStorage.SubBalance(addr, b)
}
func (s *StateObject) AddBalance(addr common.Address, b *big.Int) {
	s.stateStorage.AddBalance(addr, b)
}
func (s *StateObject) GetBalance(addr common.Address) *big.Int {
	return big.NewInt(10e+16)
}
func (s *StateObject) GetNonce(addr common.Address) uint64 {
	return 0
}
func (s *StateObject) SetNonce(addr common.Address, nonce uint64) {}

//END

func (s *StateObject) GetCodeHash(addr common.Address) common.Hash {
	return s.stateStorage.GetCodeHash(addr)
}
func (s *StateObject) GetCode(addr common.Address) []byte {
	return s.stateStorage.GetCode(addr)
}
func (s *StateObject) SetCode(addr common.Address, code []byte) {
	s.stateStorage.SetCode(addr, code)
}
func (s *StateObject) GetCodeSize(addr common.Address) int {
	return s.stateStorage.GetCodeSize(addr)
}

func (s *StateObject) AddRefund(uint64)  {}
func (s *StateObject) GetRefund() uint64 { return 0 }

func (s *StateObject) GetState(addr common.Address, h common.Hash) common.Hash {
	return s.stateStorage.GetState(addr, h)
}
func (s *StateObject) SetState(addr common.Address, k common.Hash, v common.Hash) {
	s.stateStorage.SetState(addr, k, v)
}

func (s *StateObject) Suicide(addr common.Address) bool {
	return s.stateStorage.Suicide(addr)
}
func (s *StateObject) HasSuicided(addr common.Address) bool {
	return s.stateStorage.HasSuicided(addr)
}
func (s *StateObject) Exist(addr common.Address) bool {
	return s.stateStorage.Exist(addr)
}
func (s *StateObject) Empty(addr common.Address) bool {
	return s.stateStorage.Empty(addr)
}

func (s *StateObject) RevertToSnapshot(revid int) {
	s.stateStorage.RevertToSnapshot(revid)
}

func (s *StateObject) Snapshot() int {
	return s.stateStorage.Snapshot()
}

func (s *StateObject) AddLog(log *types.Log) {
	s.stateStorage.AddLog(log)
}

func (s *StateObject) Logs() []*types.Log {
	return s.stateStorage.Logs()
}

func (s *StateObject) AddPreimage(common.Hash, []byte) {}

func (s *StateObject) ForEachStorage(addr common.Address, cb func(common.Hash, common.Hash) bool) {
	s.stateStorage.ForEachStorage(addr, cb)
}

func (s *StateObject) Commit(deleteEmptyObjects bool) (root common.Hash, err error) {
	root, err = s.stateStorage.Commit(true)
	if err != nil {
		return root, err
	}
	s.stateStorage.Finalise(deleteEmptyObjects)
	err = s.db.TrieDB().Commit(root, true)
	return
}

func (s *StateObject) Reset(root common.Hash) error {
	s.stateStorage.Reset(root)
	return nil
}
