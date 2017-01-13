package main

import (
	"fmt"
	"github.com/geobe/gostip/go/transcription"
)
var namesRu = [...]string{
	"Баитакова Аиана",
	"Нурманбетов Калыс",
	"Молдалиев Эрлан",
	"Токтосунов Джаныш",
	"Токтосунов Баиыш",
	"Акылбекова Джылдыз",
	"Асанбекова Бегаиым",
	"Исламов Салават",
	"Мухамедов Илиаз",
	"Рысбеков Максат",
	"Мелисбеков Талгат",
	"Намазов Дастан",
	"Ташкулова Диана",
	"Сарбаева Аида",
	"Асанова Аидаи",
	"Бегматов Расул",
	"Беишэкеева Гуланда",
	"Джумабаев Арген",
	"Касымбекова Аисанат",
	"Инамов Бахриддин",
	"Джоломанов Азамат",
	"Кашкарлыков Канат",
	"Шаршенбекова Бурул",
	"Талант улу Нурлан",
	"Курманбекова Аизат",
	"Омуров Алишер",
	"Кыпчаков Азим",
	"Токторбекова Джанылаи",
	"Акматова Джибек",
	"Иусупова Айджан",
	"Джунушев Арслан",
	"Мыизамов Аджибек",
	"Талант улу Мирлан",
	"Кошалы улу Нурбек",
	"Русланова Нурсель",
	"Мухамеджанова Каухар",
	"Дё Оскар",
	"Фатулаев Рустам",
	"Касымбеков Султанмырза",
	"Кудаибергенов Улукбек",
}

var namesLat = [...]string{
	"Baitakova Ayana",
	"Nurmanbetov Kalys",
	"Moldaliev Erlan",
	"Toktosunov Zhanysh",
	"Toktosunov Baiysh",
	"Akylbekova Zhyldyz",
	"Asanbekova Begaiym",
	"Islamov Salavat",
	"Muhammedov Iliaz",
	"Rysbekov Maksat",
	"Melisbekov Talgat",
	"Namazov Dastan",
	"Tashkulova Diana",
	"Sarbaeva Aida",
	"Asanova Aidai",
	"Begmatov Rasul",
	"Beishekeeva Gulanda",
	"Zhumabaev Argen",
	"Kasymbekova Aisanat",
	"Inamov Bachriddin",
	"Zholomanov Azamat",
	"Kashkarlykov Kanat",
	"Sharshenbekova Burul",
	"Talant uulu Nurlan",
	"Kurmanbekova Aizat",
	"Omurov Alisher",
	"Kypchakov Azim",
	"Toktorbekova Zhanylai",
	"Akmatova Zhibek",
	"Yusupova Aizhan",
	"Dzhunushev Arslan",
	"Mijsamov Azhibek",
	"Talant uulu Mirlan",
	"Koshaly uulu Nurbek",
	"Ruslanova Nursel",
	"Muhamedzhanova Kauhar",
	"De Oskar",
	"Fatullaev Rustam",
	"Kasymbekov Sultanmyrza",
	"Kudaibergenov Ulukbek",
}
func main() {
	//var teststrings = [...]string{
	//	"Verrückte im Taxi und heulende Wölfe in Sachsen",
	//	"Franz jagt im komplett verwahrlosten Taxi quer durch Bayern",
	//	"Daisy Duck ist die ewige Verlobte von Donald Duck, die allerdings auch oft mit dessen Vetter Gustav Gans anbändelt. Ihr Debüt fand offiziell 1937 im Kurzfilm Don Donald statt, obwohl dort eigentlich nur eine Vorläuferin namens Donna Duck auftritt. Dennoch wurde der 9. Januar 2007 als ihr 70. Geburtstag gefeiert, wobei in der Presse teilweise kritisch damit umgegangen wurde. ",
	//	"«Си́мпсоны» (англ. The Simpsons) — самый длинный мультсериал в истории американского телевидения, длящийся 28 сезонов и продлённый до 30-го. Первая мини-серия Good night была показана в «Шоу Трейси Ульман» 19 апреля 1987 года. Демонстрация полноценных серий мультсериала началась 17 декабря 1989 года на канале FOX. Выпускается «Gracie Films» для кинокомпании «20th Century Fox». Мультсериал в настоящее время известен телезрителям более чем в 100 странах",
	//}
	//
	//for i, s := range teststrings {
	//	kyr := transcription.IsKyrgyz(s)
	//	fmt.Printf("teststring %d is Kyrgyz? %t\n%s\n", i, kyr, s)
	//	fmt.Println(transcription.Transcribe(s))
	//}
	//idru := "Иван Денисович"
	//idde := "Iwan Denissowitsch"
	//
	//const repeat = 5000
	//
	//start := time.Now().UnixNano()
	//for i := 0; i < repeat; i++ {
	//	transcription.Transcribe(idru)
	//	transcription.Transcribe(idde)
	//}
	//t := (time.Now().UnixNano() - start)
	//c := 2 * repeat
	//t1 := t / (2 * repeat)
	//fmt.Printf("%d conversions took %d nsec, %d nsec/conversion\n", c, t, t1)

	cnt := 0
	for i, v := range namesRu {
		tr := transcription.Transcribe(v)
		var pr, si string
		if tr == namesLat[i] {
			pr = "OK  "
			si = "=="
		} else {
			pr = "FAIL"
			si = "!="
			cnt++
		}
		fmt.Printf("%s %s: %s %s %s\n", pr, v, namesLat[i], si, tr)
	}
	fmt.Printf("%d Differences\n", cnt)
}

func compareByRune(s1, s2 string) string {
	result := ""
	runes1 := make([]rune, len(s1))
	runes2 := make([]rune, len(s2))
	for i, v := range s1 {
		runes1[i] = v
	}
	for i, v := range s2 {
		runes2[i] = v
	}
	rlen := min(len(runes1), len(runes2))
	for i := 0; i < rlen; i++ {
		if runes1[i] != runes2[i] {
			result += fmt.Sprintf("%c", runes1[i])
		}
	}
	if len(runes1) != len(runes2) {
		result += " !!!!"
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}