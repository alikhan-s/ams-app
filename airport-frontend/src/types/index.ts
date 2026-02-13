// User & Auth Types
export interface User {
  id: number;
  full_name: string;
  email: string;
  role: 'PASSENGER' | 'STAFF' | 'ADMIN';
  created_at: string;
}

export interface RegisterRequest {
  full_name: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
}

// Flight Types
export interface Flight {
  id: number;
  flight_no: string;
  origin: string;
  destination: string;
  gate_id?: number;
  departure_time: string;
  arrival_time: string;
  status: 'SCHEDULED' | 'DELAYED' | 'BOARDING' | 'DEPARTED' | 'ARRIVED' | 'CANCELLED';
  version: number;
  created_at: string;
  updated_at: string;
  total_seats: number;
}

export interface SearchFlightsParams {
  origin?: string;
  destination?: string;
  date?: string; // YYYY-MM-DD
}

// Passenger Types
export interface Passenger {
  id: number;
  user_id: number;
  passport_no: string;
  phone: string;
  created_at: string;
  updated_at: string;
}

// Booking Types
export interface Ticket {
  id: number;
  flight_id: number;
  flight?: Flight;
  passenger_id: number;
  seat_no?: string;
  price: number;
  status: 'ACTIVE' | 'CANCELLED';
  created_at: string;
}

export interface BookingRequest {
  flight_id: number;
  passport_no?: string; // Required if passenger profile doesn't exist
  phone?: string;
}

// Airport Ops Types
export interface Gate {
  id: number;
  terminal_id: number;
  code: string;
  status: 'OPEN' | 'CLOSED' | 'MAINTENANCE';
}

export interface CreateGateRequest {
  terminal_id: number;
  code: string;
  status: 'OPEN' | 'CLOSED' | 'MAINTENANCE';
}

export interface Baggage {
  id: number;
  ticket_id: number;
  tag_code: string;
  status: 'RECEIVED' | 'LOADED' | 'IN_TRANSIT' | 'DELIVERED';
  updated_at: string;
}

export interface BaggageDetail {
  id: number;
  tag_code: string;
  status: 'RECEIVED' | 'LOADED' | 'IN_TRANSIT' | 'DELIVERED';
  updated_at: string;
  passenger_name: string;
  user_id: number;
}

export interface CreateBaggageRequest {
  ticket_id: number;
}

export interface UpdateBaggageRequest {
  status: 'RECEIVED' | 'LOADED' | 'IN_TRANSIT' | 'DELIVERED';
}

// API Error Response
export interface ApiError {
  error: string;
}