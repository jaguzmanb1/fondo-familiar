package data

import (
	"fmt"
)

// ErrValorMayor is raised when a user is not found
var ErrValorMayor = fmt.Errorf("The specified value is greater than the user has")

// ReporteGeneral describes a general report from the fondo
type ReporteGeneral struct {
	Capital           float64 `json:"capital"`
	Intereses         float64 `json:"intereses"`
	Prestado          float64 `json:"prestado"`
	Total             float64 `json:"total"`
	TotalSinIntereses float64 `json:"totalSinIntereses"`
}

// PostDescuento is the response when a discount is applied
type PostDescuento struct {
	ValorDescuento int    `json:"valorDescuento" validate:"required"`
	IDUsuario      int    `json:"idUsuario" validate:"required"`
	IDCredito      int    `json:"idCredito"`
	Antes          Estado `json:"antes"`
	Despues        Estado `json:"despues"`
}

// Estado is the amount of aportes and intereses of an user
type Estado struct {
	Aportes   int `json:"aportes"`
	Intereses int `json:"intereses"`
}

// GetReporteGeneral gives a general report of the status of the fondo
func (u *UserService) GetReporteGeneral() (ReporteGeneral, error) {
	u.l.Info("[GetReportegeneral] Getting reporte general")

	reporte := ReporteGeneral{}
	rows, err := u.DB.Query(`SELECT 
	COALESCE((SELECT SUM(valor) as valor FROM aportes), 0) as capital,
	COALESCE((SELECT SUM(valor) from creditos_intereses), 0) as intereses,
	COALESCE((SELECT SUM(totalCapital) from creditos), 0) - COALESCE((SELECT SUM(valor) from creditos_cuotas), 0) as prestado,
	COALESCE(((SELECT SUM(valor) as valor FROM aportes) + COALESCE((SELECT SUM(valor) from creditos_intereses), 0) + COALESCE((SELECT SUM(valor) from creditos_cuotas), 0) - (SELECT SUM(totalCapital) from creditos)), 0) as total,
	COALESCE(((SELECT SUM(valor) as valor FROM aportes) - (SELECT SUM(totalCapital) from creditos)), 0) + COALESCE((SELECT SUM(valor) from creditos_cuotas), 0) as totalSinIntereses`)
	if err != nil {
		return reporte, err
	}

	for rows.Next() {
		err = rows.Scan(&reporte.Capital, &reporte.Intereses, &reporte.Prestado, &reporte.Total, &reporte.TotalSinIntereses)
		if err != nil {
			return reporte, err
		}
	}

	return reporte, nil
}

// PostDescontarParaCreditoCapital discounts money on aportes to pay it to a credit from a given user
func (u *UserService) PostDescontarParaCreditoCapital(pDescuento *PostDescuento) (PostDescuento, error) {
	_, err := u.CreditExists(pDescuento.IDCredito)
	ap, err := u.GetSumAportesByID(pDescuento.IDUsuario)
	if err != nil {
		return PostDescuento{}, err
	}
	if ap.Valor >= int(pDescuento.ValorDescuento) {
		aps, err := u.GetAllAportesByID(pDescuento.IDUsuario, "", "")

		if err != nil {
			return PostDescuento{}, err
		}

		antes := Estado{Aportes: ap.Valor}
		tmp := 0
		for _, a := range aps {
			if a.Valor != 0 {
				var err error
				if (tmp + a.Valor) > pDescuento.ValorDescuento {
					sum := pDescuento.ValorDescuento - tmp
					tmp += sum
					_, err = u.DB.Query("UPDATE aportes SET valor = ? WHERE id = ?", a.Valor-sum, a.ID)

				} else {
					tmp += a.Valor
					_, err = u.DB.Query("UPDATE aportes SET valor = ? WHERE id = ?", 0, a.ID)
				}

				if err != nil {
					return PostDescuento{}, err
				}

				if tmp == pDescuento.ValorDescuento {
					pago := Pago{ValorCapital: pDescuento.ValorDescuento, IDCredito: pDescuento.IDCredito}
					ap, err := u.GetSumAportesByID(pDescuento.IDUsuario)
					err = u.CreatePago(&pago)
					despues := Estado{Aportes: ap.Valor}

					if err != nil {
						return PostDescuento{}, err
					}

					return PostDescuento{Antes: antes, Despues: despues, ValorDescuento: pDescuento.ValorDescuento}, nil
				}
			}
		}

	}
	return PostDescuento{}, ErrValorMayor
}

