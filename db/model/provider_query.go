package model

import (
	"wellnus/backend/router/http_helper/http_error"

	"database/sql"

	"github.com/lib/pq"
)

func IsProvider(user User) bool {
	return user.UserRole == "VOLUNTEER" || user.UserRole == "COUNSELLOR";
}

func AuthoriseProvider(db *sql.DB, userID int64) bool {
	user, _ := GetUser(db, userID)
	return IsProvider(user)
}

func ReadProviderSettings(rows *sql.Rows) ([]ProviderSetting, error) {
	providerSettings := make([]ProviderSetting, 0)
	for rows.Next() {
		var providerSetting ProviderSetting
		if err := rows.Scan(
			&providerSetting.UserID,
			&providerSetting.Intro,
			pq.Array(&providerSetting.Topics));
			err != nil {
				return nil, err
			}
		providerSettings = append(providerSettings, providerSetting)
	}
	return providerSettings, nil
}

func ReadProviders(rows *sql.Rows) ([]Provider, error) {
	providers := make([]Provider, 0)
	for rows.Next() {
		var provider Provider
		if err := rows.Scan(
			&provider.User.ID,
			&provider.User.FirstName, 
			&provider.User.LastName, 
			&provider.User.Gender, 
			&provider.User.Faculty, 
			&provider.User.Email,
			&provider.User.UserRole, 
			&provider.User.PasswordHash,
			&provider.Setting.UserID,
			&provider.Setting.Intro,
			pq.Array(&provider.Setting.Topics)); 
			err != nil {
				return nil, err
			}
		providers = append(providers, provider)
	}
	return providers, nil
}

func GetProviderSetting(db *sql.DB, userID int64) (ProviderSetting, error) {
	rows, err := db.Query("SELECT * FROM wn_provider_setting WHERE user_id = $1", userID)
	if err != nil { return ProviderSetting{}, err }
	defer rows.Close()
	providerSettings, err := ReadProviderSettings(rows)
	if err != nil { return ProviderSetting{}, err }
	if len(providerSettings) == 0 { return ProviderSetting{}, http_error.NotFoundError }
	return providerSettings[0], nil
}

func GetAllProviders(db *sql.DB, topics []string) ([]Provider, error) {
	var rows *sql.Rows
	var err error
	if topics == nil {
		rows, err = db.Query(
		`SELECT 
			wn_user.id,
			wn_user.first_name,
			wn_user.last_name,
			wn_user.gender,
			wn_user.faculty,
			wn_user.email,
			wn_user.user_role,
			wn_user.password_hash,
			wn_provider_setting.user_id,
			wn_provider_setting.intro,
			wn_provider_setting.topics
		FROM wn_provider_setting 
		JOIN wn_user ON wn_user.id = wn_provider_setting.user_id
		WHERE wn_user.user_role IN ('VOLUNTEER', 'COUNSELLOR')`)
	} else {
		rows, err = db.Query(
			`SELECT 
				wn_user.id,
				wn_user.first_name,
				wn_user.last_name,
				wn_user.gender,
				wn_user.faculty,
				wn_user.email,
				wn_user.user_role,
				wn_user.password_hash,
				wn_provider_setting.user_id,
				wn_provider_setting.intro,
				wn_provider_setting.topics
			FROM wn_provider_setting 
			JOIN wn_user ON wn_user.id = wn_provider_setting.user_id
			WHERE wn_user.user_role IN ('VOLUNTEER', 'COUNSELLOR') 
			AND $1 <@ wn_provider_setting.topics`,
			pq.Array(topics))
	}
	if err != nil { return nil, err }
	defer rows.Close()
	providers, err := ReadProviders(rows)
	if err != nil { return nil, err }
	return providers, nil
}

func GetProvider(db *sql.DB, userID int64) (Provider, error) {
	providerSetting, err := GetProviderSetting(db, userID)
	if err != nil { return Provider{}, err }
	provider, err := providerSetting.LoadProviderSetting(db)
	if err != nil { return Provider{}, err }
	return provider, nil
}

func GetProviderWithEvents(db *sql.DB, userID int64) (ProviderWithEvents, error) {
	provider, err := GetProvider(db, userID)
	if err != nil { return ProviderWithEvents{}, err }
	providerWithEvents, err := provider.LoadProvider(db)
	if err != nil { return ProviderWithEvents{}, err }
	return providerWithEvents, nil
}

func AddUpdateProviderSettingOfUser(db *sql.DB, providerSetting ProviderSetting, userID int64) (ProviderSetting, error) {
	providerSetting.UserID = userID
	if !AuthoriseProvider(db, userID) {
		return ProviderSetting{}, http_error.UnauthorizedError
	}
	_, err := db.Exec(
		`INSERT INTO wn_provider_setting (
			user_id,
			intro,
			topics
		) VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			intro = EXCLUDED.intro,
			topics = EXCLUDED.topics`,
		providerSetting.UserID,
		providerSetting.Intro,
		pq.Array(providerSetting.Topics))
	if err != nil { return ProviderSetting{}, err }
	return providerSetting, nil
}

func DeleteProviderSettingOfUser(db *sql.DB, userID int64) (ProviderSetting, error) {
	_, err := db.Exec("DELETE FROM wn_provider_setting WHERE user_id = $1", userID)
	if err != nil { return ProviderSetting{}, err }
	return ProviderSetting{ UserID: userID }, nil
}



