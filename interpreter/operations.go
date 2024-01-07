package interpreter

import (
	"errors"

	"github.com/nikgalushko/gan-ilox/token"
)

var ErrTypeMissmatch = errors.New("Type missmatch")

func add(left token.Literal, right token.Literal) (token.Literal, error) {
	if !((left.IsNumber() && right.IsNumber()) || (left.IsString() && right.IsString())) {
		return token.LiteralNil, ErrTypeMissmatch
	}

	if left.IsNumber() {
		if left.IsInt() && right.IsInt() {
			return token.NewLiteralInt(left.AsInt() + right.AsInt()), nil
		}
		return token.NewLiteralFloat(left.AsFloat() + right.AsFloat()), nil
	}

	return token.NewLiteralString(left.AsString() + right.AsString()), nil
}

func mul(left token.Literal, right token.Literal) (token.Literal, error) {
	if !left.IsNumber() && !right.IsNumber() {
		return token.LiteralNil, ErrTypeMissmatch
	}

	if left.IsInt() && right.IsInt() {
		return token.NewLiteralInt(left.AsInt() * right.AsInt()), nil
	}
	return token.NewLiteralFloat(left.AsFloat() * right.AsFloat()), nil
}

func div(left token.Literal, right token.Literal) (token.Literal, error) {
	if !left.IsNumber() && !right.IsNumber() {
		return token.LiteralNil, ErrTypeMissmatch
	}

	if left.IsInt() && right.IsInt() {
		return token.NewLiteralInt(left.AsInt() / right.AsInt()), nil
	}
	return token.NewLiteralFloat(left.AsFloat() / right.AsFloat()), nil
}

func sub(a, b token.Literal) (token.Literal, error) {
	if !a.IsNumber() && !b.IsNumber() {
		return token.LiteralNil, ErrTypeMissmatch
	}

	if a.IsInt() && b.IsInt() {
		return token.NewLiteralInt(a.AsInt() - b.AsInt()), nil
	}
	return token.NewLiteralFloat(a.AsFloat() - b.AsFloat()), nil
}

func less(left, right token.Literal) (token.Literal, error) {
	if !((left.IsNumber() && right.IsNumber()) || (left.IsString() && right.IsString())) {
		return token.LiteralNil, ErrTypeMissmatch
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

	return token.NewLiteralBool(ret), nil
}

func lessOrEqual(left, right token.Literal) (token.Literal, error) {
	if !((left.IsNumber() && right.IsNumber()) || (left.IsString() && right.IsString())) {
		return token.LiteralNil, ErrTypeMissmatch
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

	return token.NewLiteralBool(ret), nil
}
func graeater(left, right token.Literal) (token.Literal, error) {
	if !((left.IsNumber() && right.IsNumber()) || (left.IsString() && right.IsString())) {
		return token.LiteralNil, ErrTypeMissmatch
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

	return token.NewLiteralBool(ret), nil
}

func graeaterOrEqual(left, right token.Literal) (token.Literal, error) {
	if !((left.IsNumber() && right.IsNumber()) || (left.IsString() && right.IsString())) {
		return token.LiteralNil, ErrTypeMissmatch
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

	return token.NewLiteralBool(ret), nil
}
func equal(left, right token.Literal) (token.Literal, error) {
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

	return token.NewLiteralBool(ret), nil
}
