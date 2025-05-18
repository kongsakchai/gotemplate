package generate

import "github.com/google/uuid"

var fixedUUID = ""

func SetFixedUUID(uuidStr string) {
	fixedUUID = uuidStr
}

func ClearFixedUUID() {
	fixedUUID = ""
}

func UUID() string {
	if fixedUUID != "" {
		return fixedUUID
	}

	return uuid.NewString()
}
