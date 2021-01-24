package data

import (
	"math"
	"time"
)

// Credito describes
type Credito struct {
	FechaInicio         time.Time `json:"fechaInicio" validate:"required" schema:"date"`
	TotalCapital        int       `json:"totalCapital" validate:"required"`
	Descripcion         string    `json:"descripcion"`
	Tiempo              int       `json:"tiempo" validate:"required"`
	PorcentajeIntereses float64   `json:"porcentajeIntereses"`
	IDUsuario           int       `json:"idUsuario" validate:"required"`

	ValorCuota        int     `json:"valorCuota"`
	ValorTotalCredito int     `json:"valorTotal"`
	ID                int     `json:"id"`
	TotalIntereses    int     `json:"totalIntereses"`
	CapitalPagado     int     `json:"capitalPagado"`
	InteresPagado     int     `json:"interesPagado"`
	TotalPagado       int     `json:"totalPagado"`
	PorcentajePagado  float64 `json:"porcentajePagado"`
	DebeTotal         int     `json:"debeTotal"`
	DebeCapital       int     `json:"debeCapital"`
	DebeInteres       int     `json:"debeInteres"`
}

// CreditoExistente resumee of credit
type CreditoExistente struct {
	ValorCuota int
	ID         int
	IDUsuario  int
}

// Cuota describes the data from a single cuota
type Cuota struct {
	Capital   float64   `json:"capital"`
	Intereses float64   `json:"intereses"`
	Cuota     float64   `json:"cuota"`
	Saldo     float64   `json:"saldo"`
	Mes       time.Time `json:"mes"`
}

// Pago describes the payment of aporte or intereses
type Pago struct {
	ValorCapital    int       `json:"valorCapital"`
	ValorIntrereses int       `json:"valorIntereses"`
	Fecha           time.Time `json:"fecha" validate:"required"`
	IDCredito       int       `json:"idCredito" validate:"required"`
}

// Informe describes the data given from a specific report
type Informe struct {
	Message int `json:"capital"`
}

// Cuotas array of cuota
type Cuotas []*Cuota

// Creditos array of credito
type Creditos []*Credito

