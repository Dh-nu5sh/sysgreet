package bootstrap

import "testing"

func TestResolvePolicy_FlagOverridesEnv(t *testing.T) {
	res, err := ResolvePolicy("keep", "overwrite", true)
	if err != nil {
		t.Fatalf("ResolvePolicy error: %v", err)
	}
	if res.Value != PolicyKeep {
		t.Fatalf("expected keep policy, got %s", res.Value)
	}
	if res.Source != PolicySourceFlag {
		t.Fatalf("expected source flag, got %s", res.Source)
	}
}

func TestResolvePolicy_EnvUsedWhenFlagEmpty(t *testing.T) {
	res, err := ResolvePolicy("", "overwrite", true)
	if err != nil {
		t.Fatalf("ResolvePolicy error: %v", err)
	}
	if res.Value != PolicyOverwrite {
		t.Fatalf("expected overwrite policy, got %s", res.Value)
	}
	if res.Source != PolicySourceEnv {
		t.Fatalf("expected env source, got %s", res.Source)
	}
}

func TestResolvePolicy_DefaultPromptInteractive(t *testing.T) {
	res, err := ResolvePolicy("", "", true)
	if err != nil {
		t.Fatalf("ResolvePolicy error: %v", err)
	}
	if res.Value != PolicyPrompt {
		t.Fatalf("expected prompt policy, got %s", res.Value)
	}
	if res.Source != PolicySourceDefault {
		t.Fatalf("expected default source, got %s", res.Source)
	}
}

func TestResolvePolicy_InvalidValue(t *testing.T) {
	if _, err := ResolvePolicy("invalid", "", true); err == nil {
		t.Fatalf("expected error for invalid policy value")
	}
}

func TestResolvePolicy_NonInteractiveRequiresExplicitPolicy(t *testing.T) {
	if _, err := ResolvePolicy("", "", false); err != ErrPolicyRequired {
		t.Fatalf("expected ErrPolicyRequired, got %v", err)
	}
}
