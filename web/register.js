document.getElementById("registerForm").addEventListener("submit", function (event) {
    event.preventDefault();

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

    fetch('/api/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(registerData)
    })
        .then((response) => {
            if (response.status == 201) {
                errorMessageElement.style.display = "none";
                window.location.href = '/login.html';
            } else {
                errorMessageElement.style.display = "block";
                return response.text()
            }
        })
        .then((errorMessage) => {
            errorMessageElement.innerHTML = errorMessage;
        })
        .catch(error => console.error('Error:', error));
});