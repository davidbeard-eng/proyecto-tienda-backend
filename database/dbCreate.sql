

-- Configurar codificación UTF-8
SET client_encoding = 'UTF8';



-- 1. Eliminar tablas en orden inverso por si ya existen
DROP TABLE IF EXISTS PROD_VENTA CASCADE;
DROP TABLE IF EXISTS VENTA CASCADE;
DROP TABLE IF EXISTS PRODUCTO CASCADE;
DROP TABLE IF EXISTS TIPO_DOC CASCADE;
DROP TABLE IF EXISTS VENDEDOR CASCADE;
DROP TABLE IF EXISTS TIENDA_EMP CASCADE;
DROP TABLE IF EXISTS SUELDO CASCADE;
DROP TABLE IF EXISTS EMPLEADO CASCADE;
DROP TABLE IF EXISTS TIENDA CASCADE;
DROP TABLE IF EXISTS COMUNA CASCADE;

-- 2. Crear tablas base (sin dependencias)
CREATE TABLE COMUNA (
    id_comuna SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL
);

CREATE TABLE TIPO_DOC (
    id_tipo_doc SERIAL PRIMARY KEY,
    nombre VARCHAR(50) NOT NULL -- Aquí irá 'Boleta' o 'Factura'
);

CREATE TABLE PRODUCTO (
    id_producto SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    precio NUMERIC(12, 2) NOT NULL
);

-- 3. Crear tablas con dependencias simples
CREATE TABLE TIENDA (
    id_tienda SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    id_comuna INT REFERENCES COMUNA(id_comuna) ON DELETE CASCADE
);

CREATE TABLE EMPLEADO (
    id_empleado SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    cargo VARCHAR(100) NOT NULL,
    id_comuna INT REFERENCES COMUNA(id_comuna) ON DELETE CASCADE
);

-- 4. Crear tablas de detalles e historiales
CREATE TABLE SUELDO (
    id_sueldo SERIAL PRIMARY KEY,
    id_empleado INT REFERENCES EMPLEADO(id_empleado) ON DELETE CASCADE,
    monto NUMERIC(12, 2) NOT NULL,
    mes INT NOT NULL,
    anio INT NOT NULL
);

CREATE TABLE VENDEDOR (
    id_vendedor SERIAL PRIMARY KEY,
    id_empleado INT REFERENCES EMPLEADO(id_empleado) ON DELETE CASCADE
);

-- 5. Relación Muchos a Muchos entre Tienda y Empleado
CREATE TABLE TIENDA_EMP (
    id_tienda INT REFERENCES TIENDA(id_tienda) ON DELETE CASCADE,
    id_empleado INT REFERENCES EMPLEADO(id_empleado) ON DELETE CASCADE,
    PRIMARY KEY (id_tienda, id_empleado)
);

-- 6. Encabezado de Ventas
CREATE TABLE VENTA (
    id_venta SERIAL PRIMARY KEY,
    fecha DATE NOT NULL,
    id_tipo_doc INT REFERENCES TIPO_DOC(id_tipo_doc),
    id_tienda INT REFERENCES TIENDA(id_tienda),
    id_vendedor INT REFERENCES VENDEDOR(id_vendedor)
);

-- 7. Detalle de Ventas (Relación Muchos a Muchos entre Venta y Producto)
CREATE TABLE PROD_VENTA (
    id_venta INT REFERENCES VENTA(id_venta) ON DELETE CASCADE,
    id_producto INT REFERENCES PRODUCTO(id_producto) ON DELETE CASCADE,
    cantidad INT NOT NULL,
    PRIMARY KEY (id_venta, id_producto)
);

