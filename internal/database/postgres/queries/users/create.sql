-- name: CreateUser :exec

INSERT INTO users (id,first_name,last_name,username,email_address,hashed_password,two_factor_secret,avatar_src,user_account_status,birthday,service_role,email_address_verification_token) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);
