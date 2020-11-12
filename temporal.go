package temporal

import (
	"math/cmplx"
	"math/rand"
)

type Range struct {
	Min   float64
	Max   float64
}

type Temporal struct {
	T     float64
	C     float64
}

// Splits utilization U into N fragments, which all sum to U
func Uunifast (u_total float64, n int) []float64 {
	u := complex(u_total, 0.0)

	// Utilization components
	us := make([]complex128, n)

	// Compute r1 ... r_(n-1) ~ U(0,1)
	rs := []complex128{}
	for i := 0; i < n; i++ {
		rs = append(rs, complex(rand.Float64(), 0))
	}

	// s_n = u
	s := u

	// For i = n, ... , 2
	for i := n; i >= 2; i-- {
		exp := complex(1.0, 0.0) / complex((float64(i) - 1.0), 0.0)
		ss := s * cmplx.Pow(rs[i-1], exp)
		us[i-1] = s - ss
		s = ss
	}

	// Complete range (last element)
	sum := complex128(0.0)
	for _, v := range us {
		sum += v
	}
	us[0] = u - sum

	// Return array of float64s
	components := make([]float64, n)
	for i := 0; i < n; i++ {
		components[i] = real(us[i])
	}
	return components
}

// Gives WCET and periods for given utilizations and period range
func Make_Temporal_Data (r Range, us []float64) []Temporal {
	ts := make([]Temporal, len(us))
	for i, u := range us {
		period := (r.Min + (rand.Float64() * r.Max))
		wcet := period * u
		ts[i] = Temporal{T: period, C: wcet}
	}
	return ts
}