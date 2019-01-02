// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestHardForkMethods tries to test the congruence of HF/EIP config blocks and methods.
// It might be one enormous tautology, but at least it looks comprehensive.
func TestHardForkMethods(t *testing.T) {
	mainc := MainnetChainConfig
	// test fork/eip counts
	type testCase struct {
		name              string
		hardforkBlock     *big.Int
		hardforkForkFn    func(*big.Int) bool
		hardforkEIPBlocks []*big.Int
		hardForkEIPFns    []func(*big.Int) bool
		len               int // testing the length of the expected fork/eip blocks ensures that at least we can count properly
	}
	buildTestCasesForConfig := func(conf *ChainConfig) []testCase {
		return []testCase{
			{"Homestead", mainc.HomesteadBlock, mainc.IsHomestead, []*big.Int{mainc.EIP7Block}, []func(*big.Int) bool{
				mainc.IsEIP7,
			}, 1},
			{"Byzantium", mainc.ByzantiumBlock, mainc.IsByzantium, mainc.ByzantiumEIPBlocks(), []func(*big.Int) bool{
				mainc.IsEIP100,
				mainc.IsEIP140,
				mainc.IsEIP198,
				mainc.IsEIP211,
				mainc.IsEIP212,
				mainc.IsEIP213,
				mainc.IsEIP214,
				mainc.IsEIP649,
				mainc.IsEIP658,
			}, 9},
			{"Constantinople", mainc.ConstantinopleBlock, mainc.IsConstantinople, mainc.ConstantinopleEIPBlocks(), []func(*big.Int) bool{
				mainc.IsEIP145,
				mainc.IsEIP1014,
				mainc.IsEIP1052,
				mainc.IsEIP1283,
				mainc.IsEIP1234,
			}, 5},
		}
	}
	createTestBlockVals := func(bl *big.Int) (vals []*big.Int) {
		vals = append(vals, new(big.Int))
		if bl == nil {
			return
		}
		vals = append(vals, new(big.Int).Sub(bl, common.Big1))
		vals = append(vals, new(big.Int).Set(bl))
		vals = append(vals, new(big.Int).Add(bl, common.Big1))
		return
	}
	testHFEIPFns := func(name string, n *big.Int, isHF func(*big.Int) bool, isEIPs []func(*big.Int) bool) {
		for i, v := range createTestBlockVals(n) {
			for j, f := range isEIPs {
				if isHF(v) != f(v) {
					t.Errorf("i: %v, j: %v, want: %v, got: %v", i, j, isHF(v), f(v))
				}
			}
		}
	}
	runInterchangeabilityTest := func(tc testCase, f func(name string, n *big.Int, isHF func(*big.Int) bool, isEIPs []func(*big.Int) bool)) {
		f(tc.name, tc.hardforkBlock, tc.hardforkForkFn, tc.hardForkEIPFns)
		for k := range tc.hardforkEIPBlocks {
			tc.hardforkEIPBlocks[k] = new(big.Int).Set(tc.hardforkBlock)
		}
		f(tc.name, tc.hardforkBlock, tc.hardforkForkFn, tc.hardForkEIPFns)
		tc.hardforkBlock = nil
		f(tc.name, tc.hardforkBlock, tc.hardforkForkFn, tc.hardForkEIPFns)
	}
	for _, conf := range []*ChainConfig{MainnetChainConfig, TestnetChainConfig, RinkebyChainConfig, AllEthashProtocolChanges, AllCliqueProtocolChanges} {
		for i, c := range buildTestCasesForConfig(conf) {
			got := len(c.hardforkEIPBlocks)
			if got != c.len {
				t.Errorf("i: %d, want: %v, got: %v", i, c.len, got)
			}
			runInterchangeabilityTest(c, testHFEIPFns)
		}
	}
}

func TestCheckCompatible(t *testing.T) {
	type test struct {
		stored, new *ChainConfig
		head        uint64
		wantErr     *ConfigCompatError
	}
	tests := []test{
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 0, wantErr: nil},
		{stored: AllEthashProtocolChanges, new: AllEthashProtocolChanges, head: 100, wantErr: nil},
		{
			stored:  &ChainConfig{EIP150Block: big.NewInt(10)},
			new:     &ChainConfig{EIP150Block: big.NewInt(20)},
			head:    9,
			wantErr: nil,
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &ChainConfig{HomesteadBlock: nil},
			head:   3,
			wantErr: &ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: big.NewInt(0),
				NewConfig:    nil,
				RewindTo:     0,
			},
		},
		{
			stored: AllEthashProtocolChanges,
			new:    &ChainConfig{HomesteadBlock: big.NewInt(1)},
			head:   3,
			wantErr: &ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: big.NewInt(0),
				NewConfig:    big.NewInt(1),
				RewindTo:     0,
			},
		},
		{
			stored: &ChainConfig{HomesteadBlock: big.NewInt(30), EIP150Block: big.NewInt(10)},
			new:    &ChainConfig{HomesteadBlock: big.NewInt(25), EIP150Block: big.NewInt(20)},
			head:   25,
			wantErr: &ConfigCompatError{
				What:         "EIP150 fork block",
				StoredConfig: big.NewInt(10),
				NewConfig:    big.NewInt(20),
				RewindTo:     9,
			},
		},
	}

	for _, test := range tests {
		err := test.stored.CheckCompatible(test.new, test.head)
		if !reflect.DeepEqual(err, test.wantErr) {
			t.Errorf("error mismatch:\nstored: %v\nnew: %v\nhead: %v\nerr: %v\nwant: %v", test.stored, test.new, test.head, err, test.wantErr)
		}
	}
}
