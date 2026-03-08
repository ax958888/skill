# Requirements Document

## Introduction

本文檔定義了一個自動化工作流程系統，用於處理 GitHub skill 的分析、審查、構建和部署。該系統將接收 GitHub URL，執行安全審查和分析，生成標準操作程序（SOP），並通過 forge agent 重建、審查和部署 skill 到 cpx31 伺服器的 k3s 環境中。

## Glossary

- **Skill_Automation_Pipeline**: 完整的自動化工作流程系統，負責從 GitHub URL 輸入到最終部署的所有階段
- **Kiro_Bot**: 負責接收 GitHub URL 並執行分析階段的 agent（@Pojun_kirobot）
- **Forge_Agent**: 負責重建 skill 的 agent（@Forge_coderxbot）
- **GitHub_Skill**: 存放在 GitHub 上的程式碼專案，需要被分析和重建
- **SOP**: Standard Operating Procedure，標準操作程序，描述如何重建 skill 的詳細步驟
- **Security_Analyzer**: 安全審查模組，檢測惡意代碼和安全風險
- **Language_Detector**: 語言識別模組，判斷 skill 使用的程式語言
- **Type_Analyzer**: 類型分析模組，判斷 skill 的功能類型
- **Binary_Arsenal**: 符合零依賴、高效能、單一執行檔規範的產出物
- **Analysis_Report**: 包含安全審查、語言識別、類型分析、使用場景和 SOP 的完整分析結果
- **K3s_Cluster**: 運行在 cpx31 伺服器上的 Kubernetes 集群
- **OpenClaw_Namespace**: k3s 集群中的 openclaw 命名空間，用於部署 skill
- **Root_Directory**: skill 的最終安裝目錄（/root/workspace/agents/）

## Requirements

### Requirement 1: 接收 GitHub URL

**User Story:** 作為用戶，我想要提供 GitHub URL 給系統，以便系統能夠獲取並分析 skill

#### Acceptance Criteria

1. WHEN 用戶提供有效的 GitHub URL，THE Kiro_Bot SHALL 接受該 URL 並開始處理流程
2. WHEN 用戶提供的 GitHub URL 格式無效，THE Kiro_Bot SHALL 返回錯誤訊息並說明正確格式
3. WHEN GitHub URL 指向的儲存庫不存在或無法訪問，THE Kiro_Bot SHALL 返回錯誤訊息並記錄失敗原因
4. THE Kiro_Bot SHALL 支援 HTTPS 格式的 GitHub URL（https://github.com/user/repo）
5. THE Kiro_Bot SHALL 克隆 GitHub 儲存庫到臨時工作目錄

### Requirement 2: 執行安全審查

**User Story:** 作為系統管理員，我想要系統自動檢測惡意代碼，以便確保部署的 skill 是安全的

#### Acceptance Criteria

1. WHEN GitHub_Skill 被下載完成，THE Security_Analyzer SHALL 掃描所有原始碼檔案
2. THE Security_Analyzer SHALL 檢測已知的惡意模式（例如：任意代碼執行、未授權網路連接、檔案系統破壞）
3. THE Security_Analyzer SHALL 檢測可疑的依賴套件和外部連接
4. WHEN 檢測到高風險安全問題，THE Security_Analyzer SHALL 標記為「不安全」並終止流程
5. WHEN 檢測到中低風險問題，THE Security_Analyzer SHALL 在 Analysis_Report 中列出警告
6. THE Security_Analyzer SHALL 生成安全評分（0-100 分）

### Requirement 3: 識別程式語言

**User Story:** 作為系統，我想要自動識別 skill 使用的程式語言，以便選擇正確的構建策略

#### Acceptance Criteria

1. WHEN 安全審查通過，THE Language_Detector SHALL 分析專案檔案結構
2. THE Language_Detector SHALL 識別主要程式語言（Go、Rust、Python、JavaScript、TypeScript 等）
3. THE Language_Detector SHALL 檢測構建配置檔案（go.mod、Cargo.toml、package.json、requirements.txt）
4. WHEN 專案包含多種語言，THE Language_Detector SHALL 識別主要語言和次要語言
5. THE Language_Detector SHALL 在 Analysis_Report 中記錄語言識別結果和信心度

### Requirement 4: 分析 Skill 類型和使用場景

**User Story:** 作為系統，我想要理解 skill 的功能類型和使用場景，以便生成準確的 SOP

#### Acceptance Criteria

