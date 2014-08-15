/*created_at

https://bigquery.cloud.google.com/

*/

SELECT 
  LEFT(repository_pushed_at,7) monthx
, repository_language   
, CEIL( count(*)/1000) Tausend
FROM githubarchive:github.timeline
where 1=1
	AND  LEFT(repository_pushed_at,7) >= '2011-01'
	AND  repository_language in ('Go','go','Golang','golang','C','Java','PHP','JavaScript','C++','Python','Ruby')
	AND  type="PushEvent"
group by monthx, repository_language
order by monthx, repository_language
;


SELECT 
  repository_language   
, LEFT(repository_pushed_at,7) monthx
, CEIL( count(*)/1000) Tausend
FROM githubarchive:github.timeline
where 1=1
	AND  LEFT(repository_pushed_at,7) >= '2011-01'
	AND  repository_language in ('Go','go','Golang','golang','C','Java','PHP','JavaScript','C++','Python','Ruby')
	AND  type="PushEvent"
group by monthx, repository_language
order by repository_language   , monthx
;
