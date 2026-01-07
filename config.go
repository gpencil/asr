package asr

import "time"

// WhisperModel Whisper模型类型
type WhisperModel string

const (
	ModelTiny   WhisperModel = "tiny"   // 最小模型 ~75MB，速度最快，准确率较低
	ModelBase   WhisperModel = "base"   // 基础模型 ~142MB，速度快，准确率一般
	ModelSmall  WhisperModel = "small"  // 小型模型 ~466MB，速度中等，准确率较好
	ModelMedium WhisperModel = "medium" // 中型模型 ~1.5GB，速度较慢，准确率高（推荐长音频）
	ModelLarge  WhisperModel = "large"  // 大型模型 ~3GB，速度慢，准确率最高
)

// WhisperConfig Whisper配置
type WhisperConfig struct {
	// 模型配置
	Model    WhisperModel // 使用的模型，推荐medium或large
	Language string       // 语言代码，如"zh"中文、"en"英文，空则自动检测

	// 输出配置
	OutputFormat string // 输出格式：txt, srt, vtt, json, all（默认txt）
	Verbose      bool   // 是否显示详细日志

	// 性能配置
	Device      string        // 设备：cpu, cuda（如果有GPU）
	Threads     int           // CPU线程数，默认0（自动）
	Timeout     time.Duration // 超时时间，默认10分钟
	BeamSize    int           // Beam search大小，越大越准确但越慢，默认5
	BestOf      int           // 候选数量，默认5
	Temperature float64       // 采样温度，0-1，默认0（贪婪）

	// 长音频优化
	EnableVAD      bool // 启用语音活动检测（VAD），自动去除静音部分
	MaxSegmentLen  int  // 最大分段长度（秒），0表示不分段，推荐300（5分钟）
	SplitOnSilence bool // 在静音处分割，适合长音频
}

// DefaultConfig 返回默认配置（针对长音频优化）
func DefaultConfig() *WhisperConfig {
	return &WhisperConfig{
		Model:          ModelMedium, // 长音频推荐medium
		Language:       "zh",        // 默认中文
		OutputFormat:   "txt",
		Verbose:        false,
		Device:         "cpu",
		Threads:        0,
		Timeout:        30 * time.Minute, // 长音频需要更长超时
		BeamSize:       5,
		BestOf:         5,
		Temperature:    0,
		EnableVAD:      true, // 启用VAD节省时间
		MaxSegmentLen:  300,  // 5分钟一段
		SplitOnSilence: true, // 在静音处分割
	}
}

// FastConfig 快速配置（牺牲准确率换速度）
func FastConfig() *WhisperConfig {
	config := DefaultConfig()
	config.Model = ModelBase // 使用小模型
	config.BeamSize = 1      // 减少beam size
	config.BestOf = 1
	config.MaxSegmentLen = 180 // 3分钟一段
	return config
}

// AccurateConfig 高准确率配置（牺牲速度换准确率）
func AccurateConfig() *WhisperConfig {
	config := DefaultConfig()
	config.Model = ModelLarge // 使用大模型
	config.BeamSize = 10      // 增加beam size
	config.BestOf = 10
	config.MaxSegmentLen = 600 // 10分钟一段
	config.Timeout = 60 * time.Minute
	return config
}

// TranscriptionResult 转录结果
type TranscriptionResult struct {
	Text     string        // 完整文本
	Language string        // 检测到的语言
	Duration float64       // 音频时长（秒）
	Segments []TextSegment // 分段文本（带时间戳）
	FilePath string        // 输出文件路径（如果保存了文件）
	Model    WhisperModel  // 使用的模型
}

// TextSegment 文本片段（带时间戳）
type TextSegment struct {
	ID    int     // 片段ID
	Start float64 // 开始时间（秒）
	End   float64 // 结束时间（秒）
	Text  string  // 文本内容
}
