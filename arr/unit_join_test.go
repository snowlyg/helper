package arr

import "testing"

func TestUnitJoin(t *testing.T) {
	type args struct {
		ss  []uint
		sep string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success",
			args: struct {
				ss  []uint
				sep string
			}{ss: []uint{1, 2, 3, 4}, sep: ","},
			want: "1,2,3,4",
		},
		{
			name: "success",
			args: struct {
				ss  []uint
				sep string
			}{ss: []uint{1, 2, 3, 4}, sep: "||"},
			want: "1||2||3||4",
		},
		{
			name: "success",
			args: struct {
				ss  []uint
				sep string
			}{ss: []uint{1, 2, 3, 4}, sep: "-"},
			want: "1-2-3-4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnitJoin(tt.args.ss, tt.args.sep); got != tt.want {
				t.Errorf("UnitJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}
