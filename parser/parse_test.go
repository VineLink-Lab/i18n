package parser

import (
	"io/fs"
	"os"
	"testing"

	"github.com/VineLink-Lab/i18n/utils"
	"golang.org/x/text/language"
)

func TestParser_PreParse(t *testing.T) {
	type fields struct {
		directory string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test Case 1: Input is a Directory",
			fields: fields{
				directory: "../example/message",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				directoryPath: tt.fields.directory,
				directory:     os.DirFS(tt.fields.directory),
				bundleDir:     make(map[string]string),
				languages:     utils.NewSet[language.Tag](),
				contents:      make(map[string]BundleContent),
			}
			if err := p.PreParse(); (err != nil) != tt.wantErr {
				t.Errorf("PreParse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParser_ParseAllBundles(t *testing.T) {
	type fields struct {
		directory string
		bundleDir map[string]string
		languages utils.Set[language.Tag]
		contents  map[string]BundleContent
	}
	languages := utils.NewSet[language.Tag]()
	languages.Add(language.English)
	languages.Add(language.Japanese)
	languages.Add(language.French)
	languages.Add(language.Spanish)
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test Case 1: Parse All Bundles",
			fields: fields{
				directory: "../example/message",
				bundleDir: map[string]string{
					"main":            "../example/message",
					"with_parameters": "../example/message/with_parameters",
				},
				languages: languages,
				contents:  make(map[string]BundleContent),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				directoryPath: tt.fields.directory,
				directory:     os.DirFS(tt.fields.directory),
				bundleDir:     tt.fields.bundleDir,
				languages:     tt.fields.languages,
				contents:      tt.fields.contents,
			}
			if err := p.ParseAllBundles(); (err != nil) != tt.wantErr {
				t.Errorf("ParseAllBundles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParser_MatchLanguageTag(t *testing.T) {
	type fields struct {
		directoryPath string
		directory     fs.FS
		bundleDir     map[string]string
		languages     utils.Set[language.Tag]
		contents      map[string]BundleContent
	}
	type args struct {
		lang language.Tag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   language.Tag
	}{
		{
			name: "Test Case 1: Exact Match",
			fields: fields{
				languages: func() utils.Set[language.Tag] {
					s := utils.NewSet[language.Tag]()
					s.Add(language.English)
					s.Add(language.Japanese)
					return s
				}(),
			},
			args: args{
				lang: language.Japanese,
			},
			want: language.Japanese,
		},
		{
			name: "Test Case 2: Fallback Match",
			fields: fields{
				languages: func() utils.Set[language.Tag] {
					s := utils.NewSet[language.Tag]()
					s.Add(language.English)
					s.Add(language.Japanese)
					return s
				}(),
			},
			args: args{
				lang: language.AmericanEnglish,
			},
			want: language.English,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				directoryPath: tt.fields.directoryPath,
				directory:     tt.fields.directory,
				bundleDir:     tt.fields.bundleDir,
				languages:     tt.fields.languages,
				contents:      tt.fields.contents,
			}
			if got := p.MatchLanguage(tt.args.lang); !(got.String() == tt.want.String()) {
				t.Errorf("MatchLanguageTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
