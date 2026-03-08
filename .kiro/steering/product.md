# Product Overview

This project implements the Binary Arsenal Generator (私人兵工廠建構者) - an automated skill transformation and deployment pipeline.

The system takes GitHub repositories containing scripts or tools (Python/JS/TS) and transforms them into zero-dependency, high-performance single binary executables. It follows a complete automation workflow from analysis to deployment on k3s infrastructure.

## Core Capabilities

- Automated security analysis and code review
- Language detection and optimal target language selection (Go/Rust)
- Generation of Standard Operating Procedures (SOP) for rebuilding
- Automated compilation to static binaries with zero dependencies
- Deployment to cpx31 server k3s cluster (openclaw namespace)
- GitHub backup and version control

## Target Environment

- Platform: Ubuntu 24.04 (cpx31)
- Deployment: k3s cluster
- Output: Single static binary executables
