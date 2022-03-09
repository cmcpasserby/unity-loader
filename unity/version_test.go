package unity

import (
	"reflect"
	"testing"
)

func TestVersionFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    VersionData
		wantErr bool
	}{
		{
			name:  "no_hash",
			input: "2020.3.30f1",
			want: VersionData{
				Major:   2020,
				Minor:   3,
				Update:  30,
				VerType: "f",
				Patch:   1,
			},
			wantErr: false,
		},
		{
			name:    "invalid_ver_type",
			input:   "2020.3.30g1",
			want:    VersionData{},
			wantErr: true,
		},
		{
			name:    "invalid_number",
			input:   "202f0.3.30f1",
			want:    VersionData{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VersionFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("VersionFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VersionFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
