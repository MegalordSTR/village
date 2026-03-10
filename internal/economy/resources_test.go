package economy

import (
	"math"
	"testing"
)

func TestResourceValidate(t *testing.T) {
	tests := []struct {
		name     string
		resource Resource
		want     bool
	}{
		{
			name: "valid resource",
			resource: Resource{
				Type:     ResourceGrain,
				Quantity: 100.0,
				Quality:  QualityNormal,
			},
			want: true,
		},
		{
			name: "unknown resource type",
			resource: Resource{
				Type:     "unknown",
				Quantity: 100.0,
			},
			want: false,
		},
		{
			name: "negative quantity",
			resource: Resource{
				Type:     ResourceGrain,
				Quantity: -5.0,
			},
			want: false,
		},
		{
			name: "zero quantity",
			resource: Resource{
				Type:     ResourceGrain,
				Quantity: 0.0,
			},
			want: true,
		},
		{
			name: "NaN quantity",
			resource: Resource{
				Type:     ResourceGrain,
				Quantity: math.NaN(),
			},
			want: false,
		},
		{
			name: "positive infinity quantity",
			resource: Resource{
				Type:     ResourceGrain,
				Quantity: math.Inf(1),
			},
			want: false,
		},
		{
			name: "negative infinity quantity",
			resource: Resource{
				Type:     ResourceGrain,
				Quantity: math.Inf(-1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.resource.Validate()
			if got != tt.want {
				t.Errorf("Resource.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
