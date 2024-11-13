document.getElementById("loginForm").addEventListener("submit", function(event) {
    event.preventDefault();

    const username = document.getElementById("Username").value;
    const password = document.getElementById("Password").value;
    const errorMessageElement = document.getElementById("errorMessage");
    document.getElementById("errorMessage").style.display = "none";

    const loginData = {
        Username: username,
        Password: password
    };

    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(loginData)
    })
    .then((response) => {
        if (response.status == 200) {
            window.location.href = '/success.html';
        } else {
            errorMessageElement.style.display = "block";
            return response.text()
        }
    })
    .then((errorMessage) => {
        console.log(errorMessage);
        errorMessageElement.innerHTML = errorMessage;
    })
    .catch(error => console.error('Error:', error));
});