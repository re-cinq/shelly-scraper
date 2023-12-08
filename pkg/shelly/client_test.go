package shelly

import "testing"

func TestGetSwitchStatus(t *testing.T) {
	t.Run("get status", func(t *testing.T) {
		c := New("192.168.2.51")
		c.GetSwitchStatus("0")
	})
}
