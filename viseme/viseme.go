package viseme

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kkoch986/ai-skeletons-playback/output"
)

// Viseme represents the shape the mouth makes to emit a specific sound
// see more info here:
// https://docs.aws.amazon.com/polly/latest/dg/ph-table-english-us.html
type Viseme string

var mapping = map[Viseme]float64{
	"e": 0.2,
	"k": 0.2,
	"a": 0.8,
	"p": 0,
	"t": 0.3,
	"f": 0.1,
	"T": 0.3,
	"s": 0.2,
	"@": 0.7,
	"E": 0.6,
	"S": 0.3,
	"u": 0.25,
	"O": 0.25,
}

func (v Viseme) ToJawOpenness() float64 {
	val, ok := mapping[v]
	if !ok {
		return 0
	}
	return val
}

// Sequence represents a timestamped array of Visemes
// The keys of the map are the float second offsets from the start
// of the sequence
type Sequence map[float64]Viseme

func (r Sequence) ToJawOutputSequence(key string) output.Sequence {
	ret := output.Sequence{}
	for t, v := range r {
		ret.AddValue(t, output.Value{
			Key:   key,
			Value: v.ToJawOpenness(),
		})
	}
	return ret
}

func (r Sequence) MarshalJSON() ([]byte, error) {
	// TODO: update this so it marshals the same as it unmarshals [ [time, viseme], ...]
	m := make(map[string]Viseme, len(r))
	for k, v := range r {
		m[fmt.Sprintf("%f", k)] = v
	}
	return json.Marshal(m)
}

func (r Sequence) UnmarshalJSON(p []byte) error {
	if len(p) == 0 {
		return nil
	}
	var tmp []interface{}
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}
	for i, v := range tmp {
		va, ok := v.([]interface{})
		if !ok || len(va) != 2 {
			return errors.New("unable to parse VisemeSequence, invalid sequence entry")
		}
		ts, ok := va[0].(float64)
		if !ok {
			return fmt.Errorf("unable to parse VisemeSequence, invalid timestamp at position %d", i)
		}

		vis, ok := va[1].(string)
		if !ok {
			return fmt.Errorf("unable to parse VisemeSequence, invalid viseme at position %d", i)
		}
		r[ts] = Viseme(vis)
	}
	return nil
}
