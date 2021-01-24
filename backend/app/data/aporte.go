package data

import "database/sql"

// Aporte describes
type Aporte struct {
	Valor     int    `json:"valor" validate:"required"`
	Fecha     string `json:"fecha" validate:"required"`
	IDUsuario int    `json:"idUsuario" validate:"required"`
	ID        int    `json:"id"`
}

// SumAportes describes the sum of various aportes
type SumAportes struct {
	Valor int `json:"valor" validate:"required"`
}

// Aportes array of aportes
type Aportes []*Aporte

// CreateAporte makes
func (u *UserService) CreateAporte(id int, ap *Aporte) error {
	u.l.Info("[CreateAporte] Creating aporte", "aporte", ap)
	_, err := u.UserExists(id)

	if err != nil {
		return err
	}
	_, err = u.DB.Exec("INSERT INTO aportes (valor, idUsuario, fecha) VALUES ( ?, ?, ?)", ap.Valor, id, ap.Fecha)
	if err == nil {
		return nil
	}
	return err
}

// GetAllAportes gives all the aportes in the Fondo
func (u *UserService) GetAllAportes() (Aportes, error) {
	u.l.Info("[GetAllAportes] Getting all aportes from database")

	aportes := Aportes{}
	rows, err := u.DB.Query("SELECT valor, idUsuario, fecha FROM aportes")
	if err != nil {
		return aportes, err
	}

	for rows.Next() {
		aporte := &Aporte{}
		err = rows.Scan(&aporte.Valor, &aporte.IDUsuario, &aporte.Fecha)
		if err != nil {
			return aportes, err
		}

		aportes = append(aportes, aporte)
	}

	return aportes, nil
}

// GetAllAportesByID gives all the aportes in the Fondo for a specific user
func (u *UserService) GetAllAportesByID(id int, startDate string, endDate string) (Aportes, error) {
	u.l.Info("[GetAllAportesByID] Getting aportes from id", "userID", id)
	aportes := Aportes{}
	var (
		rows *sql.Rows
		err  error
	)

	if startDate != "" && endDate != "" {
		rows, err = u.DB.Query("SELECT valor, idUsuario, fecha, id FROM aportes where idUsuario = ? AND fecha BETWEEN ? AND ? ", id, startDate, endDate)
	} else if startDate != "" {
		rows, err = u.DB.Query("SELECT valor, idUsuario, fecha, id FROM aportes where idUsuario = ? AND fecha >= ?", id, startDate)
	} else if endDate != "" {
		rows, err = u.DB.Query("SELECT valor, idUsuario, fecha, id FROM aportes where idUsuario = ? AND fecha <= ?", id, endDate)
	} else {
		rows, err = u.DB.Query("SELECT valor, idUsuario, fecha, id FROM aportes where idUsuario = ?", id)
	}

	if err != nil {
		return aportes, err
	}

	for rows.Next() {
		aporte := &Aporte{}
		err = rows.Scan(&aporte.Valor, &aporte.IDUsuario, &aporte.Fecha, &aporte.ID)
		if err != nil {
			return aportes, err
		}

		aportes = append(aportes, aporte)
	}

	return aportes, nil
}

// GetSumAportesByID gives the sum of aportes of an specific user
func (u *UserService) GetSumAportesByID(id int) (SumAportes, error) {
	u.l.Info("[GetSumAportesByID] Getting aportes from id", "userID", id)

	sumAporte := SumAportes{}
	_, err := u.UserExists(id)

	if err != nil {
		return sumAporte, err
	}

	rows, err := u.DB.Query("SELECT COALESCE(SUM(valor), 0) FROM aportes where idUsuario = ?", id)
	if err != nil {
		return sumAporte, err
	}

	for rows.Next() {
		err = rows.Scan(&sumAporte.Valor)
		if err != nil {
			return sumAporte, err
		}
	}

	return sumAporte, nil
}
