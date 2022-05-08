#!/bin/bash
set -e

clickhouse client -n <<-EOSQL
	CREATE DATABASE sg;
	CREATE TABLE sg.events (
        device_id    String,
        device_os    String,
        session      String,
        sequence     Int32,
        event        String,
        param_int    Int32,
        param_str    String,
        ip           String,
        client_time  DateTime,
        server_time  DateTime
    ) ENGINE = Log;
EOSQL