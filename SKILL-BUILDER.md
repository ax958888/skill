# Skill: skill-builder (Forge Agent 構建工具)

## 技能宣告 (Definition)

當接收到分析報告和 SOP 時，此工具負責構建和部署階段的所有功能：按照 SOP 重建 skill、審查結果、部署到目標目錄、備份到 GitHub。

---

## 語言選擇 (Language)

**Go** - 適合 k3s/Docker 互動、自動化部署、系統操作

---

## 功能描述 (Features)

### 核心功能
1. **接收分析報告** - 讀取 skill-analyzer 產生的分析報告
2. **構建環境設置** - 建立獨立的構建環境
3. **按 SOP 構建** - 執行 Go 或 Rust 靜態編譯
4. **驗證產出物** - 檢查零依賴、介面規範、功能可用性
5. **部署到根目錄** - 複製到 /root/workspace/agents/
6. **備份到 GitHub** - 版本控制和災難恢復
7. **構建報告** - 輸出 JSON 格式的完整構建結果

### 構建引擎模組
- 支援 Go 靜態編譯（CGO_ENABLED=0）
- 支援 Rust 靜態編譯（x86_64-unknown-linux-musl）
- 按照 SOP 步驟順序執行
- 記錄構建日誌和時間
- 處理構建失敗並報告錯誤

### 驗證模組
- **單一執行檔檢查** - 確認產出物是單一檔案
- **零依賴檢查** - 使用 ldd 驗證無動態依賴
- **CLI 介面檢查** - 測試 --help 參數
- **JSON 輸出檢查** - 驗證輸出格式
- **大小檢查** - 確認檔案小於 100MB
- **執行權限檢查** - 驗證可執行性

### 部署模組
- 複製執行檔到目標目錄
- 設定檔案權限為 755
- 支援 k3s 集群部署（可選）
- 驗證部署後可執行性
- 記錄部署位置和版本資訊

### 備份模組
- 建立或使用現有 Git 儲存庫
- 提交原始碼和執行檔
- 標記版本號（時間戳或語義版本）
- 推送到 GitHub 備份儲存庫
- 重試機制（最多 3 次）

---

## CLI 介面 (Interface)

### 命令格式
```bash
skill-builder build <analysis-report> [flags]
```

### 參數說明
- `<analysis-report>` - 分析報告檔案路徑或 `-` 從 stdin 讀取（必填）
- `--output, -o` - 輸出檔案路徑（預設：stdout）
- `--work-dir, -w` - 工作目錄（預設：/tmp/skill-builder-*）
- `--deploy-dir, -d` - 部署目錄（預設：/root/workspace/agents/）
- `--skip-deploy` - 跳過部署步驟（僅構建和驗證）
- `--skip-backup` - 跳過 GitHub 備份步驟
- `--config, -c` - 配置檔案路徑
- `--verbose, -v` - 啟用詳細日誌

### 使用範例
```bash
# 從檔案讀取分析報告
skill-builder build analysis.json

# 從 stdin 讀取（管道）
cat analysis.json | skill-builder build -

# 跳過備份
skill-builder build analysis.json --skip-backup

# 自訂部署目錄
skill-builder build analysis.json -d /opt/skills/

# 僅構建和驗證，不部署
skill-builder build analysis.json --skip-deploy -v
```

---

## 輸入格式 (Input)

### 分析報告格式（JSON）
```json
{
  "workflow_id": "550e8400-e29b-41d4-a716-446655440000",
  "github_url": "https://github.com/user/repo",
  "repository": {
    "name": "repo",
    "clone_path": "/tmp/skill-analyzer-xyz/repo"
  },
  "recommendation": {
    "target_language": "go"
  },
  "sop": {
    "language": "go",
    "steps": [
      {
        "phase": "setup",
        "command": "go mod init skill-name"
      },
      {
        "phase": "build",
        "command": "CGO_ENABLED=0 go build -ldflags=\"-s -w\" -tags netgo -o skill-name"
      }
    ],
    "interface_spec": {
      "input": "CLI arguments",
      "output": "JSON to stdout",
      "errors": "stderr with non-zero exit code"
    },
    "dependencies": ["github.com/spf13/cobra"]
  }
}
```

