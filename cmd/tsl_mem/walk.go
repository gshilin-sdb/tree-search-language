// Copyright 2019 Yaacov Zamir <kobi.zamir@gmail.com>
// and other contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: 2019 Nimrod Shneor <nimrodshn@gmail.com>

package main

import (
	"fmt"
	"regexp"

	"github.com/yaacov/tsl/pkg/tsl"
)

func handleIdent(n tsl.Node, book Book) (tsl.Node, error) {
	var err error

	l := n.Left.(tsl.Node)

	switch v := book[l.Left.(string)].(type) {
	case string:
		n.Left = tsl.Node{
			Func: tsl.StringOp,
			Left: v,
		}
	case nil:
		n.Left = tsl.Node{
			Func: tsl.NullOp,
			Left: nil,
		}
	case bool:
		val := "false"
		if v {
			val = "true"
		}
		n.Left = tsl.Node{
			Func: tsl.StringOp,
			Left: val,
		}
	case float32:
		n.Left = tsl.Node{
			Func: tsl.NumberOp,
			Left: float64(v),
		}
	case float64:
		n.Left = tsl.Node{
			Func: tsl.NumberOp,
			Left: v,
		}
	case int32:
		n.Left = tsl.Node{
			Func: tsl.NumberOp,
			Left: float64(v),
		}
	case int64:
		n.Left = tsl.Node{
			Func: tsl.NumberOp,
			Left: float64(v),
		}
	case uint32:
		n.Left = tsl.Node{
			Func: tsl.NumberOp,
			Left: float64(v),
		}
	case uint64:
		n.Left = tsl.Node{
			Func: tsl.NumberOp,
			Left: float64(v),
		}
	case int:
		n.Left = tsl.Node{
			Func: tsl.NumberOp,
			Left: float64(v),
		}
	case uint:
		n.Left = tsl.Node{
			Func: tsl.NumberOp,
			Left: float64(v),
		}
	default:
		err = tsl.UnexpectedLiteralError{Literal: fmt.Sprintf("%s[%v]", l.Left.(string), v)}
	}

	return n, err
}

func handleStringOp(n tsl.Node, book Book) (bool, error) {
	l := n.Left.(tsl.Node)
	r := n.Right.(tsl.Node)

	left := l.Left.(string)
	right := r.Left.(string)

	switch n.Func {
	case tsl.EqOp:
		return left == right, nil
	case tsl.NotEqOp:
		return left != right, nil
	case tsl.LtOp:
		return left < right, nil
	case tsl.LteOp:
		return left <= right, nil
	case tsl.GtOp:
		return left > right, nil
	case tsl.GteOp:
		return left >= right, nil
	case tsl.RegexOp:
		var valid = regexp.MustCompile(right)
		return valid.MatchString(left), nil
	case tsl.NotRegexOp:
		var valid = regexp.MustCompile(right)
		return !valid.MatchString(left), nil
	}

	return false, tsl.UnexpectedLiteralError{Literal: n.Func}
}

func handleNumberOp(n tsl.Node, book Book) (bool, error) {
	l := n.Left.(tsl.Node)
	r := n.Right.(tsl.Node)

	left := l.Left.(float64)
	right := r.Left.(float64)

	switch n.Func {
	case tsl.EqOp:
		return left == right, nil
	case tsl.NotEqOp:
		return left != right, nil
	case tsl.LtOp:
		return left < right, nil
	case tsl.LteOp:
		return left <= right, nil
	case tsl.GtOp:
		return left > right, nil
	case tsl.GteOp:
		return left >= right, nil
	}

	return false, tsl.UnexpectedLiteralError{Literal: n.Func}
}

func handleStringArrayOp(n tsl.Node, book Book) (bool, error) {
	l := n.Left.(tsl.Node)
	r := n.Right.(tsl.Node)

	left := l.Left.(string)
	right := r.Right.([]tsl.Node)

	switch n.Func {
	case tsl.BetweenOp:
		begin := right[0].Left.(string)
		end := right[1].Left.(string)
		return left >= begin && left < end, nil
	case tsl.NotBetweenOp:
		begin := right[0].Left.(string)
		end := right[1].Left.(string)
		return left < begin || left >= end, nil
	case tsl.InOp:
		b := false
		for _, node := range right {
			b = b || left == node.Left.(string)
		}
		return b, nil
	case tsl.NotInOp:
		b := true
		for _, node := range right {
			b = b && left != node.Left.(string)
		}
		return b, nil
	}

	return false, tsl.UnexpectedLiteralError{Literal: n.Func}
}

func handleNumberArrayOp(n tsl.Node, book Book) (bool, error) {
	l := n.Left.(tsl.Node)
	r := n.Right.(tsl.Node)

	left := l.Left.(float64)
	right := r.Right.([]tsl.Node)

	switch n.Func {
	case tsl.BetweenOp:
		begin := right[0].Left.(float64)
		end := right[1].Left.(float64)
		return left >= begin && left < end, nil
	case tsl.NotBetweenOp:
		begin := right[0].Left.(float64)
		end := right[1].Left.(float64)
		return left < begin || left >= end, nil
	case tsl.InOp:
		b := false
		for _, node := range right {
			b = b || left == node.Left.(float64)
		}
		return b, nil
	case tsl.NotInOp:
		b := true
		for _, node := range right {
			b = b && left != node.Left.(float64)
		}
		return b, nil
	}

	return false, tsl.UnexpectedLiteralError{Literal: n.Func}
}

func handleLogicalOp(n tsl.Node, book Book) (bool, error) {
	l := n.Left.(tsl.Node)
	r := n.Right.(tsl.Node)

	right, err := Walk(r, book)
	if err != nil {
		return false, err
	}
	left, err := Walk(l, book)
	if err != nil {
		return false, err
	}

	switch n.Func {
	case tsl.AndOp:
		return right && left, nil
	case tsl.OrOp:
		return right || left, nil
	}

	return false, tsl.UnexpectedLiteralError{Literal: n.Func}
}

// Walk implements sql semantics.
func Walk(n tsl.Node, book Book) (bool, error) {
	// Check for identifiers.
	l := n.Left.(tsl.Node)
	if l.Func == tsl.IdentOp {
		newNode, err := handleIdent(n, book)
		if err != nil {
			return false, err
		}
		return Walk(newNode, book)
	}

	// Walk tree.
	switch n.Func {
	case tsl.EqOp, tsl.NotEqOp, tsl.LtOp, tsl.LteOp, tsl.GtOp, tsl.GteOp, tsl.RegexOp, tsl.NotRegexOp,
		tsl.BetweenOp, tsl.NotBetweenOp, tsl.NotInOp, tsl.InOp:
		r := n.Right.(tsl.Node)

		switch l.Func {
		case tsl.StringOp:
			if r.Func == tsl.StringOp {
				return handleStringOp(n, book)
			}
			if r.Func == tsl.ArrayOp {
				return handleStringArrayOp(n, book)
			}
		case tsl.NumberOp:
			if r.Func == tsl.NumberOp {
				return handleNumberOp(n, book)
			}
			if r.Func == tsl.ArrayOp {
				return handleNumberArrayOp(n, book)
			}
		case tsl.NullOp:
			// Any comparison operation on a null element is false.
			return false, nil
		}
	case tsl.IsNotNilOp:
		return l.Func != tsl.NullOp, nil
	case tsl.IsNilOp:
		return l.Func == tsl.NullOp, nil
	case tsl.AndOp, tsl.OrOp:
		return handleLogicalOp(n, book)
	}

	return false, tsl.UnexpectedLiteralError{Literal: n.Func}
}
