async function registerRequest() {
    const password = document.getElementById("Password").value;
    const confirmPassword = document.getElementById("confirm_password").value;
    const errorMessageElement = document.getElementById("errorMessage");
    document.getElementById("errorMessage").style.display = "none";

    // Validate password match
    if (password !== confirmPassword) {
        document.getElementById("errorMessage").style.display = "block";
        errorMessageElement.innerHTML = "Passwords do not match.";
        return;
    } else {
        document.getElementById("errorMessage").style.display = "none";
    }

    const registerData = {
        username: document.getElementById("Username").value,
        password: password
    };

    const response = await fetch('/api/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(registerData)
    })
    

    if (response.status == 201) {
        window.location.href = '/login.html';
    } else {
        const errorMessage = await response.text()
        errorMessageElement.style.display = "block";
        errorMessageElement.innerHTML = errorMessage;
    }
}

document.getElementById("registerForm").addEventListener("submit", function (event) {
    event.preventDefault();
    registerRequest();
});