// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package transactionrecord_test

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/bitmark-inc/bitmarkd/account"
	"github.com/bitmark-inc/bitmarkd/chain"
	"github.com/bitmark-inc/bitmarkd/currency"
	"github.com/bitmark-inc/bitmarkd/fault"
	"github.com/bitmark-inc/bitmarkd/mode"
	"github.com/bitmark-inc/bitmarkd/transactionrecord"
	"github.com/bitmark-inc/bitmarkd/util"
	"golang.org/x/crypto/ed25519"
	"reflect"
	"testing"
)

// to print a keypair for future tests
func TestGenerateKeypair(t *testing.T) {
	generate := false

	// generate = true // (uncomment to get a new key pair)

	if generate {
		// display key pair and fail the test
		// use the displayed values to modify data below
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if nil != err {
			t.Errorf("key pair generation error: %v", err)
			return
		}
		t.Errorf("*** GENERATED:\n%s", util.FormatBytes("publicKey", publicKey))
		t.Errorf("*** GENERATED:\n%s", util.FormatBytes("privateKey", privateKey))
		return
	}
}

// to hold a keypair for testing
type keyPair struct {
	publicKey  []byte
	privateKey []byte
}

// public/private keys from above generate

var proofedby = keyPair{
	publicKey: []byte{
		0x55, 0xb2, 0x98, 0x88, 0x17, 0xf7, 0xea, 0xec,
		0x37, 0x74, 0x1b, 0x82, 0x44, 0x71, 0x63, 0xca,
		0xaa, 0x5a, 0x9d, 0xb2, 0xb6, 0xf0, 0xce, 0x72,
		0x26, 0x26, 0x33, 0x8e, 0x5e, 0x3f, 0xd7, 0xf7,
	},
	privateKey: []byte{
		0x95, 0xb5, 0xa8, 0x0b, 0x4c, 0xdb, 0xe6, 0x1c,
		0x0f, 0x3f, 0x72, 0xcc, 0x15, 0x2d, 0x4a, 0x4f,
		0x29, 0xbc, 0xfd, 0x39, 0xc9, 0xa6, 0x7e, 0x2c,
		0x7b, 0xc6, 0xe0, 0xe1, 0x4e, 0xc7, 0xc7, 0xba,
		0x55, 0xb2, 0x98, 0x88, 0x17, 0xf7, 0xea, 0xec,
		0x37, 0x74, 0x1b, 0x82, 0x44, 0x71, 0x63, 0xca,
		0xaa, 0x5a, 0x9d, 0xb2, 0xb6, 0xf0, 0xce, 0x72,
		0x26, 0x26, 0x33, 0x8e, 0x5e, 0x3f, 0xd7, 0xf7,
	},
}

var registrant = keyPair{
	publicKey: []byte{
		0x7a, 0x81, 0x92, 0x56, 0x5e, 0x6c, 0xa2, 0x35,
		0x80, 0xe1, 0x81, 0x59, 0xef, 0x30, 0x73, 0xf6,
		0xe2, 0xfb, 0x8e, 0x7e, 0x9d, 0x31, 0x49, 0x7e,
		0x79, 0xd7, 0x73, 0x1b, 0xa3, 0x74, 0x11, 0x01,
	},
	privateKey: []byte{
		0x66, 0xf5, 0x28, 0xd0, 0x2a, 0x64, 0x97, 0x3a,
		0x2d, 0xa6, 0x5d, 0xb0, 0x53, 0xea, 0xd0, 0xfd,
		0x94, 0xca, 0x93, 0xeb, 0x9f, 0x74, 0x02, 0x3e,
		0xbe, 0xdb, 0x2e, 0x57, 0xb2, 0x79, 0xfd, 0xf3,
		0x7a, 0x81, 0x92, 0x56, 0x5e, 0x6c, 0xa2, 0x35,
		0x80, 0xe1, 0x81, 0x59, 0xef, 0x30, 0x73, 0xf6,
		0xe2, 0xfb, 0x8e, 0x7e, 0x9d, 0x31, 0x49, 0x7e,
		0x79, 0xd7, 0x73, 0x1b, 0xa3, 0x74, 0x11, 0x01,
	},
}

