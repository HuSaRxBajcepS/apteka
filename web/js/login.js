function panelForRole(role) {
    const panels = {
        patient: "/patient.html",
        employee: "/employee.html",
        doctor: "/doctor.html"
    };
    return panels[role] || "/";
}

function requireRole(expectedRole) {
    const token = localStorage.getItem("jwt");
    const role = localStorage.getItem("role");

    if (!token || role !== expectedRole) {
        window.location.href = "/";
    }
}

async function login() {
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;

    const response = await fetch("/api/login", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({ email, password })
    });

    if (!response.ok) {
        document.getElementById("loginMessage").textContent = "Nieprawidłowe dane logowania";
        return;
    }

    const data = await response.json();
    localStorage.setItem("jwt", data.token);
    localStorage.setItem("role", data.role);
    localStorage.setItem("userID", data.user_id);

    window.location.href = panelForRole(data.role);
}

function logout() {
    localStorage.clear();
    window.location.href = "/";
}

async function register() {
    const fullName = document.getElementById("registerFullName").value;
    const email = document.getElementById("registerEmail").value;
    const password = document.getElementById("registerPassword").value;
    const role = document.getElementById("registerRole").value;

    const response = await fetch("/api/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            full_name: fullName,
            email: email,
            password: password,
            role: role
        })
    });

    if (!response.ok) {
        document.getElementById("registerMessage").textContent =
            "Nie udało się zarejestrować";
        return;
    }

    document.getElementById("registerMessage").textContent =
        "Konto utworzone pomyślnie";
}

function redirectLoggedUser() {
    const token = localStorage.getItem("jwt");
    const role = localStorage.getItem("role");
    if (token && role) {
        window.location.href = panelForRole(role);
    }
}
window.login = login;
window.register = register;
window.logout = logout;
window.requireRole = requireRole;
window.redirectLoggedUser = redirectLoggedUser;