---

## 輸出格式 (Output)

### 成功輸出（stdout - JSON）
```json
{
  "workflow_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-15T10:35:00Z",
  "build": {
    "status": "success",
    "language": "go",
    "source_path": "/tmp/skill-builder-xyz/src",
    "binary_path": "/tmp/skill-builder-xyz/skill-name",
    "binary_size": 8388608,
    "build_time": "45s",
    "build_log": "go: downloading github.com/spf13/cobra v1.8.0\n..."
  },
  "validation": {
    "status": "passed",
    "checks": [
      {
        "name": "single_binary",
        "passed": true,
        "message": "Output is a single executable file"
      },
      {
        "name": "zero_dependency",
        "passed": true,
        "message": "ldd: not a dynamic executable"
      },
      {
        "name": "cli_interface",
        "passed": true,
        "message": "Accepts --help flag"
      },
      {
        "name": "json_output",
        "passed": true,
        "message": "Outputs valid JSON"
      },
      {
        "name": "size_check",
        "passed": true,
        "message": "Binary size: 8MB (< 100MB limit)"
      }
    ]
  },
  "deployment": {
    "status": "success",
    "target_path": "/root/workspace/agents/skill-name",
    "permissions": "755",
    "k3s_deployed": false,
    "deployment_time": "2024-01-15T10:35:30Z"
  },
  "backup": {
    "status": "success",
    "repository": "https://github.com/backup-org/skill-name",
    "commit_hash": "def456789abc",
    "tag": "v1.0.0-20240115",
    "backup_time": "2024-01-15T10:36:00Z"
  },
  "status": "success",
  "message": "Build, validation, deployment, and backup completed successfully"
}
```

### 錯誤輸出（stderr - JSON）
```json
{
  "error": "build failed",
  "details": "compilation error: undefined: someFunction",
  "category": "build_validation",
  "workflow_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-15T10:35:00Z",
  "retryable": false,
  "retry_count": 0,
  "suggestions": [
    "Check source code for syntax errors",
    "Verify all dependencies are available"
  ]
}
```

---

## 構建規範 (Build Specification)

### 語言和框架
- **語言**: Go 1.21+
- **CLI 框架**: Cobra
- **進程執行**: os/exec (標準庫)
- **檔案操作**: os, io (標準庫)

### 靜態編譯命令
```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo -o skill-builder cmd/build/main.go
```

### 編譯參數說明
- `CGO_ENABLED=0` - 禁用 CGO，確保靜態編譯
- `-ldflags="-s -w"` - 移除符號表和調試資訊，縮減體積
- `-tags netgo` - 使用純 Go 網路實現
- `-o skill-builder` - 輸出檔案名稱

### 依賴套件
```go
require (
    github.com/spf13/cobra v1.8.0
    github.com/google/uuid v1.5.0
    k8s.io/client-go v0.29.0  // k3s 部署（可選）
)
```

### 專案結構
```
skill-builder/
├── cmd/
│   └── build/
│       └── main.go           # 主程式入口
├── internal/
│   ├── builder/
│   │   └── engine.go         # 構建引擎
│   ├── validator/
│   │   └── validator.go      # 驗證器
│   ├── deployer/
│   │   └── deployer.go       # 部署器
│   └── backup/
│       └── github.go         # GitHub 備份
├── pkg/
│   └── models/
│       └── report.go         # 資料模型
├── go.mod
├── go.sum
└── README.md
```

---

## 構建流程 (Build Process)

