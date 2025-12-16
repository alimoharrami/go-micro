import loadConfig from './config.js';

let API_URL = '';
let currentPage = 1;
const limit = 10;

// State
let isEditing = false;
let currentUserId = null;

// DOM Elements
const usersTableBody = document.getElementById('userListBody');
const modal = document.getElementById('userModal');
const modalTitle = document.getElementById('modalTitle');
const userForm = document.getElementById('userForm');
const paginationContainer = document.getElementById('pagination');

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
    fetchUsers();
}

// Fetch Users
async function fetchUsers(page = 1) {
    try {
<<<<<<< HEAD
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/users?page=${page}&limit=${limit}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (response.status === 401) {
            localStorage.removeItem('token');
            window.location.href = 'login.html';
            return;
        }
=======
        const response = await fetch(`${API_URL}/api/users?page=${page}&limit=${limit}`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
        });
>>>>>>> 00d025c5414ad08ba72aa1709d1b1161f9169728
        if (!response.ok) throw new Error('Failed to fetch users');

        const data = await response.json();
        // Backend returns: { data: [...], pagination: { ... } }
        renderTable(data.data || []); // Handle case where data might be wrapped
        renderPagination(data.pagination || {});

        currentPage = page;
    } catch (error) {
        console.error('Error:', error);
        usersTableBody.innerHTML = `<tr><td colspan="6" style="text-align:center; color:red;">Error loading users: ${error.message}</td></tr>`;
    }
}

// Render Table
function renderTable(users) {
    usersTableBody.innerHTML = '';

    if (users.length === 0) {
        usersTableBody.innerHTML = `<tr><td colspan="6" style="text-align:center;">No users found</td></tr>`;
        return;
    }

    users.forEach(user => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>#${user.ID}</td>
            <td>${user.first_name || ''} ${user.last_name || ''}</td>
            <td>${user.Email}</td>
            <td><span class="badge badge-success">User</span></td>
            <td><span class="badge ${user.Active ? 'badge-success' : 'badge-danger'}">${user.Active ? 'Active' : 'Inactive'}</span></td>
            <td>
                <div class="action-buttons">
                    <button class="btn-icon" onclick="window.editUser(${user.ID})">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                    </button>
                    <button class="btn-icon delete" onclick="window.deleteUser(${user.ID})">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                    </button>
                </div>
            </td>
        `;
        usersTableBody.appendChild(tr);
    });
}

// Pagination
function renderPagination(pagination) {
    paginationContainer.innerHTML = '';
    if (!pagination.pages || pagination.pages <= 1) return;

    // Previous
    const prevBtn = document.createElement('button');
    prevBtn.className = 'btn btn-secondary';
    prevBtn.disabled = currentPage === 1;
    prevBtn.innerText = 'Previous';
    prevBtn.onclick = () => fetchUsers(currentPage - 1);
    paginationContainer.appendChild(prevBtn);

    // Page Info (Simple)
    const info = document.createElement('span');
    info.style.margin = '0 10px';
    info.innerText = `Page ${pagination.page} of ${pagination.pages}`;
    paginationContainer.appendChild(info);

    // Next
    const nextBtn = document.createElement('button');
    nextBtn.className = 'btn btn-secondary';
    nextBtn.disabled = currentPage === pagination.pages;
    nextBtn.innerText = 'Next';
    nextBtn.onclick = () => fetchUsers(currentPage + 1);
    paginationContainer.appendChild(nextBtn);
}

// Modal Functions - Attached to Window
window.openModal = (mode = 'create') => {
    isEditing = mode === 'edit';
    modal.classList.add('active');

    if (isEditing) {
        modalTitle.innerText = 'Edit User';
        document.getElementById('passwordGroup').style.display = 'none'; // Don't edit password here for simplicity
        document.getElementById('password').removeAttribute('required');
    } else {
        modalTitle.innerText = 'Add New User';
        userForm.reset();
        document.getElementById('passwordGroup').style.display = 'block';
        document.getElementById('password').setAttribute('required', 'true');
    }
};

window.closeModal = () => {
    modal.classList.remove('active');
    userForm.reset();
    currentUserId = null;
    isEditing = false;
};

// Form Submit
userForm.onsubmit = async (e) => {
    e.preventDefault();

    const formData = {
        first_name: document.getElementById('firstName').value,
        last_name: document.getElementById('lastName').value,
        // Using correct JSON keys backend expects
    };

    // Add only if not empty (though required fields are checked by HTML)

    // For Update (PUT)
    if (isEditing) {
        // Backend expects UpdateUserInput: FirstName, LastName, Active pointers
        // The endpoint is PUT /users/:id
        // The current UserService Update takes fields. 
        // Let's send what we have.
        try {
            // For simplicity, we assume we just update names. 
            // We need to fetch the user first to fill the form properly in real app, 
            // but here we just submit. Wait, we need to populate form first for edit.
        } catch (e) { }
    }

    const url = isEditing ? `${API_URL}/api/users/${currentUserId}` : `${API_URL}/api/users`;
    const method = isEditing ? 'PUT' : 'POST';

    const payload = {
        first_name: document.getElementById('firstName').value,
        last_name: document.getElementById('lastName').value,
    };

    if (!isEditing) {
        payload.Email = document.getElementById('email').value; // Cap E? Go struct has Email
        payload.Password = document.getElementById('password').value;
    } else {
        // Update doesn't take email/password typically in this simple implementation
        // Check backend: UpdateUserInput { FirstName, LastName, Active }
        // So we send just these.
    }

    try {
        const token = localStorage.getItem('token');
        const res = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(payload)
        });

        if (!res.ok) {
            const err = await res.json();
            throw new Error(err.error || 'Request failed');
        }

        window.closeModal();
        fetchUsers(currentPage); // Refresh
    } catch (error) {
        alert(error.message);
    }
};

// Edit User (Fetch details first)
window.editUser = async (id) => {
    currentUserId = id;
    try {
<<<<<<< HEAD
        const token = localStorage.getItem('token');
        const res = await fetch(`${API_URL}/users/${id}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (!res.ok) throw new Error("Could not fetch user");
=======
        const res = await fetch(`${API_URL}/api/users/${id}`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
        });
        if(!res.ok) throw new Error("Could not fetch user");
>>>>>>> 00d025c5414ad08ba72aa1709d1b1161f9169728
        const user = await res.json();

        // Populate form
        document.getElementById('firstName').value = user.first_name || '';
        document.getElementById('lastName').value = user.last_name || '';
        document.getElementById('email').value = user.Email || '';
        document.getElementById('userId').value = user.ID;

        // Email usually readonly on edit
        document.getElementById('email').setAttribute('readonly', 'true');

        window.openModal('edit');
    } catch (e) {
        console.error(e);
        alert("Failed to load user details");
    }
};

// Delete User
window.deleteUser = async (id) => {
    if (!confirm('Are you sure you want to delete this user?')) return;

    try {
<<<<<<< HEAD
        const token = localStorage.getItem('token');
        const res = await fetch(`${API_URL}/users/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
=======
        const res = await fetch(`${API_URL}/api/users/${id}`, {
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' }
>>>>>>> 00d025c5414ad08ba72aa1709d1b1161f9169728
        });

        if (!res.ok) throw new Error("Failed to delete");

        fetchUsers(currentPage);
    } catch (e) {
        alert(e.message);
    }
};

// Start
init();
