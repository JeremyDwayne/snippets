package main

import (
	"context"
	"os"

	"github.com/mailerlite/mailerlite-go"
)

type Mailer struct {
	client   *mailerlite.Client
	group_id string
}

func NewMailer() Mailer {
	return Mailer{
		client:   mailerlite.NewClient(os.Getenv("MAILER_LITE_API_TOKEN")),
		group_id: os.Getenv("MAILER_LITE_GROUP_ID"),
	}
}

func (app *application) NewSubscriber(ctx context.Context, email string) error {
	subscriber := &mailerlite.Subscriber{
		Email: email,
	}

	newSubscriber, _, err := app.mailer.client.Subscriber.Create(ctx, subscriber)
	if err != nil {
		app.logger.Error("new sub", err)
		return err
	}

	_, _, err = app.mailer.client.Group.Assign(ctx, app.mailer.group_id, newSubscriber.Data.ID)
	if err != nil {
		app.logger.Error("group assign", err)
		return err
	}

	return nil
}
