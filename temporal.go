package temporal

import (
	"math"
	"math/cmplx"
	"math/rand"
	"errors"
	"fmt"
)

/*
 *******************************************************************************
 *                              Type Definitions                               *
 *******************************************************************************
*/


type Range struct {
	Min   float64
	Max   float64
}

type Temporal struct {
	T     float64
	C     float64
}

/*
 *******************************************************************************
 *                         Public Function Definitions                         *
 *******************************************************************************
*/

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

// Returns the hyperperiod for a slice of floating point periods (converted to integers)
func Integral_Hyperperiod (data []Temporal) (int64, error) {

	if len(data) == 0 {
		return 0.0, errors.New("Cannot compute LCM over empty slice")
	}
	if len(data) == 1 {
		return int64(data[0].T), nil
	}

	periods := []int64{}
	for _, d := range data {
		periods = append(periods, int64(d.T))
		fmt.Printf("period (%f): %d\n", d.T, int64(d.T))
	}
	return lcm(periods), nil
}

// Gives WCET and periods for given utilizations, period range, and step
func Make_Temporal_Data (r Range, step float64, us []float64) ([]Temporal, error) {
	ts := make([]Temporal, len(us))

	// Closure: Returns value, rounded to the nearest multiple of factor
	nearest_multiple := func (value, factor float64) float64 {
		return math.Ceil(value / factor) * factor
	}

	// Ensure the step is not larger than the min
	if step >= r.Max {
		return nil, errors.New(fmt.Sprintf("Step (%f) larger than maximum value (%f)",
			step, r.Min))
	}

	for i, u := range us {
		period := nearest_multiple((r.Min + (rand.Float64() * r.Max)), step)
		wcet := period * u
		ts[i] = Temporal{T: period, C: wcet}
	}
	return ts, nil
}


/*
 *******************************************************************************
 *                        Private Function Definitions                         *
 *******************************************************************************
*/


// Euclidean GCD
func gcd(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// LCM (expecting at least two elements)
func lcm(integers []int64) int64 {

	_lcm := func (a, b int64) int64 {
		return a * b / gcd(a,b)
	}

	r := _lcm(integers[0], integers[1])
	for i := 2; i < len(integers); i++ {
		r = _lcm(r, integers[i])
	}
	return r;
}