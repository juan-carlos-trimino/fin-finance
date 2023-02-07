//

package finances

import "testing"

func TestMiscellaneous_RealInterestRate(t *testing.T) {
	type args struct {
		nominalRate   float64
		inflationRate float64
	}
	tests := []struct {
		name string
		m    *Miscellaneous
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.RealInterestRate(tt.args.nominalRate, tt.args.inflationRate); got != tt.want {
				t.Errorf("Miscellaneous.RealInterestRate() = %v, want %v", got, tt.want)
			}
		})
	}
}
