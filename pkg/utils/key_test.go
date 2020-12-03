// Copyright 2020 PingCAP, Inc. Licensed under Apache-2.0.

package utils

import (
	"encoding/hex"
	. "github.com/pingcap/check"
)

type testKeySuite struct{}

var _ = Suite(&testKeySuite{})

func (r *testKeySuite) TestParseKey(c *C) {

	testRawKey := []struct {
		rawKey string
		ans []byte
	}{
		{"1234", []byte("1234")},
		{"abcd",[]byte("abcd")},
		{"1a2b", []byte("1a2b")},
		{"AA", []byte("AA")},
		{"\a",[]byte("\a")},
		{"\\'",[]byte("\\'")},
	}

	for _, tt := range testRawKey {
		parsedKey, err := ParseKey("raw", tt.rawKey)
		c.Assert(err, IsNil)
		c.Assert(parsedKey, BytesEquals, tt.ans)
	}

	testEscapedKey := []struct {
		EscapedKey string
		ans []byte
	}{
		{"\\a\\x1", []byte("\a\x01")},
		{"\\b\\f",[]byte("\b\f")},
		{"\\n\\r", []byte("\n\r")},
		{"\\t\\v",[]byte("\t\v")},
		{"\\'",[]byte("'")},
	}

	for _, tt := range testEscapedKey {
		parsedKey, err := ParseKey("escaped", tt.EscapedKey)
		c.Assert(err, IsNil)
		c.Assert(parsedKey, BytesEquals, tt.ans)
	}

	testHexKey := []struct {
		hexKey string
		ans []byte
	}{
		{"1234", []byte("1234")},
		{"abcd",[]byte("abcd")},
		{"1a2b", []byte("1a2b")},
		{"AA", []byte("AA")},
		{"\a",[]byte("\a")},
		{"\\'",[]byte("\\'")},
		{"\x01", []byte("\x01")},
		{"\xAA", []byte("\xAA")},
	}

	for _, tt := range testHexKey {
		key := hex.EncodeToString([]byte(tt.hexKey))
		parsedKey, err := ParseKey("hex", key)
		c.Assert(err, IsNil)
		c.Assert(parsedKey, BytesEquals, tt.ans)
	}

	testNotSupportKey := []struct {
		any string
		ans []byte
	}{
		{"1234", []byte("1234")},
		{"abcd",[]byte("abcd")},
		{"1a2b", []byte("1a2b")},
		{"AA", []byte("AA")},
		{"\a",[]byte("\a")},
		{"\\'",[]byte("\\'")},
		{"\x01", []byte("\x01")},
		{"\xAA", []byte("\xAA")},
	}

	for _, tt := range testNotSupportKey {
		_, err := ParseKey("notSupport", tt.any)
		c.Assert(err, ErrorMatches, "unknown format.*")
	}
}

func (r *testKeySuite) TestCompareEndKey(c *C) {
	res := CompareEndKey([]byte("1"), []byte("2"))
	c.Assert(res, Less, 0)

	res = CompareEndKey([]byte("1"), []byte("1"))
	c.Assert(res, Equals, 0)

	res = CompareEndKey([]byte("2"), []byte("1"))
	c.Assert(res, Greater, 0)

	res = CompareEndKey([]byte("1"), []byte(""))
	c.Assert(res, Less, 0)

	res = CompareEndKey([]byte(""), []byte(""))
	c.Assert(res, Equals, 0)

	res = CompareEndKey([]byte(""), []byte("1"))
	c.Assert(res, Greater, 0)
}
