CREATE FUNCTION GetUserByID(id uuid) RETURNS TABLE(id uuid, login varchar, password varchar, admin boolean)
    AS $$ select id, login, password, admin from users where id = $1 $$
    LANGUAGE SQL;

