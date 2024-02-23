package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"watcher/domain/event"
	"watcher/entities/events/github"
	"watcher/utils"

	"github.com/gin-gonic/gin"
)

type EventController interface {
	HandleEvent(ctx *gin.Context)
}

type eventControllerImplementation struct {
	service event.Service
}

func NewEventController(service event.Service) EventController {
	return &eventControllerImplementation{
		service: service,
	}
}

func (eci *eventControllerImplementation) HandleEvent(ctx *gin.Context) {
	// check if webhook repo has push webhook active
	// if (isHookActive()) {
	ua := strings.ToLower(ctx.Request.Header.Get("User-Agent"))
	fmt.Printf("%v", ua)

	for _, val := range utils.HookSource {
		if strings.Contains(val, ua) {
			ua = val
			break
		}
	}

	fmt.Printf("ua: %v", ua)
	switch ua {
	case "github":
		event := github.PushEventRequest{
			Source: "github",
			Type:   ctx.Request.Header.Get("x-github-event"),
		}
		var body interface{}
		if err := ctx.BindJSON(&event); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		if err := ctx.BindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		bytes, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		bodyStr := string(bytes)

		fmt.Printf("event: %v", event)
		fmt.Printf("body: %v", body)
		fmt.Printf("body: %s", bodyStr)
		err = eci.service.PushGithubEvent(ctx, &event)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{})
			return
		}
	case "gitlab":
	default:
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "event pushed to kafka successfully",
	})
}
