package services

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"main/models"
)

func MatchRuleForFile(filePath string, rules []models.Rule) *models.Rule {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil
	}

	fileType := detectFileType(filePath, info)
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), ".")

	for i := range rules {
		rule := &rules[i]
		if !rule.Enabled {
			continue
		}
		if rule.AllowAllFiles {
			return rule
		}
		if ext != "" && matchExtension(ext, rule.CustomExtensions) {
			return rule
		}
		if fileType != "" && containsString(rule.FileTypes, fileType) {
			return rule
		}
	}

	return nil
}

func BuildNameFromTemplate(template []string, originalName string, aiName string, t time.Time) string {
	parts := make([]string, 0, len(template))
	originalBase := strings.TrimSuffix(originalName, filepath.Ext(originalName))
	aiBase := sanitizeName(strings.TrimSuffix(aiName, filepath.Ext(aiName)))

	for _, part := range template {
		switch part {
		case "YYYY":
			parts = append(parts, t.Format("2006"))
		case "MM":
			parts = append(parts, t.Format("01"))
		case "DD":
			parts = append(parts, t.Format("02"))
		case "HH":
			parts = append(parts, t.Format("15"))
		case "mm":
			parts = append(parts, t.Format("04"))
		case "original":
			if aiBase != "" {
				parts = append(parts, aiBase)
			} else {
				parts = append(parts, sanitizeName(originalBase))
			}
		default:
			if strings.HasPrefix(part, "separator") {
				separator := strings.TrimPrefix(part, "separator")
				parts = append(parts, separator)
			} else {
				parts = append(parts, sanitizeName(part))
			}
		}
	}

	return strings.Join(parts, "")
}

func SelectTimestamp(filePath string, dateSource string) time.Time {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Now()
	}

	switch dateSource {
	case "created", "modified":
		return info.ModTime()
	default:
		return time.Now()
	}
}

func detectFileType(filePath string, info os.FileInfo) string {
	if info.IsDir() {
		return "folder"
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff", ".heic":
		return "image"
	case ".mp4", ".mov", ".avi", ".mkv", ".wmv", ".flv", ".webm":
		return "video"
	case ".mp3", ".wav", ".aac", ".flac", ".m4a", ".ogg":
		return "audio"
	case ".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz":
		return "archive"
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt", ".md", ".rtf", ".csv":
		return "document"
	case ".go", ".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".cpp", ".c", ".h", ".rs", ".rb", ".php", ".html", ".css", ".json", ".yaml", ".yml":
		return "code"
	case ".dmg", ".pkg", ".exe", ".msi":
		return "installer"
	case ".psd", ".sketch", ".ai", ".xd", ".fig":
		return "design"
	case ".epub", ".mobi", ".azw", ".azw3":
		return "ebook"
	default:
		return ""
	}
}

func containsString(items []string, value string) bool {
	for _, item := range items {
		if strings.EqualFold(item, value) {
			return true
		}
	}
	return false
}

func matchExtension(ext string, items []string) bool {
	normalized := strings.TrimPrefix(strings.ToLower(ext), ".")
	for _, item := range items {
		if strings.TrimPrefix(strings.ToLower(item), ".") == normalized {
			return true
		}
	}
	return false
}

func sanitizeName(value string) string {
	name := strings.TrimSpace(value)
	if name == "" {
		return "untitled"
	}
	replacer := strings.NewReplacer("/", "_", "\\", "_", ":", "_")
	return replacer.Replace(name)
}
