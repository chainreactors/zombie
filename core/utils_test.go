package core

import (
	"reflect"
	"testing"
)

func TestParseUrl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *Target
		ok    bool
	}{
		{
			name:  "http url with credentials and explicit port",
			input: "http://user:pass@127.0.0.1:8080/path",
			ok:    true,
			want: &Target{
				IP:       "127.0.0.1",
				Port:     "8080",
				Service:  "http",
				Scheme:   "http",
				Username: "user",
				Password: "pass",
			},
		},
		{
			name:  "https url with default port",
			input: "https://127.0.0.1",
			ok:    true,
			want: &Target{
				IP:      "127.0.0.1",
				Port:    "443",
				Service: "https",
				Scheme:  "https",
			},
		},
		{
			name:  "service url with default port",
			input: "redis://127.0.0.1",
			ok:    true,
			want: &Target{
				IP:      "127.0.0.1",
				Port:    "6379",
				Service: "redis",
				Scheme:  "redis",
			},
		},
		{
			name:  "plain ip",
			input: "127.0.0.1",
			ok:    true,
			want: &Target{
				IP: "127.0.0.1",
			},
		},
		{
			name:  "host and port fallback",
			input: "127.0.0.1:3306",
			ok:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseUrl(tt.input)
			if ok != tt.ok {
				t.Fatalf("unexpected ok result for %q: got %v want %v", tt.input, ok, tt.ok)
			}
			if !tt.ok {
				if got != nil {
					t.Fatalf("expected nil target for %q, got %#v", tt.input, got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("unexpected target for %q:\n got: %#v\nwant: %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSimpleParseUrl(t *testing.T) {
	got := SimpleParseUrl("127.0.0.1:3306")
	want := &Target{IP: "127.0.0.1", Port: "3306"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected target: got %#v want %#v", got, want)
	}
}

func TestOptionPrepareBuildsTargetsFromInputs(t *testing.T) {
	opt := &Option{}
	opt.IP = []string{
		"https://user:pass@127.0.0.1:8443/path",
		"127.0.0.1:3306",
		"127.0.0.1",
	}
	opt.ServiceName = "mysql"
	opt.Mod = ModSniper

	runner, err := opt.Prepare()
	if err != nil {
		t.Fatal(err)
	}
	if len(runner.Targets) != 3 {
		t.Fatalf("unexpected target count: got %d want 3", len(runner.Targets))
	}

	got := runner.Targets
	if got[0].IP != "127.0.0.1" || got[0].Port != "8443" || got[0].Service != "mysql" || got[0].Username != "user" || got[0].Password != "pass" {
		t.Fatalf("unexpected first target: %#v", got[0])
	}
	if got[1].IP != "127.0.0.1" || got[1].Port != "3306" || got[1].Service != "mysql" {
		t.Fatalf("unexpected second target: %#v", got[1])
	}
	if got[2].IP != "127.0.0.1" || got[2].Port != "3306" || got[2].Service != "mysql" {
		t.Fatalf("unexpected third target: %#v", got[2])
	}
}
