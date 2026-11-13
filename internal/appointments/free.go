package appointments

import (
	"context"
	"time"
)

func (r *Repository) freeSlots(
	ctx context.Context,
	doctorID string,
	day time.Time,
	length time.Duration,
) ([]Slot, error) {
	y, m, d := day.Date()
	midnight := time.Date(y, m, d, 0, 0, 0, 0, wat)

	blocks, err := r.hours(ctx, doctorID, int(midnight.Weekday()))
	if err != nil {
		return nil, err
	}

	taken, err := r.busy(ctx, doctorID, midnight, midnight.AddDate(0, 0, 1))
	if err != nil {
		return nil, err
	}

	free := []Slot{}

	for _, b := range blocks {
		start := midnight.Add(time.Duration(b.Starts) * time.Second)
		end := midnight.Add(time.Duration(b.Ends) * time.Second)

		for t := start; !t.Add(length).After(end); t = t.Add(length) {
			s := Slot{t, t.Add(length)}

			if t.After(time.Now()) && !clashes(s, taken) {
				free = append(free, s)
			}
		}
	}

	return free, nil
}

func clashes(s Slot, taken []Slot) bool {
	for _, b := range taken {
		if s.StartsAt.Before(b.EndsAt) && b.StartsAt.Before(s.EndsAt) {
			return true
		}
	}

	return false
}
