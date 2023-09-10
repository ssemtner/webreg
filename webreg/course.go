package webreg

import "fmt"

type CourseInfo struct {
	SubjectCode string `json:"SUBJECT_CODE"`
	CourseCode  string `json:"COURSE_CODE"`
	Title       string `json:"TITLE"`
	Sections    []*SectionInfo
}

type SectionInfo struct {
	Capacity          int    `json:"SCTN_CPCTY_QTY"`
	Enrolled          int    `json:"SCTN_ENRLT_QTY"`
	SectionNumber     string `json:"SECTION_NUMBER"`
	BeginHour         int    `json:"BEGIN_HH_TIME"`
	BeginMinute       int    `json:"BEGIN_MM_TIME"`
	EndHour           int    `json:"END_HH_TIME"`
	EndMinute         int    `json:"END_MM_TIME"`
	RoomCode          string `json:"ROOM_CODE"`
	Instructor        string `json:"PERSON_FULL_NAME"`
	SectionStartDate  string `json:"SECTION_START_DATE"`
	SectionEndDate    string `json:"SECTION_END_DATE"`
	WaitlistCount     int    `json:"COUNT_ON_WAITLIST"`
	StartDate         string `json:"START_DATE"`
	DayCode           string `json:"DAY_CODE"`
	BuildingCode      string `json:"BLDG_CODE"`
	SectionCode       string `json:"SECT_CODE"`
	AvailableSeats    int    `json:"AVAIL_SEATS"`
	LongDescription   string `json:"LONG_DESCR"`
	BeforeDescription string `json:"BEFORE_DESCR"`
	PrintFlag         string `json:"PRINT_FLAG"`

	// Stuff I don't really know what it is
	SPT_EBRLT_FLAG     string `json:"STP_EBRLT_FLAG"`
	PRIMARY_INSTR_FLAG string `json:"PRIMARY_INSTR_FLAG"`
	FK_SPM_SPCL_MTG_CD string `json:"FK_SPM_SPCL_MTG_CD"`
	FK_SST_SCTN_STATCD string `json:"FK_SST_SCTN_STATCD"`
	FK_CDI_INSTR_TYPE  string `json:"FK_CDI_INSTR_TYPE"`
}

func (s *SectionInfo) Display() {
	fmt.Println("Section Number:", s.SectionNumber)
	fmt.Println("Section Code:", s.SectionCode)
	fmt.Println("Instructor:", s.Instructor)
	fmt.Println("Room Code:", s.RoomCode)
	fmt.Println("Building Code:", s.BuildingCode)
	fmt.Println("Enrolled:", s.Enrolled)
	fmt.Println("Capacity:", s.Capacity)
	fmt.Println("Available Seats:", s.AvailableSeats)
	fmt.Println("Waitlist Count:", s.WaitlistCount)
	fmt.Println("Day Code:", s.DayCode)
}
