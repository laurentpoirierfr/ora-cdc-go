# ora-cdc-go


## Demonstration

```bash
make demo
```

### Visulate Ora2Pg - demo

```yaml
oracle_dsn : dbi:Oracle:host=oracle;service_name=FREEPDB1;port=1521
user       : demo
password   : demo
```

## Administration Oracle

Dans le container docker oracle 

```bash
sh-4.4$ sqlplus "/ as sysdba"

SQL*Plus: Release 23.0.0.0.0 - for Oracle Cloud and Engineered Systems on Sun Nov 24 11:25:37 2024
Version 23.5.0.24.07

Copyright (c) 1982, 2024, Oracle.  All rights reserved.


Connected to:
Oracle Database 23ai Free Release 23.0.0.0.0 - Develop, Learn, and Run for Free
Version 23.5.0.24.07

SQL> ARCHIVE LOG LIST;
Database log mode              No Archive Mode
Automatic archival             Disabled
Archive destination            /opt/oracle/product/23ai/dbhomeFree/dbs/arch
Oldest online log sequence     21
Current log sequence           20


desc v$logfile;
select * from v$logfile;
select * from V$LOG_HISTORY;

ALTER PLUGGABLE DATABASE ALL OPEN;
ALTER DATABASE ADD SUPPLEMENTAL LOG DATA;


DBMS_LOGMNR_D.build ( dictionary_filename => 'FREE-dict.ora', dictionary_location => '/opt/oracle/oradata/FREE'); 

DBMS_LOGMNR.start_logmnr ( dictfilename => '/opt/oracle/oradata/FREE/FREE-dict.ora' );


```

## References

* https://www.oracle.com/database/sqldeveloper/technologies/download/
* https://github.com/visulate/visulate-ora2pg
* https://oracle-base.com/articles/8i/logminer
* https://docs.oracle.com/en/database/oracle/oracle-database/18/sutil/oracle-logminer-utility.html
* https://mbouayoun.developpez.com/journaux/

* https://mathiaszarick.wordpress.com/2024/05/24/oracle-to-postgresql-replication-using-debezium-and-platys/
* https://github.com/TrivadisPF/platys/tree/master