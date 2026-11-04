package keqing

import "testing"

func TestIssue14YAMLPathSafety(t *testing.T) {
	oldParent, oldChild := applicationYml, applicationChildYml
	t.Cleanup(func() { applicationYml, applicationChildYml = oldParent, oldChild })
	for _, parentNode := range []interface{}{"disabled", []interface{}{"entry"}, nil} {
		if err := LoadYml(map[string]interface{}{"service": parentNode}, map[string]interface{}{}); err != nil {
			t.Fatal(err)
		}
		if got := ValueString("${service.url}"); got != "" {
			t.Fatalf("missing scalar path returned %q", got)
		}
		if got := ValueInt("${service.port}"); got != 0 {
			t.Fatal(got)
		}
		if err := LoadYml(map[string]interface{}{"service": parentNode}, map[string]interface{}{"service": map[string]interface{}{"url": "https://example.test"}}); err != nil {
			t.Fatal(err)
		}
		if got := ValueString("${service.url}"); got != "https://example.test" {
			t.Fatal(got)
		}
	}
}

func TestIssue14YAMLOverrideAndFallback(t *testing.T) {
	oldParent, oldChild := applicationYml, applicationChildYml
	t.Cleanup(func() { applicationYml, applicationChildYml = oldParent, oldChild })
	parent := map[string]interface{}{"service": map[string]interface{}{"url": "parent", "port": 80, "enabled": true}}
	for _, childNode := range []interface{}{map[string]interface{}{}, nil, "disabled", []interface{}{}} {
		if err := LoadYml(parent, map[string]interface{}{"service": childNode}); err != nil {
			t.Fatal(err)
		}
		if got := ValueString("${service.url}"); got != "parent" {
			t.Fatalf("missing child leaf must fall back: %q", got)
		}
	}
	child := map[string]interface{}{"service": map[string]interface{}{"url": "", "port": 0, "enabled": false}}
	if err := LoadYml(parent, child); err != nil {
		t.Fatal(err)
	}
	if ValueString("${service.url}") != "" || ValueInt("${service.port}") != 0 || ValueBool("${service.enabled}") {
		t.Fatal("non-nil zero values must override parent")
	}
	child["service"].(map[string]interface{})["url"] = nil
	if ValueString("${service.url}") != "parent" {
		t.Fatal("explicit null must retain existing parent fallback")
	}
	if ValueObject("${missing.deep.key}") != nil {
		t.Fatal("missing both must be nil")
	}
	valuesMustPanic(t, func() { ValueObject("service.url") })
}
