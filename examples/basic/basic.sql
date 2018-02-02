-- define:query1
SELECT `name`, phone, email
FROM user
WHERE id=?
-- end

-- define:bar
SELECT a.name, b.name name2
FROM agent a
LEFT JOIN broker b ON a.id=b.id
ORDER BY created_at DESC
LIMIT 100
-- end