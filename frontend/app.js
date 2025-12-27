import loadConfig from './config.js';

let API_URL = '';
let currentPage = 1;
const limit = 10;

// State
let isEditing = false;
let currentUserId = null;

// DOM Elements
const usersTableBody = document.getElementById('userListBody');
const blogsTableBody = document.getElementById('blogListBody');
const userModal = document.getElementById('userModal');
const blogModal = document.getElementById('blogModal');
const modalTitle = document.getElementById('modalTitle');
const blogModalTitle = document.getElementById('blogModalTitle');
const userForm = document.getElementById('userForm');
const blogForm = document.getElementById('blogForm');
const userPaginationContainer = document.getElementById('userPagination');
const blogPaginationContainer = document.getElementById('blogPagination');

const userSection = document.getElementById('userSection');
const blogSection = document.getElementById('blogSection');
const usersNavLink = document.getElementById('usersNavLink');
const blogsNavLink = document.getElementById('blogsNavLink');

let currentView = 'users';

// Initialize
async function init() {
    const config = await loadConfig();
    API_URL = config.API_DOMAIN;
    if (API_URL === 'null' || API_URL.startsWith('file://')) {
        API_URL = 'http://localhost:8080';
    }

    API_URL = API_URL.replace(/\/$/, "");

    console.log('API Endpoint:', API_URL);
    fetchUsers();
    fetchBlogs();
}

// View Switcher
window.switchView = (view) => {
    currentView = view;
    if (view === 'users') {
        userSection.style.display = 'block';
        blogSection.style.display = 'none';
        usersNavLink.classList.add('active');
        blogsNavLink.classList.remove('active');
        fetchUsers(currentPage);
    } else {
        userSection.style.display = 'none';
        blogSection.style.display = 'block';
        usersNavLink.classList.remove('active');
        blogsNavLink.classList.add('active');
        fetchBlogs(currentPage);
    }
};

// Fetch Users
async function fetchUsers(page = 1) {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/api/users?page=${page}&limit=${limit}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (response.status === 401) {
            localStorage.removeItem('token');
            window.location.href = 'login.html';
            return;
        }

        if (!response.ok) throw new Error('Failed to fetch users');

        const data = await response.json();
        renderTable(data.data || []);
        renderPagination(data.pagination || {}, 'user');

        currentPage = page;
    } catch (error) {
        console.error('Error:', error);
        usersTableBody.innerHTML = `<tr><td colspan="6" style="text-align:center; color:red;">Error loading users: ${error.message}</td></tr>`;
    }
}

// Fetch Blogs
async function fetchBlogs(page = 1) {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/api/posts?page=${page}&limit=${limit}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) throw new Error('Failed to fetch blogs');

        const data = await response.json();
        renderBlogTable(data.data || []);
        renderPagination(data.pagination || {}, 'blog');

        currentPage = page;
    } catch (error) {
        console.error('Error:', error);
        blogsTableBody.innerHTML = `<tr><td colspan="4" style="text-align:center; color:red;">Error loading blogs: ${error.message}</td></tr>`;
    }
}

// Render User Table
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
                    <button class="btn-icon" onclick="window.subscribeChannel(${user.ID})" title="Subscribe to Main Channel">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path><path d="M13.73 21a2 2 0 0 1-3.46 0"></path></svg>
                    </button>
                </div>
            </td>
        `;
        usersTableBody.appendChild(tr);
    });
}

// Render Blog Table
function renderBlogTable(blogs) {
    blogsTableBody.innerHTML = '';
    if (blogs.length === 0) {
        blogsTableBody.innerHTML = `<tr><td colspan="4" style="text-align:center;">No blogs found</td></tr>`;
        return;
    }

    blogs.forEach(blog => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>#${blog.ID}</td>
            <td>${blog.Title}</td>
            <td>${new Date(blog.CreatedAt).toLocaleDateString()}</td>
            <td>
                <div class="action-buttons">
                    <button class="btn-icon" onclick="window.editBlog(${blog.ID})">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path></svg>
                    </button>
                    <button class="btn-icon delete" onclick="window.deleteBlog(${blog.ID})">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                    </button>
                </div>
            </td>
        `;
        blogsTableBody.appendChild(tr);
    });
}

// Pagination
function renderPagination(pagination, type) {
    const container = type === 'user' ? userPaginationContainer : blogPaginationContainer;
    container.innerHTML = '';
    if (!pagination.pages || pagination.pages <= 1) return;

    const fetchFunc = type === 'user' ? fetchUsers : fetchBlogs;

    // Previous
    const prevBtn = document.createElement('button');
    prevBtn.className = 'btn btn-secondary';
    prevBtn.disabled = currentPage === 1;
    prevBtn.innerText = 'Previous';
    prevBtn.onclick = () => fetchFunc(currentPage - 1);
    container.appendChild(prevBtn);

    // Page Info
    const info = document.createElement('span');
    info.style.margin = '0 10px';
    info.innerText = `Page ${pagination.page} of ${pagination.pages}`;
    container.appendChild(info);

    // Next
    const nextBtn = document.createElement('button');
    nextBtn.className = 'btn btn-secondary';
    nextBtn.disabled = currentPage === pagination.pages;
    nextBtn.innerText = 'Next';
    nextBtn.onclick = () => fetchFunc(currentPage + 1);
    container.appendChild(nextBtn);
}

