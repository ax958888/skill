# Skill: skill-analyzer (Kiro Bot 分析工具)

## 技能宣告 (Definition)

當用戶提供 GitHub URL 時，此工具負責分析階段的所有功能：安全審查、語言識別、類型分析、SOP 生成，並輸出完整的分析報告。

---

## 語言選擇 (Language)

**Go** - 適合網路請求、GitHub API 互動、簡單自動化

---

## 功能描述 (Features)

### 核心功能
1. **接收 GitHub URL** - 驗證並克隆儲存庫
2. **安全審查** - 檢測惡意代碼和安全風險
3. **語言識別** - 分析專案使用的程式語言
4. **類型分析** - 判斷 skill 的功能類型和使用場景
5. **SOP 生成** - 根據 Binary Arsenal 規範生成標準操作程序
6. **分析報告** - 輸出 JSON 格式的完整分析結果

### 安全審查模組
- 掃描所有原始碼檔案
- 檢測惡意模式：任意代碼執行、未授權網路連接、檔案系統破壞
- 檢測可疑依賴套件
- 生成安全評分（0-100 分）
- 高風險問題自動終止流程

### 語言識別模組
- 識別主要和次要程式語言
- 檢測構建配置檔案（go.mod、Cargo.toml、package.json、requirements.txt）
- 提取依賴列表
- 計算識別信心度

### 類型分析模組
- 分類 skill 類型：
  - network_request（網路請求）
  - cloud_api（雲端 API）
  - k3s_docker（容器互動）
  - file_processing（檔案處理）
  - encryption（加解密）
  - parsing（資料解析）
  - automation（自動化工具）
  - computation（重度計算）
- 識別使用場景和目標用戶
- 分析輸入輸出介面

### SOP 生成模組
- 應用 Binary Arsenal 決策矩陣：
  - **Go**: network_request, cloud_api, k3s_docker, automation
  - **Rust**: file_processing, encryption, parsing, computation
- 生成構建步驟（setup, build, validate）
- 定義介面規範（CLI 參數、JSON 輸出、stderr 錯誤）
- 列出依賴和預估構建時間

---

## CLI 介面 (Interface)

### 命令格式
```bash
skill-analyzer analyze <github-url> [flags]
```

### 參數說明
- `<github-url>` - GitHub 儲存庫 URL（必填）
- `--output, -o` - 輸出檔案路徑（預設：stdout）
- `--work-dir, -w` - 工作目錄（預設：/tmp/skill-analyzer-*）
- `--timeout, -t` - 分析超時秒數（預設：300）
- `--config, -c` - 配置檔案路徑
- `--verbose, -v` - 啟用詳細日誌
- `--workflow-id` - 工作流程追蹤 ID（自動生成）

### 使用範例
```bash
# 基本使用
skill-analyzer analyze https://github.com/user/repo

# 輸出到檔案
skill-analyzer analyze https://github.com/user/repo -o report.json

# 啟用詳細日誌
skill-analyzer analyze https://github.com/user/repo -v

# 自訂工作目錄
skill-analyzer analyze https://github.com/user/repo -w /tmp/my-analysis
```

---

## 輸出格式 (Output)

### 成功輸出（stdout - JSON）
```json
{
  "workflow_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-15T10:30:00Z",
  "github_url": "https://github.com/user/repo",
  "repository": {
    "name": "repo",
    "owner": "user",
    "clone_path": "/tmp/skill-analyzer-xyz/repo",
    "commit_hash": "abc123def456"
  },
  "security": {
    "score": 85,
    "status": "safe",
    "issues": []
  },
  "language": {
    "primary": "python",
    "secondary": ["bash"],
    "confidence": 0.95,
    "build_files": ["requirements.txt", "setup.py"],
    "dependencies": ["requests", "click"]
  },
  "type_analysis": {
    "category": "network_request",
    "use_cases": ["API client", "Data fetching"],
    "input_interface": "CLI arguments",
    "output_interface": "stdout text",
    "target_users": "developers"
  },
  "recommendation": {
    "target_language": "go",
    "rationale": "Network-heavy operations suit Go's concurrency model",
    "decision_matrix_rule": "network_request → Go"
  },
  "sop": {
    "language": "go",
    "steps": [
      {
        "phase": "setup",
        "command": "go mod init skill-name",
        "description": "Initialize Go module"
      },
      {
        "phase": "build",
        "command": "CGO_ENABLED=0 go build -ldflags=\"-s -w\" -tags netgo -o skill-name",
        "description": "Static compilation"
      },
      {
        "phase": "validate",
        "command": "ldd skill-name",
        "expected": "not a dynamic executable"
      }
    ],
    "interface_spec": {
      "input": "CLI arguments via Cobra",
      "output": "JSON to stdout",
      "errors": "stderr with non-zero exit code"
    },
    "dependencies": ["github.com/spf13/cobra"],
    "estimated_build_time": "2-5 minutes"
  },
  "status": "success",
  "message": "Analysis completed successfully"
}
```

