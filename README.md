## 免费语音转文字工具（基于 Whisper）

基于 OpenAI Whisper 的免费开源语音识别工具，完全本地运行，无需API密钥，支持多种语言。

## 功能特性

- ✅ **完全免费** - 基于开源Whisper，本地运行，无需付费API
- ✅ **高准确率** - 使用最先进的语音识别模型
- ✅ **多语言支持** - 支持中文、英文、粤语等99种语言
- ✅ **长音频优化** - 特别优化会议、采访等长音频处理
- ✅ **多种输出格式** - 支持txt、srt、vtt、json等格式
- ✅ **分段时间戳** - 提供精确的时间戳信息
- ✅ **批量处理** - 支持批量转录多个文件

## 环境准备

### 1. 安装 Python 和 Whisper

```bash
# 方式1：使用pip安装（推荐）
pip install -U openai-whisper

# 方式2：使用pip3
pip3 install -U openai-whisper

# 方式3：使用conda
conda install -c conda-forge openai-whisper
```

### 2. 验证安装

```bash
# 检查whisper是否安装成功
whisper --help

# 应该看到whisper的帮助信息
```

### 3. （可选）安装FFmpeg

Whisper需要FFmpeg来处理音频文件：

```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg

# Windows
# 从 https://ffmpeg.org/download.html 下载安装
```

## 快速开始

### 1. 基本使用

```go
package main

import (
    "fmt"
    "log"
    "ysgit.lunalabs.cn/products/go-common/asr"
)

func main() {
    // 创建客户端（使用默认配置）
    client, err := asr.NewWhisperClient(nil)
    if err != nil {
        log.Fatal(err)
    }

    // 转录音频文件
    result, err := client.Transcribe("meeting.mp3")
    if err != nil {
        log.Fatal(err)
    }

    // 输出结果
    fmt.Printf("语言: %s\n", result.Language)
    fmt.Printf("时长: %.2f秒\n", result.Duration)
    fmt.Printf("文本:\n%s\n", result.Text)
}
```

### 2. 长音频处理（推荐配置）

```go
// 使用默认配置（已针对长音频优化）
client, _ := asr.NewWhisperClient(nil)

// 或者使用更高准确率的配置
config := asr.AccurateConfig()
config.Language = "zh"  // 指定中文
client, _ := asr.NewWhisperClient(config)

result, _ := client.Transcribe("long_meeting.mp3")

// 获取带时间戳的分段文本
for _, segment := range result.Segments {
    fmt.Printf("[%02d:%02d] %s\n",
        int(segment.Start)/60,
        int(segment.Start)%60,
        segment.Text)
}
```

## 配置说明

### 模型选择

| 模型 | 大小 | 速度 | 准确率 | 推荐场景 |
|------|------|------|--------|----------|
| tiny | 75MB | 极快 | 较低 | 短音频、快速预览 |
| base | 142MB | 快 | 一般 | 日常使用 |
| small | 466MB | 中等 | 较好 | 平衡场景 |
| **medium** | **1.5GB** | **较慢** | **高** | **长音频、会议（推荐）** |
| large | 3GB | 慢 | 最高 | 专业转录 |

### 预设配置

```go
// 1. 默认配置（长音频优化）
config := asr.DefaultConfig()
// 模型: medium
// 语言: zh（中文）
// 启用VAD（语音活动检测）
// 5分钟自动分段

// 2. 快速配置（牺牲准确率换速度）
config := asr.FastConfig()
// 模型: base
// 3分钟分段
// 减少beam search

// 3. 高准确率配置
config := asr.AccurateConfig()
// 模型: large
// 10分钟分段
// 增加beam search
```

### 自定义配置

```go
config := &asr.WhisperConfig{
    // 模型配置
    Model:    asr.ModelMedium, // 模型选择
    Language: "zh",            // 语言："zh"中文、"en"英文、""自动检测

    // 输出配置
    OutputFormat: "json",      // txt, srt, vtt, json, all
    Verbose:      true,        // 显示详细日志

    // 性能配置
    Device:   "cpu",           // cpu 或 cuda（GPU）
    Threads:  4,               // CPU线程数
    Timeout:  30 * time.Minute, // 超时时间
    BeamSize: 5,               // Beam search大小

    // 长音频优化
    EnableVAD:      true,      // 启用语音活动检测
    MaxSegmentLen:  300,       // 最大分段长度（秒）
    SplitOnSilence: true,      // 在静音处分割
}
```

## 使用示例

### 示例1：处理会议录音

```go
config := asr.DefaultConfig()
config.Language = "zh"
config.OutputFormat = "json"  // 获取时间戳
config.Verbose = true

client, _ := asr.NewWhisperClient(config)
result, _ := client.Transcribe("meeting.mp3")

// 输出带时间戳的会议纪要
fmt.Println("会议纪要:")
for _, seg := range result.Segments {
    minutes := int(seg.Start) / 60
    seconds := int(seg.Start) % 60
    fmt.Printf("[%02d:%02d] %s\n", minutes, seconds, seg.Text)
}
```

### 示例2：生成字幕文件

```go
config := asr.DefaultConfig()
config.OutputFormat = "srt"  // SRT字幕格式

client, _ := asr.NewWhisperClient(config)
result, _ := client.Transcribe("video.mp4")

fmt.Printf("字幕文件已生成: %s\n", result.FilePath)
// 输出: 字幕文件已生成: video.srt
```

