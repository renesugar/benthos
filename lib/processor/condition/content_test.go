// Copyright (c) 2018 Ashley Jeffs
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package condition

import (
	"os"
	"testing"

	"github.com/Jeffail/benthos/lib/metrics"
	"github.com/Jeffail/benthos/lib/types"
	"github.com/Jeffail/benthos/lib/util/service/log"
)

func TestContentCheck(t *testing.T) {
	testLog := log.NewLogger(os.Stdout, log.LoggerConfig{LogLevel: "NONE"})
	testMet := metrics.DudType{}

	type fields struct {
		operator string
		part     int
		arg      string
	}
	tests := []struct {
		name   string
		fields fields
		arg    [][]byte
		want   bool
	}{
		{
			name: "equals_cs foo pos",
			fields: fields{
				operator: "equals_cs",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("foo"),
			},
			want: true,
		},
		{
			name: "equals_cs foo neg",
			fields: fields{
				operator: "equals_cs",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("not foo"),
			},
			want: false,
		},
		{
			name: "equals foo pos",
			fields: fields{
				operator: "equals",
				part:     0,
				arg:      "fOo",
			},
			arg: [][]byte{
				[]byte("foo"),
			},
			want: true,
		},
		{
			name: "equals foo pos 2",
			fields: fields{
				operator: "equals",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("fOo"),
			},
			want: true,
		},
		{
			name: "equals foo neg",
			fields: fields{
				operator: "equals",
				part:     0,
				arg:      "fOo",
			},
			arg: [][]byte{
				[]byte("f0o"),
			},
			want: false,
		},
		{
			name: "contains_cs foo pos",
			fields: fields{
				operator: "contains_cs",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("hello foo world"),
			},
			want: true,
		},
		{
			name: "contains_cs foo neg",
			fields: fields{
				operator: "contains_cs",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("hello fOo world"),
			},
			want: false,
		},
		{
			name: "contains foo pos",
			fields: fields{
				operator: "contains",
				part:     0,
				arg:      "fOo",
			},
			arg: [][]byte{
				[]byte("hello foo world"),
			},
			want: true,
		},
		{
			name: "contains foo pos 2",
			fields: fields{
				operator: "contains",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("hello fOo world"),
			},
			want: true,
		},
		{
			name: "contains foo neg",
			fields: fields{
				operator: "contains",
				part:     0,
				arg:      "fOo",
			},
			arg: [][]byte{
				[]byte("hello f0o world"),
			},
			want: false,
		},
		{
			name: "equals_cs foo pos from neg index",
			fields: fields{
				operator: "equals_cs",
				part:     -1,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("bar"),
				[]byte("foo"),
			},
			want: true,
		},
		{
			name: "equals_cs foo neg from neg index",
			fields: fields{
				operator: "equals_cs",
				part:     -2,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("bar"),
				[]byte("foo"),
			},
			want: false,
		},
		{
			name: "equals_cs neg empty msg",
			fields: fields{
				operator: "equals_cs",
				part:     0,
				arg:      "foo",
			},
			arg:  [][]byte{},
			want: false,
		},
		{
			name: "equals_cs neg oob",
			fields: fields{
				operator: "equals_cs",
				part:     1,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("foo"),
			},
			want: false,
		},
		{
			name: "equals_cs neg oob neg index",
			fields: fields{
				operator: "equals_cs",
				part:     -2,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("foo"),
			},
			want: false,
		},
		{
			name: "prefix_cs foo pos",
			fields: fields{
				operator: "prefix_cs",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("foo hello world"),
			},
			want: true,
		},
		{
			name: "prefix_cs foo neg",
			fields: fields{
				operator: "prefix_cs",
				part:     0,
				arg:      "fOo",
			},
			arg: [][]byte{
				[]byte("foo hello world"),
			},
			want: false,
		},
		{
			name: "prefix_cs foo neg 2",
			fields: fields{
				operator: "prefix_cs",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("hello foo world"),
			},
			want: false,
		},
		{
			name: "prefix foo pos",
			fields: fields{
				operator: "prefix",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("foo hello world"),
			},
			want: true,
		},
		{
			name: "prefix foo pos 2",
			fields: fields{
				operator: "prefix",
				part:     0,
				arg:      "fOo",
			},
			arg: [][]byte{
				[]byte("FoO hello world"),
			},
			want: true,
		},
		{
			name: "prefix foo neg",
			fields: fields{
				operator: "prefix",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("hello foo world"),
			},
			want: false,
		},
		{
			name: "suffix_cs foo pos",
			fields: fields{
				operator: "suffix_cs",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("hello world foo"),
			},
			want: true,
		},
		{
			name: "suffix_cs foo neg",
			fields: fields{
				operator: "suffix_cs",
				part:     0,
				arg:      "fOo",
			},
			arg: [][]byte{
				[]byte("hello world foo"),
			},
			want: false,
		},
		{
			name: "suffix_cs foo neg 2",
			fields: fields{
				operator: "suffix_cs",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("hello foo world"),
			},
			want: false,
		},
		{
			name: "suffix foo pos",
			fields: fields{
				operator: "suffix",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("hello world foo"),
			},
			want: true,
		},
		{
			name: "suffix foo pos 2",
			fields: fields{
				operator: "suffix",
				part:     0,
				arg:      "fOo",
			},
			arg: [][]byte{
				[]byte("hello world FoO"),
			},
			want: true,
		},
		{
			name: "suffix foo neg",
			fields: fields{
				operator: "suffix",
				part:     0,
				arg:      "foo",
			},
			arg: [][]byte{
				[]byte("hello foo world"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := NewConfig()
			conf.Type = "content"
			conf.Content.Operator = tt.fields.operator
			conf.Content.Part = tt.fields.part
			conf.Content.Arg = tt.fields.arg

			c, err := NewContent(conf, nil, testLog, testMet)
			if err != nil {
				t.Error(err)
				return
			}
			if got := c.Check(types.NewMessage(tt.arg)); got != tt.want {
				t.Errorf("Content.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentBadOperator(t *testing.T) {
	testLog := log.NewLogger(os.Stdout, log.LoggerConfig{LogLevel: "NONE"})
	testMet := metrics.DudType{}

	conf := NewConfig()
	conf.Type = "content"
	conf.Content.Operator = "NOT_EXIST"

	_, err := NewContent(conf, nil, testLog, testMet)
	if err == nil {
		t.Error("expected error from bad operator")
	}
}
