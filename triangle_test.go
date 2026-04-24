package main

import "testing"

const red = "\033[31m"
const reset = "\033[0m"

func Test_IsPixelInTriangle(t *testing.T) {

	tests := []struct {
		name  string
		point Vec3
		a     Vec3
		b     Vec3
		c     Vec3
		want  bool
	}{
		{
			name:  "Point on edge",
			point: Vec3{1, 1, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{1, 1, 0},
			c:     Vec3{0, 2, 0},
			want:  true,
		},
		{
			name:  "Point outside",
			point: Vec3{4, 4, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{1, 1, 0},
			c:     Vec3{0, 2, 0},
			want:  false,
		},
		{
			name:  "Point inside",
			point: Vec3{1, 2, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 2, 0},
			c:     Vec3{0, 4, 0},
			want:  true,
		},
		{
			name:  "Point inside 2",
			point: Vec3{1, 1.5, 0},
			a:     Vec3{0, 1, 0},
			b:     Vec3{2, 1, 0},
			c:     Vec3{0.5, 2, 0},
			want:  true,
		},
		{
			name:  "Point on vertex A",
			point: Vec3{0, 0, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 0, 0},
			c:     Vec3{0, 2, 0},
			want:  true,
		},
		{
			name:  "Point on vertex B",
			point: Vec3{2, 0, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 0, 0},
			c:     Vec3{0, 2, 0},
			want:  true,
		},
		{
			name:  "Point on vertex C",
			point: Vec3{0, 2, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 0, 0},
			c:     Vec3{0, 2, 0},
			want:  true,
		},
		{
			name:  "Point on edge AB",
			point: Vec3{1, 0, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 0, 0},
			c:     Vec3{0, 2, 0},
			want:  true,
		},
		{
			name:  "Point on edge AC",
			point: Vec3{0, 1, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 0, 0},
			c:     Vec3{0, 2, 0},
			want:  true,
		},
		{
			name:  "Point on edge BC",
			point: Vec3{1, 1, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 0, 0},
			c:     Vec3{0, 2, 0},
			want:  true,
		},
		{
			name:  "Just outside edge",
			point: Vec3{1, -0.0001, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 0, 0},
			c:     Vec3{0, 2, 0},
			want:  false,
		},
		{
			name:  "Just outside hypotenuse",
			point: Vec3{1.01, 1.01, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 0, 0},
			c:     Vec3{0, 2, 0},
			want:  false,
		},
		{
			name:  "Degenerate triangle (collinear)",
			point: Vec3{1, 1, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{1, 1, 0},
			c:     Vec3{2, 2, 0},
			want:  false,
		},
		{
			name:  "Reversed winding",
			point: Vec3{1, 1, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{0, 2, 0},
			c:     Vec3{2, 0, 0},
			want:  true,
		},
		{
			name:  "Large coordinates",
			point: Vec3{10001, 10001, 0},
			a:     Vec3{10000, 10000, 0},
			b:     Vec3{20000, 10000, 0},
			c:     Vec3{10000, 20000, 0},
			want:  true,
		},
		{
			name:  "Negative space",
			point: Vec3{-1, -1, 0},
			a:     Vec3{-2, -2, 0},
			b:     Vec3{0, -2, 0},
			c:     Vec3{-2, 0, 0},
			want:  true,
		},
		{
			name:  "Skinny triangle",
			point: Vec3{1, 0.001, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 0, 0},
			c:     Vec3{1, 0.01, 0},
			want:  true,
		},
		{
			name:  "Outside near vertex",
			point: Vec3{2.1, 0.1, 0},
			a:     Vec3{0, 0, 0},
			b:     Vec3{2, 0, 0},
			c:     Vec3{0, 2, 0},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, _, _ := IsPixelInTriangle(tt.point, tt.a, tt.b, tt.c)
			if got != tt.want {

				t.Errorf("\n"+red+"IsPixelInTriangle() = %v, want %v"+reset, got, tt.want)
				return
			}
		})
	}
}
