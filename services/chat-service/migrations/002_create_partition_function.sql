CREATE OR REPLACE FUNCTION create_chat_messages_partition()
RETURNS void AS $$
DECLARE
    partition_date TEXT;
    partition_name TEXT;
    start_date TEXT;
    end_date TEXT;
BEGIN
    partition_date := to_char(NOW() + INTERVAL '1 month', 'YYYY_MM');
    partition_name := 'chat_messages_' || REPLACE(partition_date, '_', '_');
    start_date := to_char(NOW() + INTERVAL '1 month', 'YYYY-MM-01');
    end_date := to_char(NOW() + INTERVAL '2 months', 'YYYY-MM-01');

    IF NOT EXISTS (
        SELECT 1
        FROM pg_class
        WHERE relname = partition_name
    ) THEN
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF chat_messages FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
    END IF;
END;
$$ LANGUAGE plpgsql;