### 階段 1：環境設置
```bash
# 建立工作目錄
mkdir -p /tmp/skill-builder-{workflow-id}

# 複製原始碼
cp -r {source-path} /tmp/skill-builder-{workflow-id}/src

# 進入工作目錄
cd /tmp/skill-builder-{workflow-id}/src
```

### 階段 2：執行 SOP 步驟
```bash
# Setup 階段
go mod init skill-name
go mod tidy

# Build 階段（Go）
CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo -o skill-name

# Build 階段（Rust）
cargo build --release --target x86_64-unknown-linux-musl
```

### 階段 3：驗證
```bash
# 檢查零依賴
ldd skill-name

# 檢查執行權限
chmod +x skill-name

# 測試 CLI 介面
./skill-name --help

# 測試 JSON 輸出
./skill-name test-command | jq .

# 檢查檔案大小
ls -lh skill-name
```

### 階段 4：部署
```bash
# 複製到目標目錄
cp skill-name /root/workspace/agents/

# 設定權限
chmod 755 /root/workspace/agents/skill-name

# 驗證部署
/root/workspace/agents/skill-name --help
```

### 階段 5：備份
```bash
# 初始化 Git（如需要）
git init
git remote add origin {backup-repo-url}

# 提交變更
git add .
git commit -m "Build skill-name v1.0.0"

# 標記版本
git tag v1.0.0-20240115

# 推送到 GitHub
git push origin main --tags
```

---

## 錯誤處理 (Error Handling)

### 錯誤分類
1. **輸入錯誤** (exit code: 1)
   - 無效的分析報告格式
   - 缺少必要欄位

2. **構建錯誤** (exit code: 4)
   - 編譯失敗
   - 依賴下載失敗
   - 不重試（需要修改代碼）

3. **驗證錯誤** (exit code: 4)
   - 零依賴檢查失敗
   - 介面規範不符
   - 檔案過大

4. **部署錯誤** (exit code: 2)
   - 檔案複製失敗
   - 權限設定失敗
   - 自動重試 3 次

5. **備份錯誤** (exit code: 2)
   - Git 操作失敗
   - GitHub 推送失敗
   - 自動重試 3 次

### 重試邏輯
僅對外部服務錯誤（部署、備份）進行重試：
- 初始嘗試：立即執行
- 重試 1：2 秒後
- 重試 2：4 秒後
- 重試 3：8 秒後

### 清理機制
```bash
# 構建失敗時清理
rm -rf /tmp/skill-builder-{workflow-id}

# 部署失敗時保留構建產物
# 允許手動部署
```

---

## 配置檔案 (Configuration)

### 配置檔案格式（JSON）
```json
{
  "work_dir": "/tmp/skill-builder",
  "deploy_dir": "/root/workspace/agents/",
  "timeout": 600,
  "max_binary_size": 104857600,
  "k3s_config": {
    "enabled": false,
    "namespace": "openclaw",
    "kubeconfig": "/etc/rancher/k3s/k3s.yaml"
  },
  "backup_repo": "https://github.com/backup-org/skills",
  "github_token": "${GITHUB_TOKEN}",
  "log_level": "info"
}
```

### 環境變數
- `GITHUB_TOKEN` - GitHub 存取令牌（備份時需要）
- `SKILL_BUILDER_CONFIG` - 配置檔案路徑
- `SKILL_BUILDER_DEPLOY_DIR` - 部署目錄路徑

---

## 驗證步驟 (Validation)

### 零依賴驗證
```bash
ldd skill-builder
# 預期輸出: "not a dynamic executable" 或 "statically linked"
```

### 執行檔大小檢查
```bash
ls -lh skill-builder
# 預期: < 20MB
```

### 功能測試
```bash
# 測試 help 命令
./skill-builder --help

# 測試構建功能（使用測試報告）
./skill-builder build test-analysis.json --skip-deploy --skip-backup -v
```

---

## 部署規範 (Deployment)

