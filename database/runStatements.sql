-- =======================================================================
-- 1. Producto más vendido por mes el 2021.
-- =======================================================================
WITH CantidadesMes AS (
    SELECT 
        EXTRACT(MONTH FROM v.fecha) AS mes,
        p.nombre AS producto,
        SUM(pv.cantidad) AS total_unidades,
        RANK() OVER (PARTITION BY EXTRACT(MONTH FROM v.fecha) ORDER BY SUM(pv.cantidad) DESC) AS ranking
    FROM VENTA v
    JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta
    JOIN PRODUCTO p ON pv.id_producto = p.id_producto
    WHERE EXTRACT(YEAR FROM v.fecha) = 2021
    GROUP BY EXTRACT(MONTH FROM v.fecha), p.nombre
)
SELECT mes, producto, total_unidades 
FROM CantidadesMes 
WHERE ranking = 1 
ORDER BY mes;

-- =======================================================================
-- 2. Producto más económico por tienda.
-- =======================================================================
WITH PreciosPorTienda AS (
    SELECT DISTINCT
        t.nombre AS tienda,
        p.nombre AS producto,
        p.precio,
        RANK() OVER (PARTITION BY t.id_tienda ORDER BY p.precio ASC) AS ranking
    FROM VENTA v
    JOIN TIENDA t ON v.id_tienda = t.id_tienda
    JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta
    JOIN PRODUCTO p ON pv.id_producto = p.id_producto
)
SELECT tienda, producto, precio 
FROM PreciosPorTienda 
WHERE ranking = 1;

-- =======================================================================
-- 3. Ventas por mes, separadas entre Boletas y Facturas.
-- =======================================================================
SELECT 
    EXTRACT(YEAR FROM v.fecha) AS anio,
    EXTRACT(MONTH FROM v.fecha) AS mes,
    td.nombre AS tipo_documento,
    COUNT(DISTINCT v.id_venta) AS cantidad_transacciones,
    SUM(pv.cantidad * p.precio) AS monto_total_recaudado
FROM VENTA v
    JOIN TIPO_DOC td ON v.id_tipo_doc = td.id_tipo_doc
    JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta
    JOIN PRODUCTO p ON pv.id_producto = p.id_producto
GROUP BY EXTRACT(YEAR FROM v.fecha), EXTRACT(MONTH FROM v.fecha), td.nombre
ORDER BY anio, mes, tipo_documento;

-- =======================================================================
-- 4. Empleado que ganó más por tienda en 2020, indicando la comuna donde vive y el cargo que tiene en la empresa.
-- =======================================================================
WITH Sueldos2020 AS (
    SELECT 
        t.nombre AS tienda,
        e.nombre AS empleado,
        e.cargo,
        c.nombre AS comuna_residencia,
        SUM(s.monto) AS total_ganado_ano,
        RANK() OVER (PARTITION BY t.id_tienda ORDER BY SUM(s.monto) DESC) AS ranking
    FROM EMPLEADO e
    JOIN SUELDO s ON e.id_empleado = s.id_empleado
    JOIN TIENDA_EMP te ON e.id_empleado = te.id_empleado
    JOIN TIENDA t ON te.id_tienda = t.id_tienda
    JOIN COMUNA c ON e.id_comuna = c.id_comuna
    WHERE s.anio = 2020
    GROUP BY t.id_tienda, t.nombre, e.id_empleado, e.nombre, e.cargo, c.nombre
)
SELECT tienda, empleado, cargo, comuna_residencia, total_ganado_ano 
FROM Sueldos2020 
WHERE ranking = 1;

-- =======================================================================
-- 5. La tienda que tiene menos empleados.
-- =======================================================================
SELECT 
    t.nombre AS tienda,
    COUNT(te.id_empleado) AS cantidad_empleados
FROM TIENDA t
LEFT JOIN TIENDA_EMP te ON t.id_tienda = te.id_tienda
GROUP BY t.id_tienda, t.nombre
ORDER BY cantidad_empleados ASC
LIMIT 1;

