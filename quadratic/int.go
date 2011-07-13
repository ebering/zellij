package quadratic

import "math"

const (
	FLOAT64_EPSILON = 2e-8
)

type Integer struct {
	a, b int64
}

func NewInteger(a, b int64) *Integer {
	return &Integer{a, b}
}

const BASE = 2

var Zero = &Integer{0, 0}

func (i *Integer) Add(j *Integer) *Integer {
	return &Integer{i.a + j.a, i.b + j.b}
}

func (i *Integer) Negation() *Integer {
	return &Integer{ -i.a,-i.b}
}

func (i *Integer) Sub(j *Integer) *Integer {
	return &Integer{i.a - j.a, i.b - j.b}
}

func (i *Integer) Mul(j *Integer) *Integer {
	return &Integer{i.a*j.a + 2*i.b*j.b, i.a*j.b + i.b*j.a}
}

func (i *Integer) MultR2On2() *Integer {
	if i.a % 2 != 0 {
		panic("cannot multiply a qint by root 2/2 if it isn't even")
	}
	return &Integer{i.b, i.a / 2}
}

func (i *Integer) Div(j *Integer) (*Integer, bool) {
	fn := j.FieldNorm()
	m := i.Mul(j.Conjugate())

	if m.a%fn == 0 && m.b%fn == 0 {
		return &Integer{m.a / fn, m.b / fn}, true
	}
	return nil, false
}

func (i *Integer) Conjugate() *Integer {
	return &Integer{i.a, -i.b}
}

func (i *Integer) Inv() (*Integer, bool) {
	if i.a%i.FieldNorm() == 0 && i.b%i.FieldNorm() == 0 {
		return &Integer{i.a / i.FieldNorm(), i.b / i.FieldNorm()}, true
	}
	return nil, false
}

func (i *Integer) Equal(j *Integer) bool {
	return i.a == j.a && i.b == j.b
}

// TODO: This is currently prone to integer overflow
func (i *Integer) Less(j *Integer) bool {
	if i.a-j.a < 0 && j.b-i.b > 0 {
		return true
	} else if i.a-j.a > 0 && j.b-i.b < 0 {
		return false
	} else if i.a-j.a >= 0 && j.b-i.b >= 0 {
		return (i.a-j.a)*(i.a-j.a) < (j.b-i.b)*(j.b-i.b)*2
	} else {
		return (i.a-j.a)*(i.a-j.a) > (j.b-i.b)*(j.b-i.b)*2
	}
	return false
}

func (i *Integer) Float64() float64 {
	return float64(i.a) + float64(i.b)*math.Sqrt(BASE)
}

func (i *Integer) Copy() *Integer {
	return &Integer{i.a, i.b}
}

func (i *Integer) FieldNorm() int64 {
	return i.a*i.a - 2*i.b*i.b
}
