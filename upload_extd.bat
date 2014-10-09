REM appcfg is an inclusion of goapp.bat
REM 
REM appcfg.py --email=peter.buchmann@web.de update_indexes .
REM appcfg.py --email=peter.buchmann@web.de update_cron    .
REM appcfg.py --email=peter.buchmann@web.de update         .
    appcfg.py --email=peter.buchmann@web.de vacuum_indexes .
    appcfg.py --email=peter.buchmann@web.de update   .