-- =======================================================================
-- 6. El vendedor con más ventas por mes.
-- =======================================================================
WITH VentasVendedor AS (
    SELECT 
        EXTRACT(YEAR FROM v.fecha) AS anio,
        EXTRACT(MONTH FROM v.fecha) AS mes,
        emp.nombre AS vendedor,
        COUNT(v.id_venta) AS cantidad_ventas,
        RANK() OVER (PARTITION BY EXTRACT(YEAR FROM v.fecha), EXTRACT(MONTH FROM v.fecha) ORDER BY COUNT(v.id_venta) DESC) AS ranking
    FROM VENTA v
    JOIN VENDEDOR vend ON v.id_vendedor = vend.id_vendedor
    JOIN EMPLEADO emp ON vend.id_empleado = emp.id_empleado
    GROUP BY EXTRACT(YEAR FROM v.fecha), EXTRACT(MONTH FROM v.fecha), emp.nombre
)
SELECT anio, mes, vendedor, cantidad_ventas 
FROM VentasVendedor 
WHERE ranking = 1 
ORDER BY anio, mes;

-- =======================================================================
-- 7. El vendedor que ha recaudado más dinero para la tienda por año.
-- =======================================================================
WITH RecaudacionAnual AS (
    SELECT 
        EXTRACT(YEAR FROM v.fecha) AS anio,
        t.nombre AS tienda,
        emp.nombre AS vendedor,
        SUM(pv.cantidad * p.precio) AS total_dinero_recaudado,
        RANK() OVER (PARTITION BY EXTRACT(YEAR FROM v.fecha), t.id_tienda ORDER BY SUM(pv.cantidad * p.precio) DESC) AS ranking
    FROM VENTA v
    JOIN TIENDA t ON v.id_tienda = t.id_tienda
    JOIN VENDEDOR vend ON v.id_vendedor = vend.id_vendedor
    JOIN EMPLEADO emp ON vend.id_empleado = emp.id_empleado
    JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta
    JOIN PRODUCTO p ON pv.id_producto = p.id_producto
    GROUP BY EXTRACT(YEAR FROM v.fecha), t.id_tienda, t.nombre, emp.nombre
)
SELECT anio, tienda, vendedor, total_dinero_recaudado 
FROM RecaudacionAnual 
WHERE ranking = 1 
ORDER BY anio, tienda;

-- =======================================================================
-- 8. El vendedor con más productos vendidos por tienda.
-- =======================================================================
WITH ProductosPorTienda AS (
    SELECT 
        t.nombre AS tienda,
        emp.nombre AS vendedor,
        SUM(pv.cantidad) AS total_unidades_vendidas,
        RANK() OVER (PARTITION BY t.id_tienda ORDER BY SUM(pv.cantidad) DESC) AS ranking
    FROM VENTA v
    JOIN TIENDA t ON v.id_tienda = t.id_tienda
    JOIN VENDEDOR vend ON v.id_vendedor = vend.id_vendedor
    JOIN EMPLEADO emp ON vend.id_empleado = emp.id_empleado
    JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta
    GROUP BY t.id_tienda, t.nombre, emp.nombre
)
SELECT tienda, vendedor, total_unidades_vendidas 
FROM ProductosPorTienda 
WHERE ranking = 1;

-- =======================================================================
-- 9. El empleado con mayor sueldo por mes.
-- =======================================================================
WITH MaxSueldoMes AS (
    SELECT 
        s.anio,
        s.mes,
        e.nombre AS empleado,
        s.monto AS sueldo_pagado,
        RANK() OVER (PARTITION BY s.anio, s.mes ORDER BY s.monto DESC) AS ranking
    FROM SUELDO s
    JOIN EMPLEADO e ON s.id_empleado = e.id_empleado
)
SELECT anio, mes, empleado, sueldo_pagado 
FROM MaxSueldoMes 
WHERE ranking = 1 
ORDER BY anio, mes;

-- =======================================================================
-- 10. La tienda con menor recaudación por mes.
-- =======================================================================
WITH RecaudacionTiendaMes AS (
    SELECT 
        EXTRACT(YEAR FROM v.fecha) AS anio,
        EXTRACT(MONTH FROM v.fecha) AS mes,
        t.nombre AS tienda,
        SUM(pv.cantidad * p.precio) AS total_mes,
        RANK() OVER (PARTITION BY EXTRACT(YEAR FROM v.fecha), EXTRACT(MONTH FROM v.fecha) ORDER BY SUM(pv.cantidad * p.precio) ASC) AS ranking
    FROM VENTA v
    JOIN TIENDA t ON v.id_tienda = t.id_tienda
    JOIN PROD_VENTA pv ON v.id_venta = pv.id_venta
    JOIN PRODUCTO p ON pv.id_producto = p.id_producto
    GROUP BY EXTRACT(YEAR FROM v.fecha), EXTRACT(MONTH FROM v.fecha), t.nombre
)
SELECT anio, mes, tienda, total_mes AS peor_recaudacion 
FROM RecaudacionTiendaMes 
WHERE ranking = 1 
ORDER BY anio, mes;