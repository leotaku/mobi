package mobi

import "golang.org/x/text/language"

var matcher language.Matcher

func matchLocale(lang language.Tag) uint32 {
	match, _, _ := matcher.Match(lang)
	return uint32(localeCodeMap[match])
}

var faroese = language.MustParse("fo")
var romansh = language.MustParse("rm")
var belarusian = language.MustParse("be")
var azerbaijani = language.MustParse("az")
var basque = language.MustParse("eu")
var sorbian = language.MustParse("wen")
var sotho = language.MustParse("st")
var tsonga = language.MustParse("ts")
var tswana = language.MustParse("tn")
var maltese = language.MustParse("mt")
var xhosa = language.MustParse("xh")
var sami = language.MustParse("smi")
var kazakh = language.MustParse("kk")
var tatar = language.MustParse("tt")
var oriya = language.MustParse("or")
var assamese = language.MustParse("as")
var sanskrit = language.MustParse("sa")
var konkani = language.MustParse("kok")

var localeCodeMap = map[language.Tag]uint16{
	language.Und:        0,
	language.Arabic:     1,
	language.Bulgarian:  2,
	language.Catalan:    3,
	language.Chinese:    4,
	language.Czech:      5,
	language.Danish:     6,
	language.German:     7,
	language.Greek:      8,
	language.English:    9,
	language.Spanish:    10,
	language.Finnish:    11,
	language.French:     12,
	language.Hebrew:     13,
	language.Hungarian:  14,
	language.Icelandic:  15,
	language.Italian:    16,
	language.Japanese:   17,
	language.Korean:     18,
	language.Dutch:      19,
	language.Norwegian:  20,
	language.Polish:     21,
	language.Portuguese: 22,
	romansh:             23,
	language.Romanian:   24,
	language.Russian:    25,
	language.Croatian:   26,
	language.Serbian:    26,
	language.Slovak:     27,
	language.Albanian:   28,
	language.Swedish:    29,
	language.Thai:       30,
	language.Turkish:    31,
	language.Urdu:       32,
	language.Indonesian: 33,
	language.Ukrainian:  34,
	belarusian:          35,
	language.Slovenian:  36,
	language.Estonian:   37,
	language.Latvian:    38,
	language.Lithuanian: 39,
	language.Persian:    41,
	language.Vietnamese: 42,
	language.Armenian:   43,
	azerbaijani:         44,
	basque:              45,
	sorbian:             46,
	language.Macedonian: 47,
	sotho:               48,
	tsonga:              49,
	tswana:              50,
	xhosa:               52,
	language.Zulu:       53,
	language.Afrikaans:  54,
	language.Georgian:   55,
	faroese:             56,
	language.Hindi:      57,
	maltese:             58,
	sami:                59,
	language.Malay:      62,
	kazakh:              63,
	language.Swahili:    65,
	language.Uzbek:      67,
	tatar:               68,
	language.Bengali:    69,
	language.Punjabi:    70,
	language.Gujarati:   71,
	oriya:               72,
	language.Tamil:      73,
	language.Telugu:     74,
	language.Kannada:    75,
	language.Malayalam:  76,
	assamese:            77,
	language.Marathi:    78,
	sanskrit:            79,
	konkani:             87,
	language.Nepali:     97,
}

// A list of locales supported by MOBI, intended for use with Go
// language matching facilities.
//
// BUG: While the MOBI language format supports regions, this
// implementation does not and explicitly removes region information
// from the locale stored in the generated MOBI file.

var SupportedLocales []language.Tag

func init() {
	for key := range localeCodeMap {
		SupportedLocales = append(SupportedLocales, key)
	}
	matcher = language.NewMatcher(SupportedLocales)
}
