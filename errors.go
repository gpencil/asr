package asr

import "fmt"

// ErrWhisperNotInstalled Whisper未安装错误
type ErrWhisperNotInstalled struct {
	Message string
}

func (e *ErrWhisperNotInstalled) Error() string {
	return fmt.Sprintf("Whisper未安装: %s\n\n安装方法：\n  pip install -U openai-whisper\n  # 或\n  pip3 install -U openai-whisper", e.Message)
}

// ErrInvalidAudioFile 音频文件无效错误
type ErrInvalidAudioFile struct {
	FilePath string
	Reason   string
}

func (e *ErrInvalidAudioFile) Error() string {
	return fmt.Sprintf("音频文件无效: %s, 原因: %s", e.FilePath, e.Reason)
}

// ErrTranscriptionFailed 转录失败错误
type ErrTranscriptionFailed struct {
	FilePath string
	Err      error
}

func (e *ErrTranscriptionFailed) Error() string {
	return fmt.Sprintf("转录失败: %s, 错误: %v", e.FilePath, e.Err)
}

func (e *ErrTranscriptionFailed) Unwrap() error {
	return e.Err
}

// ErrTimeout 超时错误
type ErrTimeout struct {
	Duration string
}

func (e *ErrTimeout) Error() string {
	return fmt.Sprintf("转录超时: %s", e.Duration)
}
