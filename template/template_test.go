package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProcessor_ProcessPath(t *testing.T) {
	processor := New("/templates", "/output")
	config := map[string]string{
		"name": "test",
		"type": "component",
	}

	tests := []struct {
		name        string
		path        string
		want        string
		wantErr     bool
		wantErrText string
	}{
		{
			name: "simple path",
