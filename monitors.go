/*
 * Datadog API for Go
 *
 * Please see the included LICENSE file for licensing information.
 *
 * Copyright 2013 by authors and contributors.
 */

package datadog

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ThresholdCount represents an object of various threshold settings applicable to metric alerts.
type ThresholdCount struct {
	Ok               *json.Number `json:"ok,omitempty"`
	Critical         *json.Number `json:"critical,omitempty"`
	Warning          *json.Number `json:"warning,omitempty"`
	CriticalRecovery *json.Number `json:"critical_recovery,omitempty"`
	WarningRecovery  *json.Number `json:"warning_recovery,omitempty"`
}

// NoDataTimeframe refers to the number of minutes before a monitor will notify when data stops reporting.
type NoDataTimeframe int

// UnmarshalJSON unmarshals the JSON object into the NoDataTimeframe type.
func (tf *NoDataTimeframe) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "false" || s == "null" {
		*tf = 0
	} else {
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return err
		}
		*tf = NoDataTimeframe(i)
	}
	return nil
}

// Options represents a dictionary of settings for the monitor.
type Options struct {
	NoDataTimeframe   NoDataTimeframe `json:"no_data_timeframe,omitempty"`
	NotifyAudit       *bool           `json:"notify_audit,omitempty"`
	NotifyNoData      *bool           `json:"notify_no_data,omitempty"`
	RenotifyInterval  *int            `json:"renotify_interval,omitempty"`
	NewHostDelay      *int            `json:"new_host_delay,omitempty"`
	EvaluationDelay   *int            `json:"evaluation_delay,omitempty"`
	Silenced          map[string]int  `json:"silenced,omitempty"`
	TimeoutH          *int            `json:"timeout_h,omitempty"`
	EscalationMessage *string         `json:"escalation_message,omitempty"`
	Thresholds        *ThresholdCount `json:"thresholds,omitempty"`
	IncludeTags       *bool           `json:"include_tags,omitempty"`
	RequireFullWindow *bool           `json:"require_full_window,omitempty"`
	Locked            *bool           `json:"locked,omitempty"`
}

// Monitor allows watching a metric or check that you care about,
// notifying your team when some defined threshold is exceeded
type Monitor struct {
	Creator *Creator `json:"creator,omitempty"`
	ID      *int     `json:"id,omitempty"`
	Type    *string  `json:"type,omitempty"`
	Query   *string  `json:"query,omitempty"`
	Name    *string  `json:"name,omitempty"`
	Message *string  `json:"message,omitempty"`
	Tags    []string `json:"tags,omitempty"`
	Options *Options `json:"options,omitempty"`
}

// Creator contains the creator of the monitor
type Creator struct {
	Email  *string `json:"email,omitempty"`
	Handle *string `json:"handle,omitempty"`
	ID     *int    `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
}

// reqMonitors receives a slice of all monitors
type reqMonitors struct {
	Monitors []Monitor `json:"monitors,omitempty"`
}

// CreateMonitor adds a new monitor to the system. This returns a pointer to a
// monitor so you can pass that to UpdateMonitor later if needed
func (client *Client) CreateMonitor(monitor *Monitor) (*Monitor, error) {
	var out Monitor
	// TODO: is this more pretty of frowned upon?
	if err := client.doJSONRequest("POST", "/v1/monitor", monitor, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateMonitor takes a monitor that was previously retrieved through some method
// and sends it back to the server
func (client *Client) UpdateMonitor(monitor *Monitor) error {
	return client.doJSONRequest("PUT", fmt.Sprintf("/v1/monitor/%d", *monitor.ID),
		monitor, nil)
}

// GetMonitor retrieves a monitor by identifier
func (client *Client) GetMonitor(id int) (*Monitor, error) {
	var out Monitor
	if err := client.doJSONRequest("GET", fmt.Sprintf("/v1/monitor/%d", id), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetMonitorsByName retrieves monitors by name
func (client *Client) GetMonitorsByName(name string) ([]Monitor, error) {
	var out reqMonitors
	query, err := url.ParseQuery(fmt.Sprintf("name=%v", name))
	if err != nil {
		return nil, err
	}

	err = client.doJSONRequest("GET", fmt.Sprintf("/v1/monitor?%v", query.Encode()), nil, &out.Monitors)
	if err != nil {
		return nil, err
	}
	return out.Monitors, nil
}

// GetMonitorsByTags retrieves monitors by a slice of tags
func (client *Client) GetMonitorsByTags(tags []string) ([]Monitor, error) {
	var out reqMonitors
	query, err := url.ParseQuery(fmt.Sprintf("monitor_tags=%v", strings.Join(tags, ",")))
	if err != nil {
		return nil, err
	}

	err = client.doJSONRequest("GET", fmt.Sprintf("/v1/monitor?%v", query.Encode()), nil, &out.Monitors)
	if err != nil {
		return nil, err
	}
	return out.Monitors, nil
}

// DeleteMonitor removes a monitor from the system
func (client *Client) DeleteMonitor(id int) error {
	return client.doJSONRequest("DELETE", fmt.Sprintf("/v1/monitor/%d", id),
		nil, nil)
}

// GetMonitors returns a slice of all monitors
func (client *Client) GetMonitors() ([]Monitor, error) {
	var out reqMonitors
	if err := client.doJSONRequest("GET", "/v1/monitor", nil, &out.Monitors); err != nil {
		return nil, err
	}
	return out.Monitors, nil
}

// MuteMonitors turns off monitoring notifications
func (client *Client) MuteMonitors() error {
	return client.doJSONRequest("POST", "/v1/monitor/mute_all", nil, nil)
}

// UnmuteMonitors turns on monitoring notifications
func (client *Client) UnmuteMonitors() error {
	return client.doJSONRequest("POST", "/v1/monitor/unmute_all", nil, nil)
}

// MuteMonitor turns off monitoring notifications for a monitor
func (client *Client) MuteMonitor(id int) error {
	return client.doJSONRequest("POST", fmt.Sprintf("/v1/monitor/%d/mute", id), nil, nil)
}

// UnmuteMonitor turns on monitoring notifications for a monitor
func (client *Client) UnmuteMonitor(id int) error {
	return client.doJSONRequest("POST", fmt.Sprintf("/v1/monitor/%d/unmute", id), nil, nil)
}
