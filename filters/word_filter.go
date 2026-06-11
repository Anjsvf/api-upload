package filters
import (
	"fmt"
	"strings"
	"unicode"
)
//LISTA DE PALAVRAS PROIBIDAS NA LEGENDA E TITULO
var bannedWords = []string{
	"idiota", "imbecil","macaco", "preto", "estupido", "lixo", "merda", "porra", "caralho",
	"viado", "arrombado", "buceta", "xoxota", "puta", "prostituta",
	"vagabunda", "safada","buceta", "cachorra", "corno",
	"nazista", "nazi", "fascista", "racista",
	"pedofilia", "pedofilo",
	"cp", "childporn", "lolita","punheteiros", "punheteiro",
	"vendo drogas", "cocaina", "heroina",
	"suicidio", "se matar",
	"lula", "luladrão", "9dedos", "nove dedos", "petista", "ptista",
	"pt", "faz o l", "faz o L", "fazol",
	"lulismo", "lulista",
	"bolsonaro", "bolsonarista", "mito", "capitão", "bozo",
	"bolsonazi", "bolsominion", "golpista",
	"ditador", "ditadura", "golpe", "golpista", "terrorista",
	"comunista", "comuna", "esquerdopata", "direitopata",
	"petralha", "coxinha", "mortadela", "globolixo", "globista",
	"corrupto", "ladrão", "ladrao", "roubo", "mamata", "mamadeira",
	"vagabundo", "vagabunda", "traidor", "vendido",
}

func normalize(s string) string {
	replacements := map[rune]rune{
		'á': 'a', 'à': 'a', 'ã': 'a', 'â': 'a',
		'é': 'e', 'è': 'e', 'ê': 'e',
		'í': 'i', 'ì': 'i', 'î': 'i',
		'ó': 'o', 'ò': 'o', 'õ': 'o', 'ô': 'o',
		'ú': 'u', 'ù': 'u', 'û': 'u',
		'ç': 'c', 'ñ': 'n',
	}

	var result strings.Builder
	for _, r := range strings.ToLower(s) {
		if mapped, ok := replacements[r]; ok {
			result.WriteRune(mapped)
		} else if unicode.IsLetter(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		} else {
			result.WriteRune(' ')
		}
	}
	return result.String()
}


func ContainsBannedWord(text string) string {
	normalized := normalize(text)
	for _, word := range bannedWords {
		if strings.Contains(normalized, normalize(word)) {
			return word
		}
	}
	return ""
}


func CheckFields(title, caption string) error {
	if ContainsBannedWord(title) != "" {
		return fmt.Errorf("o título contém conteúdo não permitido")
	}
	if ContainsBannedWord(caption) != "" {
		return fmt.Errorf("a legenda contém conteúdo não permitido")
	}
	return nil
}