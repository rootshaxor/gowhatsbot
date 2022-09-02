package texts

type QouteChar string

const (
	QuoteMono   QouteChar = "```"
	QuoteBold   QouteChar = "*"
	QuoteItalic QouteChar = "_"
	QuoteStrike QouteChar = "~"
)
