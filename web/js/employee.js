function setMessage(id, text, type) {
    const element = document.getElementById(id);
    if (!element) return;
    element.textContent = text;
    element.className = `message ${type || ""}`;
}

async function readResponseMessage(response) {
    const text = await response.text();
    try {
        const data = JSON.parse(text);
        return data.message || text;
    } catch (_) {
        return text || "Wystąpił błąd";
    }
}

async function loadMedicines() {
    const token = localStorage.getItem("jwt");
    const box = document.getElementById("medicines");
    box.innerHTML = "Ładowanie...";

    const response = await fetch("/api/employee/medicines", {
        headers: { Authorization: `Bearer ${token}` }
    });

    if (!response.ok) {
        box.innerHTML = `<p class="error">Nie udało się pobrać leków</p>`;
        return;
    }

    const medicines = await response.json();
    box.innerHTML = "";

    medicines.forEach(medicine => {
        box.innerHTML += `
            <div class="card">
                <h3>${medicine.name}</h3>
                <p>ID: ${medicine.id}</p>
                <p>Cena: ${medicine.price} PLN</p>
                <p>Stan: ${medicine.quantity ?? 0}</p>
                <p>${medicine.requires ? "Wymaga recepty" : "Bez recepty"}</p>
            </div>
        `;
    });
}

async function addStock() {
    const token = localStorage.getItem("jwt");
    const medicineID = Number(document.getElementById("stockMedicineId").value);
    const quantity = Number(document.getElementById("stockQuantity").value);

    const response = await fetch("/api/employee/stock", {
        method: "POST",
        headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ medicine_id: medicineID, quantity })
    });

    const message = await readResponseMessage(response);
    if (!response.ok) {
        setMessage("stockMessage", message, "error");
        return;
    }

    setMessage("stockMessage", "Stan leku został zwiększony", "success");
    await loadMedicines();
}

async function addMedicine() {
    const token = localStorage.getItem("jwt");
    const name = document.getElementById("newMedicineName").value;
    const price = Number(document.getElementById("newMedicinePrice").value);
    const quantity = Number(document.getElementById("newMedicineQuantity").value);
    const requiresPrescription = document.getElementById("newMedicineRequires").value === "true";

    const response = await fetch("/api/employee/medicine", {
        method: "POST",
        headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            name,
            price,
            quantity,
            requires_prescription: requiresPrescription
        })
    });

    const message = await readResponseMessage(response);
    if (!response.ok) {
        setMessage("addMedicineMessage", message, "error");
        return;
    }

    setMessage("addMedicineMessage", "Nowy lek został dodany", "success");
    await loadMedicines();
}

async function sellMedicine() {
    const token = localStorage.getItem("jwt");
    const medicineID = Number(document.getElementById("sellMedicineId").value);
    const quantity = Number(document.getElementById("sellQuantity").value);
    const patientValue = document.getElementById("sellPatient").value;
    const patientID = patientValue ? Number(patientValue) : null;
    const prescriptionCode = document.getElementById("sellPrescription").value.trim();

    const response = await fetch("/api/employee/sell", {
        method: "POST",
        headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            medicine_id: medicineID,
            quantity,
            patient_id: patientID,
            prescription_code: prescriptionCode || null
        })
    });

    const message = await readResponseMessage(response);
    if (!response.ok) {
        setMessage("sellMessage", message, "error");
        return;
    }

    setMessage("sellMessage", "Sprzedaż zakończona", "success");
    await loadMedicines();
}

window.loadMedicines = loadMedicines;
window.addStock = addStock;
window.addMedicine = addMedicine;
window.sellMedicine = sellMedicine;
