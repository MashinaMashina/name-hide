package main

const (
	LnkExt = ".lnk"
	URLExt = ".url"
)

var availableExtensions = []string{
	LnkExt,
	URLExt,
}

func ExtAvailable(ext string) bool {
	for _, extension := range availableExtensions {
		if extension == ext {
			return true
		}
	}

	return false
}
