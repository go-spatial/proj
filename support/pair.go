package support

import (
	"strconv"
	"strings"

	"github.com/go-spatial/proj4go/merror"
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
	Pairs []Pair
}

// NewPairList returns a new PairList
func NewPairList() *PairList {
	return &PairList{
		Pairs: []Pair{},
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
			return nil, merror.New(merror.BadProjStringError)
		}

		switch len(v) {
		case 0:
			pair.Key = w
			pair.Value = ""
		case 1:
			// "proj=" is okay
			pair.Key = v[0]
			pair.Value = ""
		case 2:
			pair.Key = v[0]
			pair.Value = v[1]

		default:
			// "proj=utm=bzzt"
			return nil, merror.New(merror.BadProjStringError)
		}

		ret.Add(pair)
	}

	return ret, nil

}

// Len returns the number of pairs in the list
func (pl *PairList) Len() int {
	return len(pl.Pairs)
}

// Get returns the ith pair in the list
func (pl *PairList) Get(i int) Pair {
	return pl.Pairs[i]
}

// Add adds a Pair to the end of the list
func (pl *PairList) Add(pair Pair) {
	pl.Pairs = append(pl.Pairs, pair)
}

// AddList adds a PairList's items to the end of the list
func (pl *PairList) AddList(list *PairList) {
	pl.Pairs = append(pl.Pairs, list.Pairs...)
}

// ContainsKey returns true iff the key is present in the list
func (pl *PairList) ContainsKey(key string) bool {

	for _, pair := range pl.Pairs {
		if pair.Key == key {
			return true
		}
	}

	return false
}

// CountKey returns the number of times the key is in the list
func (pl *PairList) CountKey(key string) int {

	count := 0
	for _, pair := range pl.Pairs {
		if pair.Key == key {
			count++
		}
	}

	return count
}

// get returns the (string) value of the first occurrence of the key
func (pl *PairList) get(key string) (string, bool) {

	for _, pair := range pl.Pairs {
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

// GetAsInt returns the value of the first occurrence of the key, as an int
func (pl *PairList) GetAsInt(key string) (int, bool) {
	value, ok := pl.get(key)
	if !ok {
		return 0, false
	}
	i64, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, false
	}

	return int(i64), true
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
