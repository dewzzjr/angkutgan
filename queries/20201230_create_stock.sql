CREATE TABLE IF NOT EXISTS stock (
    code VARCHAR(20) NOT NULL PRIMARY KEY,
    asset INT NOT NULL,
    inventory INT NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP,
    modified_by INT,
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL,
    FOREIGN KEY ( code ) REFERENCES items( code )
);