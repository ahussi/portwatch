package tagger_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/tagger"
)

func TestNew(t *testing.T) {
	tr := tagger.New()
	if tr == nil {
		t.Fatal("expected non-nil Tagger")
	}
}

func TestGet_BuiltIn(t *testing.T) {
	tr := tagger.New()
	tags := tr.Get(80)
	if len(tags) == 0 {
		t.Fatal("expected at least one built-in tag for port 80")
	}
	if tags[0].Label != "http" {
		t.Errorf("expected label 'http', got %q", tags[0].Label)
	}
	if tags[0].Source != "builtin" {
		t.Errorf("expected source 'builtin', got %q", tags[0].Source)
	}
}

func TestGet_Unknown(t *testing.T) {
	tr := tagger.New()
	tags := tr.Get(9999)
	if len(tags) != 0 {
		t.Errorf("expected no tags for unknown port, got %d", len(tags))
	}
}

func TestAdd_UserTag(t *testing.T) {
	tr := tagger.New()
	tr.Add(9999, "my-service")
	tags := tr.Get(9999)
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Label != "my-service" || tags[0].Source != "user" {
		t.Errorf("unexpected tag: %v", tags[0])
	}
}

func TestAdd_BuiltInAndUser(t *testing.T) {
	tr := tagger.New()
	tr.Add(80, "frontend")
	tags := tr.Get(80)
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags (builtin + user), got %d", len(tags))
	}
}

func TestRemove(t *testing.T) {
	tr := tagger.New()
	tr.Add(9999, "temp")
	tr.Remove(9999)
	tags := tr.Get(9999)
	if len(tags) != 0 {
		t.Errorf("expected no tags after Remove, got %d", len(tags))
	}
}

func TestHasTag_True(t *testing.T) {
	tr := tagger.New()
	if !tr.HasTag(443, "https") {
		t.Error("expected HasTag to return true for port 443 / 'https'")
	}
}

func TestHasTag_False(t *testing.T) {
	tr := tagger.New()
	if tr.HasTag(443, "ftp") {
		t.Error("expected HasTag to return false for port 443 / 'ftp'")
	}
}

func TestTagString(t *testing.T) {
	tag := tagger.Tag{Label: "http", Source: "builtin"}
	if tag.String() != "http(builtin)" {
		t.Errorf("unexpected String(): %q", tag.String())
	}
}
