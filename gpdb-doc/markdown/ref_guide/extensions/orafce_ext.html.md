# orafce_ext

`orafce_ext` is an extension for [orafce](https://github.com/orafce/orafce). 
It provides the external capability to support the Oracle interface and the `UTL_RAW` utility.

## Get started

Enable `orafce` and `orafce_ext` in each desired database.

```
CREATE EXTENSION orafce;
CREATE EXTENSION orafce_ext;
```

## UTL_RAW

`UTL_RAW` supports the `RAW` datatype, enabling users to store and manipulate binary data.

>**Important** The `max_length` for the `RAW` dataype is `1GB` in `orafce_ext`, compared to `32767` in Oracle.

All functions in `UTL_RAW` fall under the `utl_raw` namespace:

- `utl_raw.cast_to_raw(c IN VARCHAR2)`
- `utl_raw.cast_to_varchar2(r IN RAW)`
- `utl_raw.concat(r1 IN RAW, r2 IN RAW, r3 IN RAW, ...)`
- `utl_raw.convert(r IN RAW, to_charset IN VARCHAR2, from_charset IN VARCHAR2)`
- `utl_raw.length(r IN RAW)`
- `utl_raw.substr(r IN RAW, pos IN INTEGER, len IN INTEGER)`

### cast_to_raw

This function converts a `VARCHAR2` value into a `RAW` value.

>**Note** This example applies to `orafce` v4.9. For Greenplum Database 6, use `varchar2` instead of `oracle.varchar2`.

```
testdb=# SELECT utl_raw.cast_to_raw('abc'::oracle.varchar2);
 cast_to_raw 
-------------
 \x616263
(1 row)
```

### cast_to_varchar2

This function converts a `RAW` value into a `VARCHAR2` value.

```
testdb=# SELECT utl_raw.cast_to_varchar2('abc'::raw);
 cast_to_varchar2 
------------------
 abc
(1 row)
```

### concat

This function concatenates multiple `RAW` values into a single `RAW` value.

```
testdb=# SELECT utl_raw.concat('a'::raw, 'b'::raw, 'c'::raw);
  concat  
----------
 \x616263
(1 row)
```

### convert

This function converts `RAW r` from character set `from_charset` to character set `to_charset` and returns the resulting `RAW`.

```
testdb=# select 'abc'::raw;
      raw       
----------------
 \xe697a0e6958c
(1 row)

testdb=# SELECT utl_raw.convert('abc'::raw, 'GBK', 'UTF-8');
  convert   
------------
 \xcedeb5d0
(1 row)
```

### length

This function returns the length in bytes of a `RAW`.

```
testdb=# select utl_raw.length('abc'::raw);
 length 
--------
      6
(1 row)
```

### substr

This function returns the substring of a `RAW`.

```
utl_raw.substr(r IN RAW, pos IN INTEGER, len IN INTEGER)
```

Where:

- The value of `pos` cannot be `0` and cannot exceed `length(r)`. 

    - If `pos` is positive, `SUBSTR` counts from the beginning of `r` to find the first byte. 
    - If `pos` is negative, `SUBSTR` counts backward from the end of `r`.

- The value of `len` cannot be less than `1` and cannot exceed `length(r) - (pos - 1)`. 

    - If `len` is omitted, `SUBSTR` returns all bytes to the end of `r`. 

```
testdb=# SELECT utl_raw.substr('abc'::raw, 1, 2);
 substr 
--------
 \x6162
(1 row)
```

To return the substring without length:

```
testdb=# SELECT utl_raw.substr('abc'::raw, 1);
  substr  
----------
 \x616263
(1 row)
```
