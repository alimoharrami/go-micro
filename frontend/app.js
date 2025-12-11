import loadConfig from './config.js';

const { API_DOMAIN } = await loadConfig();

document.addEventListener('DOMContentLoaded', fetchUsers);

document.getElementById('addUserForm').addEventListener('submit', async function (e) {
    e.preventDefault();

    const btn = document.getElementById('submitBtn');
    const form = this;
    const notificationArea = document.getElementById('notificationArea');

    // Reset state
    btn.classList.add('submitting');
    notificationArea.innerHTML = '';

    // Collect Data
    const formData = new FormData(form);
    const data = {};
    formData.forEach((value, key) => {
        data[key] = value;
    });

    try {
        const response = await fetch(`${API_DOMAIN}/api/users`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        const result = await response.json();

        if (response.ok) {
            showNotification('Success', `User ${result.first_name} created successfully!`, 'success');
            form.reset();
            fetchUsers(); // Refresh the list
        } else {
            showNotification('Error', result.error || 'Failed to create user', 'error');
        }

    } catch (error) {
        showNotification('Error', 'Network error occurred. Please try again.', 'error');
        console.error('Error:', error);
    } finally {
        btn.classList.remove('submitting');
    }
});

function showNotification(title, message, type) {
    const area = document.getElementById('notificationArea');
    const div = document.createElement('div');
    div.className = `notification ${type}`;
    div.innerHTML = `<strong>${title}:</strong> ${message}`;
    area.appendChild(div);

    setTimeout(() => {
        div.style.opacity = '0';
        setTimeout(() => div.remove(), 300);
    }, 5000);
}

async function fetchUsers() {
    console.log('Fetching users...');
    try {
        const response = await fetch(`${API_DOMAIN}/api/users?page=1&limit=10`);
        if (!response.ok) {
            console.error('Failed to fetch users');
            return;
        }

        const result = await response.json();
        // The API returns { data: [...], pagination: {...} } or just [...] depending on implementation?
        // Based on user/internal/service/user_service.go: GetUsers returns map with "data" key.

        const users = result.data || [];
        const tbody = document.getElementById('userTableBody');
        tbody.innerHTML = '';

        if (users.length === 0) {
            tbody.innerHTML = '<tr><td colspan="3" style="text-align:center">No users found</td></tr>';
            return;
        }

        users.forEach(user => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>${user.ID}</td>
                <td>${user.first_name} ${user.last_name}</td>
                <td>${user.Email}</td>
            `;
            tbody.appendChild(tr);
        });

    } catch (error) {
        console.error('Error fetching users:', error);
    }
}