### 錯誤輸出（stderr - JSON）
```json
{
  "error": "failed to clone repository",
  "details": "authentication required",
  "category": "external_service",
  "workflow_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-15T10:30:00Z",
  "retryable": true,
  "retry_count": 2,
  "suggestions": [
    "Check GitHub token configuration",
    "Verify repository access permissions"
  ]
}
```

---

## 構建規範 (Build Specification)

### 語言和框架
- **語言**: Go 1.21+
- **CLI 框架**: Cobra
- **HTTP 客戶端**: net/http (標準庫)
- **JSON 處理**: encoding/json (標準庫)

### 靜態編譯命令
```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo -o skill-analyzer cmd/analyze/main.go
```

### 編譯參數說明
- `CGO_ENABLED=0` - 禁用 CGO，確保靜態編譯
- `-ldflags="-s -w"` - 移除符號表和調試資訊，縮減體積
- `-tags netgo` - 使用純 Go 網路實現
- `-o skill-analyzer` - 輸出檔案名稱

### 依賴套件
```go
require (
    github.com/spf13/cobra v1.8.0
    github.com/google/uuid v1.5.0
)
```

### 專案結構
```
skill-analyzer/
├── cmd/
│   └── analyze/
│       └── main.go           # 主程式入口
├── internal/
│   ├── github/
│   │   └── client.go         # GitHub 客戶端
│   ├── security/
│   │   └── analyzer.go       # 安全分析器
│   ├── language/
│   │   └── detector.go       # 語言檢測器
│   ├── types/
│   │   └── analyzer.go       # 類型分析器
│   └── sop/
│       └── generator.go      # SOP 生成器
├── pkg/
│   └── models/
│       └── report.go         # 資料模型
├── go.mod
├── go.sum
└── README.md
```

---

## 錯誤處理 (Error Handling)

### 錯誤分類
1. **用戶輸入錯誤** (exit code: 1)
   - 無效的 GitHub URL
   - 格式錯誤的配置

2. **外部服務錯誤** (exit code: 2)
   - GitHub API 失敗
   - 網路超時
   - 自動重試 3 次（2s, 4s, 8s 間隔）

3. **安全錯誤** (exit code: 3)
   - 檢測到惡意代碼
   - 高風險安全問題
   - 立即終止，不重試

4. **分析錯誤** (exit code: 4)
   - 語言識別失敗
   - 類型分析失敗

### 重試邏輯
- 初始嘗試：立即執行
- 重試 1：2 秒後
- 重試 2：4 秒後
- 重試 3：8 秒後
- 3 次失敗後：報告錯誤並退出

---

## 配置檔案 (Configuration)

### 配置檔案格式（JSON）
```json
{
  "work_dir": "/tmp/skill-analyzer",
  "timeout": 300,
  "max_repo_size": 104857600,
  "github_token": "${GITHUB_TOKEN}",
  "log_level": "info",
  "security_patterns": [
    {
      "name": "eval_usage",
      "regex": "eval\\s*\\(",
      "severity": "high",
      "type": "arbitrary_execution",
      "description": "Detected eval() usage"
    }
  ]
}
```

