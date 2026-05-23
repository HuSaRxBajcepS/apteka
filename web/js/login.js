async function login() {
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;
    const response = await fetch("/api/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email,
                password
            })
        });
    if (!response.ok) {
        alert("Nieprawidłowe dane logowania");
        return;
    }
    const data = await response.json();
    localStorage.setItem(
        "jwt",
        data.token
    );
    localStorage.setItem(
        "role",
        data.role
    );
    alert("Zalogowano");
    location.reload();
}

function logout() {
    localStorage.removeItem("jwt");
    localStorage.removeItem("role");
    location.reload();
}