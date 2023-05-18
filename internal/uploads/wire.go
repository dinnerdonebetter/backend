package uploads

import (
	"github.com/dinnerdonebetter/backend/internal/objectstorage"

	"github.com/google/wire"
)

var (
	// Providers represents this package's offering to the dependency manager.
	Providers = wire.NewSet(
		ProvideUploadManager,
		wire.FieldsOf(new(*Config),
			"Storage",
		),
	)
)

// ProvideUploadManager transforms a *objectstorage.Uploader into an UploadManager.
func ProvideUploadManager(u *objectstorage.Uploader) UploadManager {
	return u
}
