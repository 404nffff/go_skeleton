// +build !noasm !appengine
// Code generated by asm2asm, DO NOT EDIT.

package avx2

import (
	`github.com/bytedance/sonic/loader`
)

const (
    _entry__quote = 144
)

const (
    _stack__quote = 56
)

const (
    _size__quote = 2736
)

var (
    _pcsp__quote = [][2]uint32{
        {1, 0},
        {4, 8},
        {6, 16},
        {8, 24},
        {10, 32},
        {12, 40},
        {13, 48},
        {2687, 56},
        {2691, 48},
        {2692, 40},
        {2694, 32},
        {2696, 24},
        {2698, 16},
        {2700, 8},
        {2704, 0},
        {2736, 56},
    }
)

var _cfunc_quote = []loader.CFunc{
    {"_quote_entry", 0,  _entry__quote, 0, nil},
    {"_quote", _entry__quote, _size__quote, _stack__quote, _pcsp__quote},
}