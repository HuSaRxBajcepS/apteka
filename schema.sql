CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(30) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS medicines(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    requires_prescription BOOLEAN NOT NULL
    DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS stock(
    medicine_id INT PRIMARY KEY REFERENCES medicines(id)
    ON DELETE CASCADE,
    quantity INT NOT NULL
);

CREATE TABLE IF NOT EXISTS prescriptions(
    id SERIAL PRIMARY KEY,
    doctor_id INT REFERENCES users(id),
    patient_id INT REFERENCES users(id),
    code CHAR(4) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS prescription_items(
    id SERIAL PRIMARY KEY,
    prescription_id INT REFERENCES prescriptions(id)
    ON DELETE CASCADE,
    medicine_id INT REFERENCES medicines(id),
    quantity INT NOT NULL
);

CREATE TABLE IF NOT EXISTS prescription_status(
    prescription_id INT PRIMARY KEY REFERENCES prescriptions(id)
    ON DELETE CASCADE,
    completed BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS sales(
    id SERIAL PRIMARY KEY,
    employee_id INT REFERENCES users(id),
    patient_id INT REFERENCES users(id),
    prescription_id INT REFERENCES prescriptions(id),
    medicine_id INT REFERENCES medicines(id),
    quantity INT NOT NULL,
    unit_price NUMERIC(10,2) NOT NULL,
    total_price NUMERIC(10,2) NOT NULL,
    sold_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs(
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    action TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sales_patient
ON sales(patient_id);

CREATE INDEX idx_sales_date
ON sales(sold_at);

CREATE INDEX idx_prescription_code
ON prescriptions(code);