### 環境變數
- `GITHUB_TOKEN` - GitHub 存取令牌（可選）
- `SKILL_ANALYZER_CONFIG` - 配置檔案路徑
- `SKILL_ANALYZER_WORK_DIR` - 工作目錄路徑

---

## 驗證步驟 (Validation)

### 零依賴驗證
```bash
ldd skill-analyzer
# 預期輸出: "not a dynamic executable" 或 "statically linked"
```

### 執行檔大小檢查
```bash
ls -lh skill-analyzer
# 預期: < 20MB
```

### 功能測試
```bash
# 測試 help 命令
./skill-analyzer --help

# 測試分析功能
./skill-analyzer analyze https://github.com/user/test-repo -v
```

---

## 部署規範 (Deployment)

### 目標環境
- **作業系統**: Ubuntu 24.04
- **架構**: x86_64 (Linux)
- **安裝目錄**: /root/workspace/agents/

### 部署步驟
```bash
# 設定執行權限
chmod 755 skill-analyzer

# 複製到目標目錄
cp skill-analyzer /root/workspace/agents/

# 驗證部署
/root/workspace/agents/skill-analyzer --help
```

---

## 使用場景 (Use Cases)

### 場景 1：分析 Python 腳本
```bash
skill-analyzer analyze https://github.com/user/python-script -o analysis.json
# 輸出: 建議使用 Go 重寫（如果是網路相關）或 Rust（如果是檔案處理）
```

### 場景 2：安全審查
```bash
skill-analyzer analyze https://github.com/suspicious/repo -v
# 如果檢測到惡意代碼，立即終止並報告
```

### 場景 3：批次分析
```bash
for repo in repo1 repo2 repo3; do
  skill-analyzer analyze "https://github.com/user/$repo" -o "$repo-analysis.json"
done
```

---

## 日誌記錄 (Logging)

### 日誌格式（JSON）
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "message": "Starting repository analysis",
  "workflow_id": "550e8400-e29b-41d4-a716-446655440000",
  "github_url": "https://github.com/user/repo"
}
```

### 日誌級別
- `DEBUG` - 詳細調試資訊
- `INFO` - 一般資訊訊息
- `WARN` - 警告訊息
- `ERROR` - 錯誤訊息

---

## 正確性屬性 (Correctness Properties)

此工具實現以下可測試的正確性屬性：

1. **有效 URL 接受** - 所有有效的 GitHub HTTPS URL 都應被接受
2. **無效 URL 拒絕** - 所有無效的 URL 都應被拒絕並返回錯誤
3. **完整檔案掃描** - 安全分析器應掃描所有代碼檔案
4. **惡意模式檢測** - 應檢測已知的惡意模式
5. **安全評分範圍** - 安全評分應在 [0, 100] 範圍內
6. **語言識別準確性** - 應正確識別主要程式語言
7. **構建檔案檢測** - 應識別標準構建配置檔案
8. **類型分類** - 應將 skill 分配到有效類別
9. **報告完整性** - 分析報告應包含所有必需部分
10. **JSON 輸出格式** - 輸出應為有效的 JSON 格式
11. **決策矩陣應用** - 應根據類型正確建議語言（Go 或 Rust）
12. **SOP 完整性** - SOP 應包含所有必需部分
13. **唯一工作流程 ID** - 每個工作流程應有唯一的 ID

---

## 測試策略 (Testing)

### 單元測試
```bash
go test ./internal/... -v
```

### 屬性測試（使用 gopter）
```bash
go test ./internal/... -tags property -v
```

### 整合測試
```bash
go test ./test/integration/... -v
```

### 測試覆蓋率
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 維護和更新 (Maintenance)

### 更新安全模式
編輯配置檔案中的 `security_patterns` 陣列，添加新的惡意模式。

### 更新決策矩陣
修改 `internal/sop/generator.go` 中的決策矩陣規則。

### 版本控制
使用語義版本號：`v1.0.0`
- Major: 破壞性變更
- Minor: 新功能
- Patch: 錯誤修復

---

## 聯絡和支援 (Support)

- **文檔**: 參考 `.kiro/specs/skill-automation-pipeline/`
- **問題回報**: 通過 GitHub Issues
- **貢獻**: 歡迎提交 Pull Requests