// PostDescontarParaCreditoIntereses discounts money on aportes to pay it to a credit from a given user
func (u *UserService) PostDescontarParaCreditoIntereses(pDescuento *PostDescuento) (PostDescuento, error) {
	_, err := u.CreditExists(pDescuento.IDCredito)
	ap, err := u.GetSumAportesByID(pDescuento.IDUsuario)
	if err != nil {
		return PostDescuento{}, err
	}
	if ap.Valor >= int(pDescuento.ValorDescuento) {
		aps, err := u.GetAllAportesByID(pDescuento.IDUsuario, "", "")

		if err != nil {
			return PostDescuento{}, err
		}

		antes := Estado{Aportes: ap.Valor}
		tmp := 0
		for _, a := range aps {
			if a.Valor != 0 {
				var err error
				if (tmp + a.Valor) > pDescuento.ValorDescuento {
					sum := pDescuento.ValorDescuento - tmp
					tmp += sum
					_, err = u.DB.Query("UPDATE aportes SET valor = ? WHERE id = ?", a.Valor-sum, a.ID)

				} else {
					tmp += a.Valor
					_, err = u.DB.Query("UPDATE aportes SET valor = ? WHERE id = ?", 0, a.ID)
				}

				if err != nil {
					return PostDescuento{}, err
				}

				if tmp == pDescuento.ValorDescuento {
					pago := Pago{ValorIntrereses: pDescuento.ValorDescuento, IDCredito: pDescuento.IDCredito}
					ap, err := u.GetSumAportesByID(pDescuento.IDUsuario)
					err = u.CreatePagoInteres(&pago)
					despues := Estado{Aportes: ap.Valor}

					if err != nil {
						return PostDescuento{}, err
					}

					return PostDescuento{Antes: antes, Despues: despues, ValorDescuento: pDescuento.ValorDescuento}, nil
				}
			}
		}

	}
	return PostDescuento{}, ErrValorMayor
}

// PostDescontar discounts money on aportes
func (u *UserService) PostDescontar(pDescuento *PostDescuento) (PostDescuento, error) {
	ap, err := u.GetSumAportesByID(pDescuento.IDUsuario)
	if err != nil {
		return PostDescuento{}, err
	}
	if ap.Valor >= int(pDescuento.ValorDescuento) {
		aps, err := u.GetAllAportesByID(pDescuento.IDUsuario, "", "")

		if err != nil {
			return PostDescuento{}, err
		}

		antes := Estado{Aportes: ap.Valor}
		tmp := 0
		for _, a := range aps {
			if a.Valor != 0 {
				var err error
				if (tmp + a.Valor) > pDescuento.ValorDescuento {
					sum := pDescuento.ValorDescuento - tmp
					tmp += sum
					_, err = u.DB.Query("UPDATE aportes SET valor = ? WHERE id = ?", a.Valor-sum, a.ID)

				} else {
					tmp += a.Valor
					_, err = u.DB.Query("UPDATE aportes SET valor = ? WHERE id = ?", 0, a.ID)
				}

				if err != nil {
					return PostDescuento{}, err
				}

				if tmp == pDescuento.ValorDescuento {
					despues := Estado{Aportes: ap.Valor}

					if err != nil {
						return PostDescuento{}, err
					}

					return PostDescuento{Antes: antes, Despues: despues, ValorDescuento: pDescuento.ValorDescuento}, nil
				}
			}
		}

	}
	return PostDescuento{}, ErrValorMayor
}
