package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/joho/godotenv"
	//"fmt"
	"strings"
	"io/ioutil"
	"encoding/json"
	//"os"
)

func main() {

	//Openai chat endpoint 
	url := "https://api.openai.com/v1/chat/completions"

	envMap, mapErr := godotenv.Read(".env")
	if mapErr != nil { return }

	apiKey := envMap["API_KEY"]
	
	//Prepare struct for capturing json from post request
	type User struct {
		Message  string `json:"message"`
	}

	//Prepare struct for getting json response from chatgpt
	type OpenaiResp struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Model   string `json:"model"`
		Usage   struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
			Index        int    `json:"index"`
		} `json:"choices"`
	}

	//Prepare struct for getting json response from chatgpt
	type OpenaiImgResp struct {
			Created int `json:"created"`

			Data  []struct {
				URL string `json:"url"`
			} `json:"data"`
		
	}


	e := echo.New()
	e.POST("/chat", func(c echo.Context) error {

		u := new(User)
		if err := c.Bind(u); err != nil {return err}

		if u.Message == "" {
			return c.JSON(http.StatusCreated, "Empty message !")
		}
		
		payload := strings.NewReader(`{"model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "` + u.Message + `"}]}`)

		client := &http.Client {}

		req, err := http.NewRequest("POST", url, payload)
		if err != nil { return err}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+ apiKey )

		res, err := client.Do(req)
		if err != nil {return err}
		defer res.Body.Close()
	  
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {return err}

		openaiResp := OpenaiResp{}
		json.Unmarshal([]byte(body), &openaiResp)

		return c.JSON(http.StatusCreated, openaiResp.Choices[0].Message)

	})

	e.POST("/image", func(c echo.Context) error {
		urlImage := "https://api.openai.com/v1/images/generations"

		u := new(User)
		if err := c.Bind(u); err != nil {return err}

		if u.Message == "" {
			return c.JSON(http.StatusCreated, "Empty message !")
		}
		
		payload := strings.NewReader(`{"prompt": "` + u.Message + `", "n": 1, "size": "256x256"}`)

		client := &http.Client {}

		req, err := http.NewRequest("POST", urlImage, payload)
		if err != nil { return err}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+ apiKey )


		res, err := client.Do(req)
		if err != nil {return err}
		defer res.Body.Close()
	  
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {return err}

		openaiImgResp := OpenaiImgResp{}
		json.Unmarshal([]byte(body), &openaiImgResp)


		if openaiImgResp.Data[0].URL == "" {
			return c.JSON(http.StatusCreated, "Empty message !")
		}

		return c.JSON(http.StatusCreated, openaiImgResp.Data[0].URL)

	})




	e.Logger.Fatal(e.Start(":1323"))
}
