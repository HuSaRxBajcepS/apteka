async function sellMedicine() {
    const token = localStorage.getItem("jwt");
    const medicineID = Number(document.getElementById("sellMedicineId").value);
    const quantity =Number(document.getElementById("sellQuantity").value);
    const patientValue =document.getElementById("sellPatient").value;
    const patientID =patientValue ? Number(patientValue): null;
    const prescriptionCode = document.getElementById("sellPrescription").value;
    const body = {
        medicine_id:
            medicineID,
        quantity,
        patient_id:
            patientID,
        prescription_code:
            prescriptionCode || null
    };
    const response =await fetch("/api/employee/sell",{
                method:"POST",
                headers:{

                    Authorization:
                    `Bearer ${token}`,

                    "Content-Type":
                    "application/json"

                },

                body:
                JSON.stringify(body)

            }
        );

    const data =await response.json();
    if (!response.ok) {
        alert(
            data.message ||
            "Błąd sprzedaży"
        );
        return;
    }
    alert("Sprzedaż zakończona");
}

async function medicines() {
    const token =localStorage.getItem("jwt");
    const response =await fetch("/api/employee/medicines",{
                headers: {

                    Authorization:
                    `Bearer ${token}`

                }
            }
        );
    return await response.json();
}