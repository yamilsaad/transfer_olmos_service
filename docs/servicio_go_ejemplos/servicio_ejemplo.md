# Cómo Funciona el Servicio Go (Windows Service)


## 1 - El binario Go compilado se instala como servicio de Windows:

*Instalación (una sola vez)*
    sc create OlmosTransfer ^
    binPath="C:\Program Files\OlmosTransfer\transfer-service.exe" ^
    displayName="Olmos Transfer Service" ^
    start=auto

*Inicio*
    net start OlmosTransfer

*El servicio quedará activo indefinidamente*
*(Se reinicia automático si Windows reinicia)*

## 2 - Arquitectura Interna del Servicio Go:

´´´
    package main

    import (
        "time"
        "log"
        "database/sql"
    )

    func main() {
        // 1. Inicializar conexión a BD Don Web (reutilizable)
        db, err := sql.Open("mysql", 
            "user:pass@tcp(don-web-server:3306)/branet_gesdrims")
        if err != nil {
            log.Fatal("Error conectando a BD", err)
        }
        defer db.Close()

        // 2. Crear TICKER para polling cada 20 segundos
        ticker := time.NewTicker(20 * time.Second)
        defer ticker.Stop()

        log.Println("✅ Servicio iniciado - Polling cada 20 segundos")

        // 3. LOOP INFINITO (Se mantiene activo)
        for {
            select {
            case <-ticker.C:
                // Cada 20 segundos ejecuta esta función
                procesarFacturasPendientes(db)
            }
        }
    }

    func procesarFacturasPendientes(db *sql.DB) {
        log.Println("[CICLO] Buscando facturas nuevas...")

        // FASE 1: SELECCIÓN
        rows, err := db.Query(`
            SELECT ID, TIPCOM, PTV, NROCOM, FECCOM, HORCOM 
            FROM maecompr 
            WHERE CGOSUC = '02' 
            AND ANULAD = 0 
            AND ID NOT IN (SELECT idmaecompr FROM control_patio_olmos) 
            ORDER BY ID ASC 
            LIMIT 10
        `)
        if err != nil {
            log.Println("❌ Error en query:", err)
            return // Reintentar en 20 segundos
        }
        defer rows.Close()

        // FASE 2: CONSTRUCCIÓN
        for rows.Next() {
            var id, tipcom, ptv, nrocom, feccom, horcom string
            rows.Scan(&id, &tipcom, &ptv, &nrocom, &feccom, &horcom)

            // Procesar factura
            recordC := construirRecordC(id, tipcom, ptv, nrocom, feccom, horcom)
            recordD := construirRecordD(db, id)
            recordP := construirRecordP(db, id)

            // FASE 3: DISPARO OPERATIVO (Ejecutar CLI)
            exitCodeC := ejecutarTransfer(recordC)
            exitCodeD := ejecutarTransfer(recordD)
            exitCodeP := ejecutarTransfer(recordP)

            // FASE 4: CONFIRMACIÓN
            if exitCodeC == 0 && exitCodeD == 0 && exitCodeP == 0 {
                registrarEnControl(db, id, "OK")
                log.Printf("✅ Factura %s enviada correctamente\n", id)
            } else {
                registrarEnControl(db, id, "ERROR")
                log.Printf("❌ Factura %s falló (códigos: %d,%d,%d)\n", 
                    id, exitCodeC, exitCodeD, exitCodeP)
            }
        }
    }

    func ejecutarTransfer(record string) int {
        cmd := exec.Command("transfer.exe", record)
        err := cmd.Run()
        
        if exitErr, ok := err.(*exec.ExitError); ok {
            return exitErr.ExitCode()
        }
        return 0
    }
´´´

# 3 - Timeline Visual: Qué Hace el Servicio:

╔════════════════════════════════════════════════════════════════════╗
║          SERVICIO GO EN PC DE CAJA #1 - LÍNEA DE TIEMPO            ║
╚════════════════════════════════════════════════════════════════════╝

T=0 seg:   ✅ Servicio inicia (se instala como Windows Service)
           └─ Abre conexión a BD Don Web
           └─ Inicia TICKER de 20 segundos

