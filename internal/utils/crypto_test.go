package utils

import "testing"

func TestMd5(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "admin", args: args{str: "21232f297a57a5a743894a0e4a801fc3" + "21232f297a57a5a74"}, want: "3c3d20cf4936b81600306b09ab1f6cf4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5(tt.args.str); got != tt.want {
				t.Errorf("Md5() = %v, want %v", got, tt.want)
			}
		})
	}
}
