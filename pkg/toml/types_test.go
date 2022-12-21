package toml

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Duration(t *testing.T) {
	d := Duration(time.Minute)
	assert.Equal(t, d.Duration(), time.Minute)

	marshalF := func(duration time.Duration) string {
		d := Duration(duration)
		txt, _ := d.MarshalText()
		return string(txt)
	}
	unmarshalF := func(txt string) time.Duration {
		var d Duration
		_ = d.UnmarshalText([]byte(txt))
		return d.Duration()
	}
	assert.Equal(t, "1m0s", marshalF(time.Minute))
	assert.Equal(t, "10s", marshalF(time.Second*10))

	assert.Equal(t, time.Second, unmarshalF("1s"))
	assert.Equal(t, time.Minute, unmarshalF("1m"))
	assert.Equal(t, time.Hour, unmarshalF("3600s"))

	assert.Zero(t, unmarshalF(""))
	assert.Zero(t, unmarshalF("1fs"))
}

func Test_Duration_JSON(t *testing.T) {
	type Example struct {
		Cost Duration `json:"cost"`
		A    int      `json:"a"`
	}
	example := Example{Cost: Duration(time.Nanosecond * 1102), A: 23}
	data, err := json.Marshal(example)
	assert.NoError(t, err)
	_, err = example.Cost.MarshalJSON()
	assert.NoError(t, err)
	_, err = example.Cost.MarshalText()
	assert.NoError(t, err)

	var newExample Example
	assert.NoError(t, json.Unmarshal(data, &newExample))
	assert.Equal(t, Duration(time.Nanosecond*1102), newExample.Cost)

	txt := "{\"cost\": 322}"
	assert.NoError(t, json.Unmarshal([]byte(txt), &newExample))
	assert.Equal(t, Duration(time.Nanosecond*322), newExample.Cost)

	assert.Error(t, json.Unmarshal([]byte(`{"cost": "xxxx"}`), &newExample))
	assert.Error(t, json.Unmarshal([]byte(`{"cost": "xxxx"}`), &newExample))

	assert.Error(t, json.Unmarshal([]byte(`{"cost": null}`), &newExample))
	assert.Error(t, json.Unmarshal([]byte{1, 0}, &newExample))

	assert.Nil(t, json.Unmarshal([]byte(`{"cost": "22.265928ms"}`), &newExample))
}

func Test_Size(t *testing.T) {
	type Example struct {
		Size Size `json:"size"`
	}
	example1 := Example{Size: 10240}
	assert.Equal(t, "10 KiB", example1.Size.String())

	txt, err := example1.Size.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "10 KiB", string(txt))

	txt, err = example1.Size.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, `"10 KiB"`, string(txt))

	var s2 Example
	assert.Error(t, json.Unmarshal([]byte(``), &s2))
	assert.Error(t, json.Unmarshal([]byte(`{"size": null`), &s2))
	assert.Error(t, json.Unmarshal([]byte(`{"size": true`), &s2))
	assert.Error(t, json.Unmarshal([]byte(`{"size": "10 MiB"`), &s2))
	assert.NoError(t, json.Unmarshal([]byte(`{"size": "10 MiB"}`), &s2))
	assert.Equal(t, Size(0xa00000), s2.Size)
	assert.Error(t, json.Unmarshal([]byte(`{"size": "10 iB"}`), &s2))
	assert.NoError(t, json.Unmarshal([]byte(`{"size": 1000}`), &s2))
	assert.NoError(t, json.Unmarshal([]byte(`{"size": "969 B"}`), &s2))
	assert.NoError(t, json.Unmarshal([]byte(`{"size": 10}`), &s2))
	assert.Error(t, json.Unmarshal([]byte("{\"size\": \"\"\"\"}"), &s2))
}
