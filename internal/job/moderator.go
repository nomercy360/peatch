package job

import (
	"encoding/json"
	"fmt"
	"github.com/peatch-io/peatch/internal/db"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func getRequestBody(userProfile string) string {
	// replace newlines with spaces
	return fmt.Sprintf(`{
    "model": "gpt-4o",
    "messages": [
        {
            "role": "system",
            "content": [
                {
                    "type": "text",
                    "text": "Review user profile description and give it a score. It’s for a social network for collaboration and finding talents, so realistic and non-abstract profile is most important. Main fields are 'Title', 'First Name', 'Last Name', and 'Description'; less important is to compare Badges and Opportunities with the profile description and title, see if they correlate. The RESULT is a number from 1 to 1000, where 1 is total spam and 1000 is the most interesting profile ever seen."
                }
            ]
        },
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "Title: Project Manager First Name: Demid Last Name: Druz Description: Hello, I'm Demid Druz’! I am 29 years old, I have been doing Visual design, Art, Graffiti for more than 10 years! I work with top companies in my field and created an MOST app that I plan to 🚀 🧑🏽‍💻My portfolio - behance.net/demiddruz ☮️My app - https://apps.apple.com/app/id1566117045 Badges: Wine Lover Traveller Founder Product Designer Business Developer Visionary Photographer Graphic Designer UX Designer UI Designer Opportunities: Giving product reviews, Coaching founders, Co-founding a company, Career coaching, Beta testing new products, Advising on SEO, Advising early stage companies"
                }
            ]
        },
        {
            "role": "assistant",
            "content": [
                {
                    "type": "text",
                    "text": "734"
                }
            ]
        },
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "Title: UX/UI designer, no-code developer Readymag/Tilda First Name: Nikita Last Name: Mironov Description: 26, UX/UI designer Badges: UX Designer, Figma, Web Designer, Web Designer Opportunities: Reaching design, Designing Websites, Co-founding a company, Brand strategy consulting, Brainstorming, Teaching design"
                }
            ]
        },
        {
            "role": "assistant",
            "content": [
                {
                    "type": "text",
                    "text": "645"
                }
            ]
        },
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "Title: Trader First Name: Hopelessly Romantic Last Name: Ruby Oxler Description: When you are faced with an obstacle, remember to breathe. Always, one breath at a time Badges: Crypto, Education Expert, Product Manager, Software Tester, FinTech Opportunities: Fundraising for non-profits, Email marketing consulting, Editing books, Developing website, Designing Websites, Design projects, Developing apps"
                }
            ]
        },
        {
            "role": "assistant",
            "content": [
                {
                    "type": "text",
                    "text": "240"
                }
            ]
        },
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "Title: Ms First Name: @Manrisen Last Name //_nd Description: everything alright Badges: Telegram, Football, Saqjan,,,, Crypto, CrossFit Enjoyer, Investment Specialist Opportunities: Advising companies Answer"
                }
            ]
        },
        {
            "role": "assistant",
            "content": [
                {
                    "type": "text",
                    "text": "185"
                }
            ]
        },
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "%s"
                }
            ]
        }
    ],
    "temperature": 0.7,
    "max_tokens": 50,
    "top_p": 1,
    "frequency_penalty": 0,
    "presence_penalty": 0
}`, userProfile)
}

type OpenAIResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func sendOpenAIRequest(reqBody string, token string) (*OpenAIResponse, error) {

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(reqBody))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var openAIResponse OpenAIResponse

	if err := json.Unmarshal(body, &openAIResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &openAIResponse, nil
}

func userToString(user db.User) string {
	msg := fmt.Sprintf("Title: %s First Name: %s Last Name: %s Description: %s", *user.Title, *user.FirstName, *user.LastName, *user.Description)

	msg += fmt.Sprintf(" Badges: ")

	for _, badge := range user.Badges {
		msg += fmt.Sprintf("%s,", badge.Text)
	}

	msg += fmt.Sprintf(" Opportunities: ")

	for _, opportunity := range user.Opportunities {
		msg += fmt.Sprintf("%s,", opportunity.Text)
	}

	// replace \ with \\
	msg = strings.ReplaceAll(msg, `\`, `\\`)

	// replace all " with \"
	msg = strings.ReplaceAll(msg, `"`, `\"`)

	re := regexp.MustCompile(`\r?\n`)
	msg = re.ReplaceAllString(msg, " ")
	// replace tabs with spaces
	msg = strings.ReplaceAll(msg, "\t", " ")

	return msg
}

func (j *notifyJob) ModerateUserProfile() error {
	log.Println("Moderating user profiles")

	users, err := j.storage.ListProfilesForModeration()

	log.Printf("Found %d users for moderation", len(users))

	if err != nil {
		return fmt.Errorf("failed to list user profiles: %w", err)
	}

	for _, user := range users {
		reqBody := getRequestBody(userToString(user))

		resp, err := sendOpenAIRequest(reqBody, j.config.openAIToken)

		if err != nil {
			return fmt.Errorf("failed to send OpenAI request: %w", err)
		}

		if len(resp.Choices) == 0 {
			return fmt.Errorf("no choices in OpenAI response")
		}

		choice := resp.Choices[0]

		if choice.FinishReason != "stop" {
			return fmt.Errorf("OpenAI response did not finish")
		}

		score, err := strconv.Atoi(choice.Message.Content)
		if err != nil || score < 0 || score > 1000 {
			return fmt.Errorf("failed to parse score: %w, reqBody: %s, resp: %v", err, reqBody, resp)
		}

		if err := j.storage.UpdateProfileScore(user.ID, score); err != nil {
			return fmt.Errorf("failed to update user review status: %w", err)
		}

		log.Printf("User %s %s %s scoring status: %d", *user.FirstName, *user.Title, *user.Description, score)
		//
		//var msg, url, btnText string
		//if score < 4 {
		//	msg = fmt.Sprintf("Your profile was hidden by our moderation. Try to make it more informative and less spammy.")
		//	url = fmt.Sprintf("%s/users/edit", j.config.webappURL)
		//	btnText = "Edit Profile"
		//} else {
		//	msg = fmt.Sprintf("Your profile was approved by our moderation. Try out creating new collaboration post")
		//	url = fmt.Sprintf("%s/collaborations/edit", j.config.webappURL)
		//	btnText = "Create Post"
		//}

		//if err := j.notifier.SendTextNotification(notification.SendNotificationParams{
		//	ChatID:     user.ChatID,
		//	Message:    telegram.EscapeMarkdown(msg),
		//	WebAppURL:  url,
		//	ButtonText: btnText,
		//}); err != nil {
		//	log.Printf("Failed to send telegram notification: %s", err)
		//}
	}

	return nil
}
