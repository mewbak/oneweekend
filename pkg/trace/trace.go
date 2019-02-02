package trace

import (
	"math/rand"

	"github.com/hunterloftis/oneweekend/pkg/geom"
)

const bias = 0.001

type Rnd func() float64

// Surface represents something that can be Hit by a Ray.
type Surface interface {
	Hit(r Ray, dMin, dMax float64, rnd *rand.Rand) *Hit
	Bounds(t0, t1 float64) (bounds *AABB)
}

// Material represents a material that scatters light.
type Material interface {
	Scatter(in, norm geom.Unit, uv, p geom.Vec, rnd *rand.Rand) (out geom.Unit, attenuate Color, ok bool)
	Emit(uv, p geom.Vec) Color
}

// Hit records the details of a Ray->Surface intersection.
type Hit struct {
	Dist float64
	Norm geom.Unit
	UV   geom.Vec
	Pt   geom.Vec
	Mat  Material
}

// AABB Represents an axis-aligned bounding box
type AABB struct {
	min geom.Vec
	max geom.Vec
}

// NewAABB creates a new AABB with the given min and max positions
func NewAABB(min, max geom.Vec) *AABB {
	return &AABB{min: min, max: max} // TODO: use .min() and .max() here to remove the burden of the caller knowing which is larger
}

// Hit returns whether or not a given ray hits the AABB between dMin and dMax distances
func (a *AABB) Hit(r Ray, dMin, dMax float64) bool {
	for i := 0; i < 3; i++ {
		invD := 1 / r.Dir.E[i]
		d0 := (a.min.E[i] - r.Or.E[i]) * invD
		d1 := (a.max.E[i] - r.Or.E[i]) * invD
		if invD < 0 {
			d0, d1 = d1, d0
		}
		if d0 > dMin {
			dMin = d0
		}
		if d1 < dMax {
			dMax = d1
		}
		if dMax <= dMin {
			return false
		}
	}
	return true
}

// Plus returns a new AABB that encompasses the space of both a and b.
// If b is nil, the new AABB will be equivalent to a.
func (a *AABB) Plus(b *AABB) *AABB {
	if b == nil {
		return NewAABB(a.min, a.max)
	}
	return NewAABB(a.min.Min(b.min), a.max.Max(b.max))
}

// Corners returns a slice of the eight corners of this bounding box.
func (a *AABB) Corners() []geom.Vec {
	c := make([]geom.Vec, 0, 8)
	for i := 0.0; i < 2; i++ {
		for j := 0.0; j < 2; j++ {
			for k := 0.0; k < 2; k++ {
				x := i*a.min.X() + (1-i)*a.max.X()
				y := j*a.min.Y() + (1-j)*a.max.Y()
				z := k*a.min.Z() + (1-k)*a.max.Z()
				c = append(c, geom.NewVec(x, y, z))
			}
		}
	}
	return c
}

// Extended returns a new bounding box that encloses Vector v.
func (a *AABB) Extended(v geom.Vec) *AABB {
	return NewAABB(a.min.Min(v), a.max.Max(v))
}

// Min returns the minimum AABB vector.
func (a *AABB) Min() geom.Vec {
	return a.min
}

// Max returns the maximum AABB vector.
func (a *AABB) Max() geom.Vec {
	return a.max
}

func (a *AABB) SurfaceArea() float64 {
	dims := a.max.Minus(a.min)
	front := dims.X() * dims.Y()
	side := dims.Z() * dims.Y()
	top := dims.X() * dims.Z()
	return (front + side + top) * 2
}

func (a *AABB) Mid() geom.Vec {
	return a.min.Plus(a.max).Scaled(0.5)
}
