
-- Configurar codificación UTF-8
SET client_encoding = 'UTF8';


-- Insertar Comunas
INSERT INTO COMUNA (nombre) VALUES ('Santiago Centro'), ('Providencia'), ('Pedro Aguirre Cerda');

-- Insertar Tipos de Documento obligatorios
INSERT INTO TIPO_DOC (nombre) VALUES ('Boleta'), ('Factura');

-- Insertar Productos con precios variados
INSERT INTO PRODUCTO (nombre, precio) VALUES ('Notebook Gamer', 899990.00), ('Mouse Inalámbrico', 19990.00), ('Teclado Mecánico', 45990.00), ('Monitor 24 pulgadas', 129990.00);

-- Insertar Tiendas
INSERT INTO TIENDA (nombre, id_comuna) VALUES ('Tienda Central Alameda', 1), ('Tienda Costanera', 2), ('Tienda Outlet PAC', 3);

-- Insertar Empleados con cargos
INSERT INTO EMPLEADO (nombre, cargo, id_comuna) VALUES 
('Juan Pérez', 'Vendedor Senior', 1),
('María José', 'Vendedora Junior', 2),
('Carlos Plaza', 'Administrador', 3),
('Ana Gómez', 'Vendedora Part-Time', 1);

-- Asociar Empleados a Tiendas (Para ver quién trabaja dónde)
INSERT INTO TIENDA_EMP (id_tienda, id_empleado) VALUES (1, 1), (1, 3), (2, 2), (2, 3), (3, 4);

-- Convertir algunos empleados en Vendedores activos
INSERT INTO VENDEDOR (id_empleado) VALUES (1), (2), (4);

-- Insertar Sueldos para el año 2020 y 2021 (Meses 1 al 3 para pruebas)
INSERT INTO SUELDO (id_empleado, monto, mes, anio) VALUES 
(1, 550000.00, 1, 2020), (1, 560000.00, 2, 2020), (1, 580000.00, 1, 2021),
(2, 450000.00, 1, 2020), (2, 470000.00, 2, 2020), (2, 490000.00, 1, 2021),
(3, 900000.00, 1, 2020), (3, 950000.00, 2, 2020), (3, 980000.00, 1, 2021),
(4, 300000.00, 1, 2020), (4, 310000.00, 2, 2020), (4, 320000.00, 1, 2021);

-- Insertar Ventas en 2020 y 2021
-- Ventas 2021 (Para consultas de productos más vendidos)
INSERT INTO VENTA (fecha, id_tipo_doc, id_tienda, id_vendedor) VALUES 
('2021-01-15', 1, 1, 1), ('2021-01-20', 2, 1, 1),
('2021-02-10', 1, 2, 2), ('2021-02-18', 2, 2, 2),
('2021-03-05', 1, 3, 3);

-- Ventas 2020
INSERT INTO VENTA (fecha, id_tipo_doc, id_tienda, id_vendedor) VALUES 
('2020-05-12', 1, 1, 1), ('2020-06-22', 2, 2, 2);

-- Detalles de Ventas (Relación producto - cantidad)
INSERT INTO PROD_VENTA (id_venta, id_producto, cantidad) VALUES 
(1, 1, 2), (1, 2, 5), -- Venta 1: 2 Notebooks y 5 Mouses
(2, 3, 3),            -- Venta 2: 3 Teclados
(3, 1, 1), (3, 4, 2), -- Venta 3: 1 Notebook y 2 Monitores
(4, 2, 10),           -- Venta 4: 10 Mouses
(5, 4, 4),            -- Venta 5: 4 Monitores
(6, 1, 1),            -- Venta 6: 1 Notebook
(7, 3, 2);            -- Venta 7: 2 Teclados


