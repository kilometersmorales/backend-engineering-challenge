package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"unbabel/internal/util"

	"github.com/spf13/viper"
)

var inFormat = "2006-01-02 15:04:05.000000"
var outFormat = "2006-01-02 15:04:05"

type event struct {
	STimestamp    string `json:"timestamp"`
	Timestamp     time.Time
	TranslationID string `json:"translation_id"`
	SourceLang    string `json:"source_language"`
	TargetLang    string `json:"target_language"`
	ClientName    string `json:"client_name"`
	EventName     string `json:"event_name"`
	Duration      int    `json:"duration"`
	NRWords       int    `json:"nr_words"`
}

func Parse(config Config) {

	// Log flags
	if config.Debug {
		for key, value := range viper.GetViper().AllSettings() {
			fmt.Printf("Command flag %v = %v\n", key, value)
		}
	}

	// Read input file, line by line

	file, err := os.Open(config.InputFile)
	util.Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// A place to gather events
	events := make([]event, 0)

	for scanner.Scan() {

		line := scanner.Text()

		// Log line
		if config.Debug {
			fmt.Printf("Input line = %v\n", line)
		}

		// Unmarshal JSON to event struct
		event := event{}
		err = json.Unmarshal([]byte(line), &event)
		util.Check(err)

		// Ignore if event doesn't match the command flags
		if config.SourceLang != "all" && config.SourceLang != event.SourceLang ||
			config.TargetLang != "all" && config.TargetLang != event.TargetLang ||
			config.ClientName != "all" && config.ClientName != event.ClientName {
			fmt.Printf("Ignored input event = %v\n", event)
			continue
		}

		// Accept event

		if config.Debug {
			fmt.Printf("Accepted input event = %v\n", event)
		}

		// Parse event timestamp
		event.Timestamp, err = time.Parse(inFormat, event.STimestamp)
		util.Check(err)

		events = append(events, event)
	}

	// Log events
	if config.Debug {
		fmt.Printf("Considered events = %v\n", events)
	}

	// Assume that the events are ordered by timestamp

	if len(events) == 0 {
		fmt.Println("No events were considered")
		return
	}

	eventPool := make([]event, 0)

	// Determine the first and last "minute" blocks
	currentBlockTS := minuteBlockTS(events[0].Timestamp)
	finalBlockTS := minuteBlockTS(events[len(events)-1].Timestamp)
	finalBlockTS = finalBlockTS.Add(time.Duration(1) * time.Minute)
	if config.Debug {
		fmt.Printf("finalBlockTS = %v\n", finalBlockTS.Format(outFormat))
	}

	// Calculate average for every block

	finalOutput := ""

	for currentBlockTS.Compare(finalBlockTS) <= 0 {
		if config.Debug {
			fmt.Printf("\ncurrentBlockTS = %v\n", currentBlockTS.Format(outFormat))
		}

		// Remove events from the pool if they're not inside the time window
		for i := 0; i < len(eventPool); i++ {

			diff := minuteDifference(currentBlockTS, eventPool[i].Timestamp)
			if diff > float64(config.WindowSize) {
				// Remove from pool
				if config.Debug {
					fmt.Printf("removed event from pool = %v\n",
						eventPool[i].Timestamp.Format(outFormat))
				}
				eventPool = append(eventPool[:i], eventPool[i+1:]...)
				// Move index back because of deletion
				i--
			} else {
				// All remaining events are within the time window
				if config.Debug {
					fmt.Printf("event in window = %v\n",
						eventPool[i].Timestamp.Format(outFormat))
				}
				break
			}
		}

		// Add events to the pool if they're inside the time window
		for i := 0; i < len(events); i++ {

			// Check if event takes place after the current block
			if currentBlockTS.Compare(events[i].Timestamp) < 0 {
				// Remaining events are placed after the current block
				if config.Debug {
					fmt.Printf("future event = %v\n",
						events[i].Timestamp.Format(outFormat))
				}
				break
			}

			// Event takes place before the "minute" block

			diff := minuteDifference(currentBlockTS, events[i].Timestamp)
			if diff >= 0 && diff <= float64(config.WindowSize) {
				// Event is inside window. Add to pool
				if config.Debug {
					fmt.Printf("added event to pool = %v\n",
						events[i].Timestamp.Format(outFormat))
				}
				eventPool = append(eventPool, events[i])
				// Remove from gathered events to prevent adding it later
				events = append(events[:i], events[i+1:]...)
				// Move index back because of deletion
				i--
			}
		}

		// Calculate average from the pool of events
		var avg float64
		if len(eventPool) == 0 {
			avg = 0
		} else {
			sum := 0
			for _, event := range eventPool {
				sum += event.Duration
			}
			avg = float64(sum) / float64(len(eventPool))
		}

		// Write output
		finalOutput += fmt.Sprintf(
			"{\"date\": \"" + currentBlockTS.Format(outFormat) +
				"\", \"average_delivery_time\": " + fmt.Sprint(avg) + "}\n")

		// Move on to next "minute" block
		currentBlockTS = currentBlockTS.Add(time.Duration(1) * time.Minute)
	}

	fmt.Print(finalOutput)
}

// Returns the difference in minutes between two timestamps
func minuteDifference(t1, t2 time.Time) float64 {
	return t1.Sub(t2).Minutes()
}

// Returns the same timestamp with no information on seconds and lower fields
func minuteBlockTS(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(),
		timestamp.Hour(), timestamp.Minute(), 0, 0, timestamp.Location())
}
