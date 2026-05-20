BEGIN ;


 Alter TABLE if exists teams ADD  COLUMN created_by uuid references users(id);

COMMIT ;