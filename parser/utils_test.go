package parser

import (
	"reflect"
	"testing"

	"golang.org/x/text/language"
)

func Test_parsePathLangTag(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name        string
		args        args
		wantLangTag language.Tag
	}{
		{
			name: "Test Case 1: English Language Tag",
			args: args{
				path: "en-US.json",
			},
			wantLangTag: language.AmericanEnglish,
		},
		{
			name: "Test Case 2: Spanish Language Tag",
			args: args{
				path: "es.json",
			},
			wantLangTag: language.Spanish,
		},
		{
			name: "Test Case 3: French Language Tag",
			args: args{
				path: "fr.json",
			},
			wantLangTag: language.French,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLangTag := parseFileNameLangTag(tt.args.path); !reflect.DeepEqual(gotLangTag, tt.wantLangTag) {
				t.Errorf("parseFileNameLangTag() = %v, want %v", gotLangTag, tt.wantLangTag)
			}
		})
	}
}
