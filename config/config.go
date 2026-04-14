package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadModelConfig() (key, model, baseURL, systemPrompt string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	key = viper.GetString("openaiProvider.key")
	if key == "" {
		log.Println("key is empty, you may not be able to send messages to AI provider")
	}
	model = viper.GetString("openaiProvider.model")
	if model == "" {
		model = "gpt-5.2"
	}
	baseURL = viper.GetString("openaiProvider.baseURL")
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	systemPrompt = viper.GetString("prompt.system")
	return
}

func LoadServerConfig() string {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	port := viper.GetString("server.port")
	if port == "" {
		port = "8077"
	}
	return ":" + port
}

func LoadNapcatConfig() (apiBaseURL, expectedToken, authToken string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	apiBaseURL = viper.GetString("napcat.apiBaseURL")
	if apiBaseURL == "" {
		apiBaseURL = "http://localhost:3000"
	}
	expectedToken = viper.GetString("napcat.expectedToken")
	authToken = viper.GetString("napcat.authToken")
	if authToken == "" {
		log.Println("authToken is empty, you may not be able to send messages to napcat")
	}
	return
}

func LoadJudgeAtConfig() bool {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	return viper.GetBool("judgeAt.enable")
}
