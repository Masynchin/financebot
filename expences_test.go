package main

import (
	"log"
	"testing"
)

var s *ExpenceService

func setNewExpenceService() {
	service, err := NewExpenceService(":memory:")
	if err != nil {
		log.Fatal("Cannot create expence service")
	}
	s = service
}

func TestInsertExpence(t *testing.T) {
	setNewExpenceService()
	defer s.Close()

	userID := int64(1)
	category := "Food"
	amount := 100
	expenceID, err := s.InsertExpence(userID, category, amount)
	if err != nil {
		t.Fatalf("Error while inserting expence: %v", err.Error())
	}

	e, err := s.GetExpenceByID(expenceID)
	if err != nil {
		t.Fatalf("Cannot get just inserted expence: %v", err.Error())
	}

	if e.Id != expenceID || e.UserID != userID || e.Category != category || e.Amount != amount {
		warning := "Inserted expence and getted are different - %v, (%v, %v, %v, %v)"
		t.Fatalf(warning, e, expenceID, userID, category, amount)
	}
}

func TestGetExpence(t *testing.T) {
	setNewExpenceService()
	defer s.Close()

	userID := int64(1)
	category := "Food"
	expectedExpenceAmount := 0
	for _, amount := range [...]int{100, 200, 300} {
		s.InsertExpence(userID, category, amount)
		expectedExpenceAmount += amount
	}

	expenceAmount, err := s.GetExpence(userID, category)
	if err != nil {
		t.Fatalf("Error while getting expence: %v", err.Error())
	}

	if expenceAmount != expectedExpenceAmount {
		t.Fatalf("Unexpected expence amount - got %v, expected %v", expenceAmount, expectedExpenceAmount)
	}
}

func TestGetUserExpences(t *testing.T) {
	setNewExpenceService()
	defer s.Close()

	userID := int64(1)
	anotherUserID := int64(2)
	category := "Food"
	expencesAmounts := [...]int{100, 200, 300}
	for _, amount := range expencesAmounts {
		s.InsertExpence(userID, category, amount)
		s.InsertExpence(anotherUserID, category, amount)
	}

	userExpences, err := s.GetUserExpences(userID)
	if err != nil {
		t.Fatalf("Error while getting user expences: %v", err.Error())
	}

	if len(userExpences) != len(expencesAmounts) {
		warning := "User expences count not the same with inserted: got %v, inserted: %v"
		t.Fatalf(warning, len(userExpences), len(expencesAmounts))
	}
	for i, expence := range userExpences {
		if expence.Amount != expencesAmounts[i] {
			t.Fatalf("User expence amount not the same: got %v, must: %v", expence.Amount, expencesAmounts[i])
		} else if expence.UserID != userID {
			t.Fatalf("Got another user expence: userID %v, got: %v", userID, expence.UserID)
		}
	}
}

func TestGetUserExpencesGroupedByCategory(t *testing.T) {
	setNewExpenceService()
	defer s.Close()

	userID := int64(1)
	anotherUserID := int64(2)
	categories := [...]string{"Food", "Taxi"}
	amounts := [...]int{100, 200, 300}
	expectedAmounts := make(map[string]int)
	for _, category := range categories {
		for _, amount := range amounts {
			s.InsertExpence(userID, category, amount)
			expectedAmounts[category] += amount
			s.InsertExpence(anotherUserID, category, amount)
		}
	}

	userExpences, err := s.GetUserExpencesGroupedByCategory(userID)
	if err != nil {
		t.Fatalf("Cannot get user expences group by category: %v", err.Error())
	}

	for _, exp := range userExpences {
		expectedAmount := expectedAmounts[exp.Category]
		if expectedAmount != exp.Amount {
			warning := "Invalid user expence grouped by category: category %v, expected %v, got %v"
			t.Fatalf(warning, exp.Category, expectedAmount, exp.Amount)
		} else if exp.UserID != userID {
			t.Fatalf("Got another user expence grouped by category: got %v, must %v", exp.UserID, userID)
		}
	}
}

func TestDeleteExpence(t *testing.T) {
	setNewExpenceService()
	defer s.Close()

	userID := int64(1)
	category := "Food"
	amount := 100
	expenceID, _ := s.InsertExpence(userID, category, amount)

	if err := s.DeleteExpence(expenceID); err != nil {
		t.Fatalf("Error while deleting expence: %v", err.Error())
	}
}
