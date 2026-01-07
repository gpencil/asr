package asr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// WhisperClient Whisper客户端
type WhisperClient struct {
	config *WhisperConfig
}

// NewWhisperClient 创建Whisper客户端
func NewWhisperClient(config *WhisperConfig) (*WhisperClient, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 检查whisper是否已安装
	if err := checkWhisperInstalled(); err != nil {
		return nil, err
	}

	return &WhisperClient{
		config: config,
	}, nil
}

// checkWhisperInstalled 检查whisper是否已安装
func checkWhisperInstalled() error {
	cmd := exec.Command("whisper", "--help")
	if err := cmd.Run(); err != nil {
		return &ErrWhisperNotInstalled{
			Message: "未找到whisper命令",
		}
	}
	return nil
}

// Transcribe 转录音频文件
func (c *WhisperClient) Transcribe(audioFilePath string) (*TranscriptionResult, error) {
	return c.TranscribeWithContext(context.Background(), audioFilePath)
}

// TranscribeWithContext 转录音频文件（带上下文）
func (c *WhisperClient) TranscribeWithContext(ctx context.Context, audioFilePath string) (*TranscriptionResult, error) {
	// 检查文件是否存在
	if _, err := os.Stat(audioFilePath); os.IsNotExist(err) {
		return nil, &ErrInvalidAudioFile{
			FilePath: audioFilePath,
			Reason:   "文件不存在",
		}
	}

	// 构建命令参数
	args := c.buildWhisperArgs(audioFilePath)

	// 创建临时输出目录
	outputDir, err := os.MkdirTemp("", "whisper_output_*")
	if err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(outputDir)

	// 添加输出目录参数
	args = append(args, "--output_dir", outputDir)

	// 创建带超时的上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	// 执行whisper命令
	cmd := exec.CommandContext(timeoutCtx, "whisper", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if c.config.Verbose {
		fmt.Printf("[Whisper] 执行命令: whisper %s\n", strings.Join(args, " "))
	}

	startTime := time.Now()
	err = cmd.Run()
	elapsed := time.Since(startTime)

	if err != nil {
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return nil, &ErrTimeout{Duration: c.config.Timeout.String()}
		}
		return nil, &ErrTranscriptionFailed{
			FilePath: audioFilePath,
			Err:      fmt.Errorf("%w: %s", err, stderr.String()),
		}
	}

	if c.config.Verbose {
		fmt.Printf("[Whisper] 转录完成，耗时: %v\n", elapsed)
	}

	// 解析结果
	result, err := c.parseResult(audioFilePath, outputDir)
	if err != nil {
		return nil, err
	}

	result.Model = c.config.Model
	return result, nil
}

// buildWhisperArgs 构建whisper命令参数
func (c *WhisperClient) buildWhisperArgs(audioFilePath string) []string {
	args := []string{audioFilePath}

	// 模型
	args = append(args, "--model", string(c.config.Model))

	// 语言
	if c.config.Language != "" {
		args = append(args, "--language", c.config.Language)
	}

	// 输出格式
	args = append(args, "--output_format", c.config.OutputFormat)

	// 设备
	if c.config.Device != "" {
		args = append(args, "--device", c.config.Device)
	}

	// 线程数
	if c.config.Threads > 0 {
		args = append(args, "--threads", fmt.Sprintf("%d", c.config.Threads))
	}

	// Beam size
	if c.config.BeamSize > 0 {
		args = append(args, "--beam_size", fmt.Sprintf("%d", c.config.BeamSize))
	}

	// Best of
	if c.config.BestOf > 0 {
		args = append(args, "--best_of", fmt.Sprintf("%d", c.config.BestOf))
	}

	// Temperature
	if c.config.Temperature > 0 {
		args = append(args, "--temperature", fmt.Sprintf("%.2f", c.config.Temperature))
	}

	// 详细日志
	if c.config.Verbose {
		args = append(args, "--verbose", "True")
	}

	return args
}

// parseResult 解析转录结果
func (c *WhisperClient) parseResult(audioFilePath, outputDir string) (*TranscriptionResult, error) {
	baseName := strings.TrimSuffix(filepath.Base(audioFilePath), filepath.Ext(audioFilePath))

	result := &TranscriptionResult{
		Segments: make([]TextSegment, 0),
	}

	// 读取文本结果
	txtFile := filepath.Join(outputDir, baseName+".txt")
	if data, err := os.ReadFile(txtFile); err == nil {
		result.Text = string(data)
		result.FilePath = txtFile
	}

	// 尝试读取JSON结果（包含详细信息）
	jsonFile := filepath.Join(outputDir, baseName+".json")
	if data, err := os.ReadFile(jsonFile); err == nil {
		var jsonResult struct {
			Text     string  `json:"text"`
			Language string  `json:"language"`
			Duration float64 `json:"duration"`
			Segments []struct {
				ID    int     `json:"id"`
				Seek  int     `json:"seek"`
				Start float64 `json:"start"`
				End   float64 `json:"end"`
				Text  string  `json:"text"`
			} `json:"segments"`
		}

		if err := json.Unmarshal(data, &jsonResult); err == nil {
			result.Text = jsonResult.Text
			result.Language = jsonResult.Language
			result.Duration = jsonResult.Duration

			for _, seg := range jsonResult.Segments {
				result.Segments = append(result.Segments, TextSegment{
					ID:    seg.ID,
					Start: seg.Start,
					End:   seg.End,
					Text:  strings.TrimSpace(seg.Text),
				})
			}
		}
	}

	return result, nil
}

// TranscribeBatch 批量转录音频文件
func (c *WhisperClient) TranscribeBatch(audioFiles []string) ([]*TranscriptionResult, error) {
	results := make([]*TranscriptionResult, 0, len(audioFiles))

	for i, file := range audioFiles {
		if c.config.Verbose {
			fmt.Printf("[Whisper] 处理文件 %d/%d: %s\n", i+1, len(audioFiles), file)
		}

		result, err := c.Transcribe(file)
		if err != nil {
			fmt.Printf("[Whisper] 文件 %s 转录失败: %v\n", file, err)
			continue
		}

		results = append(results, result)
	}

	return results, nil
}

// GetModelInfo 获取模型信息
func GetModelInfo(model WhisperModel) map[string]interface{} {
	modelInfo := map[WhisperModel]map[string]interface{}{
		ModelTiny: {
			"name":        "tiny",
			"size":        "75 MB",
			"speed":       "极快",
			"accuracy":    "较低",
			"recommended": "短音频、快速预览",
		},
		ModelBase: {
			"name":        "base",
			"size":        "142 MB",
			"speed":       "快",
			"accuracy":    "一般",
			"recommended": "日常使用、快速转录",
		},
		ModelSmall: {
			"name":        "small",
			"size":        "466 MB",
			"speed":       "中等",
			"accuracy":    "较好",
			"recommended": "平衡速度和准确率",
		},
		ModelMedium: {
			"name":        "medium",
			"size":        "1.5 GB",
			"speed":       "较慢",
			"accuracy":    "高",
			"recommended": "长音频、会议记录（推荐）",
		},
		ModelLarge: {
			"name":        "large",
			"size":        "3 GB",
			"speed":       "慢",
			"accuracy":    "最高",
			"recommended": "专业转录、高要求场景",
		},
	}

	if info, ok := modelInfo[model]; ok {
		return info
	}
	return nil
}
