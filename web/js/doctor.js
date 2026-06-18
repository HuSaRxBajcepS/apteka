async function createPrescription() {
    const token = localStorage.getItem("jwt");
    const patientID = Number(document.getElementById("patientId").value);
    const medicineID = Number(document.getElementById("medicineId").value);
    const quantity = Number(document.getElementById("quantity").value);
    const message = document.getElementById("createdPrescription");

    const response = await fetch("/api/doctor/create", {
        method: "POST",
        headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            patient_id: patientID,
            medicines: [
                {
                    medicine_id: medicineID,
                    quantity: quantity
                }
            ]
        })
    });

    if (!response.ok) {
        message.textContent = await response.text() || "Błąd tworzenia recepty";
        message.className = "message error";
        return;
    }

    const result = await response.json();
    message.innerHTML = `Recepta utworzona. ID: <strong>${result.id}</strong>, kod: <strong>${result.code}</strong>`;
    message.className = "message success";
}

window.createPrescription = createPrescription;
