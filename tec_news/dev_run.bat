start  "tec-news"  dev_appserver.py --smtp_host=127.0.0.1 --smtp_port=25 ^
	--smtp_user=user_from_smtp_auth_file --smtp_password=pass_from_smtp_auth_file ^
	--port 8087  ./mod01.yaml