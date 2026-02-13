import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { opsApi } from '../../api';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Label } from '../../components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '../../components/ui/card';
import { getBaggageStatusColor } from '../../lib/utils';
import toast from 'react-hot-toast';
import type { CreateGateRequest, UpdateBaggageRequest } from '../../types';

const gateSchema = z.object({
  terminal_id: z.string().min(1, 'Terminal ID is required'),
  code: z.string().min(1, 'Gate code is required'),
  status: z.enum(['OPEN', 'CLOSED', 'MAINTENANCE']),
});

type GateFormData = z.infer<typeof gateSchema>;

export default function AdminDashboard() {
  const queryClient = useQueryClient();

  const { data: gates } = useQuery({
    queryKey: ['gates'],
    queryFn: opsApi.listGates,
  });

  const { data: baggageList } = useQuery({
    queryKey: ['all-baggage'],
    queryFn: opsApi.listAllBaggage,
  });

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<GateFormData>({
    resolver: zodResolver(gateSchema),
    defaultValues: {
      status: 'OPEN',
    },
  });

  const createGateMutation = useMutation({
    mutationFn: (data: CreateGateRequest) => opsApi.createGate(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['gates'] });
      toast.success('Gate created successfully');
      reset();
    },
  });

  const updateBaggageMutation = useMutation({
    mutationFn: ({ id, status }: { id: number; status: string }) =>
      opsApi.updateBaggage(id, { status } as UpdateBaggageRequest),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['all-baggage'] });
      toast.success('Baggage status updated');
    },
  });

  const onSubmitGate = (data: GateFormData) => {
    createGateMutation.mutate({
      ...data,
      terminal_id: parseInt(data.terminal_id, 10),
    });
  };

  const handleUpdateBaggage = (id: number, currentStatus: string) => {
    const statusFlow: Record<string, string> = {
      RECEIVED: 'LOADED',
      LOADED: 'IN_TRANSIT',
      IN_TRANSIT: 'DELIVERED',
      DELIVERED: 'DELIVERED',
    };

    const nextStatus = statusFlow[currentStatus];
    if (nextStatus) {
      updateBaggageMutation.mutate({ id, status: nextStatus });
    }
  };

  return (
    <div className="container mx-auto px-4 py-12">
      <h1 className="text-3xl font-bold mb-8">Admin Dashboard</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-12">
        {/* Gates Section */}
        <Card>
          <CardHeader>
            <CardTitle>Gates Management</CardTitle>
          </CardHeader>
          <CardContent>
            {/* Add Gate Form */}
            <form onSubmit={handleSubmit(onSubmitGate)} className="space-y-4 mb-6">
              <div className="grid grid-cols-3 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="terminal_id">Terminal</Label>
                  <Input
                    id="terminal_id"
                    type="number"
                    placeholder="1"
                    {...register('terminal_id')}
                  />
                  {errors.terminal_id && (
                    <p className="text-xs text-red-500">{errors.terminal_id.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="code">Code</Label>
                  <Input id="code" placeholder="A1" {...register('code')} />
                  {errors.code && (
                    <p className="text-xs text-red-500">{errors.code.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="status">Status</Label>
                  <select
                    id="status"
                    {...register('status')}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  >
                    <option value="OPEN">OPEN</option>
                    <option value="CLOSED">CLOSED</option>
                    <option value="MAINTENANCE">MAINTENANCE</option>
                  </select>
                </div>
              </div>

              <Button
                type="submit"
                className="w-full bg-[#3b5998] hover:bg-[#2d4373]"
                disabled={createGateMutation.isPending}
              >
                Add Gate
              </Button>
            </form>

            {/* Gates List */}
            <div className="space-y-2 max-h-64 overflow-y-auto">
              {gates && gates.length === 0 ? (
                <p className="text-sm text-gray-500">No gates available</p>
              ) : (
                gates?.map((gate) => (
                  <div
                    key={gate.id}
                    className="flex items-center justify-between p-3 border rounded-lg"
                  >
                    <div>
                      <p className="font-medium">
                        Gate {gate.code} - Terminal {gate.terminal_id}
                      </p>
                    </div>
                    <span
                      className={`px-2 py-1 rounded text-xs text-white ${gate.status === 'OPEN'
                        ? 'bg-green-500'
                        : gate.status === 'CLOSED'
                          ? 'bg-red-500'
                          : 'bg-yellow-500'
                        }`}
                    >
                      {gate.status}
                    </span>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>

        {/* Baggage Operations Section */}
        <Card>
          <CardHeader>
            <CardTitle>Baggage Operations</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2 max-h-96 overflow-y-auto">
              {baggageList && baggageList.length === 0 ? (
                <p className="text-sm text-gray-500">No baggage to display</p>
              ) : (
                baggageList?.map((bag) => (
                  <div
                    key={bag.id}
                    className="flex items-center justify-between p-3 border rounded-lg"
                  >
                    <div className="flex-1">
                      <p className="font-semibold">{bag.tag_code}</p>
                      <p className="text-sm text-gray-600">{bag.passenger_name}</p>
                      <p className="text-xs text-gray-500">ID: {bag.id}</p>
                    </div>
                    <div className="flex items-center gap-2">
                      <span
                        className={`px-2 py-1 rounded text-xs text-white ${getBaggageStatusColor(
                          bag.status
                        )}`}
                      >
                        {bag.status}
                      </span>
                      {bag.status !== 'DELIVERED' && (
                        <Button
                          size="sm"
                          onClick={() => handleUpdateBaggage(bag.id, bag.status)}
                          disabled={updateBaggageMutation.isPending}
                        >
                          Next
                        </Button>
                      )}
                    </div>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}