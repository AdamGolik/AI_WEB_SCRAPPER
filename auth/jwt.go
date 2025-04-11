package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"your_project/models"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"gorm.io/gorm"
)

type ScrapeInput struct {
	Name     string   `json:"name" binding:"required"`
	PriceMin *float32 `json:"pricemin"`
	PriceMax *float32 `json:"pricemax"`
	URLs     []string `json:"urls"`
}

type ScrapedResult struct {
	Title string  `json:"title"`
	Price float32 `json:"price"`
	URL   string  `json:"url"`
}

// AI response structure
type AIResponse struct {
	Answer     string `json:"answer"`
	BestResult string `json:"best_result"`
}

func ScrapeAndSaveHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ScrapeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Brak użytkownika"})
			return
		}
		userID := userIDVal.(int16)

		// Domyślne źródła
		urls := append(input.URLs, generateDefaultSearchURLs(input.Name)...)

		var results []ScrapedResult
		for _, url := range urls {
			col := colly.NewCollector()
			col.OnHTML("a", func(e *colly.HTMLElement) {
				title := e.Text
				href := e.Request.AbsoluteURL(e.Attr("href"))
				if title != "" && strings.Contains(href, "http") {
					results = append(results, ScrapedResult{Title: title, URL: href})
				}
			})
			col.Visit(url)
		}

		// Filtrowanie po cenie - (w praktyce trzeba scrapować cenę i parsować float z tekstu, tu uproszczenie)
		filteredResults := results // <- placeholder, tu można dodać filtrowanie

		// JSON encode do zapisania w bazie
		urlsJSON, _ := json.Marshal(urls)

		req := models.Request{
			UserID:   uint(userID),
			Name:     input.Name,
			PriceMin: input.PriceMin,
			PriceMax: input.PriceMax,
			URLs:     urlsJSON,
		}
		db.Create(&req)

		// Wysłanie danych do AI
		aiResp, err := getAIInsight(input.Name, filteredResults)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd AI"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results":     filteredResults,
			"ai_answer":   aiResp.Answer,
			"ai_best_hit": aiResp.BestResult,
		})
	}
}

func getAIInsight(topic string, results []ScrapedResult) (AIResponse, error) {
	var aiResp AIResponse

	// Składamy prompt
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Temat: %s\nOto lista wyników:\n", topic))
	for _, r := range results {
		b.WriteString(fmt.Sprintf("- %s: %.2f zł (%s)\n", r.Title, r.Price, r.URL))
	}
	b.WriteString("Oceń, podsumuj wyniki i wskaż najlepszy.")

	reqBody := map[string]string{
		"prompt": b.String(),
	}
	reqJSON, _ := json.Marshal(reqBody)

	resp, err := http.Post("https://api.llama4free.com/llama", "application/json", strings.NewReader(string(reqJSON)))
	if err != nil {
		return aiResp, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return aiResp, err
	}

	return aiResp, nil
}

func generateDefaultSearchURLs(query string) []string {
	q := strings.ReplaceAll(query, " ", "+")
	return []string{
		fmt.Sprintf("https://allegro.pl/listing?string=%s", q),
		fmt.Sprintf("https://www.olx.pl/oferty/q-%s", q),
		fmt.Sprintf("https://www.mediaexpert.pl/search?query%5Bmenu_item%5D=&query%5Bquerystring%5D=%s", q),
	}
}
