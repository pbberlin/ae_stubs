REM goapp is a wrapper around appcfg.py
REM appcfg is an inclusion of goapp.bat
REM goapp.bat   deploy
REM appcfg.py --email=peter.buchmann@web.de appcfg.py rollback .
REM appcfg.py --email=peter.buchmann@web.de update_indexes .
REM appcfg.py --email=peter.buchmann@web.de update_cron    .
REM appcfg.py --email=peter.buchmann@web.de update         .
    appcfg.py --email=peter.buchmann@web.de update   mod01/mod01.yaml mod02/mod02.yaml
    appcfg.py --email=peter.buchmann@web.de update_indexes   mod01  -A libertarian-islands
    appcfg.py --email=peter.buchmann@web.de vacuum_indexes   mod01  -A libertarian-islands
pause