### 示例3：批量处理

```go
client, _ := asr.NewWhisperClient(nil)

files := []string{
    "meeting_day1.mp3",
    "meeting_day2.mp3",
    "meeting_day3.mp3",
}

results, _ := client.TranscribeBatch(files)

for i, result := range results {
    fmt.Printf("Day %d: %s\n", i+1, result.Text)
}
```

### 示例4：多语言处理

```go
config := asr.DefaultConfig()

// 中文
config.Language = "zh"
client, _ := asr.NewWhisperClient(config)
result, _ := client.Transcribe("chinese.mp3")

// 英文
config.Language = "en"
client, _ = asr.NewWhisperClient(config)
result, _ = client.Transcribe("english.mp3")

// 自动检测
config.Language = ""
client, _ = asr.NewWhisperClient(config)
result, _ = client.Transcribe("unknown.mp3")
fmt.Println("检测到语言:", result.Language)
```

## 支持的音频格式

Whisper支持几乎所有常见音频/视频格式：

- **音频**: mp3, wav, m4a, flac, ogg, aac, wma
- **视频**: mp4, avi, mkv, mov, wmv, flv

## 性能优化建议

### 1. 长音频处理

```go
config := asr.DefaultConfig()
config.MaxSegmentLen = 300      // 5分钟一段
config.EnableVAD = true          // 启用VAD，跳过静音
config.SplitOnSilence = true     // 在静音处分割
```

### 2. 提升速度

```go
// 使用小模型
config.Model = asr.ModelBase

// 减少beam size
config.BeamSize = 1
config.BestOf = 1

// 使用GPU（如果有）
config.Device = "cuda"
```

### 3. 提升准确率

```go
// 使用大模型
config.Model = asr.ModelLarge

// 指定语言（不要自动检测）
config.Language = "zh"

// 增加beam size
config.BeamSize = 10
config.BestOf = 10
```

## 常见问题

### Q1: 提示 "Whisper未安装"？

```bash
# 安装Whisper
pip install -U openai-whisper

# 验证安装
whisper --help
```

### Q2: 转录速度慢？

**解决方法：**
1. 使用小模型（`ModelBase` 或 `ModelTiny`）
2. 启用GPU（如果有NVIDIA显卡）：`config.Device = "cuda"`
3. 减少 `BeamSize` 和 `BestOf`
4. 启用VAD：`config.EnableVAD = true`

### Q3: 准确率不高？

**解决方法：**
1. 使用大模型（`ModelMedium` 或 `ModelLarge`）
2. 明确指定语言：`config.Language = "zh"`
3. 增加 `BeamSize` 和 `BestOf`
4. 确保音频质量良好

### Q4: 如何下载模型？

Whisper会在第一次使用时自动下载模型，模型保存在：
- **macOS/Linux**: `~/.cache/whisper/`
- **Windows**: `C:\Users\<用户名>\.cache\whisper\`

手动下载：
```bash
# 下载medium模型（推荐）
whisper --model medium dummy.wav
```

### Q5: 支持离线使用吗？

是的！Whisper完全本地运行，一旦下载了模型文件就可以完全离线使用。

### Q6: 如何处理超大文件？

```go
config := asr.DefaultConfig()
config.MaxSegmentLen = 300       // 分段处理
config.Timeout = 60 * time.Minute // 增加超时
config.EnableVAD = true          // 跳过静音部分
```

### Q7: GPU加速如何启用？

```bash
# 1. 安装CUDA版本的PyTorch
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118

# 2. 配置使用GPU
config.Device = "cuda"
```

## 模型下载地址（手动下载）

如果网络环境不好，可以手动下载模型：

1. 访问：https://huggingface.co/ggerganov/whisper.cpp/tree/main
2. 下载对应模型文件
3. 放到 `~/.cache/whisper/` 目录

## 支持的语言

Whisper支持99种语言，常用的包括：

- **中文**: `zh` (普通话)
- **英语**: `en`
- **日语**: `ja`
- **韩语**: `ko`
- **法语**: `fr`
- **德语**: `de`
- **西班牙语**: `es`
- **粤语**: `yue`

完整列表：https://github.com/openai/whisper#available-models-and-languages

## 性能参考

**测试环境**: MacBook Pro M1, 16GB RAM

| 音频时长 | 模型 | 转录时间 | 准确率 |
|----------|------|----------|--------|
| 1分钟 | base | 10秒 | 85% |
| 1分钟 | medium | 30秒 | 95% |
| 10分钟 | base | 1.5分钟 | 85% |
| 10分钟 | medium | 5分钟 | 95% |
| 60分钟 | medium | 30分钟 | 95% |

## 注意事项

1. **首次使用**会自动下载模型，需要等待（medium模型约1.5GB）
2. **音频质量**越好，识别准确率越高
3. **长音频**建议使用 `MaxSegmentLen` 分段处理
4. **多人对话**可能需要手动区分说话人
5. **方言识别**准确率可能较低，建议使用专门的方言模型

## 完整示例项目

查看 `asr/example.go` 获取更多使用示例。

## License

基于 OpenAI Whisper 开源协议。
