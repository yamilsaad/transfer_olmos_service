# Procedimiento de Envío de Datos

**Versión:** Tec.P.4.2.0.0.A  
**Sistema:** TotalSale - Marck System  
**Fecha:** 01/11/2004

---

# Objetivos

Definir los diferentes métodos de captura de información de ventas utilizados por el centro comercial para interactuar con los locales comerciales.

# Alcance

Este documento está dirigido a:

- Soporte técnico del centro comercial.
- Programadores.
- Responsables informáticos de los locales.

# Responsabilidades

## Responsable del Local

- Conocer el procedimiento de envío de ventas.
- Adaptar el funcionamiento del local a los nuevos requerimientos.

## Responsable Informático

- Adaptar o desarrollar el sistema para transmitir la información de ventas.

---

# Métodos de Envío

Actualmente existen dos métodos:

1. Envío Online mediante `transfer.exe`
2. Retransmisión de archivos históricos

---

# 1. Envío Online

El sistema debe ejecutar el programa:

```bash
transfer.exe "record:[PARAMETROS]"
```

Todos los parámetros se separan mediante:

```text
|
```

Ejemplo:

```text
transfer.exe "record:C|125|1|1|2512|20041101|1350|DNI|20528522|DNI|24333125|N|1|20041101"
```

---

# Parámetros Comunes

| Campo | Tipo | Obligatorio | Descripción |
|--------|------|-------------|-------------|
| Tipo Registro | Char | Sí | C, D o P |
| Número Cliente | Numérico | Sí | Identificador del local |
| Código Comprobante | Numérico | Sí | Código AFIP |
| Punto de Venta | Numérico | Sí | Punto de venta |
| Número Comprobante | Numérico | Sí | Número del comprobante |

---

# Tipos de Registro

## Registro C - Cabecera

Incluye:

- Fecha
- Hora
- Datos del vendedor
- Datos del comprador
- Estado (Anulado/Normal)
- Terminal POS
- Fecha Operativa

### Ejemplo

```text
C|125|1|1|2512|20041101|1350|DNI|20528522|DNI|24333125|N|1|20041101
```

---

## Registro D - Detalle

Incluye:

- IVA
- Neto
- IVA calculado
- Otros impuestos
- Rubro
- Cantidad

Ejemplo:

```text
D|125|1|1|2512|21.00|150.25|31.55|0|281|5
```

---

## Registro P - Pago

Incluye:

- Forma de pago
- Tarjeta
- Cuotas
- Importe

Ejemplo

```text
P|125|1|1|2512|2|6|3|181.80
```

---

# Retransmisión

Cuando sea necesario reenviar información histórica:

```bash
transfer.exe "file:ventas.dat"
```

---

# Formas de Pago

| Código | Descripción |
|---------|-------------|
| 01 | Efectivo |
| 02 | Tarjeta |
| 03 | Cheque |
| 04 | Dólares |
| 05 | Cuenta Corriente |
| 06 | Otro |

...

---

# Tarjetas

| Código | Tarjeta |
|---------|----------|
|01|Mastercard|
|02|Visa|
|03|Visa Electron|
|04|American Express|
|05|Diners|
|06|Naranja|

...

---

# Ejemplos completos

## Ticket Factura B contado

```text
RECORD:C|...
RECORD:D|...
RECORD:P|...
```

## Ticket Factura B con Visa

```text
RECORD:C|...
RECORD:D|...
RECORD:P|...
```

---

# Referentes Técnicos

Departamento de Sistemas

Patio Olmos Shopping

- José Claudio Carrizo
- Leandro Toledo
- José Reig