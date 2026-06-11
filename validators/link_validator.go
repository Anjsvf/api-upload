package validators

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)



var (
	// Aceita: youtube.com/watch?v=ID, youtu.be/ID, youtube.com/shorts/ID
	youtubeRegex = regexp.MustCompile(
		`(?i)^(https?://)?(www\.)?(youtube\.com/(watch\?v=|shorts/)|youtu\.be/)([a-zA-Z0-9_-]{11})`)

	// Aceita: chat.whatsapp.com/INVITE_CODE
	whatsappRegex = regexp.MustCompile(
		`(?i)^(https?://)?chat\.whatsapp\.com/[A-Za-z0-9]{10,}$`)
)




func ExtractYouTubeID(rawURL string) (string, error) {
	match := youtubeRegex.FindStringSubmatch(rawURL)
	if len(match) < 6 {
		return "", fmt.Errorf("URL do YouTube inválida")
	}
	return match[5], nil
}


func ValidateYouTubeFormat(rawURL string) error {
	_, err := ExtractYouTubeID(rawURL)
	return err
}


type YouTubeAPIResponse struct {
	PageInfo struct {
		TotalResults int `json:"totalResults"`
	} `json:"pageInfo"`
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
}

/
func ValidateYouTubeExists(videoID, apiKey string) error {
	if apiKey == "" {
	
		return nil
	}

	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/youtube/v3/videos?part=id&id=%s&key=%s",
		url.QueryEscape(videoID),
		apiKey,
	)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return fmt.Errorf("erro ao consultar YouTube API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("YouTube API retornou status %d", resp.StatusCode)
	}

	var result YouTubeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("erro ao processar resposta da YouTube API: %v", err)
	}

	if result.PageInfo.TotalResults == 0 {
		return fmt.Errorf("vídeo não encontrado ou indisponível no YouTube")
	}

	return nil
}


func ValidateYouTubeLink(rawURL, apiKey string) (string, error) {
	// Normaliza o link
	rawURL = strings.TrimSpace(rawURL)
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "https://" + rawURL
	}

	
	videoID, err := ExtractYouTubeID(rawURL)
	if err != nil {
		return "", fmt.Errorf("link do YouTube inválido: precisa ser youtube.com/watch?v=... ou youtu.be/...")
	}


	if err := ValidateYouTubeExists(videoID, apiKey); err != nil {
		return "", err
	}

	// Retorna o link normalizado
	return fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID), nil
}




func ValidateWhatsAppLink(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "https://" + rawURL
	}

	if !whatsappRegex.MatchString(rawURL) {
		return "", fmt.Errorf("link do WhatsApp inválido: precisa ser chat.whatsapp.com/CODIGO")
	}

	
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("link malformado")
	}
	parsed.Scheme = "https"

	return parsed.String(), nil
}