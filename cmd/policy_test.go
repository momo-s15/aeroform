package cmd

import (
	"testing"
)

func TestPolicyCommandsRegistered(t *testing.T) {
	subs := policyCmd.Commands()
	names := make(map[string]bool)
	for _, c := range subs {
		names[c.Name()] = true
	}
	for _, want := range []string{"add", "list", "remove"} {
		if !names[want] {
			t.Errorf("policy subcommand %q not registered", want)
		}
	}
}
