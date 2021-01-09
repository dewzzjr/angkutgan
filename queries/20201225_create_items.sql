CREATE TABLE IF NOT EXISTS items (
    code VARCHAR(20) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    unit VARCHAR(255) NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP,
    modified_by INT,
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS price_sell (
    code VARCHAR(20) NOT NULL PRIMARY KEY,
    value INT NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP,
    modified_by INT,
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL,
    FOREIGN KEY ( code ) REFERENCES items( code )
);

CREATE TABLE IF NOT EXISTS price_rent (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    description VARCHAR(255) NOT NULL,
    duration INT NOT NULL,
    time_unit INT NOT NULL,
    code VARCHAR(20) NOT NULL REFERENCES items( code ),
    value INT NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by INT,
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL,
    FOREIGN KEY ( code ) REFERENCES items( code )
);
