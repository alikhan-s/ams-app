const API_BASE_URL = 'http://localhost:8080/api/v1';

class Api {
    static getHeaders() {
        const headers = {
            'Content-Type': 'application/json'
        };
        const token = localStorage.getItem('token');
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        return headers;
    }

    static async request(endpoint, method = 'GET', body = null) {
        const config = {
            method,
            headers: this.getHeaders(),
        };

        if (body) {
            config.body = JSON.stringify(body);
        }

        try {
            const response = await fetch(`${API_BASE_URL}${endpoint}`, config);

            if (response.status === 401) {
                // Token expired or invalid
                localStorage.removeItem('token');
                localStorage.removeItem('user');
                window.location.href = 'login.html';
                return;
            }

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || `Request failed with status ${response.status}`);
            }

            return await response.json();
        } catch (error) {
            console.error('API Request Error:', error);
            throw error;
        }
    }

    static get(endpoint) {
        return this.request(endpoint, 'GET');
    }

    static post(endpoint, body) {
        return this.request(endpoint, 'POST', body);
    }

    static patch(endpoint, body) {
        return this.request(endpoint, 'PATCH', body);
    }

    static delete(endpoint) {
        return this.request(endpoint, 'DELETE');
    }
}

// Export for module usage if needed, though we are using vanilla JS imports or global script loading
window.Api = Api;
