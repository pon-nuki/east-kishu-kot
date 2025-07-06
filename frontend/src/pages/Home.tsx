import React, { useState, useEffect } from 'react';
import kihoIcon from "../assets/ki-ho.png";
import { motion } from "framer-motion";
import backgroundImg from "../assets/IMG_3278.jpg";

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
    <div className="relative h-screen w-screen overflow-hidden">
      {/* 背景画像をdivで指定 */}
      <div
        className="fixed top-0 left-0 w-full h-full bg-cover bg-center brightness-75 z-0"
        style={{ backgroundImage: `url(${backgroundImg})` }}
      />

      {/* 暗めのぼかしレイヤー */}
      <div className="fixed top-0 left-0 w-full h-full bg-gray-900/60 backdrop-blur-sm z-10" />

      {/* UI全体（← z-20で前面へ） */}
      <div className="relative z-20 flex flex-col h-full w-full text-white">
        {/* タイトル部 */}
        <div className="flex flex-col items-center justify-center pt-6">
          <div className="flex items-center mb-1">
            <img
              src={kihoIcon}
              alt="きーほくん"
              className="w-8 h-8 mr-2 rounded-full shadow"
            />
            <h1 className="text-2xl font-bold text-white text-center whitespace-nowrap">
              東紀州KOT自動入力ツール（試作版）
            </h1>
          </div>
          <motion.p
            className="text-xs text-gray-300 italic"
            style={{ marginTop: "8px" }}
            initial={{ opacity: 0, y: -4 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 1.8, delay: 0.8 }}
          >
            Made with ❤️ in Kii-Nagashima, Mie
          </motion.p>
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

            {/* ボタン群 */}
            <div style={{ marginTop: "20px" }}>
              <div style={{ display: "flex", gap: "12px" }}>
                <button
                  style={{
                    flex: 1,
                    padding: "10px 16px",
                    fontSize: "14px",
                    fontWeight: 500,
                    lineHeight: "1.5",
                    borderRadius: "8px",
                    backgroundColor: "#475569",
                    color: "#fff",
                    boxShadow: "0 1px 4px rgba(0, 0, 0, 0.1)",
                    cursor: "pointer",
                    transition: "background-color 0.2s ease-in-out",
                  }}
                  onMouseOver={(e) =>
                    (e.currentTarget.style.backgroundColor = "#334155")
                  }
                  onMouseOut={(e) =>
                    (e.currentTarget.style.backgroundColor = "#475569")
                  }
                  onClick={handleChooseExcel}
                >
                  Excel選択
                </button>

                <button
                  style={{
                    flex: 1,
                    padding: "10px 16px",
                    fontSize: "14px",
                    fontWeight: 500,
                    lineHeight: "1.5",
                    borderRadius: "8px",
                    backgroundColor:
                      !excelPath || !userId || !password || isSubmitting
                        ? "#9ca3af"
                        : "#6366f1",
                    color: "#fff",
                    cursor:
                      !excelPath || !userId || !password || isSubmitting
                        ? "not-allowed"
                        : "pointer",
                    pointerEvents:
                      !excelPath || !userId || !password || isSubmitting
                        ? "none"
                        : "auto",
                    boxShadow: "0 1px 4px rgba(0, 0, 0, 0.1)",
                    transition: "background-color 0.2s ease-in-out",
                  }}
                  onMouseOver={(e) => {
                    if (excelPath && userId && password && !isSubmitting) {
                      e.currentTarget.style.backgroundColor = "#4f46e5";
                    }
                  }}
                  onMouseOut={(e) => {
                    if (excelPath && userId && password && !isSubmitting) {
                      e.currentTarget.style.backgroundColor = "#6366f1";
                    }
                  }}
                  onClick={handleRegisterFromExcel}
                  disabled={!excelPath || !userId || !password || isSubmitting}
                >
                  {isSubmitting ? "打刻中…" : "Excelから一括自動打刻！"}
                </button>
              </div>
            </div>

            {summary && (
              <div className="mt-4 text-sm text-center font-medium text-green-600">
                {summary}
              </div>
            )}
          </div>
        </main>

        {/* フッター */}
        <footer className="text-[10px] text-center text-gray-500 py-2">
          © {new Date().getFullYear()} Godspeed —{" "}
          <span className="italic">Attack from East Kishu!</span>
        </footer>
      </div>
    </div>
  );
};

export default Home;
