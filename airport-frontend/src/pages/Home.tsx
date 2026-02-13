import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { flightsApi } from '../api';
import { useAuthStore } from '../store/authStore';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { Label } from '../components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/card';
import { formatDate, formatTime, getFlightStatusColor } from '../lib/utils';
import { Plane } from 'lucide-react';
import type { SearchFlightsParams } from '../types';

const searchSchema = z.object({
  origin: z.string().min(3, 'Origin must be 3 characters').max(3).toUpperCase(),
  destination: z.string().min(3, 'Destination must be 3 characters').max(3).toUpperCase(),
  date: z.string().optional(),
});

type SearchFormData = z.infer<typeof searchSchema>;

export default function Home() {
  const [searchParams, setSearchParams] = useState<SearchFlightsParams | null>(null);
  const { isAuthenticated } = useAuthStore();
  const navigate = useNavigate();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SearchFormData>({
    resolver: zodResolver(searchSchema),
  });

  const { data: flights, isLoading } = useQuery({
    queryKey: ['flights', searchParams],
    queryFn: () => flightsApi.search(searchParams!),
    enabled: !!searchParams,
  });

  const onSearch = (data: SearchFormData) => {
    setSearchParams({
      origin: data.origin,
      destination: data.destination,
      date: data.date || undefined,
    });
  };

  const handleBookFlight = (flightId: number) => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    navigate(`/book/${flightId}`);
  };

  return (
    <div className="min-h-screen">
      {/* Hero Section with Search */}
      <div className="bg-gradient-to-br from-[#3b5998] to-[#2d4373] text-white py-20">
        <div className="container mx-auto px-4">
          <div className="text-center mb-12">
            <h1 className="text-5xl font-bold mb-4">Find your flight</h1>
            <p className="text-xl opacity-90">Search for available flights worldwide</p>
          </div>

          {/* Search Form */}
          <form onSubmit={handleSubmit(onSearch)} className="max-w-4xl mx-auto">
            <div className="bg-white rounded-lg shadow-xl p-6">
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="origin" className="text-gray-700">From</Label>
                  <Input
                    id="origin"
                    placeholder="JFK"
                    className="uppercase"
                    maxLength={3}
                    {...register('origin')}
                  />
                  {errors.origin && (
                    <p className="text-sm text-red-500">{errors.origin.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="destination" className="text-gray-700">To</Label>
                  <Input
                    id="destination"
                    placeholder="LAX"
                    className="uppercase"
                    maxLength={3}
                    {...register('destination')}
                  />
                  {errors.destination && (
                    <p className="text-sm text-red-500">{errors.destination.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="date" className="text-gray-700">Date</Label>
                  <Input
                    id="date"
                    type="date"
                    {...register('date')}
                  />
                  {errors.date && (
                    <p className="text-sm text-red-500">{errors.date.message}</p>
                  )}
                </div>

                <div className="flex items-end">
                  <Button type="submit" className="w-full bg-[#ff6b35] hover:bg-[#ff5722] text-white">
                    Find tickets
                  </Button>
                </div>
              </div>
            </div>
          </form>
        </div>
      </div>

      {/* Flight Results */}
      <div className="container mx-auto px-4 py-12">
        {isLoading && (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-[#3b5998]"></div>
            <p className="mt-4 text-gray-600">Searching flights...</p>
          </div>
        )}

        {flights && flights.length > 0 && (
          <div>
            <h2 className="text-3xl font-bold mb-8">Available flights</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              {flights.map((flight) => (
                <Card key={flight.id} className="hover:shadow-lg transition-shadow">
                  <CardHeader className="pb-4">
                    <div className="flex items-center justify-between mb-2">
                      <CardTitle className="text-2xl">{flight.flight_no}</CardTitle>
                      <span className={`px-2 py-1 rounded text-xs text-white ${getFlightStatusColor(flight.status)}`}>
                        {flight.status}
                      </span>
                    </div>
                    <div className="flex items-center text-lg font-medium text-gray-700">
                      <span>{flight.origin}</span>
                      <Plane className="mx-2 h-4 w-4" />
                      <span>{flight.destination}</span>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2 text-sm">
                      <div>
                        <p className="text-gray-500">Departure: {formatDate(flight.departure_time)}, {formatTime(flight.departure_time)}</p>
                      </div>
                      <div>
                        <p className="text-gray-500">Arrival: {formatDate(flight.arrival_time)}, {formatTime(flight.arrival_time)}</p>
                      </div>
                      <div className="pt-2">
                        <p className="text-lg font-semibold text-gray-800">Price: {flight.id * 3200} ₸</p>
                      </div>
                    </div>
                    <Button
                      className="w-full mt-4 bg-[#3b5998] hover:bg-[#2d4373]"
                      onClick={() => handleBookFlight(flight.id)}
                    >
                      Book
                    </Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        )}

        {flights && flights.length === 0 && (
          <div className="text-center py-12">
            <p className="text-xl text-gray-600">No flights found. Try different search criteria.</p>
          </div>
        )}
      </div>
    </div>
  );
}