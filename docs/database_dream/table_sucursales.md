# Tabla de SUCURSALES:

    mysql> DESCRIBE sucursal;
    +-----------------+---------------+------+-----+-----------------------+----------------+
    | Field           | Type          | Null | Key | Default               | Extra          |
    +-----------------+---------------+------+-----+-----------------------+----------------+
    | ID              | int unsigned  | NO   | PRI | NULL                  | auto_increment |
    | CGOSUC          | varchar(2)    | YES  |     | NULL                  |                |
    | NOMSUC          | varchar(35)   | YES  |     | NULL                  |                |
    | RAZON           | varchar(50)   | YES  |     | NULL                  |                |
    | DOMSUC          | varchar(55)   | YES  |     | NULL                  |                |
    | LOCSUC          | varchar(35)   | YES  |     | NULL                  |                |
    | PROSUC          | varchar(35)   | YES  |     | NULL                  |                |
    | TELSUC          | varchar(35)   | YES  |     | NULL                  |                |
    | CUISUC          | varchar(13)   | YES  |     | NULL                  |                |
    | BRUSUC          | varchar(13)   | YES  |     | NULL                  |                |
    | JUBSUC          | varchar(13)   | YES  |     | NULL                  |                |
    | INISUC          | date          | YES  |     | 0000-00-00            |                |
    | NPIE01          | varchar(80)   | YES  |     | NULL                  |                |
    | NPIE02          | varchar(80)   | YES  |     | NULL                  |                |
    | PTOSUC          | varchar(4)    | YES  |     | NULL                  |                |
    | DEPSUC          | varchar(2)    | YES  |     | NULL                  |                |
    | ULTFAA          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTNDA          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTNCA          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTFUA          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTFAB          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTNDB          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTNCB          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTFUB          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTREC          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTREM          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTINT          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTTIK          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTING          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTOPA          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTAJU          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTPTE          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTPRE          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTPED          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTNEN          | varchar(8)    | YES  |     | NULL                  |                |
    | ULTAUT          | varchar(8)    | YES  |     | NULL                  |                |
    | CLAVE1          | varchar(10)   | YES  |     | NULL                  |                |
    | CLAVE2          | varchar(10)   | YES  |     | NULL                  |                |
    | CLAVE3          | varchar(10)   | YES  |     | NULL                  |                |
    | HOMO            | int           | NO   |     | 1                     |                |
    | MAIL            | varchar(200)  | NO   |     | NULL                  |                |
    | DCTO_AUT        | decimal(10,2) | NO   |     | 0.00                  |                |
    | PTV             | varchar(5)    | NO   |     | NULL                  |                |
    | LISTAS          | int           | NO   |     | 2                     |                |
    | VTOFAC          | int           | NO   |     | 7                     |                |
    | IMPRIMEDCTO     | int           | NO   |     | 0                     |                |
    | MODFECHA        | int           | NO   |     | 0                     |                |
    | CORTO           | varchar(15)   | NO   |     | NULL                  |                |
    | PERDGR          | int           | NO   |     | 0                     |                |
    | CLVALEATORIA    | int           | NO   |     | 0                     |                |
    | AFIPCRT         | varchar(200)  | YES  |     | NULL                  |                |
    | AFIPKEY         | varchar(200)  | YES  |     | NULL                  |                |
    | CONCENTRA       | int           | NO   |     | 0                     |                |
    | CBU             | varchar(22)   | YES  |     | NULL                  |                |
    | CONDICIONIVA    | varchar(50)   | NO   |     | RESPONSABLE INSCRIPTO |                |
    | LIMITAOPERACION | int           | NO   |     | 0                     |                |
    | USATMP          | int           | NO   |     | 0                     |                |
    | DPACTADO        | int           | NO   |     | 0                     |                |
    | FORMATO_FACTURA | int           | NO   |     | 0                     |                |
    | TIPNUMRX        | int           | NO   |     | 0                     |                |
    | PDIRECTA        | int           | NO   |     | 0                     |                |
    | PRECONIVA       | int           | NO   |     | 0                     |                |
    | DEXTENDIDA      | int           | NO   |     | 0                     |                |
    | RENGLONES       | int           | NO   |     | 28                    |                |
    | REDONDEO        | int           | NO   |     | 0                     |                |
    +-----------------+---------------+------+-----+-----------------------+----------------+

    mysql>

## SUCURSALES:

    mysql> SELECT CGOSUC, NOMSUC, PTOSUC FROM sucursal;
    ERROR 4031 (HY000): The client was disconnected by the server because of inactivity. See wait_timeout and interactive_timeout for configuring this behavior.
    No connection. Trying to reconnect...
    Connection id:    1967539
    Current database: branet_gesdrims

    +--------+------------+--------+
    | CGOSUC | NOMSUC     | PTOSUC |
    +--------+------------+--------+
    | 01     | CENTRAL    | 0001   |
    | 02     | OLMOS      | 0002   |
    | 03     | DOT        | 0004   |
    | 04     | RIBERA     | 0011   |
    | 05     | RIO CUARTO | 0011   |
    | 06     | WEB        | NULL   |
    | 07     | HONDURAS   | NULL   |
    | 08     | CABILDO    | NULL   |
    | 09     | PEATONAL   | NULL   |
    +--------+------------+--------+
    9 rows in set (0,06 sec)

    mysql>
