package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"steamshopemulator/internal/domain"
	"steamshopemulator/internal/review"
	"steamshopemulator/internal/store"
	"steamshopemulator/internal/web"
)

func New() (*http.Server, error) {
	projectDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dataPath := filepath.Join(projectDir, "data", "storefront.json")
	uploadDir := filepath.Join(projectDir, "data", "uploads")

	if err := store.EnsureDefault(dataPath, defaultStorefront()); err != nil {
		return nil, err
	}

	contentStore, err := store.New(dataPath)
	if err != nil {
		return nil, err
	}

	if err := normalizeStorefront(contentStore); err != nil {
		return nil, err
	}

	reviewService := review.New(projectDir, uploadDir)
	webServer, err := web.New(contentStore, reviewService, projectDir, uploadDir)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:    ":8080",
		Handler: webServer.Routes(),
	}, nil
}

func defaultStorefront() domain.Storefront {
	return domain.Storefront{
		AppID:            "311210",
		Slug:             "Call_of_Duty_Black_Ops_III",
		Title:            "PROJECT RETAIL OPS",
		Subtitle:         "未来战术动作样片页",
		ShortDescription: "使用接近 Steam 商品页的视觉结构，验证你的标题、截图、标签与转化卖点是否足够强。",
		LongDescription:  "这是一套用于模拟 Steam 商店详情页的 Go 项目。你可以替换主视觉、胶囊图、Logo、截图、标签与卖点文案，并实时得到素材完成度与首屏吸引力建议。",
		Developer:        "Your Studio",
		Publisher:        "Your Publishing Label",
		ReleaseDate:      "Coming Soon",
		Price:            "¥ 128.00",
		DiscountText:     "-15%",
		ReviewSummary:    "特别好评",
		ReviewCount:      "1,245",
		Tags: []string{
			"动作", "第一人称射击", "多人", "合作", "赛博朋克", "剧情", "竞技", "未来战争",
		},
		Features: []string{
			"高强度战斗循环", "明确职业分工", "强烈科技感包装", "便于验证 Steam 首屏视觉",
		},
		HeroThumbs: []string{"", "", "", ""},
		AboutBlocks: []domain.AboutBlock{
			{Heading: "关于这款游戏", Body: "这是一个 Steam 风格商店页原型，用于快速替换素材并观察最终展示效果。"},
			{Heading: "核心卖点", Body: "支持后台编辑、素材上传、标签调整，以及实时审查建议。"},
		},
		DetailSections: []domain.DetailSection{
			{
				Heading: "黑暗围城，征途再启",
				Paragraphs: []string{
					"探索地图，强化能力，击败怪物，搜寻宝藏。",
					"将高密度动作体验转化为更强调节奏感、路线选择和资源管理的展示段落，让你可以验证 Steam 长详情里的图文组织能力。",
				},
			},
			{
				Heading: "滚雪球式的回合制杀戮",
				Paragraphs: []string{
					"按法力值从低到高依序打出卡牌，逐步叠加连锁效果，每一步都会放大下一张卡的威力。",
					"这里预留了大图位，方便你放入战斗截图、UI 展示图或宣传画面。",
				},
			},
			{
				Heading: "自由掌控节奏",
				Paragraphs: []string{
					"自己决定玩法风格，是谨慎规划、步步为营，或是极限手速、一气呵成一大套连招。",
					"这个区块用于承接更长的卖点说明，并模拟原始商店页在正文后半段的阅读密度。",
				},
			},
		},
		ReadMoreLabel: "继续阅读",
		SystemRequirements: domain.SystemRequirements{
			Platforms: []domain.PlatformRequirements{
				{
					Name:     "Windows",
					IsActive: true,
					Minimum: []string{
						"需要 64 位元的处理器及作业系统",
						"作业系统: Windows 10 64bit",
						"处理器: x64 architecture with SSE2",
						"记忆体: 4 GB 记忆体",
						"显示卡: DX11, DX12 capable",
						"储存空间: 2 GB 可用空间",
						"音效卡: Sound Blaster 4000",
					},
					Recommended: []string{
						"需要 64 位元的处理器及作业系统",
						"作业系统: Windows 10 64bit",
						"处理器: x64 architecture with SSE2",
						"记忆体: 4 GB 记忆体",
						"显示卡: DX11, DX12 capable",
						"储存空间: 2 GB 可用空间",
						"音效卡: Sound Mega Blaster 90000",
					},
				},
				{
					Name:     "macOS",
					IsActive: false,
					Minimum: []string{
						"Apple Silicon 或 Intel 处理器",
					},
					Recommended: []string{
						"Apple Silicon",
					},
				},
			},
		},
		SidebarHighlights: []domain.IconLabel{
			{Icon: "👤", Label: "单人"},
			{Icon: "✹", Label: "Steam 成就"},
			{Icon: "☁", Label: "Steam 云端"},
			{Icon: "👥", Label: "亲友同乐"},
		},
		ControllerSupports: []domain.IconLabel{
			{Icon: "🎮", Label: "Xbox 控制器"},
			{Icon: "🎮", Label: "PlayStation 控制器"},
		},
		LanguageSupport: []domain.LanguageSupport{
			{Name: "简体中文", Interface: true, FullAudio: false, Subtitles: true},
			{Name: "英文", Interface: true, FullAudio: true, Subtitles: true},
			{Name: "法文", Interface: true, FullAudio: false, Subtitles: true},
			{Name: "德文", Interface: true, FullAudio: false, Subtitles: true},
		},
		SteamDeck: domain.SteamDeckInfo{
			Status: "已验证",
			CTA:    "深入了解",
		},
		AchievementCount: "32",
		AchievementBadges: []domain.AchievementBadge{
			{Label: "◇"},
			{Label: "✕"},
			{Label: "☠"},
			{Label: "✦"},
		},
		MetadataEntries: []domain.KeyValue{
			{Key: "名称", Value: "PROJECT RETAIL OPS"},
			{Key: "类型", Value: "动作、独立制作、策略"},
			{Key: "开发者", Value: "Your Studio"},
			{Key: "发行商", Value: "Your Publishing Label"},
			{Key: "系列", Value: "Retail Ops Universe"},
			{Key: "发行日期", Value: "Coming Soon"},
		},
		CommunityLinks: []domain.ExternalLink{
			{Label: "X", URL: "#"},
			{Label: "YouTube", URL: "#"},
			{Label: "Discord", URL: "#"},
			{Label: "TikTok", URL: "#"},
			{Label: "Bilibili", URL: "#"},
			{Label: "QQ 群", URL: "#"},
		},
		DiscoveryLinks: []domain.ExternalLink{
			{Label: "检视更新历史记录", URL: "#"},
			{Label: "阅读相关新闻", URL: "#"},
			{Label: "检视讨论区", URL: "#"},
			{Label: "寻找社群群组", URL: "#"},
		},
	}
}

