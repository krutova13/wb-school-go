package calendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_CreateEvent(t *testing.T) {
	service := NewService()

	tests := []struct {
		name    string
		userID  string
		date    string
		text    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "успешное создание события",
			userID:  "user1",
			date:    "2023-12-31",
			text:    "Новый год",
			wantErr: false,
		},
		{
			name:    "пустой user_id",
			userID:  "",
			date:    "2023-12-31",
			text:    "Новый год",
			wantErr: true,
			errMsg:  "user_id не может быть пустым",
		},
		{
			name:    "пустой текст",
			userID:  "user1",
			date:    "2023-12-31",
			text:    "",
			wantErr: true,
			errMsg:  "текст события не может быть пустым",
		},
		{
			name:    "некорректная дата",
			userID:  "user1",
			date:    "2023-13-31",
			text:    "Новый год",
			wantErr: true,
			errMsg:  "некорректный формат даты",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := service.CreateEvent(tt.userID, tt.date, tt.text)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, tt.userID, event.UserID)
				assert.Equal(t, tt.text, event.Text)
				assert.NotEmpty(t, event.ID)
			}
		})
	}
}

func TestService_UpdateEvent(t *testing.T) {
	service := NewService()

	event, err := service.CreateEvent("user1", "2023-12-31", "Новый год")
	require.NoError(t, err)
	require.NotNil(t, event)

	tests := []struct {
		name    string
		id      string
		userID  string
		date    string
		text    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "успешное обновление события",
			id:      event.ID,
			userID:  "user1",
			date:    "2024-01-01",
			text:    "Обновленный новый год",
			wantErr: false,
		},
		{
			name:    "событие не найдено",
			id:      "несуществующий_id",
			userID:  "user1",
			date:    "2024-01-01",
			text:    "Обновленный новый год",
			wantErr: true,
			errMsg:  "событие не найдено",
		},
		{
			name:    "нет прав на обновление",
			id:      event.ID,
			userID:  "user2",
			date:    "2024-01-01",
			text:    "Обновленный новый год",
			wantErr: true,
			errMsg:  "нет прав на обновление этого события",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedEvent, err := service.UpdateEvent(tt.id, tt.userID, tt.date, tt.text)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, updatedEvent)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updatedEvent)
				assert.Equal(t, tt.text, updatedEvent.Text)
			}
		})
	}
}

func TestService_DeleteEvent(t *testing.T) {
	service := NewService()

	event, err := service.CreateEvent("user1", "2023-12-31", "Новый год")
	require.NoError(t, err)
	require.NotNil(t, event)

	tests := []struct {
		name    string
		id      string
		userID  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "успешное удаление события",
			id:      event.ID,
			userID:  "user1",
			wantErr: false,
		},
		{
			name:    "событие не найдено",
			id:      "несуществующий_id",
			userID:  "user1",
			wantErr: true,
			errMsg:  "событие не найдено",
		},
		{
			name:    "нет прав на удаление",
			id:      event.ID,
			userID:  "user2",
			wantErr: true,
			errMsg:  "нет прав на удаление этого события",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteEvent(tt.id, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "нет прав на удаление" {
					// После первого удаления событие уже не существует
					assert.Contains(t, err.Error(), "событие не найдено")
				} else {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetEventsForDay(t *testing.T) {
	service := NewService()

	_, err := service.CreateEvent("user1", "2023-12-31", "Новый год")
	require.NoError(t, err)

	_, err = service.CreateEvent("user1", "2023-12-31", "Встреча")
	require.NoError(t, err)

	_, err = service.CreateEvent("user1", "2024-01-01", "Другой день")
	require.NoError(t, err)

	_, err = service.CreateEvent("user2", "2023-12-31", "Другой пользователь")
	require.NoError(t, err)

	tests := []struct {
		name    string
		userID  string
		date    string
		wantLen int
		wantErr bool
	}{
		{
			name:    "события на день для user1",
			userID:  "user1",
			date:    "2023-12-31",
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "события на другой день",
			userID:  "user1",
			date:    "2024-01-01",
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "события другого пользователя",
			userID:  "user2",
			date:    "2023-12-31",
			wantLen: 1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, err := service.GetEventsForDay(tt.userID, tt.date)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, events, tt.wantLen)
			}
		})
	}
}

func TestService_GetEventsForWeek(t *testing.T) {
	service := NewService()

	service.CreateEvent("user1", "2023-12-25", "Понедельник")      // Понедельник
	service.CreateEvent("user1", "2023-12-26", "Вторник")          // Вторник
	service.CreateEvent("user1", "2023-12-27", "Среда")            // Среда
	service.CreateEvent("user1", "2024-01-01", "Следующая неделя") // Другая неделя

	events, err := service.GetEventsForWeek("user1", "2023-12-25")
	require.NoError(t, err)
	assert.Len(t, events, 3) // Только события текущей недели (25, 26, 27 декабря)
}

func TestService_GetEventsForMonth(t *testing.T) {
	service := NewService()

	service.CreateEvent("user1", "2023-12-01", "Декабрь")
	service.CreateEvent("user1", "2023-12-15", "Декабрь")
	service.CreateEvent("user1", "2023-12-31", "Декабрь")
	service.CreateEvent("user1", "2024-01-01", "Январь") // Другой месяц

	events, err := service.GetEventsForMonth("user1", "2023-12-15")
	require.NoError(t, err)
	assert.Len(t, events, 3) // Только события декабря (1, 15, 31 декабря)
}

func TestIsSameDay(t *testing.T) {
	date1 := time.Date(2023, 12, 31, 10, 30, 0, 0, time.UTC)
	date2 := time.Date(2023, 12, 31, 15, 45, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 1, 10, 30, 0, 0, time.UTC)

	assert.True(t, isSameDay(date1, date2))
	assert.False(t, isSameDay(date1, date3))
}

func TestGetWeekStart(t *testing.T) {
	monday := time.Date(2023, 12, 25, 10, 0, 0, 0, time.UTC)
	weekStart := getWeekStart(monday)
	assert.Equal(t, monday, weekStart)

	wednesday := time.Date(2023, 12, 27, 10, 0, 0, 0, time.UTC)
	weekStart = getWeekStart(wednesday)
	assert.Equal(t, monday, weekStart)

	sunday := time.Date(2023, 12, 31, 10, 0, 0, 0, time.UTC)
	weekStart = getWeekStart(sunday)
	assert.Equal(t, monday, weekStart)
}
