CREATE KEYSPACE IF NOT EXISTS shortener WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TABLE IF NOT EXISTS shortener.urls
    (short text, 
    long text, 
    created_at timestamp, 
    expires_at timestamp,
    last_visited timestamp, 
    PRIMARY KEY(short, created_at))
    WITH CLUSTERING ORDER BY (created_at DESC);

CREATE  INDEX IF NOT EXISTS ON shortener.urls (long); 