func normalizeStorefront(contentStore *store.Store) error {
	current := contentStore.Get()
	changed := false

	if len(current.HeroThumbs) < 4 {
		next := make([]string, 4)
		copy(next, current.HeroThumbs)
		current.HeroThumbs = next
		changed = true
	}

	if len(current.DetailSections) < 3 {
		next := make([]domain.DetailSection, 3)
		copy(next, current.DetailSections)
		current.DetailSections = next
		changed = true
	}

	if current.HeroThumbs[0] == "" && current.HeroImage != "" {
		current.HeroThumbs[0] = current.HeroImage
		changed = true
	}

	if strings.TrimSpace(current.ContentHTML) == "" {
		current.ContentHTML = legacyContentHTML(current)
		changed = true
	}

	if !changed {
		return nil
	}

	return contentStore.Save(current)
}

func legacyContentHTML(page domain.Storefront) string {
	var builder strings.Builder
	if strings.TrimSpace(page.ShortDescription) != "" {
		builder.WriteString("<p>")
		builder.WriteString(page.ShortDescription)
		builder.WriteString("</p>")
	}
	if strings.TrimSpace(page.LongDescription) != "" {
		builder.WriteString("<p>")
		builder.WriteString(page.LongDescription)
		builder.WriteString("</p>")
	}
	for _, section := range page.DetailSections {
		if strings.TrimSpace(section.Heading) != "" {
			builder.WriteString("<h3>")
			builder.WriteString(section.Heading)
			builder.WriteString("</h3>")
		}
		for _, paragraph := range section.Paragraphs {
			if strings.TrimSpace(paragraph) == "" {
				continue
			}
			builder.WriteString("<p>")
			builder.WriteString(paragraph)
			builder.WriteString("</p>")
		}
		if strings.TrimSpace(section.Image) != "" {
			builder.WriteString(fmt.Sprintf(`<p><img src="%s" alt="%s"></p>`, section.Image, section.Heading))
		}
	}
	return builder.String()
}
