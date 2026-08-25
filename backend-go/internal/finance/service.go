package finance

import (
	"math"
	"time"

	"gorm.io/gorm"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/models"
)

type Entry struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
	CreatedByID string    `json:"createdById"`
	CreatedBy   *EntryUser `json:"createdBy,omitempty"`
	CreatedAt   string    `json:"createdAt"`
	UpdatedAt   string    `json:"updatedAt"`
}

type EntryUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Summary struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

type Service struct {
	DB    *gorm.DB
	Sheet *SheetsRepository
}

func NewService(db *gorm.DB, sheet *SheetsRepository) *Service {
	return &Service{DB: db, Sheet: sheet}
}

func (s *Service) activeRows(from, to string) ([]Row, error) {
	rows, err := s.Sheet.FetchAllRows()
	if err != nil {
		return nil, err
	}

	filtered := make([]Row, 0, len(rows))
	for _, row := range rows {
		if row.ID == "" || row.DeletedAt != "" {
			continue
		}
		if from != "" && row.Date < from {
			continue
		}
		if to != "" && row.Date > to {
			continue
		}
		filtered = append(filtered, row)
	}
	return filtered, nil
}

func (s *Service) userMap(rows []Row) (map[string]EntryUser, error) {
	ids := make([]string, 0, len(rows))
	seen := map[string]bool{}
	for _, row := range rows {
		if row.CreatedByID != "" && !seen[row.CreatedByID] {
			seen[row.CreatedByID] = true
			ids = append(ids, row.CreatedByID)
		}
	}
	if len(ids) == 0 {
		return map[string]EntryUser{}, nil
	}

	var users []models.User
	if err := s.DB.Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}

	result := make(map[string]EntryUser, len(users))
	for _, u := range users {
		result[u.ID] = EntryUser{FirstName: u.FirstName, LastName: u.LastName}
	}
	return result, nil
}

func toEntry(row Row, users map[string]EntryUser) Entry {
	entry := Entry{
		ID:          row.ID,
		Type:        row.Type,
		Category:    row.Category,
		Amount:      math.Abs(row.Amount),
		Description: row.Description,
		Date:        row.Date,
		CreatedByID: row.CreatedByID,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
	if u, ok := users[row.CreatedByID]; ok {
		entry.CreatedBy = &u
	}
	return entry
}

func (s *Service) GetEntries(from, to string) ([]Entry, error) {
	rows, err := s.activeRows(from, to)
	if err != nil {
		return nil, err
	}
	users, err := s.userMap(rows)
	if err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, toEntry(row, users))
	}
	return entries, nil
}

func (s *Service) GetSummary(from, to string) (*Summary, error) {
	rows, err := s.activeRows(from, to)
	if err != nil {
		return nil, err
	}

	summary := &Summary{}
	for _, row := range rows {
		if row.Type == string(models.FinanceIncome) {
			summary.Income += math.Abs(row.Amount)
		} else if row.Type == string(models.FinanceExpense) {
			summary.Expense += math.Abs(row.Amount)
		}
	}
	summary.Balance = summary.Income - summary.Expense
	return summary, nil
}

type CreateInput struct {
	Type        string
	Category    string
	Amount      float64
	Description string
	Date        string
	CreatedByID string
}

func (s *Service) CreateEntry(input CreateInput) (*Entry, error) {
	rows, err := s.Sheet.FetchAllRows()
	if err != nil {
		return nil, err
	}

	id := s.Sheet.NextSequentialID(rows)
	now := time.Now().Format(time.RFC3339)

	signedAmount := math.Abs(input.Amount)
	if input.Type == string(models.FinanceExpense) {
		signedAmount = -signedAmount
	}

	row := Row{
		ID:          id,
		Type:        input.Type,
		Category:    input.Category,
		Amount:      signedAmount,
		Description: input.Description,
		Date:        input.Date,
		CreatedByID: input.CreatedByID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.Sheet.AppendRow(row); err != nil {
		return nil, err
	}

	users, err := s.userMap([]Row{row})
	if err != nil {
		return nil, err
	}

	entry := toEntry(row, users)
	return &entry, nil
}

func (s *Service) findRow(id string) (*Row, error) {
	rows, err := s.Sheet.FetchAllRows()
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row.ID == id && row.DeletedAt == "" {
			return &row, nil
		}
	}
	return nil, apperror.New("Finance entry tidak ditemukan", 404)
}

type UpdateInput struct {
	Type        *string
	Category    *string
	Amount      *float64
	Description *string
	Date        *string
}

func (s *Service) UpdateEntry(id string, input UpdateInput) (*Entry, error) {
	row, err := s.findRow(id)
	if err != nil {
		return nil, err
	}

	entryType := row.Type
	if input.Type != nil {
		entryType = *input.Type
	}
	if input.Category != nil {
		row.Category = *input.Category
	}
	if input.Description != nil {
		row.Description = *input.Description
	}
	if input.Date != nil {
		row.Date = *input.Date
	}

	absAmount := math.Abs(row.Amount)
	if input.Amount != nil {
		absAmount = math.Abs(*input.Amount)
	}
	signedAmount := absAmount
	if entryType == string(models.FinanceExpense) {
		signedAmount = -absAmount
	}
	row.Type = entryType
	row.Amount = signedAmount
	row.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := s.Sheet.UpdateRow(*row); err != nil {
		return nil, err
	}

	users, err := s.userMap([]Row{*row})
	if err != nil {
		return nil, err
	}

	entry := toEntry(*row, users)
	return &entry, nil
}

func (s *Service) DeleteEntry(id string) error {
	row, err := s.findRow(id)
	if err != nil {
		return err
	}
	return s.Sheet.SoftDeleteRow(row.RowNumber)
}
