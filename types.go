package main

import (
	"fmt"
	"strings"
)

const currentHelpText = `
# HELP current_power The current power of the system
# TYPE current_power gauge
`

const powerHelpText = `
# HELP total_power The amount of power generated for the day
# TYPE total_power counter
`

/*
Generated struct from my ECU's json output

Basically it's an array of timestamped data of the output for the day, plus a
field for total power. so not quite real time, but also not all dated. Annoying!
*/
type RealTimeData struct {
	Power       []Power `json:"power"`
	TodayEnergy string  `json:"today_energy"`
	Subtitle    string  `json:"subtitle"`
}

type Power struct {
	Time            int64 `json:"time"`
	EachSystemPower int   `json:"each_system_power"`
}

// Generate a Prometheus-compatible metric string, including timestamps on the
// dated data
func (r RealTimeData) ToMetrics() string {
	var b strings.Builder

	b.WriteString(currentHelpText)
	for _, c := range r.Power {
		b.WriteString(fmt.Sprintf("current_power %v %v\n", c.EachSystemPower, c.Time))
	}

	b.WriteString("\n")

	b.WriteString(powerHelpText)
	b.WriteString(fmt.Sprintf("total_power %v\n", r.TodayEnergy))

	return b.String()
}