var issuer = keyPair{
	publicKey: []byte{
		0x9f, 0xc4, 0x86, 0xa2, 0x53, 0x4f, 0x17, 0xe3,
		0x67, 0x07, 0xfa, 0x4b, 0x95, 0x3e, 0x3b, 0x34,
		0x00, 0xe2, 0x72, 0x9f, 0x65, 0x61, 0x16, 0xdd,
		0x7b, 0x01, 0x8d, 0xf3, 0x46, 0x98, 0xbd, 0xc2,
	},
	privateKey: []byte{
		0xf3, 0xf7, 0xa1, 0xfc, 0x33, 0x10, 0x71, 0xc2,
		0xb1, 0xcb, 0xbe, 0x4f, 0x3a, 0xee, 0x23, 0x5a,
		0xae, 0xcc, 0xd8, 0x5d, 0x2a, 0x80, 0x4c, 0x44,
		0xb5, 0xc6, 0x03, 0xb4, 0xca, 0x4d, 0x9e, 0xc0,
		0x9f, 0xc4, 0x86, 0xa2, 0x53, 0x4f, 0x17, 0xe3,
		0x67, 0x07, 0xfa, 0x4b, 0x95, 0x3e, 0x3b, 0x34,
		0x00, 0xe2, 0x72, 0x9f, 0x65, 0x61, 0x16, 0xdd,
		0x7b, 0x01, 0x8d, 0xf3, 0x46, 0x98, 0xbd, 0xc2,
	},
}

var ownerOne = keyPair{
	publicKey: []byte{
		0x27, 0x64, 0x0e, 0x4a, 0xab, 0x92, 0xd8, 0x7b,
		0x4a, 0x6a, 0x2f, 0x30, 0xb8, 0x81, 0xf4, 0x49,
		0x29, 0xf8, 0x66, 0x04, 0x3a, 0x84, 0x1c, 0x38,
		0x14, 0xb1, 0x66, 0xb8, 0x89, 0x44, 0xb0, 0x92,
	},
	privateKey: []byte{
		0xc7, 0xae, 0x9f, 0x22, 0x32, 0x0e, 0xda, 0x65,
		0x02, 0x89, 0xf2, 0x64, 0x7b, 0xc3, 0xa4, 0x4f,
		0xfa, 0xe0, 0x55, 0x79, 0xcb, 0x6a, 0x42, 0x20,
		0x90, 0xb4, 0x59, 0xb3, 0x17, 0xed, 0xf4, 0xa1,
		0x27, 0x64, 0x0e, 0x4a, 0xab, 0x92, 0xd8, 0x7b,
		0x4a, 0x6a, 0x2f, 0x30, 0xb8, 0x81, 0xf4, 0x49,
		0x29, 0xf8, 0x66, 0x04, 0x3a, 0x84, 0x1c, 0x38,
		0x14, 0xb1, 0x66, 0xb8, 0x89, 0x44, 0xb0, 0x92,
	},
}

var ownerTwo = keyPair{
	publicKey: []byte{
		0xa1, 0x36, 0x32, 0xd5, 0x42, 0x5a, 0xed, 0x3a,
		0x6b, 0x62, 0xe2, 0xbb, 0x6d, 0xe4, 0xc9, 0x59,
		0x48, 0x41, 0xc1, 0x5b, 0x70, 0x15, 0x69, 0xec,
		0x99, 0x99, 0xdc, 0x20, 0x1c, 0x35, 0xf7, 0xb3,
	},
	privateKey: []byte{
		0x8f, 0x83, 0x3e, 0x58, 0x30, 0xde, 0x63, 0x77,
		0x89, 0x4a, 0x8d, 0xf2, 0xd4, 0x4b, 0x17, 0x88,
		0x39, 0x1d, 0xcd, 0xb8, 0xfa, 0x57, 0x22, 0x73,
		0xd6, 0x2e, 0x9f, 0xcb, 0x37, 0x20, 0x2a, 0xb9,
		0xa1, 0x36, 0x32, 0xd5, 0x42, 0x5a, 0xed, 0x3a,
		0x6b, 0x62, 0xe2, 0xbb, 0x6d, 0xe4, 0xc9, 0x59,
		0x48, 0x41, 0xc1, 0x5b, 0x70, 0x15, 0x69, 0xec,
		0x99, 0x99, 0xdc, 0x20, 0x1c, 0x35, 0xf7, 0xb3,
	},
}

