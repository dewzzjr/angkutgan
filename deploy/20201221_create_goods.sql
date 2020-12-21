CREATE TABLE IF NOT EXISTS items (
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    unit VARCHAR(255) NOT NULL,
    PRIMARY KEY ( code )
);

CREATE TABLE IF NOT EXISTS price_sell (
    id INT NOT NULL AUTO_INCREMENT,
    code VARCHAR(20) NOT NULL,
    value INT NOT NULL,
    PRIMARY KEY ( id ),
    FOREIGN KEY ( code ) REFERENCES items( code )
);

CREATE TABLE IF NOT EXISTS price_rent (
    id INT NOT NULL AUTO_INCREMENT,
    duration INT NOT NULL,
    time_unit INT NOT NULL,
    code VARCHAR(20) NOT NULL,
    value INT NOT NULL,
    PRIMARY KEY ( id ),
    FOREIGN KEY ( code ) REFERENCES items( code )
);