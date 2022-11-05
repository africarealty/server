package domain

type TimeUnitType uint

const (
	TimeUnitDay TimeUnitType = iota
	TimeUnitMonth
	TimeUnitYear
	TimeUnitHour
)

type PeriodUnitType uint

const (
	PeriodUnitDaily PeriodUnitType = iota
	PeriodUnitMonthly
	PeriodUnitYearly
	PeriodUnitHourly
	PeriodUnitOnce
)

type Currency uint

const (
	CurrencyUSD Currency = iota
	CurrencyEUR
	CurrencyRUB
)