1. WHEN 語言識別完成，THE Type_Analyzer SHALL 分析專案的功能類型
2. THE Type_Analyzer SHALL 分類 skill 類型（網路請求、雲端 API、k3s/Docker 互動、檔案處理、加解密、資料解析、自動化工具）
3. THE Type_Analyzer SHALL 識別主要使用場景和目標用戶
4. THE Type_Analyzer SHALL 分析專案的輸入輸出介面
5. THE Type_Analyzer SHALL 在 Analysis_Report 中記錄類型分析和使用場景描述

### Requirement 5: 生成標準操作程序

**User Story:** 作為 Forge Agent，我想要獲得詳細的 SOP，以便準確地重建 skill

#### Acceptance Criteria

1. WHEN 類型分析完成，THE Kiro_Bot SHALL 根據 Binary_Arsenal 規範生成 SOP
2. WHERE 主要語言是 Go，THE Kiro_Bot SHALL 生成 Go 靜態編譯的 SOP（CGO_ENABLED=0）
3. WHERE 主要語言是 Rust，THE Kiro_Bot SHALL 生成 Rust 靜態編譯的 SOP（target-feature=+crt-static）
4. THE SOP SHALL 包含構建命令、依賴處理、測試步驟和產出物驗證
5. THE SOP SHALL 指定 CLI 參數輸入和 JSON 格式輸出的介面規範
6. THE SOP SHALL 包含錯誤處理規範（stderr 輸出）
7. THE SOP SHALL 在 Analysis_Report 中以結構化格式呈現

### Requirement 6: 生成並回傳分析報告

**User Story:** 作為用戶，我想要查看完整的分析結果，以便決定是否繼續交付流程

#### Acceptance Criteria

1. WHEN SOP 生成完成，THE Kiro_Bot SHALL 組合完整的 Analysis_Report
2. THE Analysis_Report SHALL 包含安全審查結果、語言識別、類型分析、使用場景和 SOP
3. THE Analysis_Report SHALL 以 JSON 格式輸出到 stdout
4. THE Analysis_Report SHALL 包含建議的語言選擇（基於 Binary_Arsenal 決策矩陣）
5. THE Kiro_Bot SHALL 等待用戶確認是否繼續交付

### Requirement 7: 處理用戶確認

**User Story:** 作為用戶，我想要能夠確認或拒絕交付，以便控制部署流程

#### Acceptance Criteria

1. WHEN Analysis_Report 被回傳，THE Kiro_Bot SHALL 等待用戶輸入確認指令
2. WHEN 用戶確認交付，THE Kiro_Bot SHALL 將 Analysis_Report 和 SOP 傳遞給 Forge_Agent
3. WHEN 用戶拒絕交付，THE Kiro_Bot SHALL 終止流程並記錄拒絕原因
4. WHEN 用戶在 300 秒內未回應，THE Kiro_Bot SHALL 自動終止流程
5. THE Kiro_Bot SHALL 記錄用戶的確認決策和時間戳

### Requirement 8: 交付給 Forge Agent

**User Story:** 作為 Kiro Bot，我想要將分析結果交付給 Forge Agent，以便開始重建流程

#### Acceptance Criteria

1. WHEN 用戶確認交付，THE Kiro_Bot SHALL 傳送指令給 Forge_Agent
2. THE Kiro_Bot SHALL 傳遞完整的 Analysis_Report、SOP 和原始碼位置
3. THE Kiro_Bot SHALL 設定工作流程追蹤 ID 用於後續階段的關聯
4. WHEN Forge_Agent 成功接收，THE Kiro_Bot SHALL 記錄交付成功狀態
5. WHEN Forge_Agent 無法接收或離線，THE Kiro_Bot SHALL 重試 3 次後返回錯誤

### Requirement 9: 按照 SOP 重建 Skill

**User Story:** 作為 Forge Agent，我想要按照 SOP 重建 skill，以便產出符合 Binary Arsenal 規範的執行檔

#### Acceptance Criteria

1. WHEN Forge_Agent 接收到 SOP，THE Forge_Agent SHALL 建立獨立的構建環境
2. THE Forge_Agent SHALL 按照 SOP 中的步驟順序執行構建命令
3. WHERE 語言是 Go，THE Forge_Agent SHALL 執行靜態編譯（go build -ldflags="-s -w" -tags netgo）
4. WHERE 語言是 Rust，THE Forge_Agent SHALL 執行靜態編譯（cargo build --release --target x86_64-unknown-linux-musl）
5. THE Forge_Agent SHALL 驗證產出物是單一執行檔且無外部依賴
6. WHEN 構建失敗，THE Forge_Agent SHALL 記錄錯誤訊息並通知 Kiro_Bot
7. THE Forge_Agent SHALL 在構建完成後生成構建報告

