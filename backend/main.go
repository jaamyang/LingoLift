package main

import (
	"bytes"
	// "embed"
	"encoding/json"
	"fmt"
	"io"
	// "io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExplainWordRequest is the request body for the /api/explain-word endpoint
type ExplainWordRequest struct {
	Word               string                 `json:"word"`
	Context            string                 `json:"context"`
	AudioLanguage      string                 `json:"audio_language"`
	UserNativeLanguage string                 `json:"user_native_language"`
	LlmSettings        map[string]interface{} `json:"llm_settings"`
}

// TranslateSentenceRequest is the request body for the /api/translate-sentence endpoint
type TranslateSentenceRequest struct {
	Sentence       string                 `json:"sentence"`
	TargetLanguage string                 `json:"target_language"`
	LlmSettings    map[string]interface{} `json:"llm_settings"`
}

// LLMResponse is the expected structure of the response from the LLM provider
type LLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// var frontendFS embed.FS

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	r.GET("/api/proxy", proxyHandler)
	r.POST("/api/explain-word", explainWordHandler)
	r.POST("/api/translate-sentence", translateSentenceHandler)

	// frontendDist, err := fs.Sub(frontendFS, "frontend/dist")
	// if err != nil {
	// 	log.Fatal("Failed to get frontend/dist subdirectory:", err)
	// }

	// Serve frontend static files
	r.StaticFS("/assets", http.Dir("frontend/dist/assets"))
	r.StaticFS("/models", http.Dir("frontend/dist/models"))
	r.StaticFS("/wasm", http.Dir("frontend/dist/wasm"))
  
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.File(filepath.Join("frontend/dist", "index.html"))
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func proxyHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required."})
		return
	}

	log.Printf("Proxying request for url: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to connect to the provided URL: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	io.Copy(c.Writer, resp.Body)
}

func explainWordHandler(c *gin.Context) {
	var req ExplainWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Explain word request: Word='%s', AudioLang='%s', NativeLang='%s'", req.Word, req.AudioLanguage, req.UserNativeLanguage)

	audioLangExplanation, err := fetchExplanationForLang(req.Word, req.Context, req.AudioLanguage, req.LlmSettings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get explanation from LLM service."})
		return
	}

	nativeLangExplanation := audioLangExplanation
	if req.AudioLanguage != req.UserNativeLanguage {
		nativeLangExplanation, err = fetchExplanationForLang(req.Word, req.Context, req.UserNativeLanguage, req.LlmSettings)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get explanation from LLM service."})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"audio_language_explanation":  audioLangExplanation,
		"native_language_explanation": nativeLangExplanation,
	})
}

func translateSentenceHandler(c *gin.Context) {
	var req TranslateSentenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Translate sentence request: ToLang='%s', Sentence='%.50s...'", req.TargetLanguage, req.Sentence)

	prompt := fmt.Sprintf(`Translate to %s. Provide only the translation.

Text: "%s"

Translation to %s:`, req.TargetLanguage, req.Sentence, req.TargetLanguage)

	translatedText, err := fetchLLMResponse(prompt, req.LlmSettings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to translate sentence via LLM service."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"translated_sentence": translatedText})
}

func fetchExplanationForLang(word, context, language string, llmSettings map[string]interface{}) (string, error) {
	prompt := fmt.Sprintf(`Reply in %s. Explain the word "%s" in 1-2 simple sentences based on the context. Use simple, everyday vocabulary.

Word: "%s"
Context: "%s"

Explanation (in %s):`, language, word, word, context, language)

	return fetchLLMResponse(prompt, llmSettings)
}

func fetchLLMResponse(prompt string, llmSettings map[string]interface{}) (string, error) {
	defaultAPIURL := "https://api-inference.modelscope.cn/v1/chat/completions"
	defaultModel := "qwen/qwen3-30b-a3b"

	apiURL := os.Getenv("LLM_BASE_URL")
	if apiURL == "" {
		apiURL = defaultAPIURL
	}

	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = defaultModel
	}

	apiKey := os.Getenv("LLM_API_KEY")

	if s, ok := llmSettings["baseUrl"].(string); ok {
		apiURL = s
	}
	if m, ok := llmSettings["model"].(string); ok {
		model = m
	}
	if k, ok := llmSettings["apiKey"].(string); ok {
		apiKey = k
	}

	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are an experienced language tutor, you are good at translate and explain the words and phrases in context"},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
    "enable_thinking": false,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM service returned an error: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var llmResp LLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return "", err
	}

	if len(llmResp.Choices) > 0 {
		return llmResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("could not parse LLM response")
} 