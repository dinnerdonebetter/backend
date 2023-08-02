// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_for_household_by_setting_name.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetServiceSettingConfigurationsForHouseholdBySettingName = `-- name: GetServiceSettingConfigurationsForHouseholdBySettingName :many

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
  AND service_setting_configurations.belongs_to_household = $2
`

type GetServiceSettingConfigurationsForHouseholdBySettingNameParams struct {
	Name               string `db:"name"`
	BelongsToHousehold string `db:"belongs_to_household"`
}

type GetServiceSettingConfigurationsForHouseholdBySettingNameRow struct {
	CreatedAt          time.Time      `db:"created_at"`
	CreatedAt_2        time.Time      `db:"created_at_2"`
	ArchivedAt_2       sql.NullTime   `db:"archived_at_2"`
	LastUpdatedAt_2    sql.NullTime   `db:"last_updated_at_2"`
	ArchivedAt         sql.NullTime   `db:"archived_at"`
	LastUpdatedAt      sql.NullTime   `db:"last_updated_at"`
	Name               string         `db:"name"`
	Enumeration        string         `db:"enumeration"`
	Description        string         `db:"description"`
	Type               SettingType    `db:"type"`
	ID                 string         `db:"id"`
	BelongsToUser      string         `db:"belongs_to_user"`
	BelongsToHousehold string         `db:"belongs_to_household"`
	ID_2               string         `db:"id_2"`
	Notes              string         `db:"notes"`
	Value              string         `db:"value"`
	DefaultValue       sql.NullString `db:"default_value"`
	AdminsOnly         bool           `db:"admins_only"`
}

func (q *Queries) GetServiceSettingConfigurationsForHouseholdBySettingName(ctx context.Context, db DBTX, arg *GetServiceSettingConfigurationsForHouseholdBySettingNameParams) ([]*GetServiceSettingConfigurationsForHouseholdBySettingNameRow, error) {
	rows, err := db.QueryContext(ctx, GetServiceSettingConfigurationsForHouseholdBySettingName, arg.Name, arg.BelongsToHousehold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetServiceSettingConfigurationsForHouseholdBySettingNameRow{}
	for rows.Next() {
		var i GetServiceSettingConfigurationsForHouseholdBySettingNameRow
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