### 目標環境
- **作業系統**: Ubuntu 24.04
- **架構**: x86_64 (Linux)
- **安裝目錄**: /root/workspace/agents/
- **K3s 集群**: openclaw namespace（可選）

### 部署步驟
```bash
# 設定執行權限
chmod 755 skill-builder

# 複製到目標目錄
cp skill-builder /root/workspace/agents/

# 驗證部署
/root/workspace/agents/skill-builder --help
```

### K3s 部署（可選）
```bash
# 建立 ConfigMap
kubectl create configmap skill-name-config \
  --from-file=skill-name=/root/workspace/agents/skill-name \
  -n openclaw

# 建立 Pod
kubectl run skill-name \
  --image=alpine:latest \
  --command -- /bin/sh -c "while true; do sleep 3600; done" \
  -n openclaw
```

---

## 使用場景 (Use Cases)

### 場景 1：標準構建流程
```bash
# 1. Kiro Bot 分析
skill-analyzer analyze https://github.com/user/repo -o analysis.json

# 2. 用戶確認

# 3. Forge Agent 構建
skill-builder build analysis.json -v
```

### 場景 2：僅構建不部署
```bash
skill-builder build analysis.json --skip-deploy --skip-backup -o build-report.json
# 用於測試構建流程
```

### 場景 3：管道處理
```bash
skill-analyzer analyze https://github.com/user/repo | \
  skill-builder build - --skip-backup
# 快速測試（跳過備份）
```

### 場景 4：批次構建
```bash
for analysis in *.json; do
  skill-builder build "$analysis" -o "${analysis%.json}-build.json"
done
```

---

## 驗證檢查清單 (Validation Checklist)

### 1. 單一執行檔檢查
```bash
# 確認只有一個執行檔產出
[ -f skill-name ] && [ ! -d skill-name ]
```

### 2. 零依賴檢查
```bash
# 使用 ldd 檢查
ldd skill-name 2>&1 | grep -q "not a dynamic executable"
```

### 3. CLI 介面檢查
```bash
# 測試 --help 參數
./skill-name --help
echo $?  # 應該是 0
```

### 4. JSON 輸出檢查
```bash
# 測試 JSON 輸出
./skill-name test-command | jq . > /dev/null
echo $?  # 應該是 0
```

### 5. 錯誤輸出檢查
```bash
# 測試錯誤輸出到 stderr
./skill-name invalid-command 2>&1 >/dev/null | jq .error
```

### 6. 大小檢查
```bash
# 檢查檔案大小 < 100MB
size=$(stat -f%z skill-name 2>/dev/null || stat -c%s skill-name)
[ $size -lt 104857600 ]
```

---

## 日誌記錄 (Logging)

### 日誌格式（JSON）
```json
{
  "timestamp": "2024-01-15T10:35:00Z",
  "level": "info",
  "message": "Starting build process",
  "workflow_id": "550e8400-e29b-41d4-a716-446655440000",
  "phase": "build",
  "language": "go"
}
```

### 日誌級別
- `DEBUG` - 詳細調試資訊（構建命令、輸出）
- `INFO` - 一般資訊訊息（階段開始/完成）
- `WARN` - 警告訊息（非致命錯誤）
- `ERROR` - 錯誤訊息（構建失敗、驗證失敗）

### 構建日誌
```bash
# 構建日誌保存到檔案
/tmp/skill-builder-{workflow-id}/build.log

# 包含所有構建命令的輸出
# 用於調試構建失敗
```

---

## 正確性屬性 (Correctness Properties)

此工具實現以下可測試的正確性屬性：

