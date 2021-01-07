BEGIN;

    insert into cities
        (cityname,website,mayor)
    SELECT md5(random()
    ::text) AS cityname,md5
    (random
    ()::text) AS website, md5
    (random
    ()::text) AS mayor
from generate_series
    (1,10000);


COMMIT;