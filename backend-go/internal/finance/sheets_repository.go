package finance

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"

	"nora-photobooth-backend/internal/config"
)

// Row is the raw representation of one FinanceEntries sheet row
// (columns A-J: id,type,category,amount,description,date,createdById,
// createdAt,updatedAt,deletedAt).
type Row struct {
	RowNumber   int // 1-based sheet row number (for targeted updates)
	ID          string
	Type        string
	Category    string
	Amount      float64 // signed as stored in the sheet (negative for EXPENSE)
	Description string
	Date        string
	CreatedByID string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}

type SheetsRepository struct {
	svc           *sheets.Service
	spreadsheetID string
}

func NewSheetsRepository(svc *sheets.Service, spreadsheetID string) *SheetsRepository {
	return &SheetsRepository{svc: svc, spreadsheetID: spreadsheetID}
}

func (r *SheetsRepository) FetchAllRows() ([]Row, error) {
	resp, err := r.svc.Spreadsheets.Values.Get(r.spreadsheetID, config.FinanceRange).Do()
	if err != nil {
		return nil, err
	}

	rows := make([]Row, 0, len(resp.Values))
	for i, raw := range resp.Values {
		rows = append(rows, parseRow(i+2, raw))
	}
	return rows, nil
}

func parseRow(rowNumber int, raw []interface{}) Row {
	get := func(idx int) string {
		if idx < len(raw) {
			return fmt.Sprintf("%v", raw[idx])
		}
		return ""
	}

	amount, _ := strconv.ParseFloat(get(3), 64)

	return Row{
		RowNumber:   rowNumber,
		ID:          get(0),
		Type:        get(1),
		Category:    get(2),
		Amount:      amount,
		Description: get(4),
		Date:        get(5),
		CreatedByID: get(6),
		CreatedAt:   get(7),
		UpdatedAt:   get(8),
		DeletedAt:   get(9),
	}
}

func (r *SheetsRepository) NextSequentialID(rows []Row) string {
	max := 0
	for _, row := range rows {
		trimmed := strings.TrimPrefix(row.ID, "ID-")
		if n, err := strconv.Atoi(trimmed); err == nil && n > max {
			max = n
		}
	}
	return fmt.Sprintf("ID-%d", max+1)
}

func (r *SheetsRepository) AppendRow(row Row) error {
	values := [][]interface{}{{
		row.ID, row.Type, row.Category, row.Amount, row.Description,
		row.Date, row.CreatedByID, row.CreatedAt, row.UpdatedAt, row.DeletedAt,
	}}

	_, err := r.svc.Spreadsheets.Values.Append(r.spreadsheetID, config.FinanceRange, &sheets.ValueRange{Values: values}).
		ValueInputOption("RAW").Do()
	return err
}

func (r *SheetsRepository) UpdateRow(row Row) error {
	rangeStr := fmt.Sprintf("%s!A%d:J%d", config.FinanceSheetName, row.RowNumber, row.RowNumber)
	values := [][]interface{}{{
		row.ID, row.Type, row.Category, row.Amount, row.Description,
		row.Date, row.CreatedByID, row.CreatedAt, row.UpdatedAt, row.DeletedAt,
	}}

	_, err := r.svc.Spreadsheets.Values.Update(r.spreadsheetID, rangeStr, &sheets.ValueRange{Values: values}).
		ValueInputOption("RAW").Do()
	return err
}

// SoftDeleteRow writes only column J (deletedAt), matching the old
// finance.service.ts behaviour of a single-cell update rather than removing
// the row.
func (r *SheetsRepository) SoftDeleteRow(rowNumber int) error {
	rangeStr := fmt.Sprintf("%s!J%d:J%d", config.FinanceSheetName, rowNumber, rowNumber)
	values := [][]interface{}{{time.Now().Format(time.RFC3339)}}

	_, err := r.svc.Spreadsheets.Values.Update(r.spreadsheetID, rangeStr, &sheets.ValueRange{Values: values}).
		ValueInputOption("RAW").Do()
	return err
}
