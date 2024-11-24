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

Dans le container docker oracle, création d'un user **cdc_user**.

*Notes* : 

 - You must have the **EXECUTE_CATALOG_ROLE** role to use the **DBMS_LOGMNR_D** package.
 - You must have the **EXECUTE_CATALOG_ROLE** role to use the **DBMS_LOGMNR** package.
 
* https://docs.oracle.com/en/database/oracle/oracle-database/23/arpls/DBMS_LOGMNR.html#GUID-41730EFC-C6CA-423E-834B-3E0E643346C3

```bash
sh-4.4$ sqlplus "/ as sysdba"

# New user
CREATE USER cdc_user IDENTIFIED BY password;

GRANT CONNECT, RESOURCE, DBA, EXECUTE_CATALOG_ROLE TO cdc_user;
GRANT UNLIMITED TABLESPACE TO cdc_user;

GRANT SELECT, INSERT, UPDATE, DELETE ON schema.demo TO cdc_user;

GRANT SELECT ANY TRANSACTION TO cdc_user;
GRANT LOGMINING TO cdc_user;
GRANT CREATE SESSION TO cdc_user;

```

* https://www.atlassian.com/data/admin/how-to-create-a-user-and-grant-permissions-in-oracle


## References

* https://www.oracle.com/database/sqldeveloper/technologies/download/
* https://github.com/visulate/visulate-ora2pg
* https://oracle-base.com/articles/8i/logminer
* https://docs.oracle.com/en/database/oracle/oracle-database/18/sutil/oracle-logminer-utility.html
* https://mbouayoun.developpez.com/journaux/


* https://dbaoraclesql.canalblog.com/archives/2021/06/03/38999757.html


* https://mathiaszarick.wordpress.com/2024/05/24/oracle-to-postgresql-replication-using-debezium-and-platys/
* https://github.com/TrivadisPF/platys/tree/master