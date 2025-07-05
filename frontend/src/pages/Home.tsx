import React, { useState, useEffect } from 'react';
import kihoIcon from "../assets/ki-ho.png";

const Home = () => {
  const [excelPath, setExcelPath] = useState("");
  const [userId, setUserId] = useState("");
  const [password, setPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [summary, setSummary] = useState("");
  const [meta, setMeta] = useState({ company: '', copyright: '' });

  const handleChooseExcel = async () => {
    try {
      const path = await window.go.main.App.ChooseExcelFile();
      if (path) {
        setExcelPath(path);
        setSummary("");
      }
    } catch {}
  };

  const handleRegisterFromExcel = async () => {
    if (isSubmitting || !excelPath) {
      alert("先にファイルを選択してください");
      return;
    }

    try {
      setIsSubmitting(true);
      setSummary("");
      const result = await window.go.main.App.RegisterFromExcel(userId, password, excelPath);
      if (result && typeof result.SuccessCount === "number") {
        setSummary(`${result.SuccessCount} 件登録しました`);
      } else {
        setSummary("結果が取得できませんでした");
      }
    } catch (err) {
      alert("登録処理失敗: " + err);
      setSummary("登録に失敗しました");
    } finally {
      setIsSubmitting(false);
    }
  };

  useEffect(() => {
    window.go.main.App.GetAppMeta().then(setMeta);
  }, []);

return (
  <div className="flex flex-col h-screen bg-gray-900 text-white overflow-x-hidden">

    {/* タイトル部 */}
    <div className="flex items-center justify-center pt-6">
      <img
        src={kihoIcon}
        alt="きーほくん"
        className="w-8 h-8 mr-2 rounded-full shadow"
      />
      <h1 className="text-2xl font-bold text-white text-center whitespace-nowrap">
        東紀州KOT自動入力ツール（試作版）
      </h1>
    </div>

    <hr className="mt-4 border-t border-gray-500 w-4/5 mx-auto" />

    {/* メイン部 */}
    <main className="flex-grow flex justify-center items-start py-8 px-4">
      <div className="bg-white text-black p-6 rounded-xl shadow-lg w-full max-w-[480px]">

        {/* ログインID */}
        <div className="mb-5">
          <label className="block text-sm font-semibold text-gray-800 mb-2">
            ログインID
          </label>
          <input
            type="text"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
            className="w-full border border-gray-300 rounded-lg p-3 text-sm shadow-sm focus:ring-2 focus:ring-purple-400 focus:outline-none"
            placeholder="例: lukXX-XXXXXX"
          />
        </div>

        {/* パスワード */}
        <div>
          <label className="block text-sm font-semibold text-gray-800 mb-2">
            パスワード
          </label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full border border-gray-300 rounded-lg p-3 text-sm shadow-sm focus:ring-2 focus:ring-purple-400 focus:outline-none"
          />
        </div>

        {/* ファイルパス表示 */}
        {excelPath && (
          <div className="mt-4 text-xs text-gray-700 break-all">
            選択中のファイル：<br />
            <span className="font-mono">{excelPath}</span>
          </div>
        )}

        {/* ボタン群（常に横並び、等幅） */}
        <div style={{ marginTop: '24px' }}>
          <div className="mt-6 flex flex-row gap-3">
            {/* Excelファイル選択ボタン */}
            <button
              className="w-1/2 bg-yellow-400 hover:bg-yellow-500 text-white font-semibold text-sm px-4 py-2 rounded-lg shadow-md transition duration-200 ease-in-out"
              onClick={handleChooseExcel}
            >
              Excel選択
            </button>

            {/* 自動打刻ボタン */}
            <button
              className={`w-1/2 text-sm px-4 py-2 rounded-lg font-semibold shadow-md transition duration-200 ease-in-out ${
                !excelPath || !userId || !password || isSubmitting
                  ? "bg-gray-400 cursor-not-allowed text-white"
                  : "bg-purple-600 hover:bg-purple-700 text-white cursor-pointer"
              }`}
              onClick={handleRegisterFromExcel}
              disabled={!excelPath || !userId || !password || isSubmitting}
            >
              {isSubmitting ? "打刻中…" : "Excelから一括自動打刻！"}
            </button>
          </div>
        </div>
        {/* サマリー */}
        {summary && (
          <div className="mt-4 text-sm text-center font-medium text-green-600">
            {summary}
          </div>
        )}
      </div>
    </main>

    {/* フッター */}
    <footer className="text-[10px] text-center text-gray-500 py-2">
      © {new Date().getFullYear()} Godspeed — <span className="italic">Attack from East Kishu!</span>
    </footer>
  </div>
);





};

export default Home;
