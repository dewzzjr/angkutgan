CREATE TABLE IF NOT EXISTS users (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    PRIMARY KEY ( id )
);

ALTER TABLE items ADD create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE items ADD update_time TIMESTAMP;
ALTER TABLE items ADD modified_by INT;
ALTER TABLE items ADD FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL;

ALTER TABLE price_sell ADD create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE price_sell ADD update_time TIMESTAMP;
ALTER TABLE price_sell ADD modified_by INT;
ALTER TABLE price_sell ADD FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL;

ALTER TABLE price_rent ADD create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE price_rent ADD update_time TIMESTAMP;
ALTER TABLE price_rent ADD modified_by INT;
ALTER TABLE price_rent ADD FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS user_info (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    nik VARCHAR(50),
    ktp VARCHAR(50),
    fullname VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    birthdate DATE,
    religion VARCHAR(20),
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP,
    modified_by INT,
    PRIMARY KEY ( id ),
    FOREIGN KEY ( user_id ) REFERENCES users( id ),
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS access (
    id INT NOT NULL AUTO_INCREMENT,
    code VARCHAR(20) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL,
    PRIMARY KEY ( id )
);

CREATE TABLE IF NOT EXISTS user_access (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    access_id INT NOT NULL,
    PRIMARY KEY ( id ),
    FOREIGN KEY ( user_id ) REFERENCES users( id ),
    FOREIGN KEY ( access_id ) REFERENCES access( id ),
    UNIQUE ( user_id, access_id )
);

INSERT INTO users (	id, username, password, status ) VALUES ( 0, 'admin', '$2a$10$jLquAzQsf3izuCFadOGspen0H9gzEEj/4m5INfNgzWawJwvIDrvHC', 'ACTIVE' ) ON DUPLICATE KEY UPDATE;
INSERT INTO user_info ( id, user_id, nik, ktp, fullname, address, phone, birthdate, religion ) VALUES ( 0, 0, '', '', 'Admin', '', '', CURDATE(), '' ) ON DUPLICATE KEY UPDATE;