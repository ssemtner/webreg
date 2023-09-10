package main

import "webreg/webreg"

func main() {
	client := webreg.NewWebregClient(webreg.Term{
		Code:   "FA23",
		Option: "5320:::FA23",
	})

	client.GetCookiesFromServer()

	result, err := client.GetCourseInfo("CSE", "20")
	if err != nil {
		panic(err)
	}

	for _, section := range result.Sections {
		section.Display()
	}
}
