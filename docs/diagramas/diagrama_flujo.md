┌─────────────────────────────────────────────────────────────────┐
│                    DON WEB (Ferozo - Argentina)                 │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │  MySQL: branet_gesdrims                                  │   │
│  │  ├─ Tabla: maecompr (Comprobantes/Facturas)              │   │
│  │  ├─ Tabla: movstock (Detalle de artículos)               │   │
│  │  ├─ Tabla: movpagos (Formas de pago)                     │   │
│  │  └─ Tabla: control_patio_olmos (Control - NUEVA)         │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
         ▲
         │ (Consultas SQL c/20 seg)
         │
┌────────┴───────────────────────────────────────────────────────┐
│           SUCURSAL PATIO OLMOS (18 terminales)                 │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  PC DE CAJA #1 - Windows                                 │  │
│  │  ┌────────────────────────────────────────────────────┐  │  │
│  │  │  GO/Node.js Worker Service (Daemon)                │  │  │
│  │  │  • Lee datos de BD Don Web                         │  │  │
│  │  │  • Procesa y formatea records C, D, P              │  │  │
│  │  │  • Ejecuta transfer.exe (CLI) vía OS               │  │  │
│  │  │  • Guarda auditoría en control_patio_olmos         │  │  │
│  │  └────────────────────────────────────────────────────┘  │  │
│  │           ▼                                              │  │
│  │  ┌────────────────────────────────────────────────────┐  │  │
│  │  │  transfer.exe (Marck System - TotalSale)           │  │  │
│  │  │  • Binario obligatorio del shopping                │  │  │
│  │  │  • Invocado 3 veces: C|..., D|..., P|...           │  │  │
│  │  │  • Envía a colectora del shopping                  │  │  │
│  │  └────────────────────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                │
│  [PC #2] [PC #3] ... [PC #18] (Mismo esquema)                  │
└────────────────────────────────────────────────────────────────┘
         ▼
┌─────────────────────────────────────────────────────────────────┐
│            MARCK SYSTEM - PATIO OLMOS (Colectora)               │
│  Recibe datos de ventas en tiempo real                          │
└─────────────────────────────────────────────────────────────────┘



┌─────────────────────────────────────────────────────────────────┐
│  CICLO DE POLLING (Cada 20 segundos - Configurable)             │
└─────────────────────────────────────────────────────────────────┘

1️⃣  FASE DE SELECCIÓN
   ├─ Conectar a BD Don Web (MySQL)
   ├─ Ejecutar query:
   │  SELECT ID, TIPCOM, PTV, NROCOM, FECCOM, HORCOM 
   │  FROM maecompr 
   │  WHERE CGOSUC = '02' 
   │    AND ANULAD = 0 
   │    AND ID NOT IN (SELECT idmaecompr FROM control_patio_olmos) 
   │  ORDER BY ID ASC 
   │  LIMIT 10;
   └─ Obtiene lista de facturas NUEVAS (máx 10)

2️⃣  FASE DE CONSTRUCCIÓN (Loop por cada factura)
   └─ Para cada factura encontrada:
      ├─ Recuperar datos CABECERA de maecompr
      │  (ID, TIPCOM, PTV, NROCOM, FECCOM, HORCOM)
      ├─ Recuperar DETALLE de movstock
      │  (CANTID, NETO, IVALIN, TASVIG)
      ├─ Recuperar PAGOS de movpagos
      │  (forma de pago, monto, tarjeta)
      └─ Sanitizar datos:
         • Remover ceros a izquierda
         • Formatear fechas: "DD-MM-YYYY" → "YYYYMMDD"
         • Formatear horas: "HH:MM:SS" → "HHMM"
         • Normalizar decimales con punto (.)

3️⃣  FASE DE DISPARO OPERATIVO (CLI - Ejecutar transfer.exe)
   └─ Invocar 3 comandos secuenciales:
      
      Cmd 1: transfer.exe "record:C|92|81|2|3|20260703|1430|..."
      └─ Captura Exit Code (0 = OK, otros = Error)
      
      Cmd 2: transfer.exe "record:D|1|12.00|500.00|105.00|21.00|0001|..."
      └─ Captura Exit Code
      
      Cmd 3: transfer.exe "record:P|01||500.00|..."
      └─ Captura Exit Code

4️⃣  FASE DE CONFIRMACIÓN
   ├─ Si EXIT CODE = 0 en los 3 comandos:
   │  INSERT INTO control_patio_olmos
   │  VALUES (id_factura, 1, NOW(), 'OK')
   │  └─ Marca factura como procesada ✅
   │
   └─ Si EXIT CODE ≠ 0:
      ├─ Incrementar intentos_fallidos
      ├─ Guardar respuesta_transfer (stderr)
      └─ Reintentará en siguiente ciclo (próximos 20 seg)

5️⃣  VUELVE al paso 1️⃣ (continúa indefinidamente)


FACTURA EN CAJA #1 (Sucursal Olmos, PTV 0002)
│
├─ 14:30 hrs: Cliente paga $500 (200 neto + $105 IVA 21%)
│             Se guarda en maecompr con ID=5428
│
├─ 14:50 hrs: Worker detecta ID=5428 (no está en control_patio_olmos)
│
├─ Worker construye:
│  • Record C: "record:C|92|81|2|3|20260703|1430|500|105|0"
│  • Record D: "record:D|1|12.00|200.00|105.00|21.00|0001"
│  • Record P: "record:P|01||500.00"
│
├─ Worker ejecuta (vía CLI):
│  ✅ transfer.exe "record:C|..." → EXIT 0
│  ✅ transfer.exe "record:D|..." → EXIT 0
│  ✅ transfer.exe "record:P|..." → EXIT 0
│
└─ Worker inserta en BD:
   INSERT INTO control_patio_olmos
   VALUES (5428, 1, '2026-07-03 14:50:35', 'OK')
   
   ➜ Factura AUDITABLE y NO SE DUPLICARÁ



Escenario: El shopping se cae por 2 horas

├─ Clientes siguen comprando en caja
├─ ERP sigue guardando en BD Don Web (funciona local)
├─ Worker acumula las transacciones sin procesar
│
└─ Cuando se recupera el shopping:
   ├─ Worker retoma ciclo de polling
   ├─ Ejecuta modo BATCH (opcional):
   │  transfer.exe "file:C:\ruta\olmos_batch_20260703.txt"
   │  (Envía todas pendientes de ese día)
   └─ Sincroniza automáticamente ✅