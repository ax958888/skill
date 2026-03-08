# Binary Arsenal Generator (私人兵工廠建構者)

自動化技能轉換和部署管道，將 GitHub 儲存庫中的腳本或工具（Python/JS/TS）轉換為零依賴、高效能的單一執行檔。

## 專案概述

本專案實現了完整的自動化工作流程，從分析到部署在 k3s 基礎設施上。

### 核心功能

- 自動化安全分析和代碼審查
- 語言檢測和最佳目標語言選擇（Go/Rust）
- 生成重建的標準操作程序（SOP）
- 自動編譯為零依賴的靜態執行檔
- 部署到 cpx31 伺服器 k3s 集群（openclaw namespace）
- GitHub 備份和版本控制

## 工具

### 1. skill-analyzer

分析 GitHub 儲存庫，生成完整的分析報告和構建 SOP。

**功能：**
- 接收 GitHub URL
- 安全審查（檢測惡意代碼）
- 語言識別
- 類型分析
- SOP 生成

**使用方式：**
```bash
cd skill-analyzer
./build.sh
./skill-analyzer analyze https://github.com/user/repo -o analysis.json
```

詳細文檔：[SKILL-ANALYZER.md](./SKILL-ANALYZER.md)

### 2. skill-builder

根據分析報告構建和部署技能。

**功能：**
- 按照 SOP 重建技能
- 驗證產出物
- 部署到目標目錄
- GitHub 備份

**使用方式：**
```bash
cd skill-builder
./build.sh
./skill-builder build analysis.json
```

詳細文檔：[SKILL-BUILDER.md](./SKILL-BUILDER.md)

## 完整工作流程

```bash
# 1. 分析階段（Kiro Bot）
skill-analyzer analyze https://github.com/user/repo -o analysis.json

# 2. 查看分析結果
cat analysis.json | jq .

# 3. 構建階段（Forge Agent）
skill-builder build analysis.json -v

# 4. 驗證部署
/root/workspace/agents/skill-name --help
```

## 技術棧

### 主要語言

- **Go**: 網路請求、雲端 API、k3s/Docker 互動、簡單自動化
- **Rust**: 檔案處理、加解密、高效能解析、重度計算

### 構建要求

**Go 靜態編譯：**
```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo
```

**Rust 靜態編譯：**
```bash
cargo build --release --target x86_64-unknown-linux-musl
```

### 輸出標準

所有執行檔必須符合：
- **單一執行檔**：無外部依賴
- **靜態連結**：可在 Ubuntu 24.04 上直接執行
- **標準介面**：
  - 輸入：CLI 參數
  - 輸出：成功時輸出 JSON 到 stdout
  - 錯誤：輸出到 stderr 並返回非零退出碼

## 部署環境

- 目標作業系統：Ubuntu 24.04 (cpx31)
- 架構：x86_64 (Linux)
- 容器平台：k3s
- 命名空間：openclaw
- 安裝目錄：/root/workspace/agents/

## 專案結構

```
.
├── skill-analyzer/          # 分析工具
│   ├── cmd/                 # 主程式
│   ├── pkg/models/          # 資料模型
│   ├── build.sh             # 構建腳本
│   └── README.md
├── skill-builder/           # 構建工具
│   ├── cmd/                 # 主程式
│   ├── pkg/models/          # 資料模型
│   ├── build.sh             # 構建腳本
│   └── README.md
├── .kiro/                   # Kiro 配置
│   ├── specs/               # 規格文檔
│   ├── steering/            # 專案指導規則
│   └── hooks/               # 自動化 hooks
├── SKILL-ANALYZER.md        # 分析工具詳細文檔
├── SKILL-BUILDER.md         # 構建工具詳細文檔
└── skill重塑.md             # 技能定義

```

## 快速開始

### 前置需求

- Go 1.21+
- Git
- 目標語言工具鏈（Go/Rust）用於構建技能
- ldd 命令（用於驗證零依賴）
- jq（用於 JSON 處理，可選）

### 構建工具

```bash
# 構建 skill-analyzer
cd skill-analyzer
chmod +x build.sh
./build.sh

# 構建 skill-builder
cd skill-builder
chmod +x build.sh
./build.sh
```

### 驗證

```bash
# 檢查零依賴
ldd skill-analyzer/skill-analyzer
ldd skill-builder/skill-builder

# 檢查執行檔大小
ls -lh skill-analyzer/skill-analyzer
ls -lh skill-builder/skill-builder

# 測試執行
skill-analyzer/skill-analyzer --help
skill-builder/skill-builder --help
```

## 語言決策矩陣

| 任務特徵 | 採用語言 | 理由 |
|---------|---------|------|
| 網路請求、雲端 API、k3s/Docker 互動、簡單自動化 | Go | 網路操作和系統自動化 |
| 檔案處理、加解密、高效能解析、重度計算 | Rust | 記憶體安全和效能 |

## 文檔

- [產品概述](.kiro/steering/product.md)
- [技術棧](.kiro/steering/tech.md)
- [專案結構](.kiro/steering/structure.md)
- [需求文檔](.kiro/specs/skill-automation-pipeline/requirements.md)
- [設計文檔](.kiro/specs/skill-automation-pipeline/design.md)

## 授權

MIT License

## 貢獻