// helper to make an address
func makeAccount(publicKey []byte) *account.Account {
	return &account.Account{
		AccountInterface: &account.ED25519Account{
			Test:      true,
			PublicKey: publicKey,
		},
	}
}

// test the packing/unpacking of base record
//
// ensures that pack->unpack returns the same original value
func TestPackBaseData(t *testing.T) {

	proofedbyAccount := makeAccount(proofedby.publicKey)

	r := transactionrecord.BaseData{
		Currency:       currency.Nothing,
		PaymentAddress: "nulladdress",
		Owner:          proofedbyAccount,
		Nonce:          0x12345678,
	}

	expected := []byte{
		0x01, 0x00, 0x0b, 0x6e, 0x75, 0x6c, 0x6c, 0x61,
		0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x21, 0x13,
		0x55, 0xb2, 0x98, 0x88, 0x17, 0xf7, 0xea, 0xec,
		0x37, 0x74, 0x1b, 0x82, 0x44, 0x71, 0x63, 0xca,
		0xaa, 0x5a, 0x9d, 0xb2, 0xb6, 0xf0, 0xce, 0x72,
		0x26, 0x26, 0x33, 0x8e, 0x5e, 0x3f, 0xd7, 0xf7,
		0xf8, 0xac, 0xd1, 0x91, 0x01,
	}

	expectedTxId := transactionrecord.Link{
		0x9e, 0xd1, 0x69, 0x58, 0x1f, 0xf3, 0x45, 0x02,
		0x46, 0xdc, 0xfe, 0x20, 0xf3, 0x76, 0xd8, 0x5d,
		0x56, 0xe3, 0x79, 0xc2, 0xe0, 0x97, 0xb9, 0x29,
		0xf5, 0x52, 0x4a, 0x3e, 0x6b, 0x18, 0xf4, 0x2c,
	}

	// manually sign the record and attach signature to "expected"
	signature := ed25519.Sign(proofedby.privateKey[:], expected)
	r.Signature = signature[:]
	//t.Logf("signature: %#v", r.Signature)
	l := util.ToVarint64(uint64(len(signature)))
	expected = append(expected, l...)
	expected = append(expected, signature[:]...)

	// test the packer
	packed, err := r.Pack(proofedbyAccount)
	if nil != err {
		t.Errorf("pack error: %v", err)
	}

	// if either of above fail we will have the message _without_ a signature
	if !bytes.Equal(packed, expected) {
		t.Errorf("pack record: %x  expected: %x", packed, expected)
		t.Errorf("*** GENERATED Packed:\n%s", util.FormatBytes("expected", packed))
		t.Fatal("fatal error")
	}

	// check the record type
	if transactionrecord.BaseDataTag != packed.Type() {
		t.Errorf("pack record type: %x  expected: %x", packed.Type(), transactionrecord.BaseDataTag)
	}

	t.Logf("Packed length: %d bytes", len(packed))

	// check txIds
	txId := packed.MakeLink()

	if txId != expectedTxId {
		t.Errorf("pack tx id: %#v  expected: %#v", txId, expectedTxId)
		t.Errorf("*** GENERATED tx id:\n%s", util.FormatBytes("expectedTxId", txId.Bytes()))
	}

	// =====
	// check test-network detection
	//
	// NOTE: this can only be done in the first record test since
	//       mode.Initialise may not be repeated
	if _, _, err := packed.Unpack(); err != fault.ErrWrongNetworkForPublicKey {
		t.Errorf("expected 'wrong network for public key' but got error: %v", err)
	}
	mode.Initialise(chain.Testing) // enter test mode - ONLY ALLOWED ONCE (or panic will occur
	// =====

	// test the unpacker
	unpacked, n, err := packed.Unpack()
	if nil != err {
		t.Fatalf("unpack error: %v", err)
	}

	if len(packed) != n {
		t.Errorf("did not unpack all data: only used: %d of: %d bytes", n, len(packed))
	}

	baseData, ok := unpacked.(*transactionrecord.BaseData)
	if !ok {
		t.Fatalf("did not unpack to BaseData")
	}

	// display a JSON version for information
	item := struct {
		HexTxId  string
		TxId     transactionrecord.Link
		BaseData *transactionrecord.BaseData
	}{
		HexTxId:  txId.String(),
		TxId:     txId,
		BaseData: baseData,
	}
	b, err := json.MarshalIndent(item, "", "  ")
	if nil != err {
		t.Fatalf("json error: %v", err)
	}

	t.Logf("BaseData: JSON: %s", b)

	// check that structure is preserved through Pack/Unpack
	// note reg is a pointer here
	if !reflect.DeepEqual(r, *baseData) {
		t.Errorf("different, original: %v  recovered: %v", r, *baseData)
	}
}

