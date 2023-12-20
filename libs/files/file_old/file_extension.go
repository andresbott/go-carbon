package file_old

import "strings"

type FileTypeExtInfo struct {
	MimeType string
	Type     FileType
}

func fileTypeFromExtension(ext string) FileTypeExtInfo {

	switch strings.ToLower(ext) {
	// Images
	case "jpeg":
	case "jpg":
		return FileTypeExtInfo{
			MimeType: "image/jpeg",
			Type:     IMAGE,
		}

	case "gif":
		return FileTypeExtInfo{
			MimeType: "image/gif",
			Type:     IMAGE,
		}
	case "bmp":
		return FileTypeExtInfo{
			MimeType: "image/bmp",
			Type:     IMAGE,
		}
	case "png":
		return FileTypeExtInfo{
			MimeType: "image/png",
			Type:     IMAGE,
		}
	case "svg":
		return FileTypeExtInfo{
			MimeType: "image/svg+xml",
			Type:     IMAGE,
		}
	case "webp":
		return FileTypeExtInfo{
			MimeType: "image/webp",
			Type:     IMAGE,
		}
	case "tiff":
		return FileTypeExtInfo{
			MimeType: "image/tiff",
			Type:     IMAGE,
		}
		// audio
	case "mp3":
		return FileTypeExtInfo{
			MimeType: "audio/mpeg",
			Type:     AUDIO,
		}
	case "ogg":
		return FileTypeExtInfo{
			MimeType: "audio/ogg",
			Type:     AUDIO,
		}
	case "mid":
	case "midi":
		return FileTypeExtInfo{
			MimeType: "audio/midi",
			Type:     AUDIO,
		}
	case "aac":
		return FileTypeExtInfo{
			MimeType: "audio/aac",
			Type:     AUDIO,
		}
	case "wav":
		return FileTypeExtInfo{
			MimeType: "audio/wav",
			Type:     AUDIO,
		}
	case "weba":
		return FileTypeExtInfo{
			MimeType: "audio/webm",
			Type:     AUDIO,
		}
		// video
	case "mpeg":
		return FileTypeExtInfo{
			MimeType: "video/mpeg",
			Type:     VIDEO,
		}
	case "mpg":
		return FileTypeExtInfo{
			MimeType: "video/mpeg",
			Type:     VIDEO,
		}
	case "ogv":
		return FileTypeExtInfo{
			MimeType: "video/ogg",
			Type:     VIDEO,
		}
	case "webm":
		return FileTypeExtInfo{
			MimeType: "video/webm",
			Type:     VIDEO,
		}
	case "avi":
		return FileTypeExtInfo{
			MimeType: "video/x-msvideo",
			Type:     VIDEO,
		}
	case "mp4":
		return FileTypeExtInfo{
			MimeType: "video/mp4",
			Type:     VIDEO,
		}
	case "ts":
		return FileTypeExtInfo{
			MimeType: "video/MP2T",
			Type:     VIDEO,
		}
	case "mkv":
		return FileTypeExtInfo{
			MimeType: "video/x-matroska",
			Type:     VIDEO,
		}
		// File compression
	case "bz":
	case "tar.bz":
		return FileTypeExtInfo{
			MimeType: "application/x-bzip",
			Type:     ZIP,
		}
	case "bz2":
		return FileTypeExtInfo{
			MimeType: "application/x-bzip2",
			Type:     ZIP,
		}
	case "rar":
		return FileTypeExtInfo{
			MimeType: "application/x-rar-compressed",
			Type:     ZIP,
		}
	case "tar":
		return FileTypeExtInfo{
			MimeType: "application/x-tar",
			Type:     ZIP,
		}
	case "zip":
		return FileTypeExtInfo{
			MimeType: "application/zip",
			Type:     ZIP,
		}
	case "7z":
		return FileTypeExtInfo{
			MimeType: "application/x-7z-compressed",
			Type:     ZIP,
		}
		// text Files
	case "txt":
		return FileTypeExtInfo{
			MimeType: "text/plain",
			Type:     TEXT,
		}
	case "css":
		return FileTypeExtInfo{
			MimeType: "text/css",
			Type:     TEXT,
		}
	case "csv":
		return FileTypeExtInfo{
			MimeType: "text/csv",
			Type:     TEXT,
		}
	case "html":
		return FileTypeExtInfo{
			MimeType: "text/html",
			Type:     WEB,
		}
	case "js":
		return FileTypeExtInfo{
			MimeType: "text/javascript",
			Type:     TEXT,
		}
	case "odt":
		return FileTypeExtInfo{
			MimeType: "application/vnd.oasis.opendocument.text",
			Type:     TEXT,
		}
	case "rtf":
		return FileTypeExtInfo{
			MimeType: "text/rtf",
			Type:     TEXT,
		}
	case "xml":
		return FileTypeExtInfo{
			MimeType: "text/xml",
			Type:     TEXT,
		}
		// Other
	case "exe":
		return FileTypeExtInfo{
			MimeType: "application/octet-stream",
			Type:     EXEC,
		}
	case "sh":
		return FileTypeExtInfo{
			MimeType: "text/plain",
			Type:     EXEC,
		}
	case "bash":
		return FileTypeExtInfo{
			MimeType: "text/plain",
			Type:     EXEC,
		}
	case "py":
		return FileTypeExtInfo{
			MimeType: "text/plain",
			Type:     EXEC,
		}
	case "php":
		return FileTypeExtInfo{
			MimeType: "text/plain",
			Type:     EXEC,
		}
	case "pdf":
		return FileTypeExtInfo{
			MimeType: "application/pdf",
			Type:     PDF,
		}

	default:
		return FileTypeExtInfo{
			MimeType: "application/octet-stream",
			Type:     FILE,
		}
	}

	return FileTypeExtInfo{
		MimeType: "application/octet-stream",
		Type:     FILE,
	}
}
