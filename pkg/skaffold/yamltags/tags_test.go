/*
Copyright 2019 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package yamltags

import (
	"reflect"
	"testing"
)

type otherstruct struct {
	A int `yamltags:"required"`
}

type required struct {
	A string `yamltags:"required"`
	B int    `yamltags:"required"`
	C otherstruct
}

func TestProcessStructRequired(t *testing.T) {
	type args struct {
		s interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "missing all",
			args: args{
				s: &required{},
			},
			wantErr: true,
		},
		{
			name: "all set",
			args: args{
				s: &required{
					A: "hey",
					B: 3,
					C: otherstruct{
						A: 1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missng some",
			args: args{
				s: &required{
					A: "hey",
					C: otherstruct{
						A: 1,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missng nested",
			args: args{
				s: &required{
					A: "hey",
					B: 3,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ProcessStruct(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ProcessStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type defaultTags struct {
	A string `yamltags:"default=foo"`
	B int    `yamltags:"default=3"`
}

func TestProcessStructDefault(t *testing.T) {
	type args struct {
		s interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "missing all",
			args: args{
				s: &defaultTags{},
			},
			want:    &defaultTags{A: "foo", B: 3},
			wantErr: false,
		},
		{
			name: "all set",
			args: args{
				s: &defaultTags{A: "yo", B: 1},
			},
			want:    &defaultTags{A: "yo", B: 1},
			wantErr: false,
		},
		{
			name: "some set",
			args: args{
				s: &defaultTags{A: "yo"},
			},
			want:    &defaultTags{A: "yo", B: 3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ProcessStruct(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ProcessStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.s, tt.want) {
				t.Errorf("Got %v, want %v", tt.args.s, tt.want)
			}
		})
	}
}

type oneOfStruct struct {
	A string  `yamltags:"oneOf=set1"`
	B string  `yamltags:"oneOf=set1"`
	C int     `yamltags:"oneOf=set2"`
	D *nested `yamltags:"oneOf=set2"`
	E nested  `yamltags:"oneOf=set2"`
}

type nested struct {
	F string
}

func TestOneOf(t *testing.T) {
	type args struct {
		s interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "only one",
			args: args{
				s: &oneOfStruct{
					A: "foo",
					C: 3,
				},
			},
			wantErr: false,
		},
		{
			name: "too many in one set",
			args: args{
				s: &oneOfStruct{
					A: "foo",
					B: "baz",
					C: 3,
				}},
			wantErr: true,
		},
		{
			name: "too many pointers set",
			args: args{
				s: &oneOfStruct{
					D: &nested{F: "foo"},
					E: nested{F: "foo"},
				}},
			wantErr: true,
		},
		{
			name: "too many in both sets",
			args: args{
				s: &oneOfStruct{
					A: "foo",
					B: "baz",
					C: 3,
					D: &nested{F: "foo"},
				}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ProcessStruct(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ProcessStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
