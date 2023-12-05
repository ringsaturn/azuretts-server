package main

import (
	"context"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/ringsaturn/azuretts"
)

var az = azuretts.NewClient(
	os.Getenv("SPEECH_KEY"),
	azuretts.Region(os.Getenv("SPEECH_REGION")),
)

type TTSRequest struct {
	Language   azuretts.Language  `query:"language" default:"zh-CN"`
	VoiceName  azuretts.VoiceName `query:"voice" default:"zh-CN-XiaoxiaoNeural"`
	Style      azuretts.Style     `query:"style" default:"chat"`
	Rate       float64            `query:"rate" default:"1.0"`
	Degree     int                `query:"degree" default:"0"`
	SpeechText string             `query:"speech_text" default:"你好，世界！"`
	Volume     int                `query:"volume" default:"100"`
}

func TTS(ctx context.Context, c *app.RequestContext) {
	var req TTSRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}
	opts := []azuretts.SpeakOption{}
	opts = append(opts, azuretts.WithSpeechText(req.SpeechText))

	if req.Language != "" {
		opts = append(opts, azuretts.WithLanguage(req.Language))
	}
	if req.VoiceName != "" {
		opts = append(opts, azuretts.WithVoiceName(req.VoiceName))
	}
	if req.Style != "" {
		opts = append(opts, azuretts.WithStyle(req.Style))
	}
	if req.Rate != 0 {
		opts = append(opts, azuretts.WithRate(req.Rate))
	}
	if req.Degree != 0 {
		opts = append(opts, azuretts.WithVoiceStyledegree(float64(req.Degree)))
	}

	if req.Volume != 0 {
		opts = append(opts, azuretts.WithVolume(req.Volume))
	}

	speak := azuretts.NewSpeak(opts...)
	b, err := az.GetSynthesize(ctx, &azuretts.SynthesisRequest{
		Speak:  speak,
		Output: azuretts.AudioOutputFormat_Streaming_Audio16Khz32KbitrateMonoMp3,
	})
	if err != nil {
		c.String(500, err.Error())
	}
	if err := b.Error(); err != nil {
		c.String(500, err.Error())
	}
	c.Data(200, "audio/mpeg", b.Body)
}

func main() {
	e := server.Default()
	e.GET("/tts", TTS)
	_ = e.Run()
}
