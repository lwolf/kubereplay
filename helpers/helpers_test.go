package helpers

import "testing"

func TestBlueGreenReplicas(t *testing.T) {
	tests := []struct {
		n           int32
		segmentSize int32
		blue        int32
		green       int32
	}{
		{1, 100, 1, 0},
		{1, 50, 1, 1},
		{10, 50, 5, 5},
		{10, 30, 3, 7},
		{7, 30, 2, 5},
	}

	for _, tt := range tests {
		b, g := BlueGreenReplicas(tt.n, tt.segmentSize)

		if g != tt.green || b != tt.blue {
			t.Errorf("Calculating replicas from %d with %d is wrong, want %d/%d, got %d/%d", tt.n, tt.segmentSize, tt.blue, tt.green, b, g)
		}
	}
}