// Modal Functions
window.openModal = (mode = 'create') => {
    isEditing = mode === 'edit';
    userModal.classList.add('active');
    if (isEditing) {
        modalTitle.innerText = 'Edit User';
        document.getElementById('passwordGroup').style.display = 'none';
        document.getElementById('password').removeAttribute('required');
    } else {
        modalTitle.innerText = 'Add New User';
        userForm.reset();
        document.getElementById('passwordGroup').style.display = 'block';
        document.getElementById('password').setAttribute('required', 'true');
    }
};

window.closeModal = () => {
    userModal.classList.remove('active');
    userForm.reset();
    currentUserId = null;
    isEditing = false;
};

// Blog Modal Functions
window.openBlogModal = (mode = 'create') => {
    isEditing = mode === 'edit';
    blogModal.classList.add('active');
    if (isEditing) {
        blogModalTitle.innerText = 'Edit Blog';
    } else {
        blogModalTitle.innerText = 'Add New Blog';
        blogForm.reset();
    }
};

window.closeBlogModal = () => {
    blogModal.classList.remove('active');
    blogForm.reset();
    currentUserId = null;
    isEditing = false;
};

// Form Submits
userForm.onsubmit = async (e) => {
    e.preventDefault();
    const url = isEditing ? `${API_URL}/api/users/${currentUserId}` : `${API_URL}/api/users`;
    const method = isEditing ? 'PUT' : 'POST';

    const payload = {
        first_name: document.getElementById('firstName').value,
        last_name: document.getElementById('lastName').value,
    };

    if (!isEditing) {
        payload.Email = document.getElementById('email').value;
        payload.Password = document.getElementById('password').value;
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
        fetchUsers(currentPage);
    } catch (error) {
        alert(error.message);
    }
};

blogForm.onsubmit = async (e) => {
    e.preventDefault();
    const blogId = document.getElementById('blogId').value;
    const url = isEditing ? `${API_URL}/api/posts/${blogId}` : `${API_URL}/api/posts`;
    const method = isEditing ? 'PUT' : 'POST';

    const payload = {
        Title: document.getElementById('blogTitle').value,
        Content: document.getElementById('blogContent').value,
    };

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

        window.closeBlogModal();
        fetchBlogs(currentPage);
    } catch (error) {
        alert(error.message);
    }
};

// Edit Actions
window.editUser = async (id) => {
    currentUserId = id;
    try {
        const token = localStorage.getItem('token');
        const res = await fetch(`${API_URL}/api/users/${id}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (!res.ok) throw new Error("Could not fetch user");
        const user = await res.json();
        document.getElementById('firstName').value = user.first_name || '';
        document.getElementById('lastName').value = user.last_name || '';
        document.getElementById('email').value = user.Email || '';
        document.getElementById('userId').value = user.ID;
        document.getElementById('email').setAttribute('readonly', 'true');
        window.openModal('edit');
    } catch (e) {
        alert("Failed to load user details");
    }
};

window.editBlog = async (id) => {
    try {
        const token = localStorage.getItem('token');
        const res = await fetch(`${API_URL}/api/posts/${id}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (!res.ok) throw new Error("Could not fetch blog");
        const data = await res.json();
        const blog = data.post; // Backend returns { post: ..., user: ... }
        document.getElementById('blogTitle').value = blog.Title || '';
        document.getElementById('blogContent').value = blog.Content || '';
        document.getElementById('blogId').value = blog.ID;
        window.openBlogModal('edit');
    } catch (e) {
        alert("Failed to load blog details");
    }
};

// Delete Actions
window.deleteUser = async (id) => {
    if (!confirm('Are you sure you want to delete this user?')) return;
    try {
        const token = localStorage.getItem('token');
        const res = await fetch(`${API_URL}/api/users/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (!res.ok) throw new Error("Failed to delete");
        fetchUsers(currentPage);
    } catch (e) {
        alert(e.message);
    }
};

window.deleteBlog = async (id) => {
    if (!confirm('Are you sure you want to delete this blog?')) return;
    try {
        const token = localStorage.getItem('token');
        const res = await fetch(`${API_URL}/api/posts/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        if (!res.ok) throw new Error("Failed to delete");
        fetchBlogs(currentPage);
    } catch (e) {
        alert(e.message);
    }
};

// Subscription Action
window.subscribeChannel = async (userId) => {
    try {
        const token = localStorage.getItem('token');
        const res = await fetch(`${API_URL}/api/notification/subscribe`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                channel: "main",
                userID: parseInt(userId)
            })
        });

        if (!res.ok) {
            const err = await res.json();
            throw new Error(err.error || 'Subscription failed');
        }

        alert('Subscribed successfully!');
    } catch (e) {
        alert(e.message);
    }
};

init();
