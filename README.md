# Snoozer
A utility that can create recurring  snoozes in Google Cloud monitoring parsing a YAML config.
## Run

```
go run main.go --filepath=config.yaml
```

## Config 
```yaml
---
project_id: "foobar"
snooze_display_name: "foobar  alert snooze"
policy_details:
  - "projects/foobar/alertPolicies/17623358374314480245"
snooze_schedule:
  weekday_start_date_time: "2023-10-27T00:55:00+01:00"
  weekday_end_duration_days: 5
  weekday_duration: 30
  weekend_start_time: "21:55"
```
* `project_id` - The project in which to create the snooze
* `snooze_display_name` - A display name for the Snooze
* `policy_details` - A list of Alert policies that will be snoozed
* `weekday_start_date_time` - The date and time snooze will start in RFC3339 format
* `weekday_end_duration_days` - The number of days for which the recurring snooze will be set
* `weekday_duration` - Duration of snooze on the weekday
* `weekend_start_time` - The snooze start time on a weekend (Weekend is considered Friday to Sunday)
