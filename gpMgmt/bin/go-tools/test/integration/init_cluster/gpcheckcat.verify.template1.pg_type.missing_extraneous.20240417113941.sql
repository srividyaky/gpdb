
          SELECT *
          FROM (
               SELECT relname, oid FROM pg_class WHERE reltype IN ()
               UNION ALL
               SELECT relname, oid FROM gp_dist_random('pg_class') WHERE reltype IN ()
          ) alltyprelids
          GROUP BY relname, oid ORDER BY count(*) desc
    
