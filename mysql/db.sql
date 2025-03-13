CREATE TABLE investigator (
    id INT PRIMARY KEY,
    first_name VARCHAR(64),
    last_name VARCHAR(64),
    email_address VARCHAR(128),
    start_date DATE,
    end_date DATE,
    role_code VARCHAR(64),
    nsf_id VARCHAR(64)
);

CREATE TABLE program_element (
    id INT PRIMARY KEY,
    code VARCHAR(32),
    text TEXT
);

CREATE TABLE program_reference (
    id INT PRIMARY KEY,
    code VARCHAR(32),
    text TEXT
);

CREATE TABLE foa_information (
    id INT PRIMARY KEY,
    code VARCHAR(32),
    name VARCHAR(64)
);

CREATE TABLE institution (
    id INT PRIMARY KEY,
    name VARCHAR(128),
    city_name VARCHAR(64),
    zip_code VARCHAR(32),
    phone_number VARCHAR(16),
    street_address VARCHAR(128),
    country_name VARCHAR(64),
    street_name VARCHAR(64),
    state_code VARCHAR(16)
);

CREATE TABLE directorate (
    id INT PRIMARY KEY,
    long_name TEXT,
    abbreviation VARCHAR(64)
);

CREATE TABLE division (
    id INT PRIMARY KEY,
    long_name TEXT,
    abbreviation VARCHAR(64)
);

CREATE TABLE organization (
    id INT PRIMARY KEY,
    code VARCHAR(64),
    directorate_id INT,
    division_id INT,
    FOREIGN KEY (directorate_id) REFERENCES directorate(id),
    FOREIGN KEY (division_id) REFERENCES division(id)
);

CREATE TABLE program_officer (
    id INT PRIMARY KEY,
    sign_block_name VARCHAR(128),
    email VARCHAR(128),
    phone VARCHAR(16)
);

CREATE TABLE award_instrument (
    id INT PRIMARY KEY,
    value VARCHAR(128)
);

CREATE TABLE award (
    id INT PRIMARY KEY,
    title VARCHAR(128),
    award_effective_date DATE,
    award_expiration_date DATE,
    award_amount INT,
    program_officer_id INT,
    organization_id INT,
    award_instrument_id INT,
    abstract_text TEXT,
    FOREIGN KEY (program_officer_id) REFERENCES program_officer(id),
    FOREIGN KEY (organization_id) REFERENCES organization(id),
    FOREIGN KEY (award_instrument_id) REFERENCES award_instrument(id)
);

CREATE TABLE award_foa_information (
    foa_information_id INT,
    award_id INT,
    FOREIGN KEY (foa_information_id) REFERENCES foa_information(id),
    FOREIGN KEY (award_id) REFERENCES award(id)
);

CREATE TABLE award_institution (
    institution_id INT,
    award_id INT,
    FOREIGN KEY (institution_id) REFERENCES institution(id),
    FOREIGN KEY (award_id) REFERENCES award(id)
);

CREATE TABLE award_investigator (
    investigator_id INT,
    award_id INT,
    FOREIGN KEY (investigator_id) REFERENCES investigator(id),
    FOREIGN KEY (award_id) REFERENCES award(id)
);

CREATE TABLE award_program_element (
    program_element_id INT,
    award_id INT,
    FOREIGN KEY (program_element_id) REFERENCES program_element(id),
    FOREIGN KEY (award_id) REFERENCES award(id)
);

CREATE TABLE award_program_reference (
    program_reference_id INT,
    award_id INT,
    FOREIGN KEY (program_reference_id) REFERENCES program_reference(id),
    FOREIGN KEY (award_id) REFERENCES award(id)
);
