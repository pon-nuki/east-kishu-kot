package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"east-kishu-kot/backend/automation"
	"east-kishu-kot/backend/excelparser"
)

// App struct
type App struct {
	ctx       context.Context
	company   string
	copyright string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		company:   "Godspeed",
		copyright: "© 2025 Godspeed",
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Login(id string, pass string) error {
	fmt.Println("Login called with:", id)
	_, cancel, err := automation.LoginKOT(id, pass, a.ctx)
	if cancel != nil {
		defer cancel()
	}
	return err
}

func (a *App) ChooseExcelFile() (string, error) {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "勤務表Excelファイルを選択",
		Filters: []runtime.FileFilter{{
			DisplayName: "Excel files",
			Pattern:     "*.xlsx",
		}},
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

func (a *App) RegisterFromExcel(id, pass, filePath string) (map[string]int, error) {
	runtime.LogInfo(a.ctx, "RegisterFromExcel called with: "+filePath)

	records, err := excelparser.ParseWorkRecords(filePath)
	if err != nil {
		return map[string]int{"SuccessCount": 0}, fmt.Errorf("Excel読み取り失敗: %w", err)
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("%d 件の勤務データを検出", len(records)))
	for _, r := range records {
		runtime.LogInfo(a.ctx, fmt.Sprintf("→ %s %s-%s (休憩: %v, 実働: %s)", r.Date, r.StartTime, r.EndTime, r.HasBreak, r.WorkDuration))
	}

	ctx, cancel, err := automation.LoginKOT(id, pass, a.ctx)
	if err != nil {
		return map[string]int{"SuccessCount": 0}, fmt.Errorf("ログイン失敗: %w", err)
	}
	defer cancel()

	var baseID string
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`(() => {
			const select = document.querySelector('select.htBlock-selectOther');
			if (!select) return "";
			const opt = Array.from(select.options).find(o => o.textContent.includes("打刻編集"));
			if (!opt) return "";
			return opt.value.replace("#button_", "").slice(0, -2);
		})()`, &baseID),
	)
	if err != nil || baseID == "" {
		return map[string]int{"SuccessCount": 0}, fmt.Errorf("baseIDの取得に失敗しました: %w", err)
	}

	successCount := 0

	for _, record := range records {
		t, _ := time.Parse("2006-01-02", record.Date)
		actionID := fmt.Sprintf("action_%s%02d", baseID, t.Day())
		runtime.LogInfo(a.ctx, fmt.Sprintf("[%s] 打刻編集画面を開きます (%s)", record.Date, actionID))

		err := chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(`KOT_GLOBAL.KOT_LIB.onClickActionButton("%s")`, actionID), nil),
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			runtime.LogError(a.ctx, fmt.Sprintf("[%s] 打刻編集画面呼び出し失敗: %v", record.Date, err))
			continue
		}

		var formState string
		_ = chromedp.Run(ctx,
			chromedp.Evaluate(`(() => {
				const form = document.querySelector('form[action*="/admin/"]');
				if (!form) return "not_found";
				const style = window.getComputedStyle(form);
				return (style.display === "none") ? "hidden" : "visible";
			})()`, &formState),
		)
		if formState != "visible" {
			runtime.LogError(a.ctx, fmt.Sprintf("[%s] 打刻編集画面が開かれていません (状態: %s)", record.Date, formState))
			continue
		}

		startHour, startMin := parseHM(record.StartTime)
		endHour, endMin := parseHM(record.EndTime)

		start := time.Date(0, 1, 1, startHour, startMin, 0, 0, time.UTC)
		end := time.Date(0, 1, 1, endHour, endMin, 0, 0, time.UTC)
		total := end.Sub(start)
		parsed, err := time.Parse("15:04", record.WorkDuration)
		if err != nil {
			runtime.LogError(a.ctx, fmt.Sprintf("[%s] 実働時間パース失敗: %v", record.Date, err))
			continue
		}
		actual := time.Duration(parsed.Hour())*time.Hour + time.Duration(parsed.Minute())*time.Minute
		breakDuration := total - actual

		tasks := []chromedp.Action{
			chromedp.Click(`#recording_timestamp_time_1`, chromedp.ByID),
			chromedp.SendKeys(`#recording_timestamp_time_1`, fmt.Sprintf("%02d%02d", startHour, startMin), chromedp.ByID),
			chromedp.SetAttributeValue(`#recording_timestamp_hour_1`, "value", fmt.Sprintf("%02d", startHour)),
			chromedp.SetAttributeValue(`#recording_timestamp_minute_1`, "value", fmt.Sprintf("%02d", startMin)),
			chromedp.SetValue(`#recording_type_code_1`, "1", chromedp.ByID),
		}

		if breakDuration > 0 && record.HasBreak {
			mid := start.Add(total / 2)
			breakStartTime := mid.Add(-breakDuration / 2)
			breakEndTime := mid.Add(breakDuration / 2)

			tasks = append(tasks,
				chromedp.Click(`#recording_timestamp_time_2`, chromedp.ByID),
				chromedp.SendKeys(`#recording_timestamp_time_2`, fmt.Sprintf("%02d%02d", breakStartTime.Hour(), breakStartTime.Minute()), chromedp.ByID),
				chromedp.SetAttributeValue(`#recording_timestamp_hour_2`, "value", fmt.Sprintf("%02d", breakStartTime.Hour())),
				chromedp.SetAttributeValue(`#recording_timestamp_minute_2`, "value", fmt.Sprintf("%02d", breakStartTime.Minute())),
				chromedp.SetValue(`#recording_type_code_2`, "3", chromedp.ByID),

				chromedp.Click(`#recording_timestamp_time_3`, chromedp.ByID),
				chromedp.SendKeys(`#recording_timestamp_time_3`, fmt.Sprintf("%02d%02d", breakEndTime.Hour(), breakEndTime.Minute()), chromedp.ByID),
				chromedp.SetAttributeValue(`#recording_timestamp_hour_3`, "value", fmt.Sprintf("%02d", breakEndTime.Hour())),
				chromedp.SetAttributeValue(`#recording_timestamp_minute_3`, "value", fmt.Sprintf("%02d", breakEndTime.Minute())),
				chromedp.SetValue(`#recording_type_code_3`, "4", chromedp.ByID),

				chromedp.Click(`#recording_timestamp_time_4`, chromedp.ByID),
				chromedp.SendKeys(`#recording_timestamp_time_4`, fmt.Sprintf("%02d%02d", endHour, endMin), chromedp.ByID),
				chromedp.SetAttributeValue(`#recording_timestamp_hour_4`, "value", fmt.Sprintf("%02d", endHour)),
				chromedp.SetAttributeValue(`#recording_timestamp_minute_4`, "value", fmt.Sprintf("%02d", endMin)),
				chromedp.SetValue(`#recording_type_code_4`, "2", chromedp.ByID),
			)
		} else {
			tasks = append(tasks,
				chromedp.Click(`#recording_timestamp_time_2`, chromedp.ByID),
				chromedp.SendKeys(`#recording_timestamp_time_2`, fmt.Sprintf("%02d%02d", endHour, endMin), chromedp.ByID),
				chromedp.SetAttributeValue(`#recording_timestamp_hour_2`, "value", fmt.Sprintf("%02d", endHour)),
				chromedp.SetAttributeValue(`#recording_timestamp_minute_2`, "value", fmt.Sprintf("%02d", endMin)),
				chromedp.SetValue(`#recording_type_code_2`, "2", chromedp.ByID),
			)
		}

		tasks = append(tasks, chromedp.Sleep(1*time.Second))

		err = chromedp.Run(ctx, tasks...)
		if err != nil {
			runtime.LogError(a.ctx, fmt.Sprintf("[%s] 打刻入力失敗: %v", record.Date, err))
			continue
		}

		err = chromedp.Run(ctx,
			chromedp.Evaluate(`KOT_GLOBAL.KOT_LIB.onClickActionButton("action_01")`, nil),
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			runtime.LogError(a.ctx, fmt.Sprintf("[%s] 登録失敗: %v", record.Date, err))
			continue
		}

		runtime.LogInfo(a.ctx, fmt.Sprintf("[%s] 登録完了", record.Date))
		successCount++
	}

	runtime.LogInfo(a.ctx, fmt.Sprintf("%d 件登録完了", successCount))
	return map[string]int{
		"SuccessCount": successCount,
	}, nil
}

func parseHM(hm string) (int, int) {
	parts := strings.Split(hm, ":")
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	return h, m
}

func (a *App) GetAppMeta() map[string]string {
	return map[string]string{
		"company":   a.company,
		"copyright": a.copyright,
	}
}