1. **構建環境隔離** - 並發構建使用不同的工作目錄
2. **步驟執行順序** - SOP 步驟按順序執行
3. **零依賴驗證** - 所有構建的執行檔無動態依賴
4. **構建報告生成** - 所有構建嘗試都生成報告
5. **介面驗證** - 執行檔接受 --help 並輸出 JSON
6. **大小限制** - 執行檔小於 100MB
7. **部署權限** - 部署的執行檔權限為 755
8. **部署後可執行** - 部署的執行檔可以執行
9. **部署資訊記錄** - 構建報告包含部署路徑和時間
10. **Git 備份提交** - 備份包含原始碼和執行檔
11. **版本標記** - 備份包含版本標記
12. **備份報告完整性** - 備份報告包含 commit hash 和 URL
13. **臨時檔案清理** - 失敗時清理臨時檔案
14. **錯誤輸出到 stderr** - 所有錯誤輸出到 stderr
15. **驗證失敗阻止部署** - 驗證失敗時不執行部署

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
# 端到端測試
go test ./test/integration/... -v

# 使用測試儲存庫
skill-builder build test/fixtures/test-analysis.json --skip-backup
```

### 測試覆蓋率
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 故障排除 (Troubleshooting)

### 問題 1：構建失敗
```bash
# 檢查構建日誌
cat /tmp/skill-builder-{workflow-id}/build.log

# 檢查依賴是否可用
go mod download  # Go
cargo fetch      # Rust
```

### 問題 2：零依賴檢查失敗
```bash
# 確認編譯參數
CGO_ENABLED=0 go build ...  # Go
cargo build --target x86_64-unknown-linux-musl  # Rust

# 檢查動態連結
ldd skill-name
```

### 問題 3：部署失敗
```bash
# 檢查目標目錄權限
ls -ld /root/workspace/agents/

# 檢查磁碟空間
df -h /root/workspace/agents/
```

### 問題 4：備份失敗
```bash
# 檢查 GitHub 令牌
echo $GITHUB_TOKEN

# 檢查網路連接
curl -I https://github.com

# 手動推送
cd /tmp/skill-builder-{workflow-id}/src
git push origin main --tags
```

---

## 維護和更新 (Maintenance)

### 更新構建配置
編輯配置檔案中的構建參數：
- `timeout` - 構建超時時間
- `max_binary_size` - 最大執行檔大小
- `deploy_dir` - 部署目錄

### 更新 K3s 配置
啟用或禁用 K3s 部署：
```json
{
  "k3s_config": {
    "enabled": true,
    "namespace": "openclaw",
    "kubeconfig": "/etc/rancher/k3s/k3s.yaml"
  }
}
```

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

---

## 與 Kiro Bot 整合 (Integration)

### 完整工作流程
```bash
# 1. Kiro Bot (@Pojun_kirobot) 執行分析
skill-analyzer analyze https://github.com/user/repo -o /tmp/analysis.json

# 2. 用戶確認分析結果
cat /tmp/analysis.json | jq .

# 3. Forge Agent (@Forge_coderxbot) 執行構建
skill-builder build /tmp/analysis.json -v

# 4. 驗證部署
/root/workspace/agents/skill-name --help
```

### 自動化腳本
```bash
#!/bin/bash
# auto-skill-pipeline.sh

GITHUB_URL=$1
ANALYSIS_FILE="/tmp/analysis-$(date +%s).json"
BUILD_FILE="/tmp/build-$(date +%s).json"

# 分析階段
echo "Analyzing $GITHUB_URL..."
skill-analyzer analyze "$GITHUB_URL" -o "$ANALYSIS_FILE"

if [ $? -ne 0 ]; then
  echo "Analysis failed"
  exit 1
fi

# 顯示分析結果
cat "$ANALYSIS_FILE" | jq .

# 等待用戶確認
read -p "Continue with build? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  echo "Build cancelled"
  exit 0
fi

# 構建階段
echo "Building skill..."
skill-builder build "$ANALYSIS_FILE" -o "$BUILD_FILE"

if [ $? -ne 0 ]; then
  echo "Build failed"
  exit 1
fi

# 顯示構建結果
cat "$BUILD_FILE" | jq .

echo "Skill automation pipeline completed successfully"
```
