package model

import (
	"wellnus/backend/router/http_helper/http_error"

	"database/sql"
	"errors"
)

func IsProvider(user User) bool {
	return user.UserRole == "VOLUNTEER" || user.UserRole == "COUNSELLOR";
}

func ReadProviderSettings(rows *sql.Rows) ([]ProviderSetting, error) {
	providerSettings := make([]ProviderSetting, 0)
	for rows.Next() {
		var providerSetting ProviderSetting
		if err := rows.Scan(
			&providerSetting.UserID,
			&providerSetting.Available,
			&providerSetting.Specialities);
			err != nil {
				return nil, err
			}
		providerSettings = append(providerSettings, providerSetting)
	}
	return providerSettings, nil
}

func ReadProvidersWithSetting(rows *sql.Rows) ([]ProviderWithSetting, error) {
	providersWithSetting := make([]ProviderWithSetting, 0)
	for rows.Next() {
		var providerWithSetting ProviderWithSetting
		if err := rows.Scan(
			&providerWithSetting.User.ID,
			&providerWithSetting.User.FirstName, 
			&providerWithSetting.User.LastName, 
			&providerWithSetting.User.Gender, 
			&providerWithSetting.User.Faculty, 
			&providerWithSetting.User.Email,
			&providerWithSetting.User.UserRole, 
			&providerWithSetting.User.PasswordHash,
			&providerWithSetting.Setting.UserID,
			&providerWithSetting.Setting.Available,
			&providerWithSetting.Setting.Specialities); 
			err != nil {
				return nil, err
			}
		providersWithSetting = append(providersWithSetting, providerWithSetting)
	}
	return providersWithSetting, nil
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

func GetAllProvidersWithSetting(db *sql.DB) ([]ProviderWithSetting, error) {
	rows, err := db.Query(
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
			wn_provider_setting.available,
			wn_provider_setting.specialities
		FROM wn_provider_setting 
		JOIN wn_user ON wn_user.id = wn_provider_setting.user_id
		WHERE wn_user.user_role IN ('VOLUNTEER', 'COUNSELLOR')`)
	if err != nil { return nil, err }
	defer rows.Close()
	providersWithSetting, err := ReadProvidersWithSetting(rows)
	if err != nil { return nil, err }
	return providersWithSetting, nil
}

func GetProviderWithSetting(db *sql.DB, userID int64) (ProviderWithSetting, error) {
	providerSetting, err := GetProviderSetting(db, userID)
	if err != nil { return ProviderWithSetting{}, err }
	providerWithSetting, err := providerSetting.LoadProviderSettings(db)
	if err != nil { return ProviderWithSetting{}, err }
	return providerWithSetting, nil
}

func AddUpdateProviderSettingOfProvider(db *sql.DB, providerSetting ProviderSetting, userID int64) (ProviderSetting, error) {
	providerSetting.UserID = userID
	user, err := GetUser(db, userID)
	if err != nil { return ProviderSetting{}, err }
	if !IsProvider(user) { return ProviderSetting{}, errors.New("Cannot addupdate provider settings of non-provider") }
	_, err = db.Exec(
		`INSERT INTO wn_provider_setting (
			user_id,
			available,
			specialities
		) VALUES ($1, $2, $3)
		ON CONFLICT (user_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			available = EXCLUDED.available,
			specialities = EXCLUDED.specialities`,
		providerSetting.UserID,
		providerSetting.Available,
		providerSetting.Specialities)
	if err != nil { return ProviderSetting{}, err }
	return providerSetting, nil
}

func DeleteProviderSettingOfProvider(db *sql.DB, userID int64) (ProviderSetting, error) {
	_, err := db.Exec("DELETE FROM wn_provider_setting WHERE user_id = $1", userID)
	if err != nil { return ProviderSetting{}, err }
	return ProviderSetting{ UserID: userID }, nil
}



