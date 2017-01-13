package transcription

import (
	"regexp"
)

var fromKyrgyzTrans = [...]equiv{
	{"Мухамед", "Muhamed"},
	{"Дж", "Zh"},
	{"дж", "zh"},
	{"улу", "uulu"},
	{"ыи", "ij"},
	{"Иу", "Yu"},
	{"Аиа", "Aya"},
	{"Ай", "Ai"},
	{"А", "A"},
	{"Б", "B"},
	{"В", "V"},
	{"Г", "G"},
	{"Д", "D"},
	{"Е", "E"},
	{"Ё", "Yo"},
	{"Ж", "Zh"},
	{"З", "Z"},
	{"И", "I"},
	{"Й", "Y"},
	{"К", "K"},
	{"Л", "L"},
	{"М", "M"},
	{"Н", "N"},
	{"Ң", "Ñ"},
	{"О", "O"},
	{"Ө", "Ö"},
	{"П", "P"},
	{"Р", "R"},
	{"С", "S"},
	{"Т", "T"},
	{"У", "U"},
	{"Ү", "Ü"},
	{"Ф", "F"},
	{"Х", "Ch"},
	{"Ц", "Ts"},
	{"Ч", "Ch"},
	{"Ш", "Sh"},
	{"Щ", "Shch"},
	{"Ы", "Y"},
	{"Э", "E"},
	{"Ю", "Yu"},
	{"Я", "Ya"},
	{"б", "b"},
	{"а", "a"},
	{"в", "v"},
	{"г", "g"},
	{"д", "d"},
	{"е", "e"},
	{"ё", "yo"},
	{"ж", "zh"},
	{"з", "z"},
	{"и", "i"},
	{"й", "y"},
	{"к", "k"},
	{"л", "l"},
	{"м", "m"},
	{"́м", "m"},
	{"н", "n"},
	{"ң", "ñ"},
	{"о", "o"},
	{"ө", "ö"},
	{"п", "p"},
	{"р", "r"},
	{"с", "s"},
	{"т", "t"},
	{"у", "u"},
	{"ү", "ü"},
	{"ф", "f"},
	{"х", "ch"},
	{"ц", "ts"},
	{"ч", "ch"},
	{"ш", "sh"},
	{"щ", "shch"},
	{"ы", "y"},
	{"э", "e"},
	{"ю", "yu"},
	{"я", "ya"},
	{"Ъ", ""},
	{"ъ", ""},
	{"Ь", ""},
	{"ь", ""},
}
var toKyrgyzTrans = [...]equiv{
	{"Shch", "Щ"},
	{"shch", "щ"},
	{"Tsch", "Ч"},
	{"tsch", "ч"},
	{"vich", "вич"},
	{"Ch", "Х"},
	{"Sh", "Ш"},
	{"Sch", "Ш"},
	{"Yu", "Ю"},
	{"Ya", "Я"},
	{"kh", "х"},
	{"ts", "ц"},
	{"ch", "х"},
	{"ck", "кк"},
	{"sh", "ш"},
	{"yu", "ю"},
	{"ya", "я"},
	{"Ah", "A"},
	{"Eh", "E"},
	{"Ih", "И"},
	{"Oh", "O"},
	{"Uh", "У"},
	{"ah", "a"},
	{"eh", "e"},
	{"ih", "и"},
	{"oh", "o"},
	{"uh", "У"},
	{"Qu", "Квь"},
	{"qu", "квь"},
	{"Yo", "Ё"},
	{"yo", "ё"},
	{"Zh", "Ж"},
	{"zh", "ж"},
	{"Kh", "Х"},
	{"Ts", "Ц"},
	{"H", "Г"},
	{"h", "г"},
	{"A", "А"},
	{"B", "Б"},
	{"V", "В"},
	{"W", "Вь"}, // oder Ь ?
	{"G", "Г"},
	{"D", "Д"},
	{"E", "Е"},
	{"Z", "З"},
	{"I", "И"},
	{"J", "Й"},
	{"K", "К"},
	{"L", "Л"},
	{"N", "Н"},
	{"Ñ", "Ң"},
	{"O", "О"},
	{"Ö", "Ө"},
	{"P", "П"},
	{"R", "Р"},
	{"S", "С"},
	{"T", "Т"},
	{"U", "У"},
	{"Ü", "Ү"},
	{"F", "Ф"},
	{"X", "Кc"},
	{"Y", "Ы"},
	{"Ä", "Э"},
	{"a", "а"},
	{"b", "б"},
	{"v", "в"},
	{"w", "вь"},
	{"g", "г"},
	{"d", "д"},
	{"e", "е"},
	{"z", "з"},
	{"i", "и"},
	{"j", "й"},
	{"k", "к"},
	{"l", "л"},
	{"m", "м"},
	{"n", "н"},
	{"ñ", "ң"},
	{"o", "о"},
	{"ö", "ө"},
	{"p", "п"},
	{"r", "р"},
	{"s", "с"},
	{"t", "т"},
	{"u", "у"},
	{"ü", "ү"},
	{"f", "ф"},
	{"x", "кc"},
	{"y", "ы"},
	{"ä", "э"},
}

var kyrRegex = regexp.MustCompile("[БГДЁЖЗИЙКЛНҢӨПУФЦЧШЩЫЭЮЯбвгдёжзийклмнңөптуүфцчшщыэюя]")
var useKyrRegex = regexp.MustCompile("((ky)|(ru))[a-zA-Z,-]*;q=0.[4-9]")
var regexMapKyrgyz = makeRegexKyrgyz()
var regexMapLatin = makeRegexLatin()

func IsKyrgyz(s string) bool {
	res := kyrRegex.FindStringSubmatch(s)
	return res != nil
}

func toKyrgyz(s string) string {
	for _, v := range regexMapKyrgyz {
		s = v.regex.ReplaceAllString(s, v.replacement)
	}
	return s
}

func fromKyrgyz(s string) string {
	for _, v := range regexMapLatin {
		s = v.regex.ReplaceAllString(s, v.replacement)
	}
	return s
}

func Transcribe(s string) string {
	if IsKyrgyz(s) {
		return fromKyrgyz(s)
	} else {
		return toKyrgyz(s)
	}
}

//func transcribe(in, tx string) (out string, yes bool) {
//	if in != "" && tx == "" {
//		yes = true
//		out = Transcribe(in)
//		return
//	}
//	yes = false
//	out = tx
//	return
//}

func UsesKyrillic(languageHeader []string) bool {
	for _, v := range languageHeader {
		m := useKyrRegex.FindString(v)
		if m != "" {
			return true
		}
	}
	return false
}

type trans struct {
	regex       *regexp.Regexp
	replacement string
}

type equiv struct {
	from string
	to   string
}

func makeRegexKyrgyz() (trex []trans) {
	trex = make([]trans, len(toKyrgyzTrans))
	for i, v := range toKyrgyzTrans {
		trex[i] = trans{
			regex: regexp.MustCompile(v.from),
			replacement: v.to,
		}
	}
	return
}

func makeRegexLatin() (trex []trans) {
	trex = make([]trans, len(fromKyrgyzTrans))
	for i, v := range fromKyrgyzTrans {
		trex[i] = trans{
			regex: regexp.MustCompile(v.from),
			replacement: v.to,
		}
	}
	return
}



