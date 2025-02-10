package Math

import "math"

type Mat3 struct {
	Matrix [9]float64
}

// Constructors

func NewMat3(a, b, c, d, e, f, g, h, k float64) Mat3 {
	return Mat3{[9]float64{a, b, c, d, e, f, g, h, k}}
}

func Mat3Identity() Mat3 {
	return Mat3{[9]float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	},
	}
}

func Mat3Scale(s float64) Mat3 {
	return Mat3{[9]float64{
		s, 0, 0,
		0, s, 0,
		0, 0, s,
	},
	}
}

func Mat3SeparateScale(sx, sy, sz float64) Mat3 {
	return Mat3{[9]float64{
		sx, 0, 0,
		0, sy, 0,
		0, 0, sz,
	},
	}
}

// rotation matrices

func Mat3XRotation(a float64) Mat3 {
	s := math.Sin(a)
	c := math.Cos(a)
	return Mat3{[9]float64{
		1, 0, 0,
		0, c, -s,
		0, s, c,
	},
	}
}

func Mat3YRotation(a float64) Mat3 {
	s := math.Sin(a)
	c := math.Cos(a)
	return Mat3{[9]float64{
		c, 0, -s,
		0, 1, 0,
		s, 0, c,
	},
	}
}

func Mat3ZRotation(a float64) Mat3 {
	s := math.Sin(a)
	c := math.Cos(a)
	return Mat3{[9]float64{
		c, -s, 0,
		s, c, 0,
		0, 0, 1,
	},
	}
}

func Mat3Euler(rX, rY, rZ float64) Mat3 {
	cp := math.Cos(rX)
	sp := math.Sin(rX)
	cr := math.Cos(rY)
	sr := math.Sin(rY)
	cy := math.Cos(rZ)
	sy := math.Sin(rZ)
	mat := Mat3{Matrix: [9]float64{
		cr*cy - sr*sp*sy, -cr*sy - sr*sp*cy, -sr * cp,
		cp * sy, cp * cy, -sp,
		sr*cy + cr*sp*sy, -sr*sy + cr*sp*cy, cr * cp,
	}}
	return mat
}

// Matrix multiplication

func (m Mat3) MatMul(o Mat3) Mat3 {
	// That will be a big one
	return Mat3{
		[9]float64{
			// First row
			m.Matrix[0]*o.Matrix[0] + m.Matrix[1]*o.Matrix[3] + m.Matrix[2]*o.Matrix[6],
			m.Matrix[0]*o.Matrix[1] + m.Matrix[1]*o.Matrix[4] + m.Matrix[2]*o.Matrix[7],
			m.Matrix[0]*o.Matrix[2] + m.Matrix[1]*o.Matrix[5] + m.Matrix[2]*o.Matrix[8],
			// Second row
			m.Matrix[3]*o.Matrix[0] + m.Matrix[4]*o.Matrix[3] + m.Matrix[5]*o.Matrix[6],
			m.Matrix[3]*o.Matrix[1] + m.Matrix[4]*o.Matrix[4] + m.Matrix[5]*o.Matrix[7],
			m.Matrix[3]*o.Matrix[2] + m.Matrix[4]*o.Matrix[5] + m.Matrix[5]*o.Matrix[8],
			// Third row
			m.Matrix[6]*o.Matrix[0] + m.Matrix[7]*o.Matrix[3] + m.Matrix[8]*o.Matrix[6],
			m.Matrix[6]*o.Matrix[1] + m.Matrix[7]*o.Matrix[4] + m.Matrix[8]*o.Matrix[7],
			m.Matrix[6]*o.Matrix[2] + m.Matrix[7]*o.Matrix[5] + m.Matrix[8]*o.Matrix[8],
		},
	}
}

func (m Mat3) VecMul(o Vector3) Vector3 {
	return Vector3{
		X: o.X*m.Matrix[0] + o.Y*m.Matrix[1] + o.Z*m.Matrix[2],
		Y: o.X*m.Matrix[3] + o.Y*m.Matrix[4] + o.Z*m.Matrix[5],
		Z: o.X*m.Matrix[6] + o.Y*m.Matrix[7] + o.Z*m.Matrix[8],
	}
}
