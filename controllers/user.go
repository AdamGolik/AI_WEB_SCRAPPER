package controllers

import (
	"AI_WEB_SCRAPPER/initlizers"
	"AI_WEB_SCRAPPER/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var body struct {
		Name            string `json:"name"`
		Lastname        string `json:"lastname"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"passwordconfirm"`
	}

	if err := c.Bind(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if body.Password != body.PasswordConfirm {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     body.Name,
		Lastname: body.Lastname,
		Email:    body.Email,
		Password: string(hash),
	}

	result := initlizers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Login(c *gin.Context) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	result := initlizers.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Typy
type ScrapeResult struct {
	Title string  `json:"title"`
	URL   string  `json:"url"`
	Price float32 `json:"price"`
	Store string  `json:"store"`
}

type LlamaMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type LlamaRequest struct {
	Model    string         `json:"model"`
	Messages []LlamaMessage `json:"messages"`
}

type LlamaChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type LlamaResponse struct {
	Choices []LlamaChoice `json:"choices"`
}

func ScrapeAmazon(keyword string) []ScrapeResult {
	var results []ScrapeResult
	c := colly.NewCollector()

	c.OnHTML(".s-result-item", func(e *colly.HTMLElement) {
		title := e.ChildText("h2 a span")
		priceText := e.ChildText(".a-price-whole")

		price := parsePrice(priceText)
		if title != "" && price > 0 {
			results = append(results, ScrapeResult{
				Title: title,
				URL:   "https://www.amazon.com" + e.ChildAttr("h2 a", "href"),
				Price: price,
				Store: "Amazon",
			})
		}
	})

	searchURL := "https://www.amazon.com/s?k=" + strings.ReplaceAll(keyword, " ", "+")
	c.Visit(searchURL)

	return results
}

func AskLlama(bestFrom []ScrapeResult) string {
	var textBuilder strings.Builder
	textBuilder.WriteString("Mam następujące oferty:\n")
	for _, o := range bestFrom {
		textBuilder.WriteString(fmt.Sprintf("- %s za %.2f zł (%s) %s\n", o.Title, o.Price, o.Store, o.URL))
	}
	textBuilder.WriteString("Która oferta jest najlepsza i dlaczego?")

	payload := LlamaRequest{
		Model: "meta-llama/llama-4-scout:free",
		Messages: []LlamaMessage{
			{
				Role: "user",
				Content: []map[string]string{
					{
						"type": "text",
						"text": textBuilder.String(),
					},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENROUTER_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Request error:", err)
		return "AI error"
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var llamaResp LlamaResponse
	json.Unmarshal(body, &llamaResp)

	if len(llamaResp.Choices) > 0 {
		return llamaResp.Choices[0].Message.Content
	}

	return "Brak odpowiedzi od AI"
}

func parsePrice(s string) float32 {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	var price float32
	fmt.Sscanf(s, "%f", &price)
	return price
}

func HandleScrapeAndAI(keyword string) {
	offers := ScrapeAmazon(keyword)
	aiResponse := AskLlama(offers)

	fmt.Println("Oferty znalezione:")
	for _, o := range offers {
		fmt.Printf("- %s | %.2f | %s\n", o.Title, o.Price, o.Store)
	}

	fmt.Println("\n🤖 Odpowiedź AI:")
	fmt.Println(aiResponse)
}
func Analyze(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var reqBody struct {
		Name     string    `json:"name" binding:"required"`
		PriceMin *float32  `json:"pricemin"`
		PriceMax *float32  `json:"pricemax"`
		URLs     []string  `json:"urls"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Domyślne URL, jeśli brak
	if len(reqBody.URLs) == 0 {
		reqBody.URLs = []string{
			"https://www.amazon.com/s?k=" + strings.ReplaceAll(reqBody.Name, " ", "+"),
		}
	}

	// Scraping (tu tylko Amazon — można rozszerzyć)
	var results []ScrapeResult
	results = append(results, ScrapeAmazon(reqBody.Name)...)

	// AI analiza
	aiMessage := AskLlama(results)

	// Zapisz zapytanie do DB
	request := models.Request{
		UserID:   userID,
		Name:     reqBody.Name,
		PriceMin: reqBody.PriceMin,
		PriceMax: reqBody.PriceMax,
		URLs:     reqBody.URLs,
	}
	if err := initlizers.DB.Create(&request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB save failed"})
		return
	}

	// JSON response
	c.JSON(http.StatusOK, gin.H{
		"results":     results,
		"ai_response": aiMessage,
	})
}

