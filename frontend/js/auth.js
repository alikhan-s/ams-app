document.addEventListener('DOMContentLoaded', () => {
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');

    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;

            try {
                // Adjust per API spec: /auth/login
                // Assuming it returns { access_token: "...", user: {...} } or similar
                const data = await Api.post('/auth/login', { email, password });

                // Save token and user info
                // Adjust based on actual API response structure
                if (data.token) {
                    localStorage.setItem('token', data.token);
                    if (data.user) {
                        localStorage.setItem('user', JSON.stringify(data.user));
                    }
                    window.location.href = 'index.html';
                } else {
                    alert('Login failed: No token received');
                }
            } catch (error) {
                alert('Login failed: ' + error.message);
            }
        });
    }

    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const fullName = document.getElementById('full_name').value;
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;

            try {
                // Adjust per API spec: /auth/register
                await Api.post('/auth/register', { full_name: fullName, email, password });
                alert('Registration successful! Please login.');
                window.location.href = 'login.html';
            } catch (error) {
                alert('Registration failed: ' + error.message);
            }
        });
    }
});
