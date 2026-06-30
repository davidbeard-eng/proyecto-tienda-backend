package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

// --- ESTRUCTURAS PARA LOS 10 REPORTES ---
type R1 struct {
	Mes           int    `json:"mes"`
	Producto      string `json:"producto"`
	TotalUnidades int    `json:"total_unidades"`
}
type R2 struct {
	Tienda   string  `json:"tienda"`
	Producto string  `json:"producto"`
	Precio   float64 `json:"precio"`
}
type R3 struct {
	anio                  int     `json:"anio"`
	Mes                   int     `json:"mes"`
	TipoDocumento         string  `json:"tipo_documento"`
	CantidadTransacciones int     `json:"cantidad_transacciones"`
	MontoTotalRecaudado   float64 `json:"monto_total_recaudado"`
}
type R4 struct {
	Tienda           string  `json:"tienda"`
	Empleado         string  `json:"empleado"`
	Cargo            string  `json:"cargo"`
	ComunaResidencia string  `json:"comuna_residencia"`
	TotalGanadoAno   float64 `json:"total_ganado_ano"`
}
type R5 struct {
	Tienda            string `json:"tienda"`
	CantidadEmpleados int    `json:"cantidad_empleados"`
}
type R6 struct {
	anio           int    `json:"anio"`
	Mes            int    `json:"mes"`
	Vendedor       string `json:"vendedor"`
	CantidadVentas int    `json:"cantidad_ventas"`
}
type R7 struct {
	anio                 int     `json:"anio"`
	Tienda               string  `json:"tienda"`
	Vendedor             string  `json:"vendedor"`
	TotalDineroRecaudado float64 `json:"total_dinero_recaudado"`
}
type R8 struct {
	Tienda                string `json:"tienda"`
	Vendedor              string `json:"vendedor"`
	TotalUnidadesVendidas int    `json:"total_unidades_vendidas"`
}
type R9 struct {
	anio         int     `json:"anio"`
	Mes          int     `json:"mes"`
	Empleado     string  `json:"empleado"`
	SueldoPagado float64 `json:"sueldo_pagado"`
}
type R10 struct {
	anio            int     `json:"anio"`
	Mes             int     `json:"mes"`
	Tienda          string  `json:"tienda"`
	PeorRecaudacion float64 `json:"peor_recaudacion"`
}

func main() {
	// CONFIGURA AQUÍ TU CONTRASEÑA REAL DE POSTGRES
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=tiendas_comerciales_db sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	// Configuración de CORS para permitir que React se conecte
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	})

	// --- RUTAS DE LA API ---
	r.GET("/api/reportes/r1", getR1)
	r.GET("/api/reportes/r2", getR2)
	r.GET("/api/reportes/r3", getR3)
	r.GET("/api/reportes/r4", getR4)
	r.GET("/api/reportes/r5", getR5)
	r.GET("/api/reportes/r6", getR6)
	r.GET("/api/reportes/r7", getR7)
	r.GET("/api/reportes/r8", getR8)
	r.GET("/api/reportes/r9", getR9)
	r.GET("/api/reportes/r10", getR10)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}

