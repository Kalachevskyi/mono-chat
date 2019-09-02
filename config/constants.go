package config

//Regexp patterns
const (
	DateShortRegexpPattern = `^\d{2}-\d{2}` //range of date for the current month
	DateRegexpPattern      = `\d{2}\.\d{2}\.\d{4}-\d{2}\.\d{2}\.\d{4}`
	DateTimeRegexpPattern  = `\d{2}\.\d{2}\.\d{4}T\d{2}\.\d{2}-\d{2}\.\d{2}\.\d{4}T\d{2}\.\d{2}`
)

const TimeLocation = "Europe/Kiev"
