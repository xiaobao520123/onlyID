package snowflake

import "testing"

func TestSnowFlake(t *testing.T) {
	h, err := NewHost(0)
	if err != nil {
		t.Error(err)
	}
	id := h.Generate()
	if id.ToInt64() == 0 {
		t.Errorf("generate id failed")
	}
}
