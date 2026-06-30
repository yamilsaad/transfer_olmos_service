# Middleware de Integración de Facturación - Shopping Patio Olmos

Este repositorio contiene la documentación técnica, especificaciones de bases de datos y la arquitectura propuesta para resolver la desconexión entre el sistema ERP actual de la empresa (Grupo Brant / Drims) y la colectora de datos obligatoria exigida por el shopping **Patio Olmos** (Marck System - TotalSale).

## Contexto del Problema

El cliente posee un sistema de gestión comercial web desarrollado en **PHP** y alojado en **Don Web (Base SQL en Ferozo)** implementado en 18 sucursales. Debido a un conflicto con el equipo de desarrollo original, **no se dispone de acceso al código fuente del sistema** ni a las interfaces de caja de los locales.

El shopping Patio Olmos exige contractualmente que las ventas se reporten en tiempo real a través de un ejecutable local obligatorio llamado `transfer.exe`. 

### Solución Arquitectónica
Dado que sí se cuenta con accesos de administración de base de datos (`root`), se optó por un enfoque desacoplado de **Near Real-Time Polling Worker**. Se implementará un servicio periférico liviano que correrá en background en las terminales del local, interrogará la base de datos de producción remota de forma periódica y enviará las estructuras formateadas a la CLI del shopping.

---

## Especificaciones Técnicas Clave

### 1. Datos Identificados de Producción (Sucursal)
Tras auditar las tablas del motor del ERP en producción, se determinó con exactitud la identidad del local comercial físico dentro del esquema relacional:
* **`CGOSUC` (Código de Sucursal):** `'02'`
* **`NOMSUC` (Nombre):** `OLMOS`
* **`PTOSUC` (Punto de Venta Base):** `'0002'`

### 2. Query de Fase de Selección (Algoritmo Base)
El Worker utilizará la siguiente consulta optimizada para capturar los comprobantes pendientes de envío garantizando no duplicar información:

```sql
SELECT ID, TIPCOM, PTV, NROCOM, FECCOM, HORCOM 
FROM maecompr 
WHERE CGOSUC = '02' 
  AND ANULAD = 0 
  AND ID NOT IN (SELECT idmaecompr FROM control_patio_olmos) 
ORDER BY ID ASC 
LIMIT 10;