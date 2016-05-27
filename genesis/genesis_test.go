// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package genesis_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bitmark-inc/bitmarkd/blockdigest"
	"github.com/bitmark-inc/bitmarkd/difficulty"
	"github.com/bitmark-inc/bitmarkd/genesis"
	"github.com/bitmark-inc/bitmarkd/merkle"
	"github.com/bitmark-inc/bitmarkd/util"
	"strings"
	"testing"
	"time"
)

// some constants embedded into the genesis block
const (
	genesisBlockNumber = uint64(1)
)

// some data embedded into the genesis block
// for live chain
var (
	// date -u -r $(printf '%d\n' 0x56809ab7)
	// Mon 28 Dec 2015 02:13:11 UTC
	// date -u -r $(printf '%d\n' 0x56809ab7) '+%FT%TZ'
	// 2015-12-28T02:13:11Z

	genesisLiveTimestamp = TS{0x56809ab7, "2015-12-28T02:13:11Z"}

	genesisLiveAddresses = []block.MinerAddress{
		{
			Currency: "",
			Address:  "DOWN the RABBIT hole",
		},
	}
	genesisLiveRawAddress = "\x00\x14" + "DOWN the RABBIT hole"
)

// some data embedded into the genesis block
// for test chain
var (
	// date -u -r $(printf '%d\n' 0x5478424b)
	// Fri Nov 28 09:37:15 UTC 2014
	// date -u -r $(printf '%d\n' 0x5478424b) '+%FT%TZ'
	// 2014-11-28T09:37:15Z

	genesisTestTimestamp = TS{0x5478424b, "2014-11-28T09:37:15Z"}

	// for testing chain
	genesisTestAddresses = []block.MinerAddress{
		{
			Currency: "",
			Address:  "Bitmark Testing Genesis Block",
		},
	}
	genesisTestRawAddress = "\x00\x1d" + "Bitmark Testing Genesis Block"
)

// create the live genesis block
//
// MinerUsername: "miner-arm"  JobId: "Live_efd7b4fe"  ExtraNonce2: "00000000"  Ntime: "56809ab7"  Nonce: "826f9a87"
func TestLiveGenesisAssembly(t *testing.T) {

	// fixed data used to create genesis block
	// ---------------------------------------

	// nonce provided by statum sserver
	extraNonce1 := []byte{0xef, 0xd7, 0xb4, 0xfe}

	// nonces obtained from miner
	nonce := uint32(0x826f9a87)
	extraNonce2 := []byte{0x00, 0x00, 0x00, 0x00}

	doCalc(t, "Live", genesisLiveTimestamp, extraNonce1, extraNonce2, nonce, genesisLiveAddresses, genesisLiveRawAddress, genesis.LiveGenesisDigest, genesis.LiveGenesisBlock)
}

// create the test genesis block
//
// MinerUsername: "miner-arm"  JobId: "Test_b201475b"  ExtraNonce2: "00000000"  Ntime: "5478424b"  Nonce: "1e26bad4"
func TestTestGenesisAssembly(t *testing.T) {

	// fixed data used to create genesis block
	// ---------------------------------------

	// nonce provided by statum server
	extraNonce1 := []byte{0xb2, 0x01, 0x47, 0x5b}

	// nonces obtained from miner
	nonce := uint32(0x1e26bad4)
	extraNonce2 := []byte{0x00, 0x00, 0x00, 0x00}

	doCalc(t, "Test", genesisTestTimestamp, extraNonce1, extraNonce2, nonce, genesisTestAddresses, genesisTestRawAddress, genesis.TestGenesisDigest, genesis.TestGenesisBlock)
}

// hold chain specific timestamp
type TS struct {
	ntime uint32
	utc   string
}

