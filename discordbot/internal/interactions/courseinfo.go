package interactions

import (
	"fmt"
	"fn/internal/discord"
	"fn/internal/webreg"
	"log"
	"strings"
)

var client *webreg.Client

func init() {
	term, _ := webreg.ParseTerm("FA23")

	client = webreg.NewClient(term)

	cookie, err := getCookie(false)
	if err != nil {
		panic(err)
	}

	client.SetCookie(cookie)
}

func CourseInfo(interaction discord.Interaction) discord.Response {
	courseCodeParts := strings.Split(interaction.Data.Options[0].Value.(string), " ")

	course, err := client.GetCourseInfo(courseCodeParts[0], courseCodeParts[1])
	if err != nil {
		log.Printf("Error: %v", err)

		// refresh cookie
		cookie, err := getCookie(true)
		if err != nil {
			return discordError(err)
		}

		client.SetCookie(cookie)

		return CourseInfo(interaction)
	}

	content := fmt.Sprintf("**%s %s**\n", course.SubjectCode, course.CourseCode)

	for _, section := range course.Sections {
		content += fmt.Sprintf(
			"Section %s - %d / %d (%d WL)\n",
			section.SectionCode,
			section.Enrolled,
			section.Capacity,
			section.WaitlistCount,
		)

		instructor := strings.TrimSpace(strings.Split(section.Instructor, ";")[0])

		content += fmt.Sprintf(
			"```%s - %s %s```\n",
			instructor,
			section.BuildingCode,
			section.RoomCode,
		)
	}

	response := discord.Response{
		Type: 4,
		Data: discord.ResponseData{
			Content: content,
		},
	}

	return response
}