// CreateCredito makes
func (u *UserService) CreateCredito(id int, cr *Credito) error {
	u.l.Info("[CreateCredito] Creating credito", "credito", cr)
	_, err := u.UserExists(id)

	if err != nil {
		return err
	}

	var (
		valorCuota     float64
		valorIntereses float64
	)

	if cr.PorcentajeIntereses == 0 {
		valorCuota = roundup(cr.TotalCapital / cr.Tiempo)
		valorIntereses = 0

	} else {
		valorCuota = roundup(int(float64(cr.TotalCapital) * cr.PorcentajeIntereses / (1 - math.Pow(cr.PorcentajeIntereses+1, float64(cr.Tiempo*-1)))))
		valorIntereses = float64((cr.Tiempo * int(valorCuota)) - cr.TotalCapital)

	}

	valorTotal := valorIntereses + float64(cr.TotalCapital)
	_, err = u.DB.Exec("INSERT INTO creditos (fechaInicio, descripcion, valorCuota, tiempo, idUsuario, totalIntereses, porcentajeInteres, totalCapital, valorTotalCredito, activo, visible) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		cr.FechaInicio, cr.Descripcion, valorCuota, cr.Tiempo, cr.IDUsuario, valorIntereses, cr.PorcentajeIntereses, cr.TotalCapital, valorTotal, true, true)
	if err == nil {
		return nil
	}
	return err
}

// CalcularCredito calculates the given credit without persist it
func (u *UserService) CalcularCredito(cr *Credito) Cuotas {
	u.l.Info("[CalcularCredito] Calculating quotas of credit", "credito", cr)
	cuotas := Cuotas{}
	totalCapital := float64(cr.TotalCapital)
	var (
		porcentajeInteres float64
		valorCuota        float64
	)

	if cr.PorcentajeIntereses == 0 {
		valorCuota = totalCapital / float64(cr.Tiempo)
		porcentajeInteres = 0
	} else {
		porcentajeInteres = cr.PorcentajeIntereses
		valorCuota = (float64(totalCapital) * porcentajeInteres / (1 - math.Pow(porcentajeInteres+1, float64(cr.Tiempo*-1))))
	}

	return calcularCuotas(totalCapital, cr.Tiempo, porcentajeInteres, valorCuota, 0, cuotas, cr.FechaInicio)
}

// CreatePago creates a payment in the database
func (u *UserService) CreatePago(p *Pago) error {
	u.l.Info("[CreatePago] Creating pago from credit", "aporte", p)
	_, err := u.CreditExists(p.IDCredito)

	if err != nil {
		return err
	}
	_, err = u.DB.Exec("INSERT INTO creditos_cuotas (valor, idCredito, fecha) VALUES ( ?, ?, ?)", p.ValorCapital, p.IDCredito, p.Fecha)
	if err == nil {
		return nil
	}
	return err
}

// CreatePagoInteres creates a interes payment in the database
func (u *UserService) CreatePagoInteres(p *Pago) error {
	u.l.Info("[CreatePagoInteres] Creating pago from credit", "aporte", p)
	_, err := u.CreditExists(p.IDCredito)

	if err != nil {
		return err
	}
	_, err = u.DB.Exec("INSERT INTO creditos_intereses (valor, idCredito, fecha) VALUES ( ?, ?, ?)", p.ValorIntrereses, p.IDCredito, p.Fecha)
	if err == nil {
		return nil
	}
	return err
}

func (u *UserService) darTotalPrestado(cr *Credito) Informe {
	return Informe{}
}

// GetAllCreditos gives all the aportes in the Fondo
func (u *UserService) GetAllCreditos() (Creditos, error) {
	u.l.Info("[GetAllCreditos] Getting all creditos from database")

	creditos := Creditos{}
	rows, err := u.DB.Query(`SELECT fechaInicio, descripcion, valorCuota, tiempo, id, idUsuario, totalIntereses, porcentajeInteres, totalCapital, valorTotalCredito,
				COALESCE(SUM(pagos.valor), 0) as capitalPagado,
				COALESCE(SUM(interes.valor), 0) as interesPagado,
				COALESCE(interes.valor + pagos.valor, 0) as totalPagado,
				COALESCE((interes.valor + pagos.valor) * 100, 0) / valorTotalCredito as porcentajePagado,
				COALESCE(valorTotalCredito - (COALESCE(SUM(pagos.valor), 0) + COALESCE(SUM(interes.valor), 0)), 0) as debeTotal,
				COALESCE(totalCapital - (COALESCE(SUM(pagos.valor), 0)), 0) as debeCapital,
				COALESCE(totalIntereses - COALESCE(SUM(interes.valor), 0), 0) as debeInteres
				FROM creditos 
				LEFT JOIN (SELECT idCredito, sum(valor) as valor FROM creditos_intereses GROUP BY idCredito) as interes ON creditos.id = interes.idCredito
				LEFT JOIN (SELECT idCredito, sum(valor) as valor FROM creditos_cuotas GROUP BY idCredito) as pagos ON creditos.id = pagos.idCredito
				GROUP BY creditos.id`)
	if err != nil {
		return creditos, err
	}

	for rows.Next() {
		credito := &Credito{}
		err = rows.Scan(&credito.FechaInicio, &credito.Descripcion, &credito.ValorCuota, &credito.Tiempo, &credito.ID, &credito.IDUsuario, &credito.TotalIntereses, &credito.PorcentajeIntereses, &credito.TotalCapital, &credito.ValorTotalCredito, &credito.CapitalPagado, &credito.InteresPagado, &credito.TotalPagado, &credito.PorcentajePagado, &credito.DebeTotal, &credito.DebeCapital, &credito.DebeInteres)
		if err != nil {
			return creditos, err
		}

		creditos = append(creditos, credito)
	}

	return creditos, nil
}

// GetAllCreditosByUserID gives all the creditos in the Fondo given an user id
func (u *UserService) GetAllCreditosByUserID(id int) (Creditos, error) {
	u.l.Info("[GetAllCreditos] Getting all creditos from database from user", "user", id)

	creditos := Creditos{}
	rows, err := u.DB.Query(`SELECT fechaInicio, descripcion, valorCuota, tiempo, id, idUsuario, totalIntereses, porcentajeInteres, totalCapital, valorTotalCredito,
				COALESCE(SUM(pagos.valor), 0) as capitalPagado,
				COALESCE(SUM(interes.valor), 0) as interesPagado,
				COALESCE(interes.valor + pagos.valor, 0) as totalPagado,
				COALESCE((interes.valor + pagos.valor) * 100, 0) / valorTotalCredito as porcentajePagado,
				COALESCE(valorTotalCredito - (COALESCE(SUM(pagos.valor), 0) + COALESCE(SUM(interes.valor), 0)), 0) as debeTotal,
				COALESCE(totalCapital - (COALESCE(SUM(pagos.valor), 0)), 0) as debeCapital,
				COALESCE(totalIntereses - COALESCE(SUM(interes.valor), 0), 0) as debeInteres
				FROM creditos 
				LEFT JOIN (SELECT idCredito, sum(valor) as valor FROM creditos_intereses GROUP BY idCredito) as interes ON creditos.id = interes.idCredito
				LEFT JOIN (SELECT idCredito, sum(valor) as valor FROM creditos_cuotas GROUP BY idCredito) as pagos ON creditos.id = pagos.idCredito
				WHERE idUsuario = ?
				GROUP BY creditos.id`, id)
	if err != nil {
		return creditos, err
	}

	for rows.Next() {
		credito := &Credito{}
		err = rows.Scan(&credito.FechaInicio, &credito.Descripcion, &credito.ValorCuota, &credito.Tiempo, &credito.ID, &credito.IDUsuario, &credito.TotalIntereses, &credito.PorcentajeIntereses, &credito.TotalCapital, &credito.ValorTotalCredito, &credito.CapitalPagado, &credito.InteresPagado, &credito.TotalPagado, &credito.PorcentajePagado, &credito.DebeTotal, &credito.DebeCapital, &credito.DebeInteres)
		if err != nil {
			return creditos, err
		}

		creditos = append(creditos, credito)
	}

	return creditos, nil
}

// GetCreditoByID returns a credit given an id
func (u *UserService) GetCreditoByID(id int) (Credito, error) {
	u.l.Info("[GetCreditoByID] Getting credit", "id", id)

	credito := Credito{}
	rows, err := u.DB.Query(`SELECT valorCuota, id, idUsuario FROM fondofamiliar_dev.creditos WHERE id = ?`, id)
	if err != nil {
		return Credito{}, err
	}

	for rows.Next() {
		credito := &Credito{}
		err = rows.Scan(&credito.FechaInicio, &credito.Descripcion, &credito.ValorCuota, &credito.Tiempo, &credito.ID, &credito.IDUsuario, &credito.TotalIntereses, &credito.PorcentajeIntereses, &credito.TotalCapital, &credito.ValorTotalCredito, &credito.CapitalPagado, &credito.InteresPagado, &credito.TotalPagado, &credito.PorcentajePagado)
		if err != nil {
			return Credito{}, err
		}
	}

	return credito, nil
}

func calcularCuotas(pValorTotal float64, pTiempo int, pPorcentajeInteres float64, pValorCuota float64, pNumeroCuota int, pCuotas Cuotas, pFechaInicio time.Time) Cuotas {
	if pNumeroCuota != pTiempo {
		numeroCuota := pNumeroCuota + 1
		valorInteres := roundup(int(pValorTotal * pPorcentajeInteres))
		valorCapital := roundup(int(pValorCuota - valorInteres))
		fechaSiguiente := pFechaInicio.AddDate(0, 1, 0)
		saldo := roundup(int(pValorTotal - valorCapital))
		valorCuota := roundup(int(pValorCuota))

		newCuota := Cuota{valorCapital, valorInteres, valorCuota, saldo, fechaSiguiente}
		pCuotas = append(pCuotas, &newCuota)
		return calcularCuotas(saldo, pTiempo, pPorcentajeInteres, pValorCuota, numeroCuota, pCuotas, fechaSiguiente)
	}
	return pCuotas

}

func roundup(i int) float64 {
	if i%100 == 0 {
		return float64(i)
	}
	return float64(i + 100 - i%100)
}
