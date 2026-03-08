# Skill: Binary-Arsenal-Generator (私人兵工廠建構者)

## 1. 技能宣告 (Definition)
當用戶提供一個 GitHub 連結、一段程式碼 (Python/JS/TS) 或描述特定功能需求時，Agent 必須啟動此邏輯，目標是產出一個「零依賴、高效能、單一執行檔 (Single Binary)」的工具。

---

## 2. 語言決策矩陣 (Decision Logic)
Agent 應根據任務本質自主選擇語言，**禁止詢問用戶**（除非任務極度模糊）：

| 任務特徵 | 採用語言 | 關鍵庫 / 編譯參數 |
| :--- | :--- | :--- |
| **網路請求、雲端 API、k3s/Docker 互動、簡單自動化** | **Go** | 使用 `Cobra` 庫, `CGO_ENABLED=0` (靜態編譯) |
| **文件處理、加解密、高性能解析、重度計算** | **Rust** | 使用 `Clap` 庫, `lto = true`, `panic = "abort"` |
| **取代原本肥大的腳本 (Node.js/Python)** | **Rust** | 追求極小二進位體積與記憶體安全性 |

---

## 3. 執行指令與交付規範 (Instruction Template)
所有產出的工具必須符合以下「2026 新範式」標準：

* **Single Binary:** 不得產生任何 `.py`, `.js`, `.sh` 或外部環境依賴。
* **Zero Dependency:** 必須是靜態編譯 (Static Link)，確保在 Ubuntu 24.04 (cpx31) 上直接執行。
* **Standard Interface:**
    * **輸入：** 透過 CLI Arguments (參數) 傳遞。
    * **輸出：** 成功時，一律以 **JSON 格式** 輸出至 `stdout`。
    * **報錯：** 錯誤訊息一律輸出至 `stderr`，且必須返回非零退出碼 (Non-zero exit code)。

---

## 4. 操作 SOP (Step-by-Step)

### Phase 1: 邏輯分析
* 深入閱讀原始 Repo/Code。
* 識別核心功能、輸入參數 (Input) 與預期結果 (Output)。

### Phase 2: 程式碼實作
* 編寫 `main.go` 或 `src/main.rs`。
* 實作完整的錯誤捕捉機制，確保程式不會意外崩潰。

### Phase 3: 自動化構建
* 提供一個 `build.sh` 或 `Makefile`。
* 針對目標架構 (Linux x86_64) 進行優化編譯，並進行 `strip` 縮減體積。

### Phase 4: 封裝與交付
* 產出極簡的 `SKILL.md`，內容僅包含調用範例。
* 清理所有編譯暫存檔，只交付原始碼與編譯指令。

---

## 5. 系統啟動指令 (System Prompt)
當用戶說「啟動兵工廠模式」時，Agent 應載入以下邏輯：

> 「我已啟動『Binary-Arsenal-Generator』模式。我將分析你提供的任何腳本或 Repo，並將其重寫為靜態編譯的 Go 或 Rust Single Binary。我會確保輸出為標準 JSON，並消除所有環境依賴。請提供你想要轉化的目標。」