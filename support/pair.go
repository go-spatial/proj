package support

import (
	"log"
	"strconv"
	"strings"
)

// Pair is a simple key-value pair
// Pairs use copy semantics (pass-by-value).
type Pair struct {
	Key   string
	Value string
}

// PairList is an array of Pair objects.
// (We can't use a map because order of the items is important and
// because we might have duplicate keys.)
type PairList struct {
	pairs []Pair
}

// NewPairList returns a new PairList
func NewPairList() *PairList {
	return &PairList{
		pairs: []Pair{},
	}
}

// NewPairListFromString returns a new PairList from a string
// of the form "+proj=utm +zone=11 +datum=WGS84",
// with the leading "+" is optional and ignoring extra whitespace
func NewPairListFromString(source string) (*PairList, error) {

	ret := NewPairList()

	words := strings.Fields(source)
	for _, w := range words {

		var pair Pair

		if w[0:1] == "+" {
			w = w[1:]
		}

		v := strings.Split(w, "=")

		if v[0] == "" {
			return nil, BadProjStringError
		}

		switch len(v) {
		case 0:
			pair.Key = w
			pair.Value = ""
		case 1:
			// "proj=" is okay
			pair.Key = v[0]
			pair.Value = v[1]
		case 2:
			pair.Key = v[0]
			pair.Value = v[1]

		default:
			log.Printf("BBB")
			// "proj=utm=bzzt"
			return nil, BadProjStringError
		}

		ret.Add(pair)
	}

	return ret, nil

}

// Add adds a Pair to the end of the list
func (pl *PairList) Add(pair Pair) {
	pl.pairs = append(pl.pairs, pair)
}

// AddList adds a PairList's items to the end of the list
func (pl *PairList) AddList(list *PairList) {
	pl.pairs = append(pl.pairs, list.pairs...)
}

// Contains returns true iff the key is present in the list
func (pl *PairList) Contains(key string) bool {

	for _, pair := range pl.pairs {
		if pair.Key == key {
			return true
		}
	}

	return false
}

// Count returns the number of times the key is in the list
func (pl *PairList) Count(key string) int {

	count := 0
	for _, pair := range pl.pairs {
		if pair.Key == key {
			count++
		}
	}

	return count
}

// get returns the (string) value of the first occurrence of the key
func (pl *PairList) get(key string) (string, bool) {

	for _, pair := range pl.pairs {
		if pair.Key == key {
			return pair.Value, true
		}
	}

	return "", false
}

// GetAsString returns the value of the first occurrence of the key, as a string
func (pl *PairList) GetAsString(key string) (string, bool) {

	return pl.get(key)
}

// GetAsFloat returns the value of the first occurrence of the key, as a float64
func (pl *PairList) GetAsFloat(key string) (float64, bool) {

	value, ok := pl.get(key)
	if !ok {
		return 0.0, false
	}

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, false
	}

	return f, true
}

// GetAsFloats returns the value of the first occurrence of the key,
// interpretted as comma-separated floats
func (pl *PairList) GetAsFloats(key string) ([]float64, bool) {

	value, ok := pl.get(key)
	if !ok {
		return nil, false
	}

	nums := strings.Split(value, ",")

	floats := []float64{}

	for _, num := range nums {
		f, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return nil, false
		}
		floats = append(floats, f)
	}

	return floats, true
}