### Requirement 10: 審查重建結果

**User Story:** 作為 Forge Agent，我想要審查重建的 skill，以便確認其可用性和正確性

#### Acceptance Criteria

1. WHEN 構建完成，THE Forge_Agent SHALL 執行功能測試驗證基本可用性
2. THE Forge_Agent SHALL 驗證執行檔的輸入輸出介面符合規範（CLI 參數、JSON 輸出、stderr 錯誤）
3. THE Forge_Agent SHALL 使用 ldd 命令驗證零依賴特性
4. THE Forge_Agent SHALL 測試執行檔在目標環境（Ubuntu 24.04）的相容性
5. WHEN 審查發現問題，THE Forge_Agent SHALL 記錄問題並嘗試修正或通知 Kiro_Bot
6. THE Forge_Agent SHALL 生成審查報告包含測試結果和可用性評估

### Requirement 11: 部署到根目錄

**User Story:** 作為系統管理員，我想要將審查通過的 skill 部署到指定目錄，以便 agent 可以使用

#### Acceptance Criteria

1. WHEN 審查通過，THE Forge_Agent SHALL 將執行檔複製到 Root_Directory
2. THE Forge_Agent SHALL 設定執行檔的權限為 755（可執行）
3. THE Forge_Agent SHALL 在 Root_Directory 中建立 skill 的配置檔案（如需要）
4. WHERE skill 需要在 K3s_Cluster 中運行，THE Forge_Agent SHALL 部署到 OpenClaw_Namespace
5. THE Forge_Agent SHALL 驗證部署後的 skill 可以正常啟動
6. THE Forge_Agent SHALL 記錄部署位置和版本資訊

### Requirement 12: 備份到 GitHub

**User Story:** 作為系統管理員，我想要將重建的 skill 備份到 GitHub，以便版本控制和災難恢復

#### Acceptance Criteria

1. WHEN 部署成功，THE Forge_Agent SHALL 建立 Git 儲存庫（如不存在）
2. THE Forge_Agent SHALL 提交重建的原始碼、執行檔和相關文檔
3. THE Forge_Agent SHALL 標記版本號（基於時間戳或語義版本）
4. THE Forge_Agent SHALL 推送到指定的 GitHub 備份儲存庫
5. WHEN GitHub 推送失敗，THE Forge_Agent SHALL 重試 3 次並記錄錯誤
6. THE Forge_Agent SHALL 在備份完成後生成備份報告包含 commit hash 和 URL

### Requirement 13: 工作流程狀態追蹤

**User Story:** 作為用戶，我想要追蹤整個工作流程的狀態，以便了解當前進度和問題

#### Acceptance Criteria

1. THE Skill_Automation_Pipeline SHALL 為每個工作流程實例分配唯一的追蹤 ID
2. THE Skill_Automation_Pipeline SHALL 記錄每個階段的開始時間、結束時間和狀態
3. THE Skill_Automation_Pipeline SHALL 支援查詢工作流程狀態（進行中、成功、失敗、已取消）
4. WHEN 任何階段失敗，THE Skill_Automation_Pipeline SHALL 記錄失敗原因和錯誤堆疊
5. THE Skill_Automation_Pipeline SHALL 提供工作流程歷史查詢功能
6. THE Skill_Automation_Pipeline SHALL 將狀態資訊持久化到儲存系統

### Requirement 14: 錯誤處理和恢復

**User Story:** 作為系統，我想要優雅地處理錯誤並支援恢復，以便提高系統可靠性

#### Acceptance Criteria

1. WHEN 任何階段發生錯誤，THE Skill_Automation_Pipeline SHALL 記錄詳細錯誤資訊
2. THE Skill_Automation_Pipeline SHALL 將錯誤訊息輸出到 stderr
3. WHERE 錯誤是暫時性的（網路超時、資源不足），THE Skill_Automation_Pipeline SHALL 自動重試最多 3 次
4. WHEN 重試失敗，THE Skill_Automation_Pipeline SHALL 通知用戶並提供手動恢復選項
5. THE Skill_Automation_Pipeline SHALL 支援從失敗階段恢復工作流程（不需要重新開始）
6. THE Skill_Automation_Pipeline SHALL 清理失敗工作流程的臨時檔案和資源

### Requirement 15: 日誌記錄和監控

**User Story:** 作為系統管理員，我想要完整的日誌記錄，以便診斷問題和監控系統健康狀態

#### Acceptance Criteria

