package review

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"steamshopemulator/internal/domain"
)

type Service struct {
	projectRoot string
	uploadDir   string
}

func New(projectRoot, uploadDir string) *Service {
	return &Service{
		projectRoot: projectRoot,
		uploadDir:   uploadDir,
	}
}

func (s *Service) Evaluate(page domain.Storefront) domain.ReviewReport {
	score := 100
	highlights := make([]string, 0, 6)
	warnings := make([]string, 0, 6)
	suggestions := make([]string, 0, 6)

	if len(strings.TrimSpace(page.Title)) >= 8 {
		highlights = append(highlights, "标题长度足够，商店页头部辨识度正常。")
	} else {
		score -= 10
		warnings = append(warnings, "标题过短，可能影响品牌辨识度。")
		suggestions = append(suggestions, "将标题控制在 8-48 个字符之间。")
	}

	if descLen := len([]rune(strings.TrimSpace(page.ShortDescription))); descLen >= 60 && descLen <= 180 {
		highlights = append(highlights, "短描述长度接近 Steam 商店常见有效区间。")
	} else {
		score -= 12
		warnings = append(warnings, "短描述长度不理想，首屏转化可能偏弱。")
		suggestions = append(suggestions, "短描述建议维持在 60-180 个字符，突出题材、玩法和卖点。")
	}

	if count := len(page.Tags); count >= 8 {
		highlights = append(highlights, "标签覆盖较完整，便于模拟 Steam 标签曝光。")
	} else {
		score -= 10
		warnings = append(warnings, "标签数量偏少，难以测试标签组合展示效果。")
		suggestions = append(suggestions, "至少维护 8 个标签，并覆盖题材、机制、视角、节奏四类信息。")
	}

	if len(page.Gallery) >= 4 {
		highlights = append(highlights, "图库数量足够，可以测试轮播与截图节奏。")
	} else {
		score -= 10
		warnings = append(warnings, "截图数量不足，无法模拟完整商店详情节奏。")
		suggestions = append(suggestions, "至少上传 4 张不同场景截图，分别覆盖战斗、UI、环境和奖励反馈。")
	}

	heroPath := page.HeroImage
	if len(page.HeroThumbs) > 0 && strings.TrimSpace(page.HeroThumbs[0]) != "" {
		heroPath = page.HeroThumbs[0]
	}

	s.scoreImage(&score, &highlights, &warnings, &suggestions, "hero", heroPath, 1200, 620)
	s.scoreImage(&score, &highlights, &warnings, &suggestions, "capsule", page.CapsuleImage, 616, 353)

	status := "优秀"
	switch {
	case score < 60:
		status = "需重做"
	case score < 80:
		status = "可优化"
	}

	return domain.ReviewReport{
		Score:       max(score, 0),
		Status:      status,
		Highlights:  highlights,
		Warnings:    warnings,
		Suggestions: suggestions,
	}
}

func (s *Service) scoreImage(score *int, highlights, warnings, suggestions *[]string, label, path string, wantW, wantH int) {
	if strings.TrimSpace(path) == "" {
		*score -= 12
		*warnings = append(*warnings, label+" 素材缺失。")
		*suggestions = append(*suggestions, "补齐 "+label+" 素材，并尽量贴近 Steam 常用展示比例。")
		return
	}

	fullPath := filepath.Join(s.projectRoot, strings.TrimPrefix(path, "/"))
	if strings.HasPrefix(path, "/uploads/") {
		fullPath = filepath.Join(s.uploadDir, filepath.Base(path))
	}

	file, err := os.Open(fullPath)
	if err != nil {
		*score -= 8
		*warnings = append(*warnings, label+" 素材无法读取。")
		return
	}
	defer file.Close()

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		*score -= 8
		*warnings = append(*warnings, label+" 素材不是可识别图片。")
		return
	}

	if cfg.Width >= wantW && cfg.Height >= wantH {
		*highlights = append(*highlights, label+" 素材尺寸达标。")
		return
	}

	*score -= 8
	*warnings = append(*warnings, label+" 素材尺寸偏小。")
	*suggestions = append(*suggestions, label+" 建议至少达到 "+dimensionText(wantW, wantH)+"。")
}

func dimensionText(w, h int) string {
	return strconv.Itoa(w) + "x" + strconv.Itoa(h)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