// test the packing/unpacking of registration record
//
// ensures that pack->unpack returns the same original value
func TestPackAssetData(t *testing.T) {

	registrantAccount := makeAccount(registrant.publicKey)

	r := transactionrecord.AssetData{
		Description: "Just the description",
		Name:        "Item's Name",
		Fingerprint: "0123456789abcdef",
		Registrant:  registrantAccount,
	}

	expected := []byte{
		0x02, 0x14, 0x4a, 0x75, 0x73, 0x74, 0x20, 0x74,
		0x68, 0x65, 0x20, 0x64, 0x65, 0x73, 0x63, 0x72,
		0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x0b, 0x49,
		0x74, 0x65, 0x6d, 0x27, 0x73, 0x20, 0x4e, 0x61,
		0x6d, 0x65, 0x10, 0x30, 0x31, 0x32, 0x33, 0x34,
		0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63,
		0x64, 0x65, 0x66, 0x21, 0x13, 0x7a, 0x81, 0x92,
		0x56, 0x5e, 0x6c, 0xa2, 0x35, 0x80, 0xe1, 0x81,
		0x59, 0xef, 0x30, 0x73, 0xf6, 0xe2, 0xfb, 0x8e,
		0x7e, 0x9d, 0x31, 0x49, 0x7e, 0x79, 0xd7, 0x73,
		0x1b, 0xa3, 0x74, 0x11, 0x01,
	}

	expectedTxId := transactionrecord.Link{
		0x1b, 0x01, 0x61, 0xd0, 0x0d, 0x3a, 0xfe, 0x51,
		0x6f, 0x74, 0x0c, 0x55, 0x1a, 0x72, 0x06, 0x23,
		0x6d, 0xcf, 0xc9, 0x08, 0x0c, 0x27, 0x36, 0x2d,
		0x27, 0x49, 0x6c, 0x42, 0x23, 0x0b, 0x7a, 0x2a,
	}

	expectedAssetIndex := transactionrecord.AssetIndex{
		0x59, 0xd0, 0x61, 0x55, 0xd2, 0x5d, 0xff, 0xdb,
		0x98, 0x27, 0x29, 0xde, 0x8d, 0xce, 0x9d, 0x78,
		0x55, 0xca, 0x09, 0x4d, 0x8b, 0xab, 0x81, 0x24,
		0xb3, 0x47, 0xc4, 0x06, 0x68, 0x47, 0x70, 0x56,
		0xb3, 0xc2, 0x7c, 0xcb, 0x7d, 0x71, 0xb5, 0x40,
		0x43, 0xd2, 0x07, 0xcc, 0xd1, 0x87, 0x64, 0x2b,
		0xf9, 0xc8, 0x46, 0x6f, 0x9a, 0x8d, 0x0d, 0xbe,
		0xfb, 0x4c, 0x41, 0x63, 0x3a, 0x7e, 0x39, 0xef,
	}

	// manually sign the record and attach signature to "expected"
	signature := ed25519.Sign(registrant.privateKey[:], expected)
	r.Signature = signature[:]
	//t.Logf("signature: %#v", r.Signature)
	l := util.ToVarint64(uint64(len(signature)))
	expected = append(expected, l...)
	expected = append(expected, signature[:]...)

	// test the packer
	packed, err := r.Pack(registrantAccount)
	if nil != err {
		t.Errorf("pack error: %v", err)
	}

	// if either of above fail we will have the message _without_ a signature
	if !bytes.Equal(packed, expected) {
		t.Errorf("pack record: %x  expected: %x", packed, expected)
		t.Errorf("*** GENERATED Packed:\n%s", util.FormatBytes("expected", packed))
		t.Fatal("fatal error")
	}

	// check the record type
	if transactionrecord.AssetDataTag != packed.Type() {
		t.Errorf("pack record type: %x  expected: %x", packed.Type(), transactionrecord.AssetDataTag)
	}

	t.Logf("Packed length: %d bytes", len(packed))

	// check txIds
	txId := packed.MakeLink()

	if txId != expectedTxId {
		t.Errorf("pack tx id: %#v  expected: %#v", txId, expectedTxId)
		t.Errorf("*** GENERATED tx id:\n%s", util.FormatBytes("expectedTxId", txId.Bytes()))
	}

	// check asset index
	assetIndex := r.AssetIndex()

	if assetIndex != expectedAssetIndex {
		t.Errorf("pack asset index: %#v  expected: %#v", assetIndex, expectedAssetIndex)
		t.Errorf("*** GENERATED asset index:\n%s", util.FormatBytes("expectedAssetIndex", assetIndex.Bytes()))
	}

	// test the unpacker
	unpacked, n, err := packed.Unpack()
	if nil != err {
		t.Fatalf("unpack error: %v", err)
	}
	if len(packed) != n {
		t.Errorf("did not unpack all data: only used: %d of: %d bytes", n, len(packed))
	}

	reg, ok := unpacked.(*transactionrecord.AssetData)
	if !ok {
		t.Fatalf("did not unpack to AssetData")
	}

	// display a JSON version for information
	item := struct {
		HexTxId   string
		TxId      transactionrecord.Link
		HexAsset  string
		Asset     transactionrecord.AssetIndex
		AssetData *transactionrecord.AssetData
	}{
		HexTxId:   txId.String(),
		TxId:      txId,
		HexAsset:  assetIndex.String(),
		Asset:     assetIndex,
		AssetData: reg,
	}
	b, err := json.MarshalIndent(item, "", "  ")
	if nil != err {
		t.Fatalf("json error: %v", err)
	}

	t.Logf("AssetData: JSON: %s", b)

	// check that structure is preserved through Pack/Unpack
	// note reg is a pointer here
	if !reflect.DeepEqual(r, *reg) {
		t.Fatalf("different, original: %v  recovered: %v", r, *reg)
	}
}

