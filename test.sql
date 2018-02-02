-- define:stmt1
SELECT *
FROM a
WHERE b=?
--end

-- define:stmt2
SELECT a, b, c
FROM table1
JOIN table2 ON a=b
ORDER BY c DESC
-- end