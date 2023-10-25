package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v2"
)

func main() {

	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var snoozeConfig SnoozeConfig
	if err := yaml.Unmarshal(f, &snoozeConfig); err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%+v\n", snoozeConfig)

	start, err := time.ParseInLocation(time.RFC3339, snoozeConfig.SnoozeSchedule.WeekdayStartDateTime, loc)
	if err != nil {
		log.Fatal(err)
	}

	end := start.AddDate(0, 0, snoozeConfig.SnoozeSchedule.WeekdayEndDurationDays)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		fmt.Printf("Processing snooze for - %s \n", d.Format("2006-01-02"))
		dayOfWeek := d.Weekday()

		switch dayOfWeek {
		case time.Friday: // Weekend snooze if from Friday 10pm to Sunday 10pm
			fmt.Printf("Processing snooze for - %s, %s \n", dayOfWeek, d.Format("2006-01-02"))
			weekendStart, err := time.ParseInLocation(time.RFC3339, fmt.Sprintf("%sT%s:00+01:00", d.Format("2006-01-02"), snoozeConfig.SnoozeSchedule.WeekendStartTime), loc)
			if err != nil {
				log.Fatal(err)
			}

			weekendEnd := weekendStart.AddDate(0, 0, snoozeConfig.SnoozeSchedule.WeekendDurationDays)
			err = createCustomSnooze(weekendStart, weekendEnd.Add(time.Minute*30), snoozeConfig)
			if err != nil {
				log.Fatal(err)
			}

		case time.Saturday, time.Sunday: // Skip snooze for Sat and Sun as covered above
			fmt.Printf("Skipping snooze for - %s, %s \n", dayOfWeek, d.Format("2006-01-02"))
			continue
		default: // Snooze for weekday
			fmt.Printf("Processing snooze for - %s, %s \n", dayOfWeek, d.Format("2006-01-02"))
			err = createCustomSnooze(d, d.Add(time.Minute*time.Duration(snoozeConfig.SnoozeSchedule.WeekdayDuration)), snoozeConfig)
			if err != nil {
				log.Fatal(err)
			}

		}

	}

}

// createCustomSnooze creates a custom snooze specified by the config.
func createCustomSnooze(start time.Time, end time.Time, snoozeConfig SnoozeConfig) error {

	ctx := context.Background()
	c, err := monitoring.NewSnoozeClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	req := &monitoringpb.CreateSnoozeRequest{
		Parent: fmt.Sprintf("projects/%s", snoozeConfig.ProjectID),
		Snooze: &monitoringpb.Snooze{
			DisplayName: fmt.Sprintf("%s - %s", snoozeConfig.SnoozeDisplayName, start.Format("2006-01-02")),
			Criteria: &monitoringpb.Snooze_Criteria{
				Policies: snoozeConfig.PolicyDetails,
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

type SnoozeConfig struct {
	ProjectID         string   `yaml:"project_id"`
	SnoozeDisplayName string   `yaml:"snooze_display_name"`
	PolicyDetails     []string `yaml:"policy_details"`
	SnoozeSchedule    struct {
		WeekdayStartDateTime   string `yaml:"weekday_start_date_time"`
		WeekdayEndDurationDays int    `yaml:"weekday_end_duration_days"`
		WeekdayDuration        int    `yaml:"weekday_duration"`
		WeekendStartTime       string `yaml:"weekend_start_time"`
		WeekendDurationDays    int    `yaml:"weekend_duration_days"`
	} `yaml:"snooze_schedule"`
}
