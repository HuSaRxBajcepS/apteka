async function createPrescription() {
    const token = localStorage.getItem("jwt");
    const patientID = Number(document.getElementById("patientId").value);
    const medicineID = Number(document.getElementById("medicineId").value);
    const quantity = Number(document.getElementById("quantity").value);
    const response = await fetch("/api/doctor/create",{
                method: "POST",
                headers: {
                    Authorization:
                        `Bearer ${token}`,
                    "Content-Type":
                        "application/json"
                },
                body: JSON.stringify({
                    patient_id:
                        patientID,
                    medicines: [
                        {
                            medicine_id:
                                medicineID,
                            quantity:
                                quantity
                        }
                    ]
                })
            }
        );

    if (!response.ok) {
        alert("Błąd tworzenia recepty");
        return;
    }

    const result = await response.json();
    alert(
`Recepta utworzona

ID:
${result.id}

Kod:
${result.code}`
    );

}