/*
Copyright 2023 eatmoreapple

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package expr

import (
	"errors"
	"fmt"
	"github.com/eatmoreapple/juice/internal/reflectlite"
	"go/token"
	"reflect"
)

var (
	// nilValue represents the nil value
	nilValue = reflect.ValueOf(nil)

	// invalidValue represents the invalid value
	invalidValue = reflect.Value{}
)

// BinaryExprExecutor is the interface for binary expression executor
type BinaryExprExecutor interface {
	// Exec execute the binary expression
	// right is the right value of the binary expression
	// next is the function to get the left value of the binary expression
	// return the result of the binary expression
	Exec(x, y func() (reflect.Value, error)) (reflect.Value, error)
}

// OperatorExecutor is the executor for operator
type OperatorExecutor struct {
	Operator
}

// Exec execute the binary expression
func (c OperatorExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	left, err := x()
	if err != nil {
		return invalidValue, err
	}
	right, err := y()
	if err != nil {
		return invalidValue, err
	}
	return c.Operator.Operate(left, right)
}

// EQLExprExecutor is the executor for ==
type EQLExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (EQLExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Eq}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// NEQExprExecutor is the executor for !=
type NEQExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (NEQExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Ne}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// LSSExprExecutor is the executor for <
type LSSExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (LSSExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Lt}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// LEQExprExecutor is the executor for <=
type LEQExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (LEQExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Le}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// GTRExprExecutor is the executor for >
type GTRExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (GTRExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Gt}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// GEQExprExecutor is the executor for >=
type GEQExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (GEQExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Ge}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// LANDExprExecutor is the executor for &&
type LANDExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (LANDExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	left, err := x()
	if err != nil {
		return invalidValue, err
	}
	left = reflectlite.Unwrap(left)
	if left.Kind() != reflect.Bool {
		return invalidValue, fmt.Errorf("expected bool, got %v", left.Kind())
	}
	if !left.Bool() {
		return left, nil
	}
	right, err := y()
	if err != nil {
		return invalidValue, err
	}
	right = reflectlite.Unwrap(right)
	if right.Kind() != reflect.Bool {
		return invalidValue, fmt.Errorf("expected bool, got %v", right.Kind())
	}
	return right, nil
}

// LORExprExecutor is the executor for ||
type LORExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (LORExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	left, err := x()
	if err != nil {
		return invalidValue, err
	}
	left = reflectlite.Unwrap(left)
	if left.Kind() != reflect.Bool {
		return invalidValue, fmt.Errorf("expected bool, got %v", left.Kind())
	}
	if left.Bool() {
		return left, nil
	}
	right, err := y()
	if err != nil {
		return invalidValue, err
	}
	right = reflectlite.Unwrap(right)
	if right.Kind() != reflect.Bool {
		return invalidValue, fmt.Errorf("expected bool, got %v", right.Kind())
	}
	return right, nil
}

// ADDExprExecutor is the executor for +
type ADDExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (ADDExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Add}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// SUBExprExecutor is the executor for -
type SUBExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (SUBExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Sub}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// MULExprExecutor is the executor for *
type MULExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (MULExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Mul}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// QUOExprExecutor is the executor for /
type QUOExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (QUOExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Quo}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// REMExprExecutor is the executor for %
type REMExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (REMExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	var operator = GenericOperator{OperatorExpr: Rem}
	executor := OperatorExecutor{Operator: operator}
	return executor.Exec(x, y)
}

// LPARENExprExecutor is the executor for (
type LPARENExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (LPARENExprExecutor) Exec(_, y func() (reflect.Value, error)) (reflect.Value, error) {
	return y()
}

// RPARENExprExecutor is the executor for )
// it just return the value
type RPARENExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (RPARENExprExecutor) Exec(x, _ func() (reflect.Value, error)) (reflect.Value, error) {
	return x()
}

type COMMENTExprExecutor struct{}

func (COMMENTExprExecutor) Exec(_, _ func() (reflect.Value, error)) (reflect.Value, error) {
	return reflect.ValueOf(true), nil
}

// NOTExprExecutor is the executor for !
type NOTExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (NOTExprExecutor) Exec(_, y func() (reflect.Value, error)) (reflect.Value, error) {
	right, err := y()
	if err != nil {
		return invalidValue, err
	}
	right = reflectlite.Unwrap(right)
	if right.Kind() != reflect.Bool {
		return invalidValue, fmt.Errorf("expected bool, got %v", right.Kind())
	}
	return reflect.ValueOf(!right.Bool()), nil
}

// ANDExprExecutor is the executor for &&
type ANDExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (ANDExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	executor := LANDExprExecutor{}
	return executor.Exec(x, y)
}

// ORExprExecutor is the executor for ||
type ORExprExecutor struct{}

// Exec execute the binary expression
// implement BinaryExprExecutor interface
func (ORExprExecutor) Exec(x, y func() (reflect.Value, error)) (reflect.Value, error) {
	executor := LORExprExecutor{}
	return executor.Exec(x, y)
}

// ErrUnsupportedBinaryExpr is the error that the binary expression is unsupported
var ErrUnsupportedBinaryExpr = errors.New("unsupported binary expression")

// binaryExprExecutors is a map from token to BinaryExprExecutor
var binaryExprExecutors = map[token.Token]BinaryExprExecutor{
	token.EQL:     EQLExprExecutor{},
	token.NEQ:     NEQExprExecutor{},
	token.LSS:     LSSExprExecutor{},
	token.LEQ:     LEQExprExecutor{},
	token.GTR:     GTRExprExecutor{},
	token.GEQ:     GEQExprExecutor{},
	token.LAND:    LANDExprExecutor{},
	token.LOR:     LORExprExecutor{},
	token.ADD:     ADDExprExecutor{},
	token.SUB:     SUBExprExecutor{},
	token.MUL:     MULExprExecutor{},
	token.QUO:     QUOExprExecutor{},
	token.REM:     REMExprExecutor{},
	token.LPAREN:  LPARENExprExecutor{},
	token.RPAREN:  RPARENExprExecutor{},
	token.COMMENT: COMMENTExprExecutor{},
	token.NOT:     NOTExprExecutor{},
	token.AND:     ANDExprExecutor{},
	token.OR:      ORExprExecutor{},
}

// FromToken returns the BinaryExprExecutor from the token
func FromToken(t token.Token) (BinaryExprExecutor, error) {
	executor, ok := binaryExprExecutors[t]
	if !ok {
		return nil, ErrUnsupportedBinaryExpr
	}
	return executor, nil
}