1. THE Skill_Automation_Pipeline SHALL 記錄所有階段的詳細操作日誌
2. THE Skill_Automation_Pipeline SHALL 使用結構化日誌格式（JSON）
3. THE Skill_Automation_Pipeline SHALL 記錄日誌級別（DEBUG、INFO、WARN、ERROR）
4. THE Skill_Automation_Pipeline SHALL 將日誌輸出到檔案和 stdout
5. THE Skill_Automation_Pipeline SHALL 支援日誌輪轉（每日或達到大小限制）
6. THE Skill_Automation_Pipeline SHALL 提供監控指標（處理時間、成功率、失敗率）

### Requirement 16: 配置管理

**User Story:** 作為系統管理員，我想要靈活配置系統參數，以便適應不同環境和需求

#### Acceptance Criteria

1. THE Skill_Automation_Pipeline SHALL 從配置檔案讀取系統參數
2. THE Skill_Automation_Pipeline SHALL 支援環境變數覆蓋配置檔案設定
3. THE Skill_Automation_Pipeline SHALL 驗證配置參數的有效性
4. THE Skill_Automation_Pipeline SHALL 提供預設配置值
5. WHERE 配置無效或缺失，THE Skill_Automation_Pipeline SHALL 使用預設值並記錄警告
6. THE Skill_Automation_Pipeline SHALL 支援熱重載配置（不需要重啟系統）

### Requirement 17: 並行處理支援

**User Story:** 作為系統，我想要支援並行處理多個 skill，以便提高整體吞吐量

#### Acceptance Criteria

1. THE Skill_Automation_Pipeline SHALL 支援同時處理多個 GitHub_Skill
2. THE Skill_Automation_Pipeline SHALL 為每個工作流程實例分配獨立的工作目錄
3. THE Skill_Automation_Pipeline SHALL 限制並行工作流程數量（可配置，預設 3）
4. WHEN 達到並行限制，THE Skill_Automation_Pipeline SHALL 將新請求加入佇列
5. THE Skill_Automation_Pipeline SHALL 確保不同工作流程之間的資源隔離
6. THE Skill_Automation_Pipeline SHALL 提供佇列狀態查詢功能

### Requirement 18: 語言決策矩陣應用

**User Story:** 作為系統，我想要根據 Binary Arsenal 決策矩陣建議最佳語言，以便優化 skill 效能

#### Acceptance Criteria

1. WHEN Type_Analyzer 識別 skill 類型，THE Kiro_Bot SHALL 應用 Binary_Arsenal 決策矩陣
2. WHERE skill 類型是網路請求、雲端 API、k3s/Docker 互動或簡單自動化，THE Kiro_Bot SHALL 建議使用 Go
3. WHERE skill 類型是檔案處理、加解密、高效能解析或重度計算，THE Kiro_Bot SHALL 建議使用 Rust
4. WHERE 原始語言與建議語言不同，THE Kiro_Bot SHALL 在 Analysis_Report 中說明原因和權衡
5. THE Kiro_Bot SHALL 在 SOP 中包含語言選擇的理由
6. THE Kiro_Bot SHALL 允許用戶覆蓋語言建議

### Requirement 19: 產出物驗證

**User Story:** 作為系統，我想要驗證最終產出物符合所有規範，以便確保品質

#### Acceptance Criteria

1. WHEN 構建完成，THE Forge_Agent SHALL 驗證產出物是單一執行檔
2. THE Forge_Agent SHALL 驗證執行檔大小合理（小於 100MB）
3. THE Forge_Agent SHALL 驗證執行檔可以在目標平台執行
4. THE Forge_Agent SHALL 驗證執行檔接受 CLI 參數並輸出 JSON 格式
5. THE Forge_Agent SHALL 驗證錯誤情況下輸出到 stderr
6. WHEN 任何驗證失敗，THE Forge_Agent SHALL 拒絕部署並返回詳細錯誤報告

### Requirement 20: 整合測試

**User Story:** 作為開發者，我想要端到端的整合測試，以便驗證整個工作流程正確運作

#### Acceptance Criteria

1. THE Skill_Automation_Pipeline SHALL 提供整合測試模式
2. THE Skill_Automation_Pipeline SHALL 使用測試用的 GitHub 儲存庫執行完整流程
3. THE Skill_Automation_Pipeline SHALL 驗證每個階段的輸出符合預期
4. THE Skill_Automation_Pipeline SHALL 在測試模式下不執行實際部署和 GitHub 備份
5. THE Skill_Automation_Pipeline SHALL 生成整合測試報告
6. FOR ALL 有效的測試輸入，執行兩次整合測試 SHALL 產生一致的結果（冪等性）
