package job

import (
	"encoding/json"
	"fmt"
	telegram "github.com/go-telegram/bot"
	"github.com/peatch-io/peatch/internal/db"
	"github.com/peatch-io/peatch/internal/notification"
	"io"
	"log"
	"net/http"
	"strings"
)

func getRequestBody(userProfile string) string {
	// replace newlines with spaces
	userProfile = strings.ReplaceAll(userProfile, "\n", " ")

	return fmt.Sprintf(`{
    "model": "gpt-4o",
    "messages": [
        {
            "role": "system",
            "content": [
                {
                    "type": "text",
                    "text": "Review user profile and decide if its a spam or not. Its a social network for collaboration and finding talents. So people profiles are important The more important fields are the title, first_name, last_name, and description Less important is to compare the badges and opportunities with profile description and title, see if they correlate The RESULT is EXACTLY on of - spam, not_spam, not_sure"
                }
            ]
        },
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "Title: Project Manager First Name: German Last Name: Ermolenko Description: 21 y.o. Junior Project/Product manager with marketing skills, student Badges: Business Developer, Content Creator, Project Manager, Business Assistant, Pizza Lover, Manager Opportunities: Giving product reviews, Coaching founders, Co-founding a company, Career coaching, Beta testing new products, Advising on SEO, Advising early stage companies"
                }
            ]
        },
        {
            "role": "assistant",
            "content": [
                {
                    "type": "text",
                    "text": "not_spam"
                }
            ]
        },
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "Title: UX/UI designer, no-code developer Readymag/Tilda First Name: Nikita Last NameL Mironov Description: 26, UX/UI designer Badges: UX Designer, Figma, Web Designer, Web Designer Opportunities: Reaching design, Designing Websites, Co-founding a company, Brand strategy consulting, Brainstorming, Teaching design"
                }
            ]
        },
        {
            "role": "assistant",
            "content": [
                {
                    "type": "text",
                    "text": "not_spam"
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
                    "text": "spam"
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
                    "text": "spam"
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
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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

	for _, badge := range user.Badges {
		msg += fmt.Sprintf(" Badges: %s", badge.Text)
	}

	for _, opportunity := range user.Opportunities {
		msg += fmt.Sprintf(" Opportunities: %s", opportunity.Text)
	}

	return msg
}

func (j *notifyJob) ModerateUserProfile() error {
	users, err := j.storage.ListProfilesForModeration()

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

		status := choice.Message.Content
		if status != "not_spam" && status != "spam" && status != "unsure" {
			log.Printf("Unknown status from OpenAI: %s", status)
			status = "unsure"
		}

		if err := j.storage.UpdateUserReviewStatus(user.ID, status); err != nil {
			return fmt.Errorf("failed to update user review status: %w", err)
		}

		log.Printf("User %d review status: %s", user.ID, status)

		var msg, url, btnText string
		if status == "spam" {
			msg = fmt.Sprintf("Your profile was hidden by our moderation. Try to make it more informative and less spammy.")
			url = fmt.Sprintf("%s/users/edit", j.config.webappURL)
			btnText = "Edit Profile"
		} else if status == "not_spam" {
			msg = fmt.Sprintf("Your profile was approved by our moderation. You can create your first post now.")
			url = fmt.Sprintf("%s/collaborations/edit", j.config.webappURL)
			btnText = "Create Post"
		} else {
			continue
		}

		if err := j.notifier.SendTextNotification(notification.SendNotificationParams{
			ChatID:     j.config.groupChatID,
			Message:    telegram.EscapeMarkdown(msg),
			WebAppURL:  url,
			ButtonText: btnText,
		}); err != nil {
			log.Printf("Failed to send telegram notification: %s", err)
		}
	}

	return nil
}
