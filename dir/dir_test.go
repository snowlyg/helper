package dir

import "testing"

func TestSelfPath(t *testing.T) {
	t.Run("test selfpath", func(t *testing.T) {
		selfPath, _ := RealPath("air")
		t.Log(selfPath)
	})
}
