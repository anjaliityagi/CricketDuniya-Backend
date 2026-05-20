BEGIN ;

Alter table if exists matches add column match_date timestamp;


COMMIT ;