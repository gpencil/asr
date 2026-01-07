package asr

import (
	"fmt"
	"log"
)

// Example1_BasicUsage 基本使用示例
func Example1_BasicUsage() {
	// 创建客户端（使用默认配置）
	client, err := NewWhisperClient(nil)
	if err != nil {
		log.Fatalf("创建Whisper客户端失败: %v", err)
	}

	// 转录音频文件
	result, err := client.Transcribe("meeting.mp3")
	if err != nil {
		log.Fatalf("转录失败: %v", err)
	}

	// 输出结果
	fmt.Printf("检测语言: %s\n", result.Language)
	fmt.Printf("音频时长: %.2f 秒\n", result.Duration)
	fmt.Printf("文本内容:\n%s\n", result.Text)
}

// Example2_CustomConfig 自定义配置
func Example2_CustomConfig() {
	// 创建自定义配置（长音频、高准确率）
	config := &WhisperConfig{
		Model:          ModelMedium, // 使用medium模型
		Language:       "zh",        // 指定中文
		OutputFormat:   "json",      // 输出JSON格式（包含时间戳）
		Verbose:        true,        // 显示详细日志
		Device:         "cpu",       // 使用CPU（如果有GPU可改为"cuda"）
		EnableVAD:      true,        // 启用VAD
		MaxSegmentLen:  300,         // 5分钟一段
		SplitOnSilence: true,        // 在静音处分割
	}

	client, err := NewWhisperClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	result, err := client.Transcribe("long_meeting.mp3")
	if err != nil {
		log.Fatalf("转录失败: %v", err)
	}

	// 输出带时间戳的分段文本
	fmt.Println("分段文本:")
	for _, segment := range result.Segments {
		fmt.Printf("[%.2f - %.2f] %s\n", segment.Start, segment.End, segment.Text)
	}
}

// Example3_FastMode 快速模式（牺牲准确率换速度）
func Example3_FastMode() {
	// 使用快速配置
	config := FastConfig()
	config.Verbose = true

	client, err := NewWhisperClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	result, err := client.Transcribe("voice_message.mp3")
	if err != nil {
		log.Fatalf("转录失败: %v", err)
	}

	fmt.Printf("快速转录结果:\n%s\n", result.Text)
}

// Example4_AccurateMode 高准确率模式
func Example4_AccurateMode() {
	// 使用高准确率配置
	config := AccurateConfig()
	config.Language = "zh" // 指定中文
	config.Verbose = true

	client, err := NewWhisperClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	result, err := client.Transcribe("important_interview.mp3")
	if err != nil {
		log.Fatalf("转录失败: %v", err)
	}

	fmt.Printf("高精度转录结果:\n%s\n", result.Text)
	fmt.Printf("使用模型: %s\n", result.Model)
}

// Example5_BatchTranscription 批量转录
func Example5_BatchTranscription() {
	config := DefaultConfig()
	config.Verbose = true

	client, err := NewWhisperClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	// 批量转录多个文件
	files := []string{
		"meeting1.mp3",
		"meeting2.mp3",
		"meeting3.mp3",
	}

	results, err := client.TranscribeBatch(files)
	if err != nil {
		log.Fatalf("批量转录失败: %v", err)
	}

	// 输出每个文件的结果
	for i, result := range results {
		fmt.Printf("\n文件 %d:\n", i+1)
		fmt.Printf("时长: %.2f秒\n", result.Duration)
		fmt.Printf("文本: %s\n", result.Text)
	}
}

// Example6_MultiLanguage 多语言支持
func Example6_MultiLanguage() {
	config := DefaultConfig()

	// 英文转录
	config.Language = "en"
	client, _ := NewWhisperClient(config)
	result, _ := client.Transcribe("english_audio.mp3")
	fmt.Printf("English: %s\n", result.Text)

	// 中文转录
	config.Language = "zh"
	client, _ = NewWhisperClient(config)
	result, _ = client.Transcribe("chinese_audio.mp3")
	fmt.Printf("中文: %s\n", result.Text)

	// 自动检测语言
	config.Language = ""
	client, _ = NewWhisperClient(config)
	result, _ = client.Transcribe("unknown_language.mp3")
	fmt.Printf("检测到语言: %s, 文本: %s\n", result.Language, result.Text)
}

// Example7_ModelComparison 不同模型对比
func Example7_ModelComparison() {
	audioFile := "test.mp3"

	models := []WhisperModel{ModelTiny, ModelBase, ModelMedium}

	for _, model := range models {
		config := DefaultConfig()
		config.Model = model
		config.Verbose = false

		client, _ := NewWhisperClient(config)

		// 获取模型信息
		info := GetModelInfo(model)
		fmt.Printf("\n测试模型: %s\n", info["name"])
		fmt.Printf("  大小: %s\n", info["size"])
		fmt.Printf("  速度: %s\n", info["speed"])
		fmt.Printf("  准确率: %s\n", info["accuracy"])

		result, err := client.Transcribe(audioFile)
		if err != nil {
			fmt.Printf("  转录失败: %v\n", err)
			continue
		}

		fmt.Printf("  结果: %s\n", result.Text[:min(100, len(result.Text))])
	}
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Example8_SaveToFile 保存到文件
func Example8_SaveToFile() {
	config := DefaultConfig()
	config.OutputFormat = "all" // 保存所有格式：txt, srt, vtt, json

	client, err := NewWhisperClient(config)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}

	result, err := client.Transcribe("meeting.mp3")
	if err != nil {
		log.Fatalf("转录失败: %v", err)
	}

	fmt.Printf("结果已保存到: %s\n", result.FilePath)
	fmt.Println("同时生成了以下文件:")
	fmt.Println("  - meeting.txt (纯文本)")
	fmt.Println("  - meeting.srt (字幕文件)")
	fmt.Println("  - meeting.vtt (字幕文件)")
	fmt.Println("  - meeting.json (详细JSON)")
}
