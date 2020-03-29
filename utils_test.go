/*
 * @Package captcha
 * @Author Quan Chen
 * @Date 2020/3/29
 * @Description
 *
 */
package captcha

import "testing"

func TestRandText(t *testing.T) {
	type args struct {
		num   int
		chars []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"1",
			args{
				num:   20,
				chars: []string{"232"},
			},
			"1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandText(tt.args.num, tt.args.chars...); got != tt.want {
				t.Errorf("RandText() = %v, want %v", got, tt.want)
			}
		})
	}
}
