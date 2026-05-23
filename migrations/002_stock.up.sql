INSERT INTO medicines(name, price, requires_prescription)
VALUES ('Paracetamol', 12.99, false),
('Ibuprofen', 15.50, false),
('Amoxicillin', 35.00, true);

INSERT INTO stock(medicine_id, quantity)
VALUES(1, 150),
(2, 120),
(3, 40);