package interpreter

import (
	"errors"

	"github.com/nikgalushko/gan-ilox/internal"
)

var ErrTypeMissmatch = errors.New("type missmatch")

func add(left internal.Literal, right internal.Literal) (internal.Literal, error) {
	if !((left.IsNumber() && right.IsNumber()) || (left.IsString() && right.IsString())) {
		return internal.LiteralNil, ErrTypeMissmatch
	}

	if left.IsNumber() {
		if left.IsInt() && right.IsInt() {
			return internal.NewLiteralInt(left.AsInt() + right.AsInt()), nil
		}
		return internal.NewLiteralFloat(left.AsFloat() + right.AsFloat()), nil
	}

	return internal.NewLiteralString(left.AsString() + right.AsString()), nil
}

func mul(left internal.Literal, right internal.Literal) (internal.Literal, error) {
	if !left.IsNumber() || !right.IsNumber() {
		return internal.LiteralNil, ErrTypeMissmatch
	}

	if left.IsInt() && right.IsInt() {
		return internal.NewLiteralInt(left.AsInt() * right.AsInt()), nil
	}
	return internal.NewLiteralFloat(left.AsFloat() * right.AsFloat()), nil
}

func div(left internal.Literal, right internal.Literal) (internal.Literal, error) {
	if !left.IsNumber() || !right.IsNumber() {
		return internal.LiteralNil, ErrTypeMissmatch
	}

	if left.IsInt() && right.IsInt() {
		return internal.NewLiteralInt(left.AsInt() / right.AsInt()), nil
	}
	return internal.NewLiteralFloat(left.AsFloat() / right.AsFloat()), nil
}

func sub(a, b internal.Literal) (internal.Literal, error) {
	if !a.IsNumber() || !b.IsNumber() {
		return internal.LiteralNil, ErrTypeMissmatch
	}

	if a.IsInt() && b.IsInt() {
		return internal.NewLiteralInt(a.AsInt() - b.AsInt()), nil
	}
	return internal.NewLiteralFloat(a.AsFloat() - b.AsFloat()), nil
}

func less(left, right internal.Literal) (internal.Literal, error) {
	if !((left.IsNumber() && right.IsNumber()) || (left.IsString() && right.IsString())) {
		return internal.LiteralNil, ErrTypeMissmatch
	}

	var ret bool
	if left.IsNumber() {
		if left.IsInt() && right.IsInt() {
			ret = left.AsInt() < right.AsInt()
		} else {
			ret = left.AsFloat() < right.AsFloat()
		}
	} else {
		ret = left.AsString() < right.AsString()
	}

	return internal.NewLiteralBool(ret), nil
}

func lessOrEqual(left, right internal.Literal) (internal.Literal, error) {
	if !((left.IsNumber() && right.IsNumber()) || (left.IsString() && right.IsString())) {
		return internal.LiteralNil, ErrTypeMissmatch
	}

	var ret bool
	if left.IsNumber() {
		if left.IsInt() && right.IsInt() {
			ret = left.AsInt() <= right.AsInt()
		} else {
			ret = left.AsFloat() <= right.AsFloat()
		}
	} else {
		ret = left.AsString() <= right.AsString()
	}

	return internal.NewLiteralBool(ret), nil
}
func graeater(left, right internal.Literal) (internal.Literal, error) {
	if !((left.IsNumber() && right.IsNumber()) || (left.IsString() && right.IsString())) {
		return internal.LiteralNil, ErrTypeMissmatch
	}

	var ret bool
	if left.IsNumber() {
		if left.IsInt() && right.IsInt() {
			ret = left.AsInt() > right.AsInt()
		} else {
			ret = left.AsFloat() > right.AsFloat()
		}
	} else {

		ret = left.AsString() > right.AsString()
	}

	return internal.NewLiteralBool(ret), nil
}

func graeaterOrEqual(left, right internal.Literal) (internal.Literal, error) {
	if !((left.IsNumber() && right.IsNumber()) || (left.IsString() && right.IsString())) {
		return internal.LiteralNil, ErrTypeMissmatch
	}

	var ret bool
	if left.IsNumber() {
		if left.IsInt() && right.IsInt() {
			ret = left.AsInt() >= right.AsInt()
		} else {
			ret = left.AsFloat() >= right.AsFloat()
		}
	} else {

		ret = left.AsString() >= right.AsString()
	}

	return internal.NewLiteralBool(ret), nil
}
func equal(left, right internal.Literal) (internal.Literal, error) {
	var ret bool

	if left.IsNumber() && right.IsNumber() {
		if left.IsInt() && right.IsInt() {
			ret = left.AsInt() == right.AsInt()
		} else {
			ret = left.AsFloat() == right.AsFloat()
		}
	} else if left.IsString() && right.IsString() {
		ret = left.AsString() == right.AsString()
	} else if left.IsBool() && right.IsBool() {
		ret = left.AsBool() == right.AsBool()
	}

	return internal.NewLiteralBool(ret), nil
}
