package simulation

import (
	"math/rand"
	"testing"
)

// mockSystem is a test double that records calls to Update.
type mockSystem struct {
	calls    [][]interface{}
	onUpdate func(week int, state *GameState, rng *rand.Rand) []Event
}

func (m *mockSystem) Update(week int, state *GameState, rng *rand.Rand) []Event {
	m.calls = append(m.calls, []interface{}{week, state, rng})
	if m.onUpdate != nil {
		return m.onUpdate(week, state, rng)
	}
	return nil
}

func TestNewTurnProcessor(t *testing.T) {
	tp := NewTurnProcessor()
	if tp == nil {
		t.Fatal("NewTurnProcessor returned nil")
	}
	if tp.systems == nil {
		t.Error("systems map not initialized")
	}
	if len(tp.order) == 0 {
		t.Error("order slice empty")
	}
	// Ensure order matches SystemOrder
	if len(tp.order) != len(SystemOrder) {
		t.Errorf("order length mismatch: got %d, want %d", len(tp.order), len(SystemOrder))
	}
	for i, typ := range tp.order {
		if typ != SystemOrder[i] {
			t.Errorf("order mismatch at index %d: got %v, want %v", i, typ, SystemOrder[i])
		}
	}
}

func TestRegisterSystem(t *testing.T) {
	tp := NewTurnProcessor()
	mock := &mockSystem{}
	tp.RegisterSystem(SystemEnvironment, mock)
	if tp.systems[SystemEnvironment] != mock {
		t.Error("RegisterSystem did not store system")
	}
	// Replace existing system
	mock2 := &mockSystem{}
	tp.RegisterSystem(SystemEnvironment, mock2)
	if tp.systems[SystemEnvironment] != mock2 {
		t.Error("RegisterSystem did not replace system")
	}
}

func TestTurnProcessor(t *testing.T) {
	tp := NewTurnProcessor()
	// Record call order
	callOrder := []SystemType{}
	// Create a mock that records its type when called
	mockFactory := func(typ SystemType) System {
		return &mockSystem{
			onUpdate: func(week int, state *GameState, rng *rand.Rand) []Event {
				callOrder = append(callOrder, typ)
				return nil
			},
		}
	}
	for _, typ := range SystemOrder {
		tp.RegisterSystem(typ, mockFactory(typ))
	}

	state := NewGameState("test", 123)
	_ = tp.ProcessWeek(state) // events not used yet

	// Verify call order matches SystemOrder
	if len(callOrder) != len(SystemOrder) {
		t.Fatalf("called %d systems, expected %d", len(callOrder), len(SystemOrder))
	}
	for i, typ := range callOrder {
		if typ != SystemOrder[i] {
			t.Errorf("call order mismatch at position %d: got %v, want %v", i, typ, SystemOrder[i])
		}
	}
}

func TestProcessWeekAddsEventsToHistory(t *testing.T) {
	tp := NewTurnProcessor()
	// Create a mock system that returns two events
	event1 := Event{
		ID:        "event-1",
		Type:      "test",
		Timestamp: "2025-W1",
		Data:      map[string]interface{}{"foo": "bar"},
	}
	event2 := Event{
		ID:        "event-2",
		Type:      "test",
		Timestamp: "2025-W1",
		Data:      map[string]interface{}{"baz": 42},
	}
	mock := &mockSystem{
		onUpdate: func(week int, state *GameState, rng *rand.Rand) []Event {
			return []Event{event1, event2}
		},
	}
	tp.RegisterSystem(SystemEnvironment, mock)

	state := NewGameState("test", 999)
	initialHistoryLen := len(state.History)

	events := tp.ProcessWeek(state)

	// Verify returned events
	if len(events) != 2 {
		t.Errorf("ProcessWeek returned %d events, expected 2", len(events))
	}
	// Verify events added to state history
	if len(state.History) != initialHistoryLen+2 {
		t.Errorf("history length %d, expected %d", len(state.History), initialHistoryLen+2)
	}
	// Verify the added events match
	if len(state.History) >= 2 {
		// Events should be appended in order
		if state.History[initialHistoryLen].ID != event1.ID {
			t.Errorf("first history event ID mismatch: got %v, want %v", state.History[initialHistoryLen].ID, event1.ID)
		}
		if state.History[initialHistoryLen+1].ID != event2.ID {
			t.Errorf("second history event ID mismatch: got %v, want %v", state.History[initialHistoryLen+1].ID, event2.ID)
		}
	}
}

func TestProcessWeekAdvancesCalendar(t *testing.T) {
	tp := NewTurnProcessor()
	state := NewGameState("test", 456)
	initialWeek := state.Calendar.Week
	initialYear := state.Calendar.Year

	tp.ProcessWeek(state)

	if state.Calendar.Week != initialWeek+1 {
		t.Errorf("week not advanced: got %d, want %d", state.Calendar.Week, initialWeek+1)
	}
	// Year should not change unless week wraps
	if state.Calendar.Week <= 52 && state.Calendar.Year != initialYear {
		t.Errorf("year changed unexpectedly: got %d, want %d", state.Calendar.Year, initialYear)
	}
}

func TestProcessWeekWeekWrap(t *testing.T) {
	tp := NewTurnProcessor()
	state := NewGameState("test", 789)
	// Set calendar to week 52, year 1
	state.Calendar.Week = 52
	state.Calendar.Year = 1
	state.Calendar.Month = 12 // December

	tp.ProcessWeek(state)

	if state.Calendar.Week != 1 {
		t.Errorf("week after wrap: got %d, want 1", state.Calendar.Week)
	}
	if state.Calendar.Year != 2 {
		t.Errorf("year after wrap: got %d, want 2", state.Calendar.Year)
	}
	if state.Calendar.Month != 1 {
		t.Errorf("month after wrap: got %d, want 1", state.Calendar.Month)
	}
}
