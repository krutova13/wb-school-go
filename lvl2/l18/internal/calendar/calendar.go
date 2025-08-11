package calendar

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"calendar/internal/types"
)

// Service представляет сервис календаря
type Service struct {
	events map[string]*types.Event
	mutex  sync.RWMutex
}

// NewService создает новый экземпляр сервиса календаря
func NewService() *Service {
	return &Service{
		events: make(map[string]*types.Event),
	}
}

// CreateEvent создает новое событие
func (s *Service) CreateEvent(userID, dateStr, text string) (*types.Event, error) {
	if userID == "" {
		return nil, errors.New("user_id не может быть пустым")
	}
	if text == "" {
		return nil, errors.New("текст события не может быть пустым")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат даты: %v", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	eventID := fmt.Sprintf("%s_%s_%d", userID, dateStr, time.Now().UnixNano())
	event := &types.Event{
		ID:     eventID,
		UserID: userID,
		Date:   date,
		Text:   text,
	}

	s.events[eventID] = event
	return event, nil
}

// UpdateEvent обновляет существующее событие
func (s *Service) UpdateEvent(id, userID, dateStr, text string) (*types.Event, error) {
	if id == "" {
		return nil, errors.New("id события не может быть пустым")
	}
	if userID == "" {
		return nil, errors.New("user_id не может быть пустым")
	}
	if text == "" {
		return nil, errors.New("текст события не может быть пустым")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат даты: %v", err)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	event, exists := s.events[id]
	if !exists {
		return nil, errors.New("событие не найдено")
	}

	if event.UserID != userID {
		return nil, errors.New("нет прав на обновление этого события")
	}

	event.Date = date
	event.Text = text

	return event, nil
}

// DeleteEvent удаляет событие
func (s *Service) DeleteEvent(id, userID string) error {
	if id == "" {
		return errors.New("id события не может быть пустым")
	}
	if userID == "" {
		return errors.New("user_id не может быть пустым")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	event, exists := s.events[id]
	if !exists {
		return errors.New("событие не найдено")
	}

	if event.UserID != userID {
		return errors.New("нет прав на удаление этого события")
	}

	delete(s.events, id)
	return nil
}

// GetEventsForDay возвращает события на конкретный день
func (s *Service) GetEventsForDay(userID, dateStr string) ([]*types.Event, error) {
	if userID == "" {
		return nil, errors.New("user_id не может быть пустым")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат даты: %v", err)
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var events []*types.Event
	for _, event := range s.events {
		if event.UserID == userID && isSameDay(event.Date, date) {
			events = append(events, event)
		}
	}

	return events, nil
}

// GetEventsForWeek возвращает события на неделю
func (s *Service) GetEventsForWeek(userID, dateStr string) ([]*types.Event, error) {
	if userID == "" {
		return nil, errors.New("user_id не может быть пустым")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат даты: %v", err)
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	weekStart := getWeekStart(date)
	weekEnd := weekStart.AddDate(0, 0, 7)

	var events []*types.Event
	for _, event := range s.events {
		if event.UserID == userID && (event.Date.Equal(weekStart) || event.Date.After(weekStart)) && event.Date.Before(weekEnd) {
			events = append(events, event)
		}
	}

	return events, nil
}

// GetEventsForMonth возвращает события на месяц
func (s *Service) GetEventsForMonth(userID, dateStr string) ([]*types.Event, error) {
	if userID == "" {
		return nil, errors.New("user_id не может быть пустым")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат даты: %v", err)
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	var events []*types.Event
	for _, event := range s.events {
		if event.UserID == userID && (event.Date.Equal(monthStart) || event.Date.After(monthStart)) && event.Date.Before(monthEnd) {
			events = append(events, event)
		}
	}

	return events, nil
}

func isSameDay(date1, date2 time.Time) bool {
	return date1.Year() == date2.Year() && date1.YearDay() == date2.YearDay()
}

func getWeekStart(date time.Time) time.Time {
	weekday := date.Weekday()
	daysToSubtract := int(weekday)
	if weekday == time.Sunday {
		daysToSubtract = 6
	} else {
		daysToSubtract = int(weekday - 1)
	}
	return date.AddDate(0, 0, -daysToSubtract)
}
