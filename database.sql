/*Tabel User_credential*/ -- fixed
CREATE TABLE user_credential
(
  id VARCHAR(225) PRIMARY KEY NOT NULL,
  username VARCHAR(100) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(225) NOT NULL,
  role VARCHAR(50) NOT NULL,
  virtual_account_number VARCHAR(225) DEFAULT '',
  is_active BOOLEAN
);

-- fixed
INSERT INTO user_credential (id, username, email, password, role, is_active) VALUES ('456', 'admin', 'compani.mail.yo', '$2a$10$FTqRPKh1IrHzvzi1YbhTbOY0pk.zQPAnh7OxJxK7D4YEih2GG2DqK','admin', true);
INSERT INTO user_credential (id, username, email, password, role, is_active) VALUES ('taxacc', 'tax', 'compani.tax.yo', '$2a$10$FTqRPKh1IrHzvzi1YbhTbOY0pk.zQPAnh7OxJxK7D4YEih2GG2DqK','admin', true);

-- Account Table
CREATE TABLE biodata (
    id VARCHAR(55) PRIMARY KEY NOT NULL,
    user_credential_id VARCHAR(55) NOT NULL,
    full_name VARCHAR(255),
    nik VARCHAR(20),
    phone_number VARCHAR(20),
    occupation VARCHAR(255),
    place_of_birth VARCHAR(255),
    date_of_birth DATE,
    postal_code VARCHAR(10),
    is_eglible BOOLEAN,
	status_update BOOLEAN,
	additional_information TEXT NULL DEFAULT 'biodata is not updated',
    FOREIGN KEY (user_credential_id) REFERENCES public.user_credential (id)
);

-- Account Table
CREATE TABLE saldo (
    id VARCHAR(55) PRIMARY KEY NOT NULL,
    user_credential_id VARCHAR(55) NOT NULL,
    total_saving INT,
    FOREIGN KEY (user_credential_id) REFERENCES public.user_credential (id)
);
-- fixed
INSERT INTO saldo (id, user_credential_id, total_saving) VALUES ('789', '456', 100000000);
INSERT INTO saldo (id, user_credential_id, total_saving) VALUES ('tax', 'taxacc', 0);

-- Deposit Interest Table
CREATE TABLE deposit_interest (
    id VARCHAR(55) PRIMARY KEY NOT NULL,
    created_date DATE,
    interest_rate DECIMAL,
	tax_rate DECIMAL,
    duration_mounth INT NOT NULL
);

-- Top Up Table
CREATE TABLE top_up (
    id VARCHAR(55) PRIMARY KEY NOT NULL,
    user_credential_id VARCHAR(55) NOT NULL,
    top_up_amount INT,
    maturity_time TIMESTAMP,
	accepted_time TIMESTAMP,
    accepted_status BOOLEAN,
    status_information TEXT, -- accepted/waiting/canceled
    transfer_confirmation_recipt BOOLEAN,
	recipt_file VARCHAR(55),
    FOREIGN KEY (user_credential_id) REFERENCES user_credential (id)
);

CREATE TABLE deposit (
    id VARCHAR(55) PRIMARY KEY NOT NULL,
    user_credential_id VARCHAR(55) NOT NULL,
    deposit_amount INT,
    interest_rate DECIMAL,
    tax_rate DECIMAL,
    duration int,
    created_date DATE,
    maturity_date DATE,
    status BOOLEAN,
    gross_profit int,
    tax int,
    net_profit int,
    total_return int,
    FOREIGN KEY (user_credential_id) REFERENCES user_credential (id)
);

-- application cost Table
CREATE TABLE application_handling_cost (
    id VARCHAR(55) PRIMARY KEY NOT NULL,
    name VARCHAR(100) NOT NULL,
    nominal DECIMAL NOT  NULL,
    unit VARCHAR(100) NOT NULL
);

-- Loan duration Table
CREATE TABLE loan_interest (
    id VARCHAR(100) PRIMARY KEY NOT NULL,
    duration_months INT NOT NULL,
    loan_interest_rate DECIMAL NOT NULL
);

-- Loan Table
CREATE TABLE loan (
    id VARCHAR(55) PRIMARY KEY NOT NULL,
    user_credential_id VARCHAR(55) NOT NULL,
    loan_amount INT,
    loan_duration INT,
    loan_interest_rate DECIMAL NOT NULL,
	loan_interest_nominal INT,
    total_amount_of_dept INT NOT NULL,
    application_handling_cost_nominal INT NOT NULL,
    application_handling_cost_unit VARCHAR(55) NULL,
    loan_date_created DATE,
    is_payed BOOLEAN,
    status TEXT,
    FOREIGN KEY (user_credential_id) REFERENCES user_credential (id)
);

-- Installenment Loan
CREATE TABLE installenment_loan (
    id VARCHAR(55) PRIMARY KEY NOT NULL,
    loan_id VARCHAR(55) NOT NULL,
    is_payed BOOLEAN,
    payment_installenment_cost INT,
    payment_deadline DATE,
    total_amount_of_dept INT NOT NULL,
    late_payment_fee_nominal DECIMAL,
    late_payment_fee_unit VARCHAR(55),
    late_payment_fee_day INT,
    late_payment_fee_total INT,
    payment_date DATE,
    status TEXT,
    transfer_confirmation_recipt BOOLEAN,
	recipt_file VARCHAR(55),
    FOREIGN KEY (loan_id) REFERENCES loan (id)
);

CREATE TABLE late_payment_fee (
    id VARCHAR(55) PRIMARY KEY NOT NULL,
    name VARCHAR(100) NOT NULL,
    nominal DECIMAL NOT NULL,
    unit VARCHAR(100) NOT NULL
);


DROP TABLE saldo;
DROP TABLE biodata;
DROP TABLE user_credential;
DROP TABLE deposit_interest;
DROP TABLE top_up;
DROP TABLE deposit;
DROP TABLE application_handling_cost;
DROP TABLE loan_interest;
DROP TABLE installenment_loan;
DROP TABLE loan;
DROP TABLE late_payment_fee;
