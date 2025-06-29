package detect

var mimeTypeToExt = map[string]string{
	"image/jpeg": ".jpeg",
	"image/png":  ".png",
	"image/gif":  ".gif",
}

func ExtensionFromMimeType(mimeType string) (string, bool) {
	ext, ok := mimeTypeToExt[mimeType]
	return ext, ok
}
