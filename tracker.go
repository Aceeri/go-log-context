package logContext

import (
	"fmt"
	"sync"
	"time"
)

type Trackers struct {
	markers sync.Map
	profile sync.Map
}

type TrackerInstant struct {
	Timestamp time.Time
	State     string
	Tags      []string
}

type ProfileInstant struct {
	Timestamp time.Duration
	State     string
	Tags      []string
}

type ProfileState struct {
	Top      []ProfileInstant
	Previous time.Time
	lock     *sync.Mutex
}

func (state *ProfileState) SetTop(name string, location ...string) {
	state.lock.Lock()
	defer state.lock.Unlock()

	difference := time.Now().Sub(state.Previous)
	state.Previous = time.Now()

	if len(state.Top) > 0 {
		state.Top[len(state.Top)-1].Timestamp = difference
	}

	state.Top = append(state.Top, ProfileInstant{
		State: location[0],
		Tags:  location[1:],
	})
}

func (trackers *Trackers) Register(markers []string) {
	for _, tracker := range markers {
		trackers.SetTracker(tracker, "register")
	}
}

func (trackers *Trackers) RegisterProfile(markers []string) {
	for _, tracker := range markers {
		trackers.SetProfile(tracker, "register")
	}
}

func (trackers *Trackers) Unregister(markers []string) {
	for _, tracker := range markers {
		trackers.UnsetTracker(tracker)
	}
}

func (trackers *Trackers) SetTracker(name string, location ...string) {
	trackers.markers.Store(name, TrackerInstant{
		Timestamp: time.Now(),
		State:     location[0],
		Tags:      location[1:],
	})

	trackers.SetProfile(name, location...)
}

func (trackers *Trackers) SetProfile(name string, location ...string) {
	any, ok := trackers.profile.Load(name)
	profile, cast := any.(*ProfileState)
	if !ok || !cast {
		return
	}

	profile.SetTop(name, location...)
}

func (trackers *Trackers) GetProfile(name string) ProfileState {
	any, ok := trackers.profile.Load(name)
	profile, cast := any.(*ProfileState)
	if !ok || !cast {
		var lock sync.Mutex
		return ProfileState{
			lock: &lock,
		}
	}

	return *profile
}

func (trackers *Trackers) Reset(name string) {
	var lock sync.Mutex
	state := ProfileState{
		lock: &lock,
	}
	trackers.profile.Store(name, &state)
}

func (trackers *Trackers) UnsetTracker(name string) {
	trackers.markers.Delete(name)
	trackers.profile.Delete(name)
}

func (trackers *Trackers) GetTracker(name string) (TrackerInstant, bool) {
	value, ok := trackers.markers.Load(name)
	if !ok {
		return TrackerInstant{}, false
	}

	casted, ok := value.(TrackerInstant)
	return casted, ok
}

func (trackers *Trackers) Overdue(name string, by int64) (bool, bool) {
	value, exists := trackers.GetTracker(name)
	if !exists {
		return false, false
	}

	overdue := value.Timestamp.Unix() < (time.Now().Unix() - by)
	return overdue, true
}

func (trackers *Trackers) AnyOverdue(threshold int64) []string {
	var overdueTrackers []string
	for _, tracker := range trackers.Tracked() {
		overdue, exists := trackers.Overdue(tracker, threshold)
		if exists && overdue {
			overdueTrackers = append(overdueTrackers, tracker)
		}
	}

	return overdueTrackers
}

func (trackers *Trackers) Tracked() []string {
	var markers []string
	trackers.markers.Range(func(key, value interface{}) bool {
		marker, _ := key.(string)
		markers = append(markers, marker)
		return true
	})

	return markers
}

func (trackers *Trackers) FormattedStates(names []string) string {
	states := make(map[string]string)
	for _, name := range names {
		value, ok := trackers.GetTracker(name)
		if !ok {
			states[name] = "not registered"
			continue
		}

		states[name] = fmt.Sprintf("%+v", value)
	}

	return fmt.Sprintf("%+v", states)
}

func (tracker *Trackers) Child(name string) Tracker {
	return Tracker{
		name: name,
		ref:  tracker,
	}
}

type Tracker struct {
	name string
	ref  *Trackers
}

func (tracker *Tracker) Set(location ...string) {
	if tracker.ref != nil {
		tracker.ref.SetTracker(tracker.name, location...)
	}
}

func (tracker *Tracker) Reset() {
	if tracker.ref != nil {
		tracker.ref.Reset(tracker.name)
	}
}

func (tracker *Tracker) GetProfile() ProfileState {
	if tracker.ref != nil {
		return tracker.ref.GetProfile(tracker.name)
	}

	var lock sync.Mutex
	return ProfileState{
		lock: &lock,
	}
}
