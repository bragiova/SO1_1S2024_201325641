CREATE DATABASE proyecto1_so1

CREATE TABLE cpu_info
(
	id_cpu INTEGER AUTO_INCREMENT
	,porcentaje DECIMAL(4,2)
	,fecha DATETIME NULL
	,CONSTRAINT pk_cpu_info PRIMARY KEY(id_cpu)
)


CREATE TABLE ram_info
(
	id_ram INTEGER AUTO_INCREMENT
	,porcentaje DECIMAL(4,2)
	,fecha DATETIME NULL
	,CONSTRAINT pk_ram_info PRIMARY KEY(id_ram)
)


SELECT * FROM cpu_info
SELECT * FROM ram_info