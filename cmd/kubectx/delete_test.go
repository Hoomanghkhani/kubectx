package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/ahmetb/kubectx/internal/env"
)

func TestDeleteOpWarnsWhenDeletingCurrentContext(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		wantWarn bool
	}{
		{name: "by name", target: "ctx-a", wantWarn: true},
		{name: "by shorthand", target: ".", wantWarn: true},
		{name: "inactive context", target: "ctx-b", wantWarn: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("KUBECONFIG", writeTestKubeconfig(t))
			t.Setenv(env.EnvIsolatedShell, "")

			var stderr bytes.Buffer
			if err := (DeleteOp{Contexts: []string{tt.target}}).Run(io.Discard, &stderr); err != nil {
				t.Fatalf("delete context: %v", err)
			}
			if gotWarn := strings.Contains(stderr.String(), "You deleted the current context."); gotWarn != tt.wantWarn {
				t.Errorf("warning present = %t, want %t; stderr = %q", gotWarn, tt.wantWarn, stderr.String())
			}
		})
	}
}
