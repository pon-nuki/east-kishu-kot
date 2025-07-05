package excelparser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type WorkRecord struct {
	Date               string // "2025-07-01"
	StartTime          string // "10:00"
	EndTime            string // "19:00"
	HasBreak           bool   // 勤務時間 > 実働時間
	WorkDuration       string // 表示用の実働時間（例 "2:00"）
	WorkDurationMinute int    // 実働時間（分単位）
}

// Excelファイルをパースして勤務データを返す
func ParseWorkRecords(filePath string) ([]WorkRecord, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	sheet := f.GetSheetName(0)

	// A3 から日付シリアル or 文字列を取得
	rawDate, err := f.GetCellValue(sheet, "A3")
	if err != nil || rawDate == "" {
		return nil, fmt.Errorf("A3セルの日付が取得できません")
	}

	var baseDate time.Time

	// シリアル値 or 日付文字列として解析
	if floatVal, err := strconv.ParseFloat(rawDate, 64); err == nil {
		baseDate, err = excelize.ExcelDateToTime(floatVal, false)
		if err != nil {
			return nil, fmt.Errorf("A3のシリアル値を日付に変換できません: %v", err)
		}
	} else {
		layouts := []string{
			"2006/01/02", "2006-01-02", "2006年1月2日", "2006年1月", "2006年",
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, rawDate); err == nil {
				baseDate = t
				break
			}
		}
		if baseDate.IsZero() {
			return nil, fmt.Errorf("A3の日付形式が不正です: \"%s\"", rawDate)
		}
	}

	var records []WorkRecord

	// C列: 開始時刻, D列: 終了時刻, E列: 実働時間, 9行目から下
	for i := 9; ; i++ {
		startCell, err1 := f.GetCellValue(sheet, fmt.Sprintf("C%d", i))
		endCell, err2 := f.GetCellValue(sheet, fmt.Sprintf("D%d", i))
		totalCellRaw, _ := f.GetCellValue(sheet, fmt.Sprintf("E%d", i))

		if err1 != nil && err2 != nil {
			break
		}

		start := strings.TrimSpace(startCell)
		end := strings.TrimSpace(endCell)
		total := strings.TrimSpace(totalCellRaw)

		if start == "" || end == "" {
			continue
		}

		startTime, err := time.Parse("15:04", start)
		if err != nil {
			continue
		}
		endTime, err := time.Parse("15:04", end)
		if err != nil {
			continue
		}

		duration := endTime.Sub(startTime)
		actualMinutes := int(duration.Minutes())

		// 実働時間（E列）を文字列＋分数で取得
		var totalText string
		var totalMinutes int
		if total == "" {
			totalText = ""
			totalMinutes = 0
		} else if floatVal, err := strconv.ParseFloat(total, 64); err == nil {
			// シリアル値（例: 0.0833）
			totalDuration := time.Duration(floatVal * 24 * float64(time.Hour))
			totalMinutes = int(totalDuration.Minutes())
			h := totalMinutes / 60
			m := totalMinutes % 60
			totalText = fmt.Sprintf("%d:%02d", h, m)
		} else {
			if t, err := time.Parse("15:04", total); err == nil {
				totalMinutes = t.Hour()*60 + t.Minute()
				totalText = total
			} else {
				totalText = ""
				totalMinutes = 0
			}
		}

		hasBreak := totalMinutes > 0 && totalMinutes < actualMinutes

		dayOffset := i - 9
		workDate := baseDate.AddDate(0, 0, dayOffset)

		records = append(records, WorkRecord{
			Date:               workDate.Format("2006-01-02"),
			StartTime:          start,
			EndTime:            end,
			HasBreak:           hasBreak,
			WorkDuration:       totalText,
			WorkDurationMinute: totalMinutes,
		})
	}

	return records, nil
}
