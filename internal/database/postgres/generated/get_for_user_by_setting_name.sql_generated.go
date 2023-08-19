// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_for_user_by_setting_name.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const getServiceSettingConfigurationsForUserBySettingName = `-- name: GetServiceSettingConfigurationsForUserBySettingName :many

SELECT
	service_setting_configurations.id,
    service_setting_configurations.value,
    service_setting_configurations.notes,
	service_settings.id,
    service_settings.name,
    service_settings.type,
    service_settings.description,
    service_settings.default_value,
    service_settings.enumeration,
    service_settings.admins_only,
    service_settings.created_at,
    service_settings.last_updated_at,
    service_settings.archived_at,
    service_setting_configurations.belongs_to_user,
    service_setting_configurations.belongs_to_household,
    service_setting_configurations.created_at,
    service_setting_configurations.last_updated_at,
    service_setting_configurations.archived_at
FROM service_setting_configurations
 JOIN service_settings ON service_setting_configurations.service_setting_id=service_settings.id
WHERE service_settings.archived_at IS NULL
  AND service_setting_configurations.archived_at IS NULL
  AND service_settings.name = $1
  AND service_setting_configurations.belongs_to_user = $2
`

type GetServiceSettingConfigurationsForUserBySettingNameParams struct {
	Name          string
	BelongsToUser string
}

type GetServiceSettingConfigurationsForUserBySettingNameRow struct {
	ID                 string
	Value              string
	Notes              string
	ID_2               string
	Name               string
	Type               SettingType
	Description        string
	DefaultValue       sql.NullString
	Enumeration        string
	AdminsOnly         bool
	CreatedAt          time.Time
	LastUpdatedAt      sql.NullTime
	ArchivedAt         sql.NullTime
	BelongsToUser      string
	BelongsToHousehold string
	CreatedAt_2        time.Time
	LastUpdatedAt_2    sql.NullTime
	ArchivedAt_2       sql.NullTime
}

func (q *Queries) GetServiceSettingConfigurationsForUserBySettingName(ctx context.Context, db DBTX, arg *GetServiceSettingConfigurationsForUserBySettingNameParams) ([]*GetServiceSettingConfigurationsForUserBySettingNameRow, error) {
	rows, err := db.QueryContext(ctx, getServiceSettingConfigurationsForUserBySettingName, arg.Name, arg.BelongsToUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetServiceSettingConfigurationsForUserBySettingNameRow{}
	for rows.Next() {
		var i GetServiceSettingConfigurationsForUserBySettingNameRow
		if err := rows.Scan(
			&i.ID,
			&i.Value,
			&i.Notes,
			&i.ID_2,
			&i.Name,
			&i.Type,
			&i.Description,
			&i.DefaultValue,
			&i.Enumeration,
			&i.AdminsOnly,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.BelongsToUser,
			&i.BelongsToHousehold,
			&i.CreatedAt_2,
			&i.LastUpdatedAt_2,
			&i.ArchivedAt_2,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
