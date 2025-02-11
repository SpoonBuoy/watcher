package eventService

import (
	"fmt"
	"watcher/domain/event"
	"watcher/entities/events/github"
	"watcher/utils"

	"github.com/gin-gonic/gin"
)

type eventServiceImplementation struct {
	token string
}

func NewEventServiceImplementation(config *utils.Config) event.Service {
	return &eventServiceImplementation{
		token: config.Token,
	}
}

func (osi *eventServiceImplementation) PushGithubEvent(ctx *gin.Context, event *github.GithubHookRequest) error {
	// push to kafka with proper payload
	fmt.Println(event)
	return nil
}
