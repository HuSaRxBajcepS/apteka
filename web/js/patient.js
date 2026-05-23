async function getPrescription() {
    const code = document.getElementById("code").value;
    const response = await fetch(`/api/patient/prescription?code=${code}`);
    if (!response.ok) {
        alert("Nie znaleziono recepty");
        return;
    }
    const medicines = await response.json();
    const box = document.getElementById("prescription");
    box.innerHTML = "";
    medicines.forEach(medicine => {
        box.innerHTML += `
            <div class="card">
                <h3>${medicine.name}</h3>

                <p>
                    Ilość:
                    ${medicine.quantity}
                </p>
            </div>
        `;
    });
}

async function loadOTC() {
    const response = await fetch("/api/patient/otc");
    const medicines = await response.json();
    const box = document.getElementById("otc");
    box.innerHTML = "";
    medicines.forEach(medicine => {
        box.innerHTML += `
            <div class="card">

                <h3>
                    ${medicine.name}
                </h3>

                <div class="price">
                    ${medicine.price} PLN
                </div>

            </div>
        `;
    });
}