package detect

var mimeTypeToExt = map[string]string{
	"image/jpeg": ".jpeg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

func ExtensionFromMimeType(mimeType string) (string, bool) {
	ext, ok := mimeTypeToExt[mimeType]
	return ext, ok
}
