package notification

import (
	"context"
	"fmt"

	telegram "github.com/go-telegram/bot"
	"github.com/peatch-io/peatch/internal/db"
)

type Notifier struct {
	bot         *telegram.Bot
	adminChatID int64
	botWebApp   string
	webappURL   string
}

func NewNotifier(botToken string, adminChatID int64, botWebApp, webappURL string) (*Notifier, error) {
	bot, err := telegram.New(botToken)
	if err != nil {
		return nil, err
	}
	return &Notifier{bot: bot, adminChatID: adminChatID, botWebApp: botWebApp, webappURL: webappURL}, nil
}

func (n *Notifier) NotifyUserVerified(user db.User) error {
	msgText := fmt.Sprintf("🎉 Congratulations! Your profile has been verified.\nYou can now access all features and be visible to other users.\n\nCheck your profile: %s/users/%s",
		n.webappURL, user.ID)

	_, err := n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{ChatID: fmt.Sprintf("%d", user.ChatID), Text: msgText})

	return err
}

func (n *Notifier) NotifyCollaborationVerified(collab db.Collaboration) error {
	msgText := fmt.Sprintf("🎉 Congratulations! Your collaboration \"%s\" has been verified.\nIt is now visible to all Peatch users.\n\nView collaboration: %s/collaborations/%s",
		collab.Title, n.webappURL, collab.ID)

	_, err := n.bot.SendMessage(context.Background(), &telegram.SendMessageParams{
		ChatID: fmt.Sprintf("%d", collab.User.ChatID),
		Text:   msgText,
	})

	return err
}

func (n *Notifier) NotifyNewPendingUser(user db.User) error {
	firstName := ""
	if user.FirstName != nil {
		firstName = *user.FirstName
	}
	lastName := ""
	if user.LastName != nil {
		lastName = *user.LastName
	}

	msgText := fmt.Sprintf("🔔 New user pending verification:\nName: %s %s\nUsername: @%s\n\nVerify them in the admin panel: %s/admin/users",
		firstName, lastName, user.Username, n.webappURL)

	params := &telegram.SendMessageParams{ChatID: fmt.Sprintf("%d", n.adminChatID), Text: msgText}

	_, err := n.bot.SendMessage(context.Background(), params)

	return err
}

func (n *Notifier) NotifyNewPendingCollaboration(user db.User, collab db.Collaboration) error {
	firstName := ""
	if user.FirstName != nil {
		firstName = *user.FirstName
	}
	lastName := ""
	if user.LastName != nil {
		lastName = *user.LastName
	}

	msgText := fmt.Sprintf("🔔 New collaboration pending verification:\nTitle: %s\nBy: %s %s (@%s)\n\nVerify it in the admin panel: %s/admin/collaborations",
		collab.Title, firstName, lastName, user.Username, n.webappURL)

	params := &telegram.SendMessageParams{
		ChatID: fmt.Sprintf("%d", n.adminChatID),
		Text:   msgText,
	}

	_, err := n.bot.SendMessage(context.Background(), params)

	return err
}
