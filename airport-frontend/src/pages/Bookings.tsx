import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { bookingsApi } from '../api';
import { Button } from '../components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '../components/ui/card';
import { formatDate, formatTime, getFlightStatusColor, getBaggageStatusColor } from '../lib/utils';
import { Plane, Package } from 'lucide-react';
import toast from 'react-hot-toast';
import { useState } from 'react';

export default function Bookings() {
  const queryClient = useQueryClient();
  const [showBaggage, setShowBaggage] = useState(false);

  const { data: bookings, isLoading } = useQuery({
    queryKey: ['my-bookings'],
    queryFn: bookingsApi.getMyBookings,
  });

  const { data: baggage } = useQuery({
    queryKey: ['my-baggage'],
    queryFn: bookingsApi.getMyBaggage,
    enabled: showBaggage,
  });

  const cancelMutation = useMutation({
    mutationFn: (ticketId: number) => bookingsApi.cancel(ticketId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['my-bookings'] });
      toast.success('Ticket cancelled successfully');
    },
  });

  const handleCancel = (ticketId: number) => {
    if (window.confirm('Are you sure you want to cancel this ticket?')) {
      cancelMutation.mutate(ticketId);
    }
  };

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-12">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-[#3b5998]"></div>
          <p className="mt-4 text-gray-600">Loading bookings...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-12">
      <div className="flex items-center justify-between mb-8">
        <h1 className="text-3xl font-bold">My bookings</h1>
        <Button
          variant={showBaggage ? 'default' : 'outline'}
          onClick={() => setShowBaggage(!showBaggage)}
          className="flex items-center gap-2"
        >
          <Package className="h-4 w-4" />
          {showBaggage ? 'Hide Baggage' : 'Check Baggage'}
        </Button>
      </div>

      {/* Baggage Section */}
      {showBaggage && baggage && (
        <Card className="mb-8 border-[#3b5998]/20">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Package className="h-5 w-5" />
              My Baggage
            </CardTitle>
          </CardHeader>
          <CardContent>
            {baggage.length === 0 ? (
              <p className="text-gray-500">No baggage checked in yet.</p>
            ) : (
              <div className="space-y-3">
                {baggage.map((bag) => (
                  <div
                    key={bag.id}
                    className="flex items-center justify-between p-4 border rounded-lg"
                  >
                    <div>
                      <p className="font-semibold">{bag.tag_code}</p>
                      <p className="text-sm text-gray-500">
                        Ticket ID: {bag.ticket_id}
                      </p>
                    </div>
                    <span
                      className={`px-3 py-1 rounded-full text-xs text-white ${getBaggageStatusColor(bag.status)}`}
                    >
                      {bag.status}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Bookings List */}
      {bookings && bookings.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-xl text-gray-600">No bookings found.</p>
            <p className="text-gray-500 mt-2">Book a flight to see it here.</p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {bookings?.map((ticket) => (
            <Card key={ticket.id} className="hover:shadow-lg transition-shadow">
              <CardHeader className="pb-4">
                <div className="flex items-center justify-between mb-2">
                  <CardTitle className="text-xl">{ticket.flight?.flight_no}</CardTitle>
                  <span
                    className={`px-2 py-1 rounded text-xs text-white ${
                      ticket.status === 'ACTIVE' ? 'bg-green-500' : 'bg-gray-500'
                    }`}
                  >
                    {ticket.status}
                  </span>
                </div>
                {ticket.flight && (
                  <>
                    <div className="flex items-center text-base font-medium text-gray-700">
                      <span>{ticket.flight.origin}</span>
                      <Plane className="mx-2 h-4 w-4" />
                      <span>{ticket.flight.destination}</span>
                    </div>
                    <span
                      className={`inline-block mt-2 px-2 py-1 rounded text-xs text-white ${getFlightStatusColor(
                        ticket.flight.status
                      )}`}
                    >
                      {ticket.flight.status}
                    </span>
                  </>
                )}
              </CardHeader>
              <CardContent>
                {ticket.flight && (
                  <div className="space-y-2 text-sm">
                    <div>
                      <p className="text-gray-500">
                        Departure: {formatDate(ticket.flight.departure_time)},{' '}
                        {formatTime(ticket.flight.departure_time)}
                      </p>
                    </div>
                    <div>
                      <p className="text-gray-500">
                        Arrival: {formatDate(ticket.flight.arrival_time)},{' '}
                        {formatTime(ticket.flight.arrival_time)}
                      </p>
                    </div>
                    {ticket.seat_no && (
                      <div>
                        <p className="text-gray-500">Seat: {ticket.seat_no}</p>
                      </div>
                    )}
                    <div className="pt-2">
                      <p className="text-lg font-semibold text-gray-800">
                        Price: {ticket.price} ₸
                      </p>
                    </div>
                  </div>
                )}

                {ticket.status === 'ACTIVE' && (
                  <Button
                    variant="destructive"
                    className="w-full mt-4"
                    onClick={() => handleCancel(ticket.id)}
                    disabled={cancelMutation.isPending}
                  >
                    Cancel Ticket
                  </Button>
                )}
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}