// test the packing/unpacking of Bitmark issue record
//
// ensures that pack->unpack returns the same original value
func TestPackBitmarkIssue(t *testing.T) {

	issuerAccount := makeAccount(issuer.publicKey)

	var asset transactionrecord.AssetIndex
	_, err := fmt.Sscan("BMA159d06155d25dffdb982729de8dce9d7855ca094d8bab8124b347c40668477056b3c27ccb7d71b54043d207ccd187642bf9c8466f9a8d0dbefb4c41633a7e39ef", &asset)
	if nil != err {
		t.Fatalf("hex to link error: %v", err)
	}

	r := transactionrecord.BitmarkIssue{
		AssetIndex: asset,
		Owner:      issuerAccount,
		Nonce:      99,
	}

	expected := []byte{
		0x03, 0x40, 0xef, 0x39, 0x7e, 0x3a, 0x63, 0x41,
		0x4c, 0xfb, 0xbe, 0x0d, 0x8d, 0x9a, 0x6f, 0x46,
		0xc8, 0xf9, 0x2b, 0x64, 0x87, 0xd1, 0xcc, 0x07,
		0xd2, 0x43, 0x40, 0xb5, 0x71, 0x7d, 0xcb, 0x7c,
		0xc2, 0xb3, 0x56, 0x70, 0x47, 0x68, 0x06, 0xc4,
		0x47, 0xb3, 0x24, 0x81, 0xab, 0x8b, 0x4d, 0x09,
		0xca, 0x55, 0x78, 0x9d, 0xce, 0x8d, 0xde, 0x29,
		0x27, 0x98, 0xdb, 0xff, 0x5d, 0xd2, 0x55, 0x61,
		0xd0, 0x59, 0x21, 0x13, 0x9f, 0xc4, 0x86, 0xa2,
		0x53, 0x4f, 0x17, 0xe3, 0x67, 0x07, 0xfa, 0x4b,
		0x95, 0x3e, 0x3b, 0x34, 0x00, 0xe2, 0x72, 0x9f,
		0x65, 0x61, 0x16, 0xdd, 0x7b, 0x01, 0x8d, 0xf3,
		0x46, 0x98, 0xbd, 0xc2, 0x63,
	}

	expectedTxId := transactionrecord.Link{
		0xbb, 0x82, 0x7a, 0xf2, 0x01, 0xdf, 0x8d, 0xfd,
		0x14, 0x76, 0xfb, 0x23, 0x50, 0xef, 0xec, 0x35,
		0x3e, 0x92, 0xf0, 0x9c, 0xc3, 0xe2, 0xd1, 0x6c,
		0x3e, 0x3d, 0x9f, 0x15, 0x9c, 0x90, 0xac, 0x25,
	}

	// manually sign the record and attach signature to "expected"
	signature := ed25519.Sign(issuer.privateKey[:], expected)
	r.Signature = signature[:]
	l := util.ToVarint64(uint64(len(signature)))
	expected = append(expected, l...)
	expected = append(expected, signature[:]...)

	// test the packer
	packed, err := r.Pack(issuerAccount)
	if nil != err {
		t.Errorf("pack error: %v", err)
	}

	// if either of above fail we will have the message _without_ a signature
	if !bytes.Equal(packed, expected) {
		t.Errorf("pack record: %x  expected: %x", packed, expected)
		t.Errorf("*** GENERATED Packed:\n%s", util.FormatBytes("expected", packed))
		t.Fatal("fatal error")
	}

	t.Logf("Packed length: %d bytes", len(packed))

	// check txId
	txId := packed.MakeLink()

	if txId != expectedTxId {
		t.Errorf("pack tx id: %#v  expected: %x", txId, expectedTxId)
		t.Errorf("*** GENERATED tx id:\n%s", util.FormatBytes("expectedTxId", txId.Bytes()))
		t.Fatal("fatal error")
	}

	// test the unpacker
	unpacked, n, err := packed.Unpack()
	if nil != err {
		t.Fatalf("unpack error: %v", err)
	}
	if len(packed) != n {
		t.Errorf("did not unpack all data: only used: %d of: %d bytes", n, len(packed))
	}

	bmt, ok := unpacked.(*transactionrecord.BitmarkIssue)
	if !ok {
		t.Fatalf("did not unpack to BitmarkIssue")
	}

	// display a JSON version for information
	item := struct {
		HexTxId      string
		TxId         transactionrecord.Link
		BitmarkIssue *transactionrecord.BitmarkIssue
	}{
		txId.String(),
		txId,
		bmt,
	}
	b, err := json.MarshalIndent(item, "", "  ")
	if nil != err {
		t.Fatalf("json error: %v", err)
	}

	t.Logf("Bitmark Issue: JSON: %s", b)

	// check that structure is preserved through Pack/Unpack
	// note reg is a pointer here
	if !reflect.DeepEqual(r, *bmt) {
		t.Fatalf("different, original: %v  recovered: %v", r, *bmt)
	}
}

