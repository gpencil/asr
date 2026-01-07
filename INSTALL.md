# Whisper 安装指南

## 快速安装（3步）

### 第1步：安装 Python

确保你的系统已安装 Python 3.8+

```bash
# 检查Python版本
python --version
# 或
python3 --version
```

### 第2步：安装 Whisper

```bash
# 推荐方式：使用pip
pip install -U openai-whisper

# 或使用pip3
pip3 install -U openai-whisper
```

### 第3步：验证安装

```bash
# 检查whisper命令
whisper --help

# 应该看到类似输出：
# usage: whisper [-h] [--model {tiny,base,small,medium,large}] ...
```

## 完整安装步骤

### macOS

```bash
# 1. 安装Homebrew（如果未安装）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. 安装Python
brew install python

# 3. 安装FFmpeg（可选但推荐）
brew install ffmpeg

# 4. 安装Whisper
pip3 install -U openai-whisper

# 5. 验证
whisper --help
```

### Ubuntu/Debian Linux

```bash
# 1. 更新包列表
sudo apt update

# 2. 安装Python和pip
sudo apt install python3 python3-pip

# 3. 安装FFmpeg（可选但推荐）
sudo apt install ffmpeg

# 4. 安装Whisper
pip3 install -U openai-whisper

# 5. 验证
whisper --help
```

### Windows

```bash
# 1. 下载并安装Python
# 访问 https://www.python.org/downloads/
# 下载Python 3.8或更高版本
# 安装时勾选 "Add Python to PATH"

# 2. 打开命令提示符（CMD）或PowerShell

# 3. 安装Whisper
pip install -U openai-whisper

# 4. （可选）安装FFmpeg
# 从 https://ffmpeg.org/download.html 下载
# 解压并添加到系统PATH

# 5. 验证
whisper --help
```

## 模型下载

Whisper会在第一次使用时自动下载模型。

### 手动预下载模型

```bash
# 下载推荐的medium模型（约1.5GB）
# 运行一次whisper命令，它会自动下载
echo "test" > test.txt
whisper test.txt --model medium

# 模型会保存到：
# macOS/Linux: ~/.cache/whisper/
# Windows: C:\Users\<用户名>\.cache\whisper\
```

### 加速下载（国内用户）

如果下载速度慢，可以手动下载模型：

```bash
# 1. 访问 Hugging Face 镜像
# https://hf-mirror.com/ggerganov/whisper.cpp/tree/main

# 2. 下载对应模型文件（例如medium.pt）

# 3. 放到缓存目录
# macOS/Linux:
mkdir -p ~/.cache/whisper
mv ~/Downloads/medium.pt ~/.cache/whisper/

# Windows:
# 放到 C:\Users\<用户名>\.cache\whisper\
```

## GPU加速（可选）

如果你有NVIDIA显卡，可以启用GPU加速：

```bash
# 1. 卸载CPU版本的PyTorch
pip uninstall torch

# 2. 安装CUDA版本（需要先安装CUDA驱动）
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118

# 3. 使用时指定GPU
# 在Go代码中设置：config.Device = "cuda"
```

## 常见问题

### Q: 提示 "command not found: whisper"？

**原因**: Python脚本路径不在系统PATH中

**解决方法**:

```bash
# macOS/Linux
# 添加到 ~/.zshrc 或 ~/.bashrc
export PATH="$HOME/.local/bin:$PATH"

# 或使用完整路径
~/.local/bin/whisper --help

# Windows
# 将 C:\Users\<用户名>\AppData\Local\Programs\Python\Python3X\Scripts
# 添加到系统环境变量PATH
```

### Q: 安装失败，提示权限错误？

```bash
# 不要使用sudo，使用用户安装
pip install --user openai-whisper

# 或使用虚拟环境
python3 -m venv whisper_env
source whisper_env/bin/activate  # Windows: whisper_env\Scripts\activate
pip install openai-whisper
```

### Q: 下载模型失败？

**方法1**: 使用镜像源

```bash
# 配置pip镜像（清华源）
pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
```

**方法2**: 手动下载模型（见上方"加速下载"）

### Q: 如何卸载？

```bash
pip uninstall openai-whisper
```

## 测试安装

创建测试文件 `test.go`：

```go
package main

import (
    "fmt"
    "log"
    "ysgit.lunalabs.cn/products/go-common/asr"
)

func main() {
    // 测试Whisper是否安装
    client, err := asr.NewWhisperClient(nil)
    if err != nil {
        log.Fatal("Whisper未正确安装:", err)
    }

    fmt.Println("✓ Whisper安装成功！")

    // 显示模型信息
    info := asr.GetModelInfo(asr.ModelMedium)
    fmt.Printf("推荐模型: %s (大小: %s)\n", info["name"], info["size"])
}
```

运行测试：

```bash
go run test.go
```

如果看到 "✓ Whisper安装成功！"，说明环境已就绪。

## 升级Whisper

```bash
# 升级到最新版本
pip install -U openai-whisper

# 查看版本
pip show openai-whisper
```

## 更多帮助

- Whisper官方文档: https://github.com/openai/whisper
- 问题反馈: https://github.com/openai/whisper/issues
