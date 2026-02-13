import { useParams, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useMutation, useQuery } from '@tanstack/react-query';
import { flightsApi, bookingsApi } from '../api';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { Label } from '../components/ui/label';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/card';
import { Plane, ArrowLeft } from 'lucide-react';
import { formatDate, formatTime, getFlightStatusColor } from '../lib/utils';
import toast from 'react-hot-toast';
import type { BookingRequest } from '../types';

const bookingSchema = z.object({
  passport_no: z.string().min(5, 'Passport number is required'),
  phone: z.string().min(10, 'Phone number is required'),
});

type BookingFormData = z.infer<typeof bookingSchema>;

export default function BookFlight() {
  const { flightId } = useParams<{ flightId: string }>();
  const navigate = useNavigate();

  const { data: flights } = useQuery({
    queryKey: ['flight-search'],
    queryFn: () => flightsApi.search({}),
  });

  const flight = flights?.find((f) => f.id === Number(flightId));

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<BookingFormData>({
    resolver: zodResolver(bookingSchema),
  });

  const bookMutation = useMutation({
    mutationFn: (data: BookingRequest) => bookingsApi.create(data),
    onSuccess: () => {
      toast.success('Flight booked successfully!');
      navigate('/bookings');
    },
  });

  const onSubmit = (data: BookingFormData) => {
    if (!flightId) return;

    bookMutation.mutate({
      flight_id: Number(flightId),
      passport_no: data.passport_no,
      phone: data.phone,
    });
  };

  if (!flight) {
    return (
      <div className="container mx-auto px-4 py-12">
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-xl text-gray-600">Flight not found</p>
            <Button className="mt-4" onClick={() => navigate('/')}>
              Go Back
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-12">
      <Button
        variant="ghost"
        className="mb-6"
        onClick={() => navigate('/')}
      >
        <ArrowLeft className="mr-2 h-4 w-4" />
        Back to Search
      </Button>

      <div className="max-w-2xl mx-auto space-y-6">
        {/* Flight Details */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Plane className="h-5 w-5" />
              Flight Details
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-2xl font-bold">{flight.flight_no}</p>
                  <div className="flex items-center text-lg text-gray-700 mt-1">
                    <span>{flight.origin}</span>
                    <Plane className="mx-2 h-4 w-4" />
                    <span>{flight.destination}</span>
                  </div>
                </div>
                <span
                  className={`px-3 py-1 rounded text-sm text-white ${getFlightStatusColor(flight.status)}`}
                >
                  {flight.status}
                </span>
              </div>

              <div className="grid grid-cols-2 gap-4 pt-4 border-t">
                <div>
                  <p className="text-sm text-gray-500">Departure</p>
                  <p className="font-medium">
                    {formatDate(flight.departure_time)}
                  </p>
                  <p className="text-gray-700">{formatTime(flight.departure_time)}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-500">Arrival</p>
                  <p className="font-medium">{formatDate(flight.arrival_time)}</p>
                  <p className="text-gray-700">{formatTime(flight.arrival_time)}</p>
                </div>
              </div>

              <div className="pt-4 border-t">
                <p className="text-2xl font-bold text-[#3b5998]">
                  Price: {flight.id * 3200} ₸
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Booking Form */}
        <Card>
          <CardHeader>
            <CardTitle>Passenger Information</CardTitle>
            <CardDescription>
              Please provide your passport and contact information
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="passport_no">Passport Number</Label>
                <Input
                  id="passport_no"
                  placeholder="AB1234567"
                  {...register('passport_no')}
                />
                {errors.passport_no && (
                  <p className="text-sm text-red-500">{errors.passport_no.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="phone">Phone Number</Label>
                <Input
                  id="phone"
                  type="tel"
                  placeholder="+7 700 123 4567"
                  {...register('phone')}
                />
                {errors.phone && (
                  <p className="text-sm text-red-500">{errors.phone.message}</p>
                )}
              </div>

              <Button
                type="submit"
                className="w-full bg-[#3b5998] hover:bg-[#2d4373]"
                disabled={bookMutation.isPending}
              >
                {bookMutation.isPending ? 'Booking...' : 'Confirm Booking'}
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}