// test the packing/unpacking of Bitmark transfer record
//
// transfer from issue
// ensures that pack->unpack returns the same original value
func TestPackBitmarkTransferOne(t *testing.T) {

	issuerAccount := makeAccount(issuer.publicKey)
	ownerOneAccount := makeAccount(ownerOne.publicKey)

	var link transactionrecord.Link
	_, err := fmt.Sscan("BMK1bb827af201df8dfd1476fb2350efec353e92f09cc3e2d16c3e3d9f159c90ac25", &link)
	if nil != err {
		t.Fatalf("hex to link error: %v", err)
	}

	r := transactionrecord.BitmarkTransfer{
		Link:  link,
		Owner: ownerOneAccount,
	}

	expected := []byte{
		0x04, 0x20, 0x25, 0xac, 0x90, 0x9c, 0x15, 0x9f,
		0x3d, 0x3e, 0x6c, 0xd1, 0xe2, 0xc3, 0x9c, 0xf0,
		0x92, 0x3e, 0x35, 0xec, 0xef, 0x50, 0x23, 0xfb,
		0x76, 0x14, 0xfd, 0x8d, 0xdf, 0x01, 0xf2, 0x7a,
		0x82, 0xbb, 0x00, 0x21, 0x13, 0x27, 0x64, 0x0e,
		0x4a, 0xab, 0x92, 0xd8, 0x7b, 0x4a, 0x6a, 0x2f,
		0x30, 0xb8, 0x81, 0xf4, 0x49, 0x29, 0xf8, 0x66,
		0x04, 0x3a, 0x84, 0x1c, 0x38, 0x14, 0xb1, 0x66,
		0xb8, 0x89, 0x44, 0xb0, 0x92,
	}

	expectedTxId := transactionrecord.Link{
		0x1c, 0xcf, 0x4b, 0x31, 0xd1, 0xe0, 0xb6, 0x1b,
		0x6b, 0x64, 0x93, 0xd2, 0xc1, 0x8c, 0xe5, 0x3a,
		0x75, 0x8e, 0x5f, 0xc3, 0x65, 0x70, 0x97, 0xb1,
		0x77, 0x35, 0x9e, 0x52, 0xed, 0x4c, 0xa3, 0x49,
	}

	// manually sign the record and attach signature to "expected"
	signature := ed25519.Sign(issuer.privateKey[:], expected)
	r.Signature = signature[:]
	l := util.ToVarint64(uint64(len(signature)))
	expected = append(expected, l...)
	expected = append(expected, signature[:]...)

	// test the packer
	packed, err := r.Pack(issuerAccount)
	if nil != err {
		t.Errorf("pack error: %v", err)
	}

	// if either of above fail we will have the message _without_ a signature
	if !bytes.Equal(packed, expected) {
		t.Errorf("pack record: %x  expected: %x", packed, expected)
		t.Errorf("*** GENERATED Packed:\n%s", util.FormatBytes("expected", packed))
		t.Fatal("fatal error")
	}

	t.Logf("Packed length: %d bytes", len(packed))

	// check txId
	txId := packed.MakeLink()

	if txId != expectedTxId {
		t.Errorf("pack txId: %#v  expected: %x", txId, expectedTxId)
		t.Errorf("*** GENERATED txId:\n%s", util.FormatBytes("expectedTxId", txId.Bytes()))
		t.Fatal("fatal error")
	}

	// test the unpacker
	unpacked, n, err := packed.Unpack()
	if nil != err {
		t.Fatalf("unpack error: %v", err)
	}
	if len(packed) != n {
		t.Errorf("did not unpack all data: only used: %d of: %d bytes", n, len(packed))
	}

	bmt, ok := unpacked.(*transactionrecord.BitmarkTransfer)
	if !ok {
		t.Fatalf("did not unpack to BitmarkTransfer")
	}

	// display a JSON version for information
	item := struct {
		HexTxId         string
		TxId            transactionrecord.Link
		BitmarkTransfer *transactionrecord.BitmarkTransfer
	}{
		txId.String(),
		txId,
		bmt,
	}
	b, err := json.MarshalIndent(item, "", "  ")
	if nil != err {
		t.Fatalf("json error: %v", err)
	}

	t.Logf("Bitmark Transfer: JSON: %s", b)

	// check that structure is preserved through Pack/Unpack
	// note reg is a pointer here
	if !reflect.DeepEqual(r, *bmt) {
		t.Fatalf("different, original: %v  recovered: %v", r, *bmt)
	}
}

