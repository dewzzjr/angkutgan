CREATE TABLE IF NOT EXISTS users (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    PRIMARY KEY ( id )
);

CREATE TABLE IF NOT EXISTS user_info (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
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
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL,
    FOREIGN KEY ( user_id ) REFERENCES users( id )
);

CREATE TABLE IF NOT EXISTS access (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(20) NOT NULL UNIQUE,
    description VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_access (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    access_id INT NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ( id ),
    FOREIGN KEY ( user_id ) REFERENCES users( id ),
    FOREIGN KEY ( access_id ) REFERENCES access( id ),
    UNIQUE ( user_id, access_id )
);

REPLACE INTO users ( id, username, password, status ) VALUES ( 0, 'admin', '$2a$10$jLquAzQsf3izuCFadOGspen0H9gzEEj/4m5INfNgzWawJwvIDrvHC', 'ACTIVE' );
REPLACE INTO user_info ( id, user_id, nik, ktp, fullname, address, phone, birthdate, religion ) VALUES ( 0, 0, '', '', 'Admin', '', '', CURRENT_DATE, '' );