歡迎提交 Pull Requests 和 Issues。
=======
# skill
用Rust/go 語言重塑skill
>>>>>>> af9e4e47d5b46c04d21d04a9bd1aeaf697ad0d1dnary Arsenal Generator (私人兵工廠建構者)

自動化技能轉換和部署管道，將 GitHub 儲存庫中的腳本或工具（Python/JS/TS）轉換為零依賴、高效能的單一執行檔。

## 專案概述

本專案實現了完整的自動化工作流程，從分析到部署在 k3s 基礎設施上。

### 核心功能

- 自動化安全分析和代碼審查
- 語言檢測和最佳目標語言選擇（Go/Rust）
- 生成重建的標準操作程序（SOP）
- 自動編譯為零依賴的靜態執行檔
- 部署到 cpx31 伺服器 k3s 集群（openclaw namespace）
- GitHub 備份和版本控制

## 工具

### 1. skill-analyzer

分析 GitHub 儲存庫，生成完整的分析報告和構建 SOP。

**功能：**
- 接收 GitHub URL
- 安全審查（檢測惡意代碼）
- 語言識別
- 類型分析
- SOP 生成

**使用方式：**
```bash
cd skill-analyzer
./build.sh
./skill-analyzer analyze https://github.com/user/repo -o analysis.json
```

詳細文檔：[SKILL-ANALYZER.md](./SKILL-ANALYZER.md)

### 2. skill-builder

根據分析報告構建和部署技能。

**功能：**
- 按照 SOP 重建技能
- 驗證產出物
- 部署到目標目錄
- GitHub 備份

**使用方式：**
```bash
cd skill-builder
./build.sh
./skill-builder build analysis.json
```

詳細文檔：[SKILL-BUILDER.md](./SKILL-BUILDER.md)

## 完整工作流程

```bash
# 1. 分析階段（Kiro Bot）
skill-analyzer analyze https://github.com/user/repo -o analysis.json

# 2. 查看分析結果
cat analysis.json | jq .

# 3. 構建階段（Forge Agent）
skill-builder build analysis.json -v

# 4. 驗證部署
/root/workspace/agents/skill-name --help
```

## 技術棧

### 主要語言

- **Go**: 網路請求、雲端 API、k3s/Docker 互動、簡單自動化
- **Rust**: 檔案處理、加解密、高效能解析、重度計算

### 構建要求

**Go 靜態編譯：**
```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -tags netgo
```

**Rust 靜態編譯：**
```bash
cargo build --release --target x86_64-unknown-linux-musl
```

### 輸出標準

所有執行檔必須符合：
- **單一執行檔**：無外部依賴
- **靜態連結**：可在 Ubuntu 24.04 上直接執行
- **標準介面**：
  - 輸入：CLI 參數
  - 輸出：成功時輸出 JSON 到 stdout
  - 錯誤：輸出到 stderr 並返回非零退出碼

## 部署環境

- 目標作業系統：Ubuntu 24.04 (cpx31)
- 架構：x86_64 (Linux)
- 容器平台：k3s
- 命名空間：openclaw
- 安裝目錄：/root/workspace/agents/

## 專案結構

```
.
├── skill-analyzer/          # 分析工具
│   ├── cmd/                 # 主程式
│   ├── pkg/models/          # 資料模型
│   ├── build.sh             # 構建腳本
│   └── README.md
├── skill-builder/           # 構建工具
│   ├── cmd/                 # 主程式
│   ├── pkg/models/          # 資料模型
│   ├── build.sh             # 構建腳本
│   └── README.md
├── .kiro/                   # Kiro 配置
│   ├── specs/               # 規格文檔
│   ├── steering/            # 專案指導規則
│   └── hooks/               # 自動化 hooks
├── SKILL-ANALYZER.md        # 分析工具詳細文檔
├── SKILL-BUILDER.md         # 構建工具詳細文檔
└── skill重塑.md             # 技能定義

```

## 快速開始

### 前置需求

- Go 1.21+
- Git
- 目標語言工具鏈（Go/Rust）用於構建技能

### 構建工具

```bash
# 構建 skill-analyzer
cd skill-analyzer
chmod +x build.sh
./build.sh

# 構建 skill-builder
cd skill-builder
chmod +x build.sh
./build.sh
```

### 驗證

```bash
# 檢查零依賴
ldd skill-analyzer/skill-analyzer
ldd skill-builder/skill-builder

# 測試執行
skill-analyzer/skill-analyzer --help
skill-builder/skill-builder --help
```

## 語言決策矩陣

| 任務特徵 | 採用語言 | 理由 |
|---------|---------|------|
| 網路請求、雲端 API、k3s/Docker 互動、簡單自動化 | Go | 網路操作和系統自動化 |
| 檔案處理、加解密、高效能解析、重度計算 | Rust | 記憶體安全和效能 |

## 文檔

- [產品概述](.kiro/steering/product.md)
- [技術棧](.kiro/steering/tech.md)
- [專案結構](.kiro/steering/structure.md)
- [需求文檔](.kiro/specs/skill-automation-pipeline/requirements.md)
- [設計文檔](.kiro/specs/skill-automation-pipeline/design.md)

## 授權

MIT License

## 貢獻

歡迎提交 Pull Requests 和 Issues。
