/***************************************************************************
 *
 * Copyright (c) 2017 Baidu.com, Inc. All Rights Reserved
 * @author duanbing(duanbing@baidu.com)
 *
 **************************************************************************/

/**
 * @filename main_statedb.go
 * @desc
 * @create time 2018-05-07 09:40:41
**/
package main

import (
	//"github.com/duanbing/go-evm/state"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	mdb, err := ethdb.NewMemDatabase()
	tr, err := trie.New(common.Hash{}, mdb)
	must(err)
	key := []byte("abc")
	value := []byte("123")
	must(tr.TryUpdate(key, value))

	key = []byte("ab")
	value = []byte("abc")
	must(tr.TryUpdate(key, value))

	fmt.Println("----------")
	root, err := tr.Commit()
	must(err)
	tr2, err := trie.New(root, mdb)
	must(err)
	key = []byte("abcd")
	value = []byte("你好")
	must(tr2.TryUpdate(key, value))
}