// 1. Producto más vendido por mes el 2021
func getR1(c *gin.Context) {
	query := `WITH CantidadesMes AS (SELECT EXTRACT(MONTH FROM v.fecha) AS mes, p.nombre AS producto, SUM(pv.cantidad) AS total_unidades, RANK() OVER (PARTITION BY EXTRACT(MONTH FROM v.fecha) ORDER BY SUM(pv.cantidad) DESC) AS ranking FROM VENTA v JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta JOIN PRODUCTO p ON pv.id_producto = p.id_producto WHERE EXTRACT(YEAR FROM v.fecha) = 2021 GROUP BY EXTRACT(MONTH FROM v.fecha), p.nombre) SELECT CAST(mes AS INT), producto, CAST(total_unidades AS INT) FROM CantidadesMes WHERE ranking = 1 ORDER BY mes;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var res []R1
	for rows.Next() {
		var r R1
		rows.Scan(&r.Mes, &r.Producto, &r.TotalUnidades)
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}

// 2. Producto más económico por tienda
func getR2(c *gin.Context) {
	query := `WITH PreciosPorTienda AS (SELECT DISTINCT t.nombre AS tienda, p.nombre AS producto, p.precio, RANK() OVER (PARTITION BY t.id_tienda ORDER BY p.precio ASC) AS ranking FROM VENTA v JOIN TIENDA t ON v.id_tienda = t.id_tienda JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta JOIN PRODUCTO p ON pv.id_producto = p.id_producto) SELECT tienda, producto, CAST(precio AS DOUBLE PRECISION) FROM PreciosPorTienda WHERE ranking = 1;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var res []R2
	for rows.Next() {
		var r R2
		rows.Scan(&r.Tienda, &r.Producto, &r.Precio)
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}

// 3. Ventas por mes, separadas entre Boletas y Facturas
func getR3(c *gin.Context) {
	query := `SELECT CAST(EXTRACT(YEAR FROM v.fecha) AS INT) AS anio, CAST(EXTRACT(MONTH FROM v.fecha) AS INT) AS mes, td.nombre AS tipo_documento, CAST(COUNT(DISTINCT v.id_venta) AS INT) AS cantidad_transacciones, CAST(SUM(pv.cantidad * p.precio) AS DOUBLE PRECISION) AS monto_total_recaudado FROM VENTA v JOIN TIPO_DOC td ON v.id_tipo_doc = td.id_tipo_doc JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta JOIN PRODUCTO p ON pv.id_producto = p.id_producto GROUP BY EXTRACT(YEAR FROM v.fecha), EXTRACT(MONTH FROM v.fecha), td.nombre ORDER BY anio, mes, tipo_documento;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var res []R3
	for rows.Next() {
		var r R3
		rows.Scan(&r.anio, &r.Mes, &r.TipoDocumento, &r.CantidadTransacciones, &r.MontoTotalRecaudado)
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}

// 4. Empleado que ganó más por tienda en 2020
func getR4(c *gin.Context) {
	query := `WITH Sueldos2020 AS (SELECT t.nombre AS tienda, e.nombre AS empleado, e.cargo, c.nombre AS comuna_residencia, SUM(s.monto) AS total_ganado_ano, RANK() OVER (PARTITION BY t.id_tienda ORDER BY SUM(s.monto) DESC) AS ranking FROM EMPLEADO e JOIN SUELDO s ON e.id_empleado = s.id_empleado JOIN TIENDA_EMP te ON e.id_empleado = te.id_empleado JOIN TIENDA t ON te.id_tienda = t.id_tienda JOIN COMUNA c ON e.id_comuna = c.id_comuna WHERE s.anio = 2020 GROUP BY t.id_tienda, t.nombre, e.id_empleado, e.nombre, e.cargo, c.nombre) SELECT tienda, empleado, cargo, comuna_residencia, CAST(total_ganado_ano AS DOUBLE PRECISION) FROM Sueldos2020 WHERE ranking = 1;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var res []R4
	for rows.Next() {
		var r R4
		rows.Scan(&r.Tienda, &r.Empleado, &r.Cargo, &r.ComunaResidencia, &r.TotalGanadoAno)
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}

// 5. La tienda que tiene menos empleados
func getR5(c *gin.Context) {
	query := `SELECT t.nombre AS tienda, CAST(COUNT(te.id_empleado) AS INT) AS cantidad_empleados FROM TIENDA t LEFT JOIN TIENDA_EMP te ON t.id_tienda = te.id_tienda GROUP BY t.id_tienda, t.nombre ORDER BY cantidad_empleados ASC LIMIT 1;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var res []R5
	for rows.Next() {
		var r R5
		rows.Scan(&r.Tienda, &r.CantidadEmpleados)
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}

// 6. El vendedor con más ventas por mes
func getR6(c *gin.Context) {
	query := `WITH VentasVendedor AS (SELECT CAST(EXTRACT(YEAR FROM v.fecha) AS INT) AS anio, CAST(EXTRACT(MONTH FROM v.fecha) AS INT) AS mes, emp.nombre AS vendedor, CAST(COUNT(v.id_venta) AS INT) AS cantidad_ventas, RANK() OVER (PARTITION BY EXTRACT(YEAR FROM v.fecha), EXTRACT(MONTH FROM v.fecha) ORDER BY COUNT(v.id_venta) DESC) AS ranking FROM VENTA v JOIN VENDEDOR vend ON v.id_vendedor = vend.id_vendedor JOIN EMPLEADO emp ON vend.id_empleado = emp.id_empleado GROUP BY EXTRACT(YEAR FROM v.fecha), EXTRACT(MONTH FROM v.fecha), emp.nombre) SELECT anio, mes, vendedor, cantidad_ventas FROM VentasVendedor WHERE ranking = 1 ORDER BY anio, mes;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var res []R6
	for rows.Next() {
		var r R6
		rows.Scan(&r.anio, &r.Mes, &r.Vendedor, &r.CantidadVentas)
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}

// 7. El vendedor que ha recaudado más dinero para la tienda por anio
func getR7(c *gin.Context) {
	query := `WITH RecaudacionAnual AS (SELECT CAST(EXTRACT(YEAR FROM v.fecha) AS INT) AS anio, t.nombre AS tienda, emp.nombre AS vendedor, SUM(pv.cantidad * p.precio) AS total_dinero_recaudado, RANK() OVER (PARTITION BY EXTRACT(YEAR FROM v.fecha), t.id_tienda ORDER BY SUM(pv.cantidad * p.precio) DESC) AS ranking FROM VENTA v JOIN TIENDA t ON v.id_tienda = t.id_tienda JOIN VENDEDOR vend ON v.id_vendedor = vend.id_vendedor JOIN EMPLEADO emp ON vend.id_empleado = emp.id_empleado JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta JOIN PRODUCTO p ON pv.id_producto = p.id_producto GROUP BY EXTRACT(YEAR FROM v.fecha), t.id_tienda, t.nombre, emp.nombre) SELECT anio, tienda, vendedor, CAST(total_dinero_recaudado AS DOUBLE PRECISION) FROM RecaudacionAnual WHERE ranking = 1 ORDER BY anio, tienda;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var res []R7
	for rows.Next() {
		var r R7
		rows.Scan(&r.anio, &r.Tienda, &r.Vendedor, &r.TotalDineroRecaudado)
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}

// 8. El vendedor con más productos vendidos por tienda
func getR8(c *gin.Context) {
	query := `WITH ProductosPorTienda AS (SELECT t.nombre AS tienda, emp.nombre AS vendedor, CAST(SUM(pv.cantidad) AS INT) AS total_unidades_vendidas, RANK() OVER (PARTITION BY t.id_tienda ORDER BY SUM(pv.cantidad) DESC) AS ranking FROM VENTA v JOIN TIENDA t ON v.id_tienda = t.id_tienda JOIN VENDEDOR vend ON v.id_vendedor = vend.id_vendedor JOIN EMPLEADO emp ON vend.id_empleado = emp.id_empleado JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta GROUP BY t.id_tienda, t.nombre, emp.nombre) SELECT tienda, vendedor, total_unidades_vendidas FROM ProductosPorTienda WHERE ranking = 1;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var res []R8
	for rows.Next() {
		var r R8
		rows.Scan(&r.Tienda, &r.Vendedor, &r.TotalUnidadesVendidas)
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}

// 9. El empleado con mayor sueldo por mes
func getR9(c *gin.Context) {
	query := `WITH MaxSueldoMes AS (SELECT s.anio, s.mes, e.nombre AS empleado, s.monto AS sueldo_pagado, RANK() OVER (PARTITION BY s.anio, s.mes ORDER BY s.monto DESC) AS ranking FROM SUELDO s JOIN EMPLEADO e ON s.id_empleado = e.id_empleado) SELECT anio, mes, empleado, CAST(sueldo_pagado AS DOUBLE PRECISION) FROM MaxSueldoMes WHERE ranking = 1 ORDER BY anio, mes;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var res []R9
	for rows.Next() {
		var r R9
		rows.Scan(&r.anio, &r.Mes, &r.Empleado, &r.SueldoPagado)
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}

// 10. La tienda con menor recaudación por mes
func getR10(c *gin.Context) {
	query := `WITH RecaudacionTiendaMes AS (SELECT CAST(EXTRACT(YEAR FROM v.fecha) AS INT) AS anio, CAST(EXTRACT(MONTH FROM v.fecha) AS INT) AS mes, t.nombre AS tienda, SUM(pv.cantidad * p.precio) AS total_mes, RANK() OVER (PARTITION BY EXTRACT(YEAR FROM v.fecha), EXTRACT(MONTH FROM v.fecha) ORDER BY SUM(pv.cantidad * p.precio) ASC) AS ranking FROM VENTA v JOIN TIENDA t ON v.id_tienda = t.id_tienda JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta JOIN PRODUCTO p ON pv.id_producto = p.id_producto GROUP BY EXTRACT(YEAR FROM v.fecha), EXTRACT(MONTH FROM v.fecha), t.nombre) SELECT anio, mes, tienda, CAST(total_mes AS DOUBLE PRECISION) AS peor_recaudacion FROM RecaudacionTiendaMes WHERE ranking = 1 ORDER BY anio, mes;`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var res []R10
	for rows.Next() {
		var r R10
		rows.Scan(&r.anio, &r.Mes, &r.Tienda, &r.PeorRecaudacion)
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}