func doCalc(t *testing.T, title string, ts TS, extraNonce1 []byte, extraNonce2 []byte, nonce uint32, addresses []block.MinerAddress, rawAddress string, gDigest blockdigest.Digest, gBlock blockrecord.PackedHeader) {

	timestamp, err := time.Parse(time.RFC3339, ts.utc)
	if nil != err {
		t.Fatalf("failed to parse time: err = %v", err)
	}
	timeUint64 := uint64(timestamp.UTC().Unix())
	ntime := ts.ntime
	if timeUint64 != uint64(ntime) {
		t.Fatalf("time converted to: 0x%08x  expectd: %08x", timeUint64, ntime)
	}

	// some common static data
	version := uint32(1) // snapshot of version number
	previousBlock := blockdigest.Digest{}

	// Just calculations after this point
	// ----------------------------------

	coinbase := block.NewFullCoinbase(genesisBlockNumber, timestamp, append(extraNonce1, extraNonce2...), addresses)
	cDigest := merkle.NewDigest(coinbase)
	coinbaseLength := len(coinbase)

	transactionCount := 1

	// merkle tree
	tree := merkle.FullMerkleTree(cDigest, []merkle.Digest{})
	if tree[len(tree)-1] != cDigest {
		t.Fatalf("failed to compute tree: actual: %#v  expected: %#v", tree[len(tree)-1], cDigest)
	}

	// default difficulty
	bits := difficulty.New() // defaults to 1

	// block header
	h := blockrecord.Header{
		Version:       version,
		PreviousBlock: previousBlock,
		MerkleRoot:    tree[len(tree)-1],
		Time:          ntime,
		Bits:          *bits,
		Nonce:         nonce,
	}

	header := h.Pack()
	hDigest := header.Digest()

	// ok - log the header and coinbase data
	t.Logf("Title: %s", title)
	t.Logf("header: %#v\n", h)
	t.Logf("packed header: %x\n", header)
	t.Logf("coinbase: %x\n", coinbase)
	t.Logf("coinbase digest: %#v\n", cDigest)
	t.Logf("merkle tree: %#v\n", tree)
	t.Logf("merkle root little endian hex: %x\n", [blockdigest.Length]byte(tree[0]))
	t.Logf("hDigest: %#v\n", hDigest)
	t.Logf("hDigest little endian hex: %x\n", [blockdigest.Length]byte(hDigest))

	t.Log(util.FormatBytes(title+"ProposedLEhash", []byte(hDigest[:])))

	// chack that it matches
	if hDigest != gDigest {
		t.Errorf("digest mismatch actual: %#v  expected: %#v", hDigest, gDigest)

		hexExtraNonce1 := fmt.Sprintf("%08x", extraNonce1)
		hexCoinbase := hex.EncodeToString(coinbase)
		n1 := strings.Index(hexCoinbase, hexExtraNonce1)
		n2 := n1 + 2*(len(extraNonce1)+len(extraNonce2)) // since 2 hex chars = 1 byte
		login := []struct {
			ID     interface{}   `json:"id"`
			Method string        `json:"method"`
			Result interface{}   `json:"result"`
			Error  []interface{} `json:"error"`
		}{
			{
				ID:     1,
				Method: "mining.subscribe",
				Result: []interface{}{
					[][]string{
						{"mining.set_difficulty", "1357"},
						{"mining.notify", "1234"},
					},
					hexExtraNonce1,
					len(extraNonce2),
				},
			},
			{
				ID:     2,
				Method: "mining.authorize",
				Result: true,
			},
			{
				ID:     3,
				Method: "mining.extranonce.subscribe",
				Result: true,
			},
		}

		requests := []struct {
			Method string        `json:"method"`
			Params []interface{} `json:"params"`
		}{
			{
				Method: "mining.set_difficulty",
				Params: []interface{}{
					difficulty.Current.Reciprocal(),
				},
			},
			{
				Method: "mining.set_extranonce",
				Params: []interface{}{
					"08000002",
					4,
				},
			},
			{
				Method: "mining.notify",
				Params: []interface{}{
					title, // [0] job_id
					fmt.Sprintf("%s", previousBlock), // [1] previous link
					hexCoinbase[:n1],                 // [2] coinbase 1
					hexCoinbase[n2:],                 // [3] coinbase 2
					[]interface{}{},                  // [4] minimised merkle tree (empty)
					fmt.Sprintf("%08x", version),     // [5] version
					difficulty.Current.String(),      // [6] bits
					fmt.Sprintf("%08x", ntime),       // [7] time
					true, // [8] clean_jobs
				},
			},
		}

		for _, r := range login {
			b, err := json.Marshal(r)
			if nil != err {
				t.Errorf("json error: %v", err)
				return
			}

			t.Logf("JSON: %s", b)
		}
		for _, r := range requests {
			b, err := json.Marshal(r)
			if nil != err {
				t.Errorf("json error: %v", err)
				return
			}

			t.Logf("JSON: %s", b)
		}
	}

	// check difficulty
	if hDigest.Cmp(bits.BigInt()) > 0 {
		t.Errorf("difficulty NOT met\n")
	}

	// compute block size
	blockSize := len(header) + 2 + coinbaseLength + 2 + len(tree)*merkle.DigestLength

	// pack the block
	blk := make([]byte, 0, blockSize)
	blk = append(blk, header...)
	blk = append(blk, byte(coinbaseLength&0xff))
	blk = append(blk, byte(coinbaseLength>>8))
	blk = append(blk, coinbase...)
	blk = append(blk, byte(transactionCount&0xff))
	blk = append(blk, byte(transactionCount>>8))

	buffer := new(bytes.Buffer)
	err = binary.Write(buffer, binary.LittleEndian, tree)
	if nil != err {
		t.Fatalf("binary.Write: err = %v", err)
	}

	blk = append(blk, buffer.Bytes()...)

	if len(blk) != blockSize {
		t.Fatalf("block size mismatch: actual: %d, expected: %d", len(blk), blockSize)
	}

	if !bytes.Equal(blk, gBlock) {
		t.Errorf("initial block assembly mismatch actual: %x  expected: %x", blk, gBlock)
		t.Log(util.FormatBytes(title+"GenesisBlock", blk))
	}

	// unpack the block
	var unpacked block.Block
	err = blockrecord.PackedHeader(blk).Unpack(&unpacked)
	if nil != err {
		t.Fatalf("unpack block failed: err = %v", err)
	}

	if unpacked.Header.Time != ntime {
		t.Fatalf("block ntime mismatch: actual 0x%08x  expected 0x%08x", unpacked.Header.Time, ntime)
	}

	if unpacked.Timestamp != timestamp {
		t.Fatalf("block timestamp mismatch: actual %v  expected %v", unpacked.Timestamp, timestamp)
	}

	t.Logf("unpacked block: %#v", unpacked)

	// re-pack
	reDigest, rePacked, ok := block.Pack(unpacked.Number, timestamp, &unpacked.Header.Bits, unpacked.Header.Time, unpacked.Header.Nonce, append(extraNonce1, extraNonce2...), unpacked.Addresses, unpacked.TxIds)

	if !ok {
		t.Fatal("block.Pack failed")
	}

	if reDigest != gDigest {
		t.Fatalf("re-digest mismatch actual: %#v  expected: %#v", reDigest, gDigest)
	}

	if !bytes.Equal(rePacked, blk) {
		t.Fatalf("re-packed mismatch actual: %x  expected: %x", rePacked, blk)
	}

	// log the final result
	if verboseTesting { // turn on in all_test.go
		t.Logf("Genesis digest: %#v", reDigest)
		t.Logf("Genesis block:  %x", rePacked)
	}

	// hex dumps for genesis.go
	t.Log(formatBytes(title+"GenesisBlock", rePacked))
	t.Log(formatBytes(title+"GenesisDigest", reDigest[:]))

	// check that these match the current genesis block/digest
	if reDigest != gDigest {
		t.Fatalf("re-digest/Genesis mismatch actual: %#v  expected: %#v", reDigest, gDigest)
	}

	if !bytes.Equal(rePacked, gBlock) {
		t.Fatalf("re-packed/Genesis mismatch actual: %x  expected: %x", rePacked, gBlock)
	}
}

