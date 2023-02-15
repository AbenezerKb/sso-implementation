-- name: GetKafkaOffset :one
SELECT offset_val
from kafka_offsets
where true
LIMIT 1;
-- name: SetKafkaOffset :exec
UPDATE kafka_offsets
SET offset_val = $1
WHERE true;