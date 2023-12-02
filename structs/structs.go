package structs

type Campaign struct {
	ID           uint
	Name         string
	MgTemplate   string
	DefaultLang  string
	Translations []Translation
}

type Translation struct {
	Lang       string
	From       string
	Subject    string
	Recipients []Recipient
}

type SendStat struct {
	CampID   uint
	Lang     string
	Email    string
	ExtID    string
	Success  bool
	ErrorMsg string
}

type Recipient struct {
	Name  string
	Email string
	Lang  string
	ExtID string
}

type RecipientExcluded struct {
	Recipient   Recipient
	SendStatsID uint
}
