-- name: CreateHousehold :exec
INSERT INTO households (id,name,billing_status,contact_email,contact_phone,time_zone,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6,$7);