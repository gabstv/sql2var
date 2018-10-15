-- define:base
SELECT
  u.name, u.score
FROM user u
-- end

-- define: sqluser1
-- sql: base
WHERE u.id=?
-- end

-- define: sqluser2
-- sql: base
WHERE u.age > 40
ORDER BY u.id ASC
-- end