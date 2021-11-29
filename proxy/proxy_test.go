package proxy

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnescapeQueryString(t *testing.T) {
	type args struct {
		qs map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "unescape",
			args: args{
				qs: map[string]string{"hoge%5B%5D": "%7B%22foo%22%3A%20%22bar%22%7D", "abc": "def"},
			},
			want:    map[string]string{"hoge[]": `{"foo": "bar"}`, "abc": "def"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnescapeQueryString(tt.args.qs)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnescapeQueryString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("UnescapeQueryString() = %v, want %v\ndiff=%v", got, tt.want, diff)
			}
		})
	}
}

func TestUnescapeMultiValueQueryString(t *testing.T) {
	type args struct {
		qs map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][]string
		wantErr bool
	}{
		{
			name: "unescape",
			args: args{
				qs: map[string][]string{"hoge%5B%5D": {"%7B%22foo%22%3A%20%22bar%22%7D", "baz"}, "abc": {"def"}},
			},
			want:    map[string][]string{"hoge[]": {`{"foo": "bar"}`, "baz"}, "abc": {"def"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnescapeMultiValueQueryString(tt.args.qs)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnescapeMultiValueQueryString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			diff := cmp.Diff(got, tt.want)
			if diff != "" {
				t.Errorf("UnescapeMultiValueQueryString() = %v, want %v\ndiff=%v", got, tt.want, diff)
			}
		})
	}
}