// test the real genesis block
func TestGenesisBlock(t *testing.T) {
	doReal(t, "Live", block.LiveGenesisDigest, block.LiveGenesisBlock)
	doReal(t, "Test", block.TestGenesisDigest, block.TestGenesisBlock)
}

func doReal(t *testing.T, title string, gDigest blockdigest.Digest, gBlock blockrecord.PackedHeader) {

	// unpack the block
	var unpacked block.Block
	err := block.LiveGenesisBlock.Unpack(&unpacked)
	if nil != err {
		t.Fatalf("unpack block failed: err = %v", err)
	}

	if verboseTesting { // turn on in all_test.go
		t.Logf("unpacked block: %v", unpacked)
	}

	// check current genesis digest matches
	if unpacked.Digest != block.LiveGenesisDigest {
		t.Fatalf("digest/Genesis mismatch actual: %#v  expected: %#v", unpacked.Digest, block.LiveGenesisDigest)
	}

	// check block number
	if unpacked.Number != genesisBlockNumber {
		t.Fatalf("block number: %d  expected %d", unpacked.Number, genesisBlockNumber)
	}

	// check the address matches
	if 1 != len(unpacked.Addresses) {
		t.Fatalf("Addresses: found: %d  expected: %d", len(unpacked.Addresses), 1)
	}
	if unpacked.Addresses[0].String() != genesisLiveRawAddress {
		t.Fatalf("RawAddress: %q  expected: %q", unpacked.Addresses[0].String(), genesisLiveRawAddress)
	}
}