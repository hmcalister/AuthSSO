async function makeAuthenticatedRequest() {
    const token = localStorage.getItem('token');

    if (!token) {
        window.location.href = '/login.html';
    }

    const headerElement = document.getElementById("header");
    const infoElement = document.getElementById("info");

    const response = await fetch("/api/authenticate", {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}`,
        },
    })
    if (response.ok) {
        const userInformation = await response.json();
        headerElement.innerHTML = `Welcome, ${userInformation.Username}.`;
        infoElement.innerHTML = "You have logged in successfully.";
    } else {
        headerElement.innerHTML = "You are not logged in.";
        infoElement.innerHTML = '<a href="/login.html" style="border-radius: 0.5em; background-color: var(--pico-primary-background); color: var(--pico-contrast); text-decoration: none; padding: 0.2em;">Log in</a></small>';
    }
}

makeAuthenticatedRequest();

