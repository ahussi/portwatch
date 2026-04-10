package filter_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/filter"
)

func baseConfig() *config.Config {
	cfg := config.Default()
	cfg.WatchPorts = []int{8080, 9090, 3000}
	cfg.AllowPorts = []int{9090}
	return cfg
}

func TestResultString(t *testing.T) {
	cases := []struct {
		r    filter.Result
		want string
	}{
		{filter.Allowed, "allowed"},
		{filter.DeniedNotWatched, "denied:not-watched"},
		{filter.DeniedExplicit, "denied:explicit"},
		{filter.Result(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.r.String(); got != tc.want {
			t.Errorf("Result(%d).String() = %q, want %q", tc.r, got, tc.want)
		}
	}
}

func TestCheckAllowed(t *testing.T) {
	f := filter.New(baseConfig())
	if got := f.Check(8080); got != filter.Allowed {
		t.Errorf("expected Allowed for 8080, got %s", got)
	}
}

func TestCheckDeniedExplicit(t *testing.T) {
	f := filter.New(baseConfig())
	if got := f.Check(9090); got != filter.DeniedExplicit {
		t.Errorf("expected DeniedExplicit for 9090, got %s", got)
	}
}

func TestCheckDeniedNotWatched(t *testing.T) {
	f := filter.New(baseConfig())
	if got := f.Check(5432); got != filter.DeniedNotWatched {
		t.Errorf("expected DeniedNotWatched for 5432, got %s", got)
	}
}

func TestShouldAlert(t *testing.T) {
	f := filter.New(baseConfig())

	if !f.ShouldAlert(8080) {
		t.Error("expected ShouldAlert=true for watched, non-allowed port 8080")
	}
	if f.ShouldAlert(9090) {
		t.Error("expected ShouldAlert=false for explicitly allowed port 9090")
	}
	if f.ShouldAlert(5432) {
		t.Error("expected ShouldAlert=false for unwatched port 5432")
	}
}

func TestNew(t *testing.T) {
	f := filter.New(baseConfig())
	if f == nil {
		t.Fatal("New returned nil")
	}
}
