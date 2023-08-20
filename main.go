package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/earlgray283/slack-event-notification-bot/config"
	"github.com/slack-go/slack"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

func main() {
	envConf, err := config.LoadEnvConfig()
	if err != nil {
		panic(err)
	}
	yamlConf, err := config.LoadYamlConfig(envConf.YamlConfigPath)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	now := time.Now()
	slackClient := slack.New(envConf.SlackBotUserOAuthToken)
	calendarSrv, err := calendar.NewService(ctx, option.WithScopes(calendar.CalendarReadonlyScope))
	if err != nil {
		panic(err)
	}

	for _, entry := range yamlConf.Calendars {
		eventList, err := calendarSrv.Events.List(entry.ID).Do(
			googleapi.QueryParameter("singleEvents", "true"),
			googleapi.QueryParameter("orderBy", "startTime"),
			googleapi.QueryParameter("timeMin", now.Format(time.RFC3339)),
		)
		if err != nil {
			panic(err)
		}

		notifyEvents := make([]*calendar.Event, 0)
		for _, event := range eventList.Items {
			startAt, _ := time.Parse(time.RFC3339, event.Start.DateTime)
			if !now.Before(startAt) {
				continue
			}
			matched, err := regexp.MatchString(entry.Event.Summary, event.Summary)
			if err != nil {
				log.Println(err)
				break
			}
			log.Println(entry.Event.Summary, event.Summary, matched)
			if !matched {
				continue
			}
			if startAt.Sub(now) <= time.Hour+10*time.Second {
				notifyEvents = append(notifyEvents, event)
			} else {
				break
			}
		}

		attachments := make([]slack.Attachment, 0)
		for _, event := range notifyEvents {
			startAt, _ := time.Parse(time.RFC3339, event.Start.DateTime)
			endAt, _ := time.Parse(time.RFC3339, event.End.DateTime)
			attachments = append(attachments, slack.Attachment{
				Text:  fmt.Sprintf("%s ~ %s: %s", startAt.Format("15:04"), endAt.Format("15:04"), event.Summary),
				Color: "#99b7dc",
			})
		}
		if len(attachments) > 0 {
			if err := postMessageMulti(slackClient, entry.Channels, slack.MsgOptionText("<!channel>\n1時間以内の予定一覧だよ", false), slack.MsgOptionAttachments(attachments...)); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func postMessageMulti(c *slack.Client, channelIDs []string, options ...slack.MsgOption) error {
	for _, channelID := range channelIDs {
		_, _, err := c.PostMessage(channelID, options...)
		if err != nil {
			return err
		}
	}
	return nil
}
