package main

import (
	"context"
	"fmt"
	"log"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {

	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		log.Fatal(err)
	}

	start, err := time.ParseInLocation(time.RFC3339, "2023-10-27T00:55:00+01:00", loc)
	if err != nil {
		log.Fatal(err)
	}

	end := start.AddDate(0, 0, 5)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		fmt.Printf("Processing snooze for - %s \n", d.Format("2006-01-02"))
		dayOfWeek := d.Weekday()

		switch dayOfWeek {
		case time.Friday: // Weekend snooze if from Friday 10pm to Sunday 10pm
			fmt.Printf("Processing snooze for - %s, %s \n", dayOfWeek, d.Format("2006-01-02"))
			weekendStart, err := time.ParseInLocation(time.RFC3339, fmt.Sprintf("%sT21:55:00+01:00", d.Format("2006-01-02")), loc)
			if err != nil {
				log.Fatal(err)
			}

			weekendEnd := weekendStart.AddDate(0, 0, 2)
			err = createCustomSnooze(weekendStart, weekendEnd.Add(time.Minute*30))
			if err != nil {
				log.Fatal(err)
			}

		case time.Saturday, time.Sunday: // Skip snooze for Sat and Sun as covered above
			fmt.Printf("Skipping snooze for - %s, %s \n", dayOfWeek, d.Format("2006-01-02"))
			continue
		default: // Snooze for weekday
			fmt.Printf("Processing snooze for - %s, %s \n", dayOfWeek, d.Format("2006-01-02"))
			err = createCustomSnooze(d, d.Add(time.Minute*30))
			if err != nil {
				log.Fatal(err)
			}

		}

	}

}

// createCustomSnooze creates a custom snooze specified by the config.
func createCustomSnooze(start time.Time, end time.Time) error {

	ctx := context.Background()
	c, err := monitoring.NewSnoozeClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	req := &monitoringpb.CreateSnoozeRequest{
		Parent: "projects/<<CHANGE_ME>>",
		Snooze: &monitoringpb.Snooze{
			DisplayName: fmt.Sprintf("<<CHANGE_ME>>  Alert snooze - %s", start.Format("2006-01-02")),
			Criteria: &monitoringpb.Snooze_Criteria{
				Policies: []string{"projects/<<CHANGE_ME>>/alertPolicies/1793074354357463121"},
			},
			Interval: &monitoringpb.TimeInterval{
				StartTime: timestamppb.New(start),
				EndTime:   timestamppb.New(end),
			},
		},
	}
	//fmt.Printf("%+v", req)
	m, err := c.CreateSnooze(ctx, req)
	if err != nil {
		return fmt.Errorf("\ncould not create custom snooze: %w", err)
	}

	fmt.Printf("Created %s\n", m.GetName())
	return nil
}
