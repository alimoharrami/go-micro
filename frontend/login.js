
import loadConfig from './config.js';

let API_URL = '';
let currentPage = 1;
const limit = 10;

const loginForm = document.getElementById('loginForm');
const errorMessage = document.getElementById('errorMessage');

// Initialize
async function init() {
    // const token = localStorage.getItem('token');
    // if (!token) {
    //     window.location.href = 'login.html';
    //     return;
    // }

    const config = await loadConfig();
    // Default to localhost:8080 (Gateway) if origin is file:// or similar
    API_URL = config.API_DOMAIN;
    if (API_URL === 'null' || API_URL.startsWith('file://')) {
        API_URL = 'http://localhost:8080';
    }

    // Remove trailing slash
    API_URL = API_URL.replace(/\/$/, "");

    console.log('API Endpoint:', API_URL);
}
    console.log(loginForm);
loginForm.onsubmit = async (e) => {
    e.preventDefault();
    console.log('here');
console.log('Constructed API_URL:', API_URL);
    errorMessage.style.display = 'none';

    if (!API_URL) {
        errorMessage.textContent = 'Configuration not loaded. Please refresh.';
        errorMessage.style.display = 'block';
        return;
    }

    const payload = {
        email: document.getElementById('email').value,
        password: document.getElementById('password').value
    };

    const submitBtn = loginForm.querySelector('button[type="submit"]');
    const originalBtnText = submitBtn.innerHTML;
    submitBtn.innerHTML = 'Signing in...';
    submitBtn.disabled = true;

    try {
        const res = await fetch(`${API_URL}/api/auth`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        });

        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.error || 'Login failed');
        }

        // Save auth data
        localStorage.setItem('token', data.token);
        localStorage.setItem('user', JSON.stringify(data.user));

        window.location.href = 'index.html';

    } catch (error) {
        console.error('Login error:', error);
        errorMessage.textContent = error.message || 'Invalid email or password';
        errorMessage.style.display = 'block';

        submitBtn.innerHTML = originalBtnText;
        submitBtn.disabled = false;
    }
};

init();