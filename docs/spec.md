# Documento de Especificaciones Técnicas (Specs V1)
## Middleware de Integración de Facturación - Patio Olmos Shopping (Grupo Brant / Drims)

**Estado:** Inicial / Propuesta de Arquitectura  
**Fecha:** Junio de 2026  
**Versión de Lineamientos de Interfaz:** Tec.P.4.2.0.0.A (TotalSale de Marck System)

---
## 0. Contexto del Proyecto
El sistema de facturación utilizado por Grupo Brant fue desarrollado por un proveedor externo sobre PHP y actualmente se encuentra en producción.

Debido a conflictos contractuales con el proveedor original, no es posible realizar modificaciones sobre el código fuente de la aplicación.

Sin embargo, Patio Olmos Shopping exige la transmisión online de las ventas utilizando la Interface de Comunicación "transfer.exe" desarrollada por Marck System.

El presente proyecto implementa un middleware independiente encargado de leer la información directamente desde la base de datos del ERP y transmitirla utilizando el protocolo homologado por Patio Olmos sin alterar el sistema existente.

## 1. Arquitectura de la Solución y Modelo de Despliegue
El sistema comercial actual está desarrollado en PHP, alojado en Don Web y utiliza una base de datos centralizada MySQL (esquema `branet_gesdrims`). Al no poseer acceso al código fuente del ERP, la captura de transacciones se realiza de forma externa y desacoplada mediante consultas cíclicas directas sobre la base de datos (Near Real-Time Polling Worker).

### Modelo Propuesto: Worker Local en cada Punto de Venta (Opción Recomendada)
Se propone compilar un servicio liviano (Daemon/Servicio de Windows) escrito en Go o Node.js que se instale individualmente en cada computadora de caja física de la sucursal del Patio Olmos.

| Componente / Criterio | Descripción y Ventajas | Mitigación de Riesgos |
| :--- | :--- | :--- |
| **Desacoplamiento** | Cero modificaciones en el servidor de producción Don Web o en el código de la aplicación PHP. | No altera el circuito habitual de facturación ni degrada la performance global del backend web. |
| **Ejecución de CLI** | Al estar instalado de forma local, el servicio tiene permisos directos del Sistema Operativo para invocar el binario `transfer.exe`. | Evita configuraciones complejas de ejecución remota de comandos o brechas de seguridad por SSH inverso. |
| **Resiliencia a Fallas de Red** | Si la sucursal pierde conectividad a Internet momentáneamente, los tickets se siguen guardando en Don Web por el ERP, y al restablecerse la red, el Worker local sincroniza los pendientes de forma automática secuencial. | Si falla el servidor local del shopping, se implementa almacenamiento de auditoría local en archivos planos temporales de contingencia ("file:"). |

---

## 2. Esquema de Datos de Control Interno (Idempotencia)
Para asegurar que ninguna factura sea procesada por duplicado (idempotencia) y permitir auditorías exhaustivas del estado de comunicación con la colectora del shopping, se requiere la creación de una tabla de control permanente.

Dado que contamos con acceso total de administrador (`root`) a la base de datos MySQL de producción, se ejecutará el siguiente bloque DDL estructural dentro del esquema `branet_gesdrims`:

```sql
CREATE TABLE IF NOT EXISTS control_patio_olmos (
    idmaecompr INT UNSIGNED NOT NULL,
    enviado TINYINT(1) DEFAULT 0,
    fecha_envio DATETIME NULL,
    linea_c_enviada TEXT NULL,
    respuesta_transfer VARCHAR(255) NULL,
    intentos_fallidos INT DEFAULT 0,
    PRIMARY KEY (idmaecompr)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

## 3. Matriz de Mapeo y Traducción de Campos
La recolección y conversión de tipos de datos desde las tablas relacionales origen hacia la interfaz posicional delimitada por barras (|) de la colectora TotalSale se estructurará bajo el siguiente esquema lógico:

### 3.1. Registro C - Cabecera del Comprobante
Origen de Datos Principal: Tabla maecompr. Filtro crítico de sucursal mediante el campo CGOSUC (código numérico asignado a Patio Olmos).

	-NroCliente (Código de Local del Shopping): Valor estático configurado por variable de entorno (Ej. 92).

	-CodComprobanteAFIP: Traducción condicional del string guardado en TIPCOM / ATIPCOM (Ej. Si "FB", mapear a 82 [Factura B]; si "FA", mapear a 81 [Factura A]; Notas de Crédito correspondientes).

	-Punto de Venta / NroComprobante: Campos PTV y NROCOM. Se les debe aplicar un parseo numérico/expresión regular para remover de forma estricta los ceros a la izquierda (Ej: "00003" se transforma en "3").

	-Fecha / Hora del Comprobante: Campos FECCOM y HORCOM. Conversión estricta de cadenas para remover guiones y dos puntos. Formatos resultantes requeridos: Fecha = AAAAMMDD (8 caracteres continuos), Hora = HHMM (4 caracteres continuos).

### 3.2. Registro D - Detalle de Artículos
Origen de Datos Principal: Tabla movstock (Enlace: movstock.idmaecompr = maecompr.ID).

	-Alícuota IVA: Extraído del campo TASVIG (Ej: 21.00). Los valores decimales se formatean forzando el separador de punto (.) y eliminando comas o separadores de miles.
	-Neto / IVA de Línea: Campos NETO e IVALIN respectivamente.
	-Cantidad: Campo CANTID casteado a flotante/entero limpio.
	-Código de Rubro: Fijado por la administración del Patio Olmos para el rubro general del local (Valor fijo hardcodeado: 0001).

### 3.3. Registro P - Formas de Pago
Origen de Datos Principal: Tabla movpagos o datosfpagos (Enlace: movpagos.id_movVentas = maecompr.ID).

	-Forma de Pago: Mapeo numérico homologado basado en el campo id_formaPago y la cadena DescripcionPago (Efectivo = 01, Tarjeta = 02).
	-Código de Tarjeta de Crédito: Evaluación por condicional sobre DescripcionPago para mapear con la lista oficial de Patio Olmos (Ej. Si incluye "VISA" $\rightarrow$ Código 2; si incluye "MASTER" $\rightarrow$ Código 1; si incluye "NARANJA" $\rightarrow$ Código 6). Si la forma de pago es efectivo, la posición del string se envía vacía.
	-Importe Cobrado: Campo monto parseado con punto decimal.

## 4.Algoritmo de Procesamiento de Negocio y Ciclo de Polling
El servicio middleware ejecutará una rutina en segundo plano con intervalos de tiempo parametrizables (Polling cíclico sugerido cada 20 segundos) siguiendo la siguiente lógica secuencial:

	1-Fase de Selección: Consultar las facturas de la base de datos de producción que cumplan con los filtros de sucursal del shopping, no se encuentren marcadas como anuladas, y que no existan en el universo registrado dentro de la tabla de control local:

		SELECT ID, TIPCOM, PTV, NROCOM, FECCOM, HORCOM 
		FROM maecompr 
		WHERE CGOSUC = '02' 
		AND ANULAD = 0 
		AND ID NOT IN (SELECT idmaecompr FROM control_patio_olmos) 
		ORDER BY ID ASC 
		LIMIT 10;

	2-Fase de Construcción (Loop por Factura):

		-Recuperar datos de la cabecera, del set de registros de artículos en movstock y de los pagos en movpagos.

		-Aplicar sanitización de datos (remover ceros, caracteres especiales, formatear strings de fecha/hora y normalizar el punto decimal).

		-Concatenar las cadenas agregando los prefijos requeridos para el modo en línea ("record:C|...", "record:D|...", "record:P|...").

	3-Fase de Disparo Operativo: Invocar de forma secuencial, síncrona y ordenada las tres ejecuciones del binario mediante comandos CLI nativos:

	exec.Command("transfer.exe", "record:C|...")

	Capturar la salida del flujo estándar (stdout/stderr) devuelto por el sistema operativo.

	4-Fase de Confirmación: Si las ejecuciones CLI devolvieron códigos de salida normales (Exit Code 0), realizar el guardado en la tabla de control:

	INSERT INTO control_patio_olmos (idmaecompr, enviado, fecha_envio, respuesta_transfer) VALUES (id_actual, 1, NOW(), 'OK');

## 5. Estrategia de Contingencia para Envío por Lotes (Bajo Demanda)
En el escenario de fallos prolongados en la interfaz de red o interrupciones en la colectora del shopping, la especificación técnica requiere soporte de transmisión masiva asíncrona mediante archivos físicos.

	-El middleware implementará una bandera o flag secundario ejecutable vía CLI (Ej: --mode=batch --date=AAAAMMDD).

	-Este módulo consultará las facturas históricas no procesadas de la fecha indicada, generará de forma consolidada todas las líneas C, D y P en un único archivo de texto plano denominado de forma interna (Ej: olmos_batch.txt).

	-Finalmente, el script invocará la herramienta de Marck System utilizando el puntero de archivo masivo: transfer.exe "file:C:\Ruta\olmos_batch.txt", asegurando el cumplimiento fiscal de retransmisión bajo demanda del shopping.