package summarize

import "math"

type WeightFunc int

const (
	// all weights are equal (unit value 1)
	Equal WeightFunc = iota
	// product of attribute weight and exponentially decaying relevance
	Exponential
	// product of attribute weight and linearly decaying relevance
	Linear
	// only attribute weight
	OnlyAttribute
)

type Assessor struct {
	weights   []float64
	function  WeightFunc
	NumTuples int
	weightEnd float64
}

// Weight computes the cover weight of a cell
func (a Assessor) Weight(attribute *Attribute, rank int) float64 {
	switch a.function {
	case Equal:
		return 1.0
	case Exponential:
		if a.NumTuples < 0 {
			panic("set numTuples")
		}
		// exponential function with f(0) = 1 and f(n) = weightEnd
		return a.weights[attribute.index] * math.Pow(a.weightEnd, float64(rank)/float64(a.NumTuples))
	case Linear:
		if a.NumTuples < 0 {
			panic("set numTuples")
		}
		// linear function with f(0) = 1 and f(n) = weightEnd
		return a.weights[attribute.index] * (1.0 - (a.weightEnd * float64(rank) / float64(a.NumTuples)))
	case OnlyAttribute:
		return a.weights[attribute.index]
	default:
		panic("invalid weighting function")
	}
}

// Weights returns the weight of a list of tuples for a single attribute
func (a Assessor) Weights(tuples TupleCover) float64 {
	// all weights are 1 so we can just return the length
	if a.function == Equal {
		return float64(len(tuples))
	}

	// sum up the weights
	sum := 0.0
	for _, cover := range tuples {
		sum += cover.weight
	}
	return sum
}

func MakeEqualWeightAssessor() Assessor {
	weights := make([]float64, 0)
	return Assessor{weights, Equal, -1, -1}
}

func MakeExponentialAssessor(weights []float64) Assessor {
	return Assessor{weights, Exponential, -1, 0.5}
}

func MakeOnlyAttributeAssessor(weights []float64) Assessor {
	return Assessor{weights, OnlyAttribute, -1, -1}
}
