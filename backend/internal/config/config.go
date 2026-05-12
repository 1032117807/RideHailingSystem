package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ServerAddress  string
	MySQLDSN       string
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	RedisKeyPrefix string
	JWTSecret      string
	JWTExpireTime  time.Duration
	CodeTTL        time.Duration

	OpenAIAPIKey  string
	OpenAIBaseURL string
	OpenAIModel   string
	OpenAITimeout time.Duration

	EmbeddingAPIKey       string
	EmbeddingBaseURL      string
	EmbeddingModel        string
	RerankAPIKey          string
	RerankBaseURL         string
	RerankModel           string
	KnowledgeMemoryDir    string
	KnowledgeChunkSize    int
	KnowledgeChunkOverlap int
	KnowledgeRecallLimit  int
	KnowledgeTopK         int

	EnableMockPayment    bool
	AIRateLimitPrefix    string
	AIPassengerChatLimit int
	AIDriverDraftLimit   int
	AIRateLimitWindow    time.Duration
}

func Load() *Config {
	return &Config{
		ServerAddress:         getEnv("SERVER_ADDR", ":8080"),
		MySQLDSN:              getEnv("MYSQL_DSN", "root:hyh@tcp(127.0.0.1:3306)/ridehailing_demo?charset=utf8mb4&parseTime=True&loc=Local"),
		RedisAddr:             getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:         getEnv("REDIS_PASSWORD", ""),
		RedisDB:               getEnvInt("REDIS_DB", 0),
		RedisKeyPrefix:        getEnv("REDIS_KEY_PREFIX", "auth:email_code:"),
		JWTSecret:             getEnv("JWT_SECRET", "replace-with-your-secret"),
		JWTExpireTime:         time.Duration(getEnvInt("JWT_EXPIRE_HOURS", 168)) * time.Hour,
		CodeTTL:               time.Duration(getEnvInt("EMAIL_CODE_TTL_SECONDS", 300)) * time.Second,
		OpenAIAPIKey:          getEnv("OPENAI_API_KEY", "sk-d174a51266e24d95aee858d18d3be2cf"),
		OpenAIBaseURL:         getEnv("OPENAI_BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1"),
		OpenAIModel:           getEnv("OPENAI_MODEL", "qwen-plus"),
		OpenAITimeout:         time.Duration(getEnvInt("OPENAI_TIMEOUT_SECONDS", 20)) * time.Second,
		EmbeddingAPIKey:       getEnv("EMBEDDING_API_KEY", getEnv("OPENAI_API_KEY", "sk-d174a51266e24d95aee858d18d3be2cf")),
		EmbeddingBaseURL:      getEnv("EMBEDDING_BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1"),
		EmbeddingModel:        getEnv("EMBEDDING_MODEL", "text-embedding-v3"),
		RerankAPIKey:          getEnv("RERANK_API_KEY", getEnv("OPENAI_API_KEY", "sk-d174a51266e24d95aee858d18d3be2cf")),
		RerankBaseURL:         getEnv("RERANK_BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1"),
		RerankModel:           getEnv("RERANK_MODEL", "rerank-v1"),
		KnowledgeMemoryDir:    getEnv("KNOWLEDGE_MEMORY_DIR", "E:\\ZHE\\WrokSpace\\AI\\RideHailingSystem\\backend\\memory"),
		KnowledgeChunkSize:    getEnvInt("KNOWLEDGE_CHUNK_SIZE", 500),
		KnowledgeChunkOverlap: getEnvInt("KNOWLEDGE_CHUNK_OVERLAP", 80),
		KnowledgeRecallLimit:  getEnvInt("KNOWLEDGE_RECALL_LIMIT", 60),
		KnowledgeTopK:         getEnvInt("KNOWLEDGE_TOP_K", 4),
		EnableMockPayment:     getEnvBool("ENABLE_MOCK_PAYMENT", false),
		AIRateLimitPrefix:     getEnv("AI_RATE_LIMIT_PREFIX", "ai:rate:"),
		AIPassengerChatLimit:  getEnvInt("AI_PASSENGER_CHAT_LIMIT", 12),
		AIDriverDraftLimit:    getEnvInt("AI_DRIVER_DRAFT_LIMIT", 6),
		AIRateLimitWindow:     time.Duration(getEnvInt("AI_RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}

func getEnvBool(key string, fallback bool) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if v == "" {
		return fallback
	}

	switch v {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