T=20 seg:  🔄 CICLO 1
           ├─ Lee 10 facturas nuevas de BD
           ├─ Ejecuta transfer.exe (3 comandos por factura)
           ├─ Registra resultado en control_patio_olmos
           └─ Espera 20 segundos

T=40 seg:  🔄 CICLO 2
           ├─ Lee 10 facturas nuevas (las anteriores están en control)
           ├─ Ejecuta transfer.exe
           ├─ Registra resultado
           └─ Espera 20 segundos

T=60 seg:  🔄 CICLO 3
           └─ ... (continúa indefinidamente)

...

T=∞:       ✅ Servicio sigue corriendo hasta:
           • Reinicio manual: net stop OlmosTransfer
           • Reinicio de Windows (se reinicia automático)
           • Error grave no recuperable (logging + reinicio automático)

# 4 - Manejo de Errores & Resiliencia:

´´´
    func procesarFacturasPendientes(db *sql.DB) {
        // Si la BD se desconecta
        err := db.Ping()
        if err != nil {
            log.Println("⚠️  BD no responde, reintentando en próximo ciclo...")
            return // No falla, espera 20 seg y reintenta
        }

        // Si transfer.exe no existe
        cmd := exec.Command("transfer.exe", record)
        err = cmd.Run()
        if err != nil {
            log.Println("⚠️  transfer.exe error:", err)
            // Se registra como fallido, reintenta en próximo ciclo
            return
        }

        // Reconexión automática
        if rows.Err() != nil {
            log.Println("⚠️  Error leyendo filas, reconectando...")
            db.Close()
            db, _ = sql.Open("mysql", connString)
        }
    }
´´´

# 5 - Configuración Típica en Windows

Servicio: OlmosTransfer
├─ Tipo: Servicio en background (nunca visible)
├─ Inicia: Automático (al arrancar Windows)
├─ Ejecutable: C:\Program Files\OlmosTransfer\transfer-service.exe
├─ Archivo de logs: C:\Program Files\OlmosTransfer\logs\service.log
└─ Estado: RUNNING (permanentemente)

*Ver estado en PowerShell:*
    Get-Service OlmosTransfer
    // Status   Name                      DisplayName
    // ------   ----                      -----------
    // Running  OlmosTransfer             Olmos Transfer Service


# 6 - Diagrama de Ejecución Interna:

┌─────────────────────────────────────────────────────────┐
│        SERVICIO GO - LOOP PRINCIPAL INFINITO            │
└─────────────────────────────────────────────────────────┘

        ╔════════════════════════════════════╗
        ║   main() - Inicia conexión BD      ║
        ║   - Conexión a Don Web (abierta)   ║
        ║   - Crea ticker 20 segundos        ║
        ║   - Inicia loop infinito           ║
        ╚════════════════════════════════════╝
                        │
                        ▼
        ┌──────────────────────────────────┐
        │  for { select { <-ticker.C } }   │  ◄── LOOP INFINITO
        │                                  │      Se repite cada 20 seg
        └──────────────────────────────────┘
                        │
         ┌──────────────┴──────────────┐
         ▼                             ▼
        [Fase 1]                    [Fase 2]
        Seleccionar                 Construir
        Facturas BD                 Records C,D,P
         │                             │
         └──────────────┬──────────────┘
                        ▼
                   [Fase 3]
                   Ejecutar
                   transfer.exe (CLI)
                        │
         ┌──────────────┴──────────────┐
         ▼                             ▼
    [EXIT CODE 0]               [EXIT CODE ≠ 0]
    Registra ÉXITO          Registra FALLO
    en control_patio_olmos  e incrementa intentos
         │                             │
         └──────────────┬──────────────┘
                        ▼
            Espera 20 segundos
                        │
                        ▼
                 ⬅️ VUELVE AL LOOP


# 7 - Consumo de Recursos (Típico):

Memoria:      ~50-100 MB (conexión abierta a BD + variables)
CPU:          ~0-2% (durmiendo 20 segundos)
Conexión BD:  1 sola (reutilizada en cada ciclo)
Logs:         ~1-2 MB/día (eventos de procesamiento)

