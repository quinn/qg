package config

import (
	"testing"
)

func TestConfig_FindGenerator(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		gName       string
		wantErr     bool
		wantErrText string
	}{
		{
			name: "finds existing generator",
			config: Config{
				Version: "1.0",
				Generators: []Generator{
					{
						Name: "test-gen",
						Args: []string{"arg1", "arg2"},
					},
				},
			},
			gName:   "test-gen",
			wantErr: false,
		},
		{
			name: "returns error for non-existent generator",
			config: Config{
				Version:    "1.0",
				Generators: []Generator{},
			},
			gName:       "missing-gen",
			wantErr:     true,
			wantErrText: "generator not found: missing-gen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.FindGenerator(tt.gName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.FindGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err.Error() != tt.wantErrText {
					t.Errorf("Config.FindGenerator() error = %v, wantErrText %v", err, tt.wantErrText)
				}
				return
			}
			if got.Name != tt.gName {
				t.Errorf("Config.FindGenerator() = %v, want %v", got.Name, tt.gName)
			}
		})
	}
}
