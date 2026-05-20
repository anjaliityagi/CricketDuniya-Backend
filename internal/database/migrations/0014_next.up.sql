BEGIN ;
Alter table if exists matches add column team_a_id uuid references teams(id);
Alter table if exists matches add column team_b_id uuid references teams(id);


COMMIT ;