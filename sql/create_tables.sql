CREATE TABLE IF NOT EXISTS vectors (
    vector_id varchar(250) NOT NULL,
    database_id varchar(250) NOT NULL,
    value varchar(250) NOT NULL,
    vector varchar(250) NOT NULL,
    metadata varchar(250),
    PRIMARY KEY(vector_id)
)