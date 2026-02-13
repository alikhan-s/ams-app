document.addEventListener('DOMContentLoaded', () => {
    updateNavbar();

    // Load initial flights or handle search
    const searchForm = document.getElementById('search-form');
    if (searchForm) {
        searchForm.addEventListener('submit', handleSearch);
    }

    // Initial load (optional, maybe showing all upcoming flights?)
    // loadFlights(); 
});

function updateNavbar() {
    const token = localStorage.getItem('token');
    const authLinks = document.getElementById('auth-links');

    if (token) {
        const user = JSON.parse(localStorage.getItem('user') || '{}');
        let links = `
            <li class="nav-item">
                <a class="nav-link" href="bookings.html">My Bookings</a>
            </li>
        `;

        if (user.role === 'ADMIN' || user.role === 'STAFF') {
            links += `
            <li class="nav-item">
                <a class="nav-link" href="admin.html">Admin Panel</a>
            </li>
        `;
        }

        links += `
            <li class="nav-item">
                <a class="nav-link" href="#" id="logout-btn">Logout</a>
            </li>
        `;
        authLinks.innerHTML = links;

        document.getElementById('logout-btn').addEventListener('click', (e) => {
            e.preventDefault();
            logout();
        });
    } else {
        authLinks.innerHTML = `
            <li class="nav-item">
                <a class="nav-link" href="login.html">Login</a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="register.html">Register</a>
            </li>
        `;
    }
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = 'login.html';
}

async function handleSearch(e) {
    e.preventDefault();
    const origin = document.getElementById('origin').value;
    const destination = document.getElementById('destination').value;
    const date = document.getElementById('date').value;

    // Build query string
    const params = new URLSearchParams();
    if (origin) params.append('origin', origin);
    if (destination) params.append('destination', destination);
    if (date) params.append('date', date);

    try {
        const flights = await Api.get(`/flights?${params.toString()}`);
        renderFlights(flights);
    } catch (error) {
        console.error('Search failed:', error);
        alert('Failed to search flights');
    }
}

function renderFlights(flights) {
    const container = document.getElementById('flights-container');
    container.innerHTML = '';

    if (!flights || flights.length === 0) {
        container.innerHTML = '<div class="col-12 text-center text-muted">No flights found</div>';
        return;
    }

    flights.forEach(flight => {
        const card = document.createElement('div');
        card.className = 'col-md-4 mb-4';

        const statusBadge = getStatusBadge(flight.status);
        const date = new Date(flight.departure_time).toLocaleDateString();
        const time = new Date(flight.departure_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });

        card.innerHTML = `
            <div class="card h-100">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-center mb-3">
                        <h5 class="card-title mb-0 fw-bold">${flight.flight_number}</h5>
                        <span class="badge ${statusBadge}">${flight.status}</span>
                    </div>
                    <div class="d-flex justify-content-between mb-2">
                        <div>
                            <small class="text-muted">From</small>
                            <div class="fw-medium">${flight.origin}</div>
                        </div>
                        <div class="text-end">
                            <small class="text-muted">To</small>
                            <div class="fw-medium">${flight.destination}</div>
                        </div>
                    </div>
                    <div class="mb-3">
                        <small class="text-muted">Departure</small>
                        <div class="fw-medium">${date} at ${time}</div>
                    </div>
                    <div class="d-grid">
                        <button class="btn btn-accent btn-sm" onclick="openBookingModal(${flight.id}, '${flight.flight_number}')">
                            Book Flight
                        </button>
                    </div>
                </div>
            </div>
        `;
        container.appendChild(card);
    });
}

function getStatusBadge(status) {
    switch (status) {
        case 'SCHEDULED': return 'badge-scheduled';
        case 'DELAYED': return 'badge-delayed';
        case 'CANCELLED': return 'badge-cancelled';
        default: return 'bg-secondary';
    }
}

// Booking Modal Logic
let currentFlightId = null;

window.openBookingModal = (flightId, flightNumber) => {
    const token = localStorage.getItem('token');
    if (!token) {
        window.location.href = 'login.html';
        return;
    }

    currentFlightId = flightId;
    document.getElementById('bookingFlightNumber').textContent = flightNumber;
    const modal = new bootstrap.Modal(document.getElementById('bookingModal'));
    modal.show();
};

document.getElementById('confirmBookingBtn').addEventListener('click', async () => {
    const passport = document.getElementById('passport').value;
    const phone = document.getElementById('phone').value;

    if (!passport || !phone) {
        alert('Please fill in all details');
        return;
    }

    try {
        await Api.post('/bookings', {
            flight_id: currentFlightId,
            passport_number: passport,
            phone: phone
            // Add other fields if required by backend, e.g. seat_number? API spec says passport & phone.
        });
        alert('Booking confirmed!');
        bootstrap.Modal.getInstance(document.getElementById('bookingModal')).hide();
        // clear form
        document.getElementById('passport').value = '';
        document.getElementById('phone').value = '';
    } catch (error) {
        alert('Booking failed: ' + error.message);
    }
});
