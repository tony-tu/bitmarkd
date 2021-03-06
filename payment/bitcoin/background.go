// Copyright (c) 2014-2017 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package bitcoin

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/bitmark-inc/bitmarkd/currency"
	"github.com/bitmark-inc/bitmarkd/pay"
	"github.com/bitmark-inc/bitmarkd/storage"
	"github.com/bitmark-inc/bitmarkd/util"
	"github.com/bitmark-inc/logger"
	"time"
)

const (
	saveModulus          = 16     // to reduce fequency of rewrites of currency record
	hardForkBlockCount   = 6 * 24 // back one day in case of hard fork
	bitcoinConfirmations = 3      // stop processing this many blocks back from most recent block
	maximumBlockCount    = 500    // total blocks in one download
	maximumBlockRate     = 20.0   // blocks per second
)

// wait for new blocks
func (state *bitcoinData) Run(args interface{}, shutdown <-chan struct{}) {

	log := state.log

loop:
	for {
		log.Info("waiting…")
		select {
		case <-shutdown:
			break loop

		case <-time.After(60 * time.Second):
			var blockNumber uint64
			err := bitcoinCall("getblockcount", []interface{}{}, &blockNumber)
			if nil != err {
				continue loop
			}
			log.Infof("block number: %d", blockNumber)

			if blockNumber <= bitcoinConfirmations {
				continue loop
			}
			blockNumber -= bitcoinConfirmations
			if blockNumber <= state.latestBlockNumber {
				continue loop
			}
			n, hash := process(log, state.latestBlockNumber, blockNumber, state.latestBlockHash)
			if 0 == n || "" == hash {
				continue loop
			}

			state.saveCount += n - state.latestBlockNumber
			state.latestBlockNumber = n
			state.latestBlockHash = hash
			if state.saveCount >= saveModulus {

				state.saveCount = 0

				key := make([]byte, 8)
				value := make([]byte, 8+len(hash))
				binary.BigEndian.PutUint64(key, currency.Bitcoin.Uint64())
				binary.BigEndian.PutUint64(value, n)
				copy(value[8:], hash)
				storage.Pool.Currency.Put(key, value)
			}
		}
	}
}

const (
	bitcoin_OP_RETURN_HEX_CODE      = "6a30" // op code with 48 byte parameter
	bitcoin_OP_RETURN_PREFIX_LENGTH = len(bitcoin_OP_RETURN_HEX_CODE)
	bitcoin_OP_RETURN_PAY_ID_OFFSET = bitcoin_OP_RETURN_PREFIX_LENGTH
	bitcoin_OP_RETURN_RECORD_LENGTH = bitcoin_OP_RETURN_PREFIX_LENGTH + 2*48
)

func process(log *logger.L, startBlockNumber uint64, endBlockNumber uint64, lastHash string) (uint64, string) {
	//func processBlocks(log *logger.L, startBlockNumber uint64) {

	var hash string
	log.Infof("starting from block: %d", startBlockNumber)
	err := bitcoinGetBlockHash(startBlockNumber, &hash)
	if nil != err {
		log.Errorf("failed to get hash for block: %d", startBlockNumber)
		return 0, ""
	}
	// block rescan in case of hard fork
	if startBlockNumber >= hardForkBlockCount && lastHash != hash {
		startBlockNumber -= hardForkBlockCount
		log.Infof("fork detected: old hash: %q  hash: %q", lastHash, hash)
		log.Infof("fork detected: restarting from block: %d", startBlockNumber)
		err := bitcoinGetBlockHash(startBlockNumber, &hash)
		if nil != err {
			log.Errorf("failed to get hash for block: %d", startBlockNumber)
			return 0, ""
		}
	}

	// to record last block processed
	n := uint64(0)
	startTime := time.Now()
	counter := 0
loop:
	for {
		var block bitcoinBlock
		err = bitcoinGetBlock(hash, &block)
		if nil != err {
			log.Errorf("failed to get block for hash: %q", hash)
			break loop
		}

		log.Infof("block: %d  hash: %q", block.Height, block.Hash)
		log.Tracef("block contents: %#v", block)

		transationCount := len(block.Tx) // first is the coinbase and can be ignored
		if transationCount > 1 {
		txLoop:
			for i, txId := range block.Tx[1:] {
				// fetch transaction and decode
				var reply bitcoinTransaction
				err := bitcoinGetRawTransaction(txId, &reply)
				if nil != err {
					log.Errorf("failed to get block: %d  transaction[%d] for: %q", block.Height, i, txId)
					continue
				}
				for j, vout := range reply.Vout {
					if bitcoin_OP_RETURN_RECORD_LENGTH == len(vout.ScriptPubKey.Hex) && bitcoin_OP_RETURN_HEX_CODE == vout.ScriptPubKey.Hex[0:4] {
						var payId pay.PayId
						payIdBytes := []byte(vout.ScriptPubKey.Hex[bitcoin_OP_RETURN_PAY_ID_OFFSET:])
						err := payId.UnmarshalText(payIdBytes)
						if nil != err {
							log.Errorf("failed to get pay id: %d  transaction[%d] id: %q", block.Height, i, txId)
							continue txLoop
						}
						log.Infof("possible transaction: %#v", reply)
						scanTx(log, payId, j, &reply)
						continue txLoop
					}
				}

			}
		}

		counter += 1 // number of blocks processed

		n = block.Height // current block number
		if "" == block.NextBlockHash {
			break loop
		}
		if n >= endBlockNumber || counter >= maximumBlockCount {
			break loop
		}

		// rate limit
		timeTaken := time.Since(startTime)
		rate := float64(counter) / timeTaken.Seconds()
		if rate > maximumBlockRate {
			time.Sleep(2 * time.Second)
		}

		// set up to get next block
		hash = block.NextBlockHash
	}
	return n, hash
}

func scanTx(log *logger.L, payId pay.PayId, payIdIndex int, tx *bitcoinTransaction) {

	amounts := make(map[string]uint64)

	// extract payments, skipping already determine OP_RETURN vout
	for i, vout := range tx.Vout {
		log.Infof("vout[%d]: %v ", i, vout)
		if payIdIndex == i {
			continue
		}
		if 1 == len(vout.ScriptPubKey.Addresses) {
			amounts[vout.ScriptPubKey.Addresses[0]] += convertToSatoshi(vout.Value)
		}
	}

	if 0 == len(amounts) {
		log.Warnf("found pay id bu no payments in tx id: %s", tx.TxId)
		return
	}

	// create packed structure to store payment details
	packed := util.ToVarint64(currency.Bitcoin.Uint64())

	// transaction ID
	txId, err := hex.DecodeString(tx.TxId)
	if nil != err {
		log.Errorf("decode bitcoin tx id error: %s", err)
		return
	}
	packed = append(packed, util.ToVarint64(uint64(len(txId)))...)
	packed = append(packed, txId...)

	// number of Vout payments
	packed = append(packed, util.ToVarint64(uint64(len(amounts)))...)

	// individual payments
	for address, value := range amounts {
		packed = append(packed, util.ToVarint64(uint64(len(address)))...)
		packed = append(packed, address...)
		packed = append(packed, util.ToVarint64(value)...)
	}

	storage.Pool.Payment.Put(payId[:], packed)
}
