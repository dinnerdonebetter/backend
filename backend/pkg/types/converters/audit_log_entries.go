package converters

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertAuditLogEntryToAuditLogEntryDatabaseCreationInput builds a AuditLogEntryDatabaseCreationInput from a AuditLogEntry.
func ConvertAuditLogEntryToAuditLogEntryDatabaseCreationInput(x *types.AuditLogEntry) *types.AuditLogEntryDatabaseCreationInput {
	v := &types.AuditLogEntryDatabaseCreationInput{
		Changes:            x.Changes,
		BelongsToHousehold: x.BelongsToHousehold,
		ID:                 x.ID,
		ResourceType:       x.ResourceType,
		RelevantID:         x.RelevantID,
		EventType:          x.EventType,
		BelongsToUser:      x.BelongsToUser,
	}

	return v
}
