import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}

export function formatTime(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
  });
}

export function formatDateTime(dateString: string): string {
  return `${formatDate(dateString)} ${formatTime(dateString)}`;
}

export function getFlightStatusColor(status: string): string {
  const colors: Record<string, string> = {
    SCHEDULED: 'bg-green-500',
    BOARDING: 'bg-blue-500',
    DEPARTED: 'bg-gray-500',
    DELAYED: 'bg-yellow-500',
    CANCELLED: 'bg-red-500',
    ARRIVED: 'bg-purple-500',
  };
  return colors[status] || 'bg-gray-400';
}

export function getBaggageStatusColor(status: string): string {
  const colors: Record<string, string> = {
    RECEIVED: 'bg-blue-500',
    LOADED: 'bg-green-500',
    IN_TRANSIT: 'bg-yellow-500',
    DELIVERED: 'bg-purple-500',
  };
  return colors[status] || 'bg-gray-400';
}