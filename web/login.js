async function loginRequest() {
    const username = document.getElementById("Username").value;
    const password = document.getElementById("Password").value;
    const errorMessageElement = document.getElementById("errorMessage");
    document.getElementById("errorMessage").style.display = "none";

    const loginData = {
        Username: username,
        Password: password
    };

    const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(loginData)
    })
    const responseText = await response.text()
    if (response.status == 200) {
        localStorage.setItem('token', responseText);
        window.location.href = '/authenticated.html';
    } else {
        errorMessageElement.style.display = "block";
        errorMessageElement.innerHTML = responseText;
    }
}

document.getElementById("loginForm").addEventListener("submit", function(event) {
    event.preventDefault();
    loginRequest();
});