// test the packing/unpacking of Bitmark transfer record
//
// test transfer to transfer
// ensures that pack->unpack returns the same original value
func TestPackBitmarkTransferTwo(t *testing.T) {

	ownerOneAccount := makeAccount(ownerOne.publicKey)
	ownerTwoAccount := makeAccount(ownerTwo.publicKey)

	var link transactionrecord.Link
	_, err := fmt.Sscan("BMK1f61f5cdb0757cdee36c0ae9514f6b87d6306475d578efbc191980a63323b6ab6", &link)
	if nil != err {
		t.Fatalf("hex to link error: %v", err)
	}

	r := transactionrecord.BitmarkTransfer{
		Link: link,
		Payment: &transactionrecord.Payment{
			Currency: currency.Bitcoin,
			Address:  "some-payment-address",
			Amount:   5000,
		},
		Owner: ownerTwoAccount,
	}

	expected := []byte{
		0x04, 0x20, 0xb6, 0x6a, 0x3b, 0x32, 0x63, 0x0a,
		0x98, 0x91, 0xc1, 0xfb, 0x8e, 0x57, 0x5d, 0x47,
		0x06, 0x63, 0x7d, 0xb8, 0xf6, 0x14, 0x95, 0xae,
		0xc0, 0x36, 0xee, 0xcd, 0x57, 0x07, 0xdb, 0x5c,
		0x1f, 0xf6, 0x01, 0x01, 0x14, 0x73, 0x6f, 0x6d,
		0x65, 0x2d, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
		0x74, 0x2d, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
		0x73, 0x88, 0x27, 0x21, 0x13, 0xa1, 0x36, 0x32,
		0xd5, 0x42, 0x5a, 0xed, 0x3a, 0x6b, 0x62, 0xe2,
		0xbb, 0x6d, 0xe4, 0xc9, 0x59, 0x48, 0x41, 0xc1,
		0x5b, 0x70, 0x15, 0x69, 0xec, 0x99, 0x99, 0xdc,
		0x20, 0x1c, 0x35, 0xf7, 0xb3,
	}

	expectedTxId := transactionrecord.Link{
		0xf4, 0x1e, 0xe0, 0xc7, 0xd4, 0x17, 0x99, 0xbd,
		0x90, 0x47, 0x7e, 0x66, 0xce, 0x4c, 0xc4, 0xf8,
		0xa7, 0x66, 0xb5, 0x13, 0xd6, 0xd2, 0x93, 0x07,
		0x9c, 0x47, 0x32, 0xe5, 0x58, 0x8f, 0x95, 0xec,
	}

	// manually sign the record and attach signature to "expected"
	signature := ed25519.Sign(ownerOne.privateKey[:], expected)
	r.Signature = signature[:]
	l := util.ToVarint64(uint64(len(signature)))
	expected = append(expected, l...)
	expected = append(expected, signature[:]...)

	// test the packer
	packed, err := r.Pack(ownerOneAccount)
	if nil != err {
		t.Errorf("pack error: %v", err)
	}

	// if either of above fail we will have the message _without_ a signature
	if !bytes.Equal(packed, expected) {
		t.Errorf("pack record: %x  expected: %x", packed, expected)
		t.Errorf("*** GENERATED Packed:\n%s", util.FormatBytes("expected", packed))
		t.Fatal("fatal error")
	}

	t.Logf("Packed length: %d bytes", len(packed))

	// check txId
	txId := packed.MakeLink()

	if txId != expectedTxId {
		t.Errorf("pack txId: %#v  expected: %x", txId, expectedTxId)
		t.Errorf("*** GENERATED txId:\n%s", util.FormatBytes("expectedTxId", txId.Bytes()))
		t.Fatal("fatal error")
	}

	// test the unpacker
	unpacked, n, err := packed.Unpack()
	if nil != err {
		t.Fatalf("unpack error: %v", err)
	}
	if len(packed) != n {
		t.Errorf("did not unpack all data: only used: %d of: %d bytes", n, len(packed))
	}

	bmt, ok := unpacked.(*transactionrecord.BitmarkTransfer)
	if !ok {
		t.Fatalf("did not unpack to BitmarkTransfer")
	}

	// display a JSON version for information
	item := struct {
		HexTxId         string
		TxId            transactionrecord.Link
		BitmarkTransfer *transactionrecord.BitmarkTransfer
	}{
		txId.String(),
		txId,
		bmt,
	}
	b, err := json.MarshalIndent(item, "", "  ")
	if nil != err {
		t.Fatalf("json error: %v", err)
	}

	t.Logf("Bitmark Transfer: JSON: %s", b)

	// check that structure is preserved through Pack/Unpack
	// note reg is a pointer here
	if !reflect.DeepEqual(r, *bmt) {
		t.Fatalf("different, original: %v  recovered: %v", r, *bmt)
	}
}
