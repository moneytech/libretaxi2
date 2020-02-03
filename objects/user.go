package objects

type MenuId int

const (
	Menu_Init        MenuId = 100
	Menu_AskLocation MenuId = 200
	Menu_Feed        MenuId = 300
	Menu_Post        MenuId = 400
	Menu_Ban 	     MenuId = 999999
)

type User struct {
	UserId int64
	MenuId MenuId
	Username string
	FirstName string
	LastName string
	Lon float64
	Lat float64
	LanguageCode string
	ReportCnt int
	ShadowBanned bool
}
