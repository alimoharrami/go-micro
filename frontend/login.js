const loginForm = document.getElementById('loginForm');
const errorMessage = document.getElementById('errorMessage');

loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    errorMessage.style.display = 'none';

    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    const submitBtn = loginForm.querySelector('button[type="submit"]');
    const originalBtnText = submitBtn.innerHTML;
    submitBtn.innerHTML = 'Signing in...';
    submitBtn.disabled = true;

    try {
        const response = await fetch('http://localhost:8080/api/auth', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });

        const data = await response.json();

        if (response.ok) {
            // Store token and user info
            localStorage.setItem('token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user));

            // Redirect to dashboard
            window.location.href = 'index.html';
        } else {
            throw new Error(data.error || 'Login failed');
        }
    } catch (error) {
        console.error('Login error:', error);
        errorMessage.textContent = error.message || 'Invalid email or password. Please try again.';
        errorMessage.style.display = 'block';

        // Reset button
        submitBtn.innerHTML = originalBtnText;
        submitBtn.disabled = false;
    }
});
