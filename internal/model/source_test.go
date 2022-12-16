package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetSources(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		error      bool
	}{
		{
			name:       "ok",
			configPath: "../../configs/event-sources.toml",
			error:      false,
		},
		{
			name:       "wrong config path",
			configPath: "wrongPath/configs/event-sources.toml",
			error:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSources(tt.configPath)

			if tt.error == true {
				assert.Error(t, err)
				return
			} else {
				if !assert.NoError(t, err) {
					return
				}
			}

			if assert.True(t, len(got) > 0, "no sources in file") {
				hasWantName := false
				for _, src := range got {
					hasWantName = hasWantName == true || isWantName(src.Name)
					assert.NotEmpty(t, src.Name, "empty name")
					assert.NotEmpty(t, src.Url, "empty url")
				}
				assert.True(t, hasWantName, "not found want source name")
			}
		})
	}
}

func isWantName(name string) bool {
	namesWant := [...]string{"porto"}

	for _, nameWant := range namesWant {
		if nameWant == name {
			return true
		}
	}
	return false
}
