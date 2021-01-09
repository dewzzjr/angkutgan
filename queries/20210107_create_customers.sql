CREATE TABLE IF NOT EXISTS customers (
    code VARCHAR(20) NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type INT NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    nik VARCHAR(255),
    role VARCHAR(255),
    group_name VARCHAR(255),
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_time TIMESTAMP,
    modified_by INT,
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS projects (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255) NOT NULL,
    create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by INT,
    FOREIGN KEY ( modified_by ) REFERENCES users( id ) ON DELETE SET NULL,
    FOREIGN KEY ( code ) REFERENCES customers( code )
);