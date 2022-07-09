package lookup

import "github.com/microcosm-cc/bluemonday"

type Source struct {
	Name            string
	ConstantParams  LookupParams
	TransformParams func(map[string]string) (any, error)
	DoLookup        func(SourceEnv, string, any) (any, error)
}

type SourceEnv struct {
	UserAgent string
}

type LookupParams struct {
	SourceLang *string `nt:"sourceLang" json:"sourceLang"`
	TargetLang *string `nt:"targetLang" json:"targetLang"`
}

// Sad
var pl = "pl"
var it = "it"

var StripTagsPolicy = bluemonday.StrictPolicy()

var AllowFormattedHTMLPolicy = bluemonday.NewPolicy().AllowElements("em")
