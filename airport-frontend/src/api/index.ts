import apiClient from './axios';
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  SearchFlightsParams,
  Flight,
  BookingRequest,
  Ticket,
  Gate,
  CreateGateRequest,
  Baggage,
  BaggageDetail,
  CreateBaggageRequest,
  UpdateBaggageRequest,
} from '../types';

// Auth API
export const authApi = {
  register: async (data: RegisterRequest) => {
    const response = await apiClient.post<{ message: string }>('/auth/register', data);
    return response.data;
  },

  login: async (data: LoginRequest) => {
    const response = await apiClient.post<AuthResponse>('/auth/login', data);
    return response.data;
  },
};

// Flights API
export const flightsApi = {
  search: async (params: SearchFlightsParams) => {
    const response = await apiClient.get<Flight[]>('/flights', { params });
    return response.data;
  },
};

// Bookings API
export const bookingsApi = {
  create: async (data: BookingRequest) => {
    const response = await apiClient.post<Ticket>('/bookings', data);
    return response.data;
  },

  getMyBookings: async () => {
    const response = await apiClient.get<Ticket[]>('/bookings/my');
    return response.data;
  },

  cancel: async (ticketId: number) => {
    const response = await apiClient.post<{ message: string }>(`/bookings/${ticketId}/cancel`);
    return response.data;
  },

  getMyBaggage: async () => {
    const response = await apiClient.get<Baggage[]>('/bookings/baggage');
    return response.data;
  },
};

// Airport Operations API (Staff/Admin)
export const opsApi = {
  // Gates
  createGate: async (data: CreateGateRequest) => {
    const response = await apiClient.post<Gate>('/ops/gates', data);
    return response.data;
  },

  listGates: async () => {
    const response = await apiClient.get<Gate[]>('/ops/gates');
    return response.data;
  },

  // Baggage
  checkInBaggage: async (data: CreateBaggageRequest) => {
    const response = await apiClient.post<Baggage>('/ops/baggage', data);
    return response.data;
  },

  listAllBaggage: async () => {
    const response = await apiClient.get<BaggageDetail[]>('/ops/baggage');
    return response.data;
  },

  updateBaggage: async (id: number, data: UpdateBaggageRequest) => {
    const response = await apiClient.patch<Baggage>(`/ops/baggage/${id}`, data);
    return response.data;
  },
};