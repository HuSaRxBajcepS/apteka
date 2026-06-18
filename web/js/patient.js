function formatDate(value) {
    if (!value) return "-";
    return new Date(value).toLocaleDateString("pl-PL");
}

async function loadPatientPanel() {
    await loadAccount();
    await loadMyPrescriptions();
}

async function loadAccount() {
    const token = localStorage.getItem("jwt");
    const box = document.getElementById("account");

    const response = await fetch("/api/patient/me", {
        headers: { Authorization: `Bearer ${token}` }
    });

    if (!response.ok) {
        box.textContent = "Nie udało się pobrać danych konta";
        box.className = "error";
        return;
    }

    const patient = await response.json();
    box.innerHTML = `
        <div class="card">
            <h3>${patient.full_name}</h3>
            <p>Email: ${patient.email}</p>
            <p>ID pacjenta: ${patient.id}</p>
        </div>
    `;
}

async function loadMyPrescriptions() {
    const token = localStorage.getItem("jwt");
    const box = document.getElementById("prescriptions");
    box.innerHTML = "Ładowanie...";

    const response = await fetch("/api/patient/prescriptions", {
        headers: { Authorization: `Bearer ${token}` }
    });

    if (!response.ok) {
        box.innerHTML = `<p class="error">Nie udało się pobrać recept</p>`;
        return;
    }

    const prescriptions = await response.json();
    box.innerHTML = "";

    if (!prescriptions || prescriptions.length === 0) {
        box.innerHTML = `<p>Brak recept.</p>`;
        return;
    }

    prescriptions.forEach(prescription => {
        const medicines = (prescription.medicines || [])
            .map(medicine => `<li>${medicine.name} — ilość: ${medicine.quantity}</li>`)
            .join("");

        box.innerHTML += `
            <div class="card">
                <h3>Recepta ${prescription.code}</h3>
                <p>ID: ${prescription.id}</p>
                <p>Wystawiona: ${formatDate(prescription.created_at)}</p>
                <p>Ważna do: ${formatDate(prescription.expires_at)}</p>
                <p>Status: ${prescription.completed ? "Zrealizowana" : "Aktywna"}</p>
                <ul>${medicines}</ul>
            </div>
        `;
    });
}

window.loadPatientPanel = loadPatientPanel;
window.loadMyPrescriptions = loadMyPrescriptions;
