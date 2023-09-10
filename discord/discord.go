package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"webreg/webreg"

	"github.com/bwmarrin/discordgo"
)

func Run(token string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	command := discordgo.ApplicationCommand{
		Name:        "courseinfo",
		Description: "Get course information",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "code",
				Description: "Course code",
				Required:    true,
			},
		},
	}

	if err = dg.Open(); err != nil {
		panic(err)
	}

	defer dg.Close()

	client := webreg.NewWebregClient(webreg.Term{
		Code:   "FA23",
		Option: "5320:::FA23",
	})
	client.GetCookiesFromServer()

	dg.ApplicationCommandCreate(dg.State.User.ID, "912837595065614387", &command)
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name == "courseinfo" {
			options := i.ApplicationCommandData().Options
			log.Printf("Start courseinfo request for %s\n", options[0].StringValue())

			split := strings.Split(options[0].StringValue(), " ")
			subject := split[0]
			course := split[1]

			// send inital response
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})

			log.Println("Initial response sent, getting course info")

			result, err := client.GetCourseInfo(subject, course)
			if err != nil {
				panic(err)
			}

			log.Println("Course info retrieved, creating embeds")

			// fields := []*discordgo.MessageEmbedField{}

			contents := []string{}

			for _, section := range result.Sections {
				content := ""
				content += fmt.Sprintf("**Section %s** - ", section.SectionCode)
				content += fmt.Sprintf("**%d** / **%d** | **%d** Waitlist\n", section.Enrolled, section.Capacity, section.WaitlistCount)
				content += fmt.Sprintf("```%s | %s %s```", strings.TrimSpace(strings.Split(section.Instructor, ";")[0]), section.BuildingCode, section.RoomCode)

				contents = append(contents, content)
			}

			log.Println("Sending followup")

			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: fmt.Sprintf("%s %s\n%s", result.SubjectCode, result.CourseCode, strings.Join(contents, "")),
			})

			log.Println("Followup sent")
		}
	})

	fmt.Println("Bot running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
