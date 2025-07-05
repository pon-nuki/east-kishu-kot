package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/sys/windows/registry"
)

// Chrome または Edge の実行パスを取得（Chrome優先）
func getBrowserExecPath() (string, error) {
	// Chrome優先
	if path, err := getRegistryBrowserPath(`SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\chrome.exe`); err == nil {
		return path, nil
	}
	// 次点：Edge
	if path, err := getRegistryBrowserPath(`SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\msedge.exe`); err == nil {
		return path, nil
	}
	// 両方なし
	return "", fmt.Errorf("ChromeまたはEdgeが見つかりませんでした")
}

// レジストリからブラウザパスを取得するヘルパー
func getRegistryBrowserPath(keyPath string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	path, _, err := k.GetStringValue("")
	if err != nil {
		return "", err
	}
	return path, nil
}

func LoginKOT(id, pass string, wailsCtx context.Context) (context.Context, context.CancelFunc, error) {
	// ブラウザ実行パスを取得（Chrome→Edge）
	browserPath, err := getBrowserExecPath()
	if err != nil {
		runtime.LogError(wailsCtx, "対応ブラウザが見つかりません: "+err.Error())
		return nil, nil, err
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(browserPath),
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1000),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, taskCancel := chromedp.NewContext(allocCtx)

	runtime.LogInfo(wailsCtx, "[chromedp] Step 1: Navigating to login page")
	err = chromedp.Run(ctx,
		chromedp.Navigate("https://login.ta.kingoftime.jp/admin"),
	)
	if err != nil {
		runtime.LogError(wailsCtx, "Step 1 failed: "+err.Error())
		goto cancel
	}

	runtime.LogInfo(wailsCtx, "[chromedp] Step 2: Waiting for login ID field")
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`#login_id`, chromedp.ByID),
	)
	if err != nil {
		runtime.LogError(wailsCtx, "Step 2 failed: "+err.Error())
		goto cancel
	}

	runtime.LogInfo(wailsCtx, "[chromedp] Step 3: Entering ID and password")
	err = chromedp.Run(ctx,
		chromedp.Click(`#login_id`, chromedp.ByID),
		chromedp.SendKeys(`#login_id`, id, chromedp.ByID),
		chromedp.Click(`#login_password`, chromedp.ByID),
		chromedp.SendKeys(`#login_password`, pass, chromedp.ByID),
	)
	if err != nil {
		runtime.LogError(wailsCtx, "Step 3 failed: "+err.Error())
		goto cancel
	}

	runtime.LogInfo(wailsCtx, "[chromedp] Step 4: Clicking login button")
	err = chromedp.Run(ctx,
		chromedp.Click(`#login_button`, chromedp.ByID),
	)
	if err != nil {
		runtime.LogError(wailsCtx, "Step 4 failed: "+err.Error())
		goto cancel
	}

	runtime.LogInfo(wailsCtx, "[chromedp] Step 5: Waiting for post-login body")
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		runtime.LogError(wailsCtx, "Step 5 failed: "+err.Error())
		goto cancel
	}

	runtime.LogInfo(wailsCtx, "[chromedp] Step 6: Skipped - no 勤務表リンクをクリックする必要なし")

	runtime.LogInfo(wailsCtx, "[chromedp] Step 7: Waiting for 勤務表 table (.htBlock-adjastableTableF)")
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`.htBlock-adjastableTableF`, chromedp.ByQuery),
	)
	if err != nil {
		runtime.LogError(wailsCtx, "Step 7 failed: "+err.Error())
		goto cancel
	}

	runtime.LogInfo(wailsCtx, "LoginKOT 完了")
	return ctx, func() {
		taskCancel()
		allocCancel()
	}, nil

cancel:
	runtime.LogError(wailsCtx, "LoginKOT 全体失敗")
	allocCancel()
	taskCancel()
	return nil, nil, err
}
