'use client'

import { useState, useEffect } from 'react'
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog"
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage } from "@/components/ui/form"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"

type ServiceData = {
  [key: string]: any[]
}

type TableConfig = {
  columns: {
    key: string
    label: string
    type: 'string' | 'number'
  }[]
}

const SERVICE_TABLES: Record<string, TableConfig> = {
  user: {
    columns: [
      { key: 'id', label: 'ID', type: 'string' },
      { key: 'name', label: 'Name', type: 'string' },
      { key: 'email', label: 'Email', type: 'string' },
    ]
  },
  product: {
    columns: [
      { key: 'id', label: 'ID', type: 'string' },
      { key: 'Name', label: 'Name', type: 'string' },
      { key: 'Category', label: 'Category', type: 'string' },
      { key: 'Price', label: 'Price', type: 'number' },
    ]
  },  
  order: {
    columns: [
      { key: 'id', label: 'ID', type: 'string' },
      { key: 'user_id', label: 'User ID', type: 'string' },
      { key: 'product_id', label: 'Product ID', type: 'string' },
      { key: 'quantity', label: 'Quantity', type: 'number' },
      { key: 'status', label: 'Status', type: 'string' },
    ]
  },
  payment: {
    columns: [
      { key: 'id', label: 'ID', type: 'string' },
      { key: 'order_id', label: 'Order ID', type: 'string' },
      { key: 'amount', label: 'Amount', type: 'number' },
      { key: 'status', label: 'Status', type: 'string' },
    ]
  },
  notification: {
    columns: [
      { key: 'id', label: 'ID', type: 'string' },
      { key: 'user_id', label: 'User ID', type: 'string' },
      { key: 'message', label: 'Message', type: 'string' },
      { key: 'status', label: 'Status', type: 'string' },
    ]
  },
}

function ResultsTable({ service, data, onEdit, onDelete }: { service: string; data: any[]; onEdit: (item: any) => void; onDelete: (id: string) => void }) {
  const config = SERVICE_TABLES[service]
  
  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            {config.columns.map((column) => (
              <TableHead key={column.key}>{column.label}</TableHead>
            ))}
            <TableHead>Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data.map((item) => (
            <TableRow key={item.id}>
              {config.columns.map((column) => (
                <TableCell key={column.key}>{String(item[column.key] ?? '-')}</TableCell>
              ))}
              <TableCell>
                <Button variant="outline" size="sm" className="mr-2" onClick={() => onEdit(item)}>Edit</Button>
                <Button variant="destructive" size="sm" onClick={() => onDelete(item.id)}>Delete</Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}

async function getJwtToken(userId: string, password: string): Promise<string> {
  const response = await fetch('/api/proxy/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ user_id: userId, password }),
  });

  if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

  const data = await response.json();
  return data.token;
}

async function fetchData(serviceName: string, jwtToken: string): Promise<any[]> {
  const response = await fetch(`/api/proxy/${serviceName}/${serviceName}s`, {
    headers: { 'Authorization': `${jwtToken}` },
  });

  if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

  return response.json();
}

async function createItem(serviceName: string, jwtToken: string, item: any): Promise<any> {
  const response = await fetch(`/api/proxy/${serviceName}/${serviceName}s`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `${jwtToken}`
    },
    body: JSON.stringify(item)
  });

  if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

  return response;
}

async function updateItem(serviceName: string, jwtToken: string, item: any): Promise<any> {
  const response = await fetch(`/api/proxy/${serviceName}/${serviceName}s/${item.id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `${jwtToken}`
    },
    body: JSON.stringify(item)
  });

  if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

  return response;
}

async function deleteItem(serviceName: string, jwtToken: string, id: string): Promise<void> {
  const response = await fetch(`/api/proxy/${serviceName}/${serviceName}s/${id}`, {
    method: 'DELETE',
    headers: { 'Authorization': `${jwtToken}` },
  });

  if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
}

async function resetService(serviceName: string, jwtToken: string): Promise<void> {
  const items = await fetchData(serviceName, jwtToken);
  for (const item of items) {
    if (Number(item.id) > 2) {
      await deleteItem(serviceName, jwtToken, item.id);
    }
  }
}

export default function MicroservicesTester() {
  const [serviceData, setServiceData] = useState<ServiceData>({})
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [userId, setUserId] = useState('')
  const [password, setPassword] = useState('')
  const [jwtToken, setJwtToken] = useState<string | null>(null)
  const [currentService, setCurrentService] = useState('user')
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingItem, setEditingItem] = useState<any | null>(null)

  const formSchema = z.object({
    ...Object.fromEntries(
      SERVICE_TABLES[currentService].columns
        .filter(column => column.key !== 'id')
        .map(column => {
          switch (column.type) {
            case 'string':
              return [column.key, z.string().min(1, { message: `${column.label} is required` })];
            case 'number':
              return [column.key, z.number().min(1, { message: `${column.label} is required` })];
            default:
              return [column.key, z.string().min(1, { message: `${column.label} is required` })];
          }
        })
      )
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: Object.fromEntries(
      SERVICE_TABLES[currentService].columns
        .filter(column => column.key !== 'id')
        .map(column => [column.key, column.type === 'number' ? 0 : ''])
    ),
  })

  useEffect(() => {
    if (jwtToken) {
      loadServiceData(currentService);
    }
  }, [jwtToken, currentService]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    try {
      const token = await getJwtToken(userId, password)
      setJwtToken(token)
    } catch (error: any) {
      setError(`Failed to login: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }

  const loadServiceData = async (serviceName: string) => {
    setIsLoading(true)
    setError(null)

    try {
      const data = await fetchData(serviceName, jwtToken!)
      const sortedData = data.sort((a: any, b: any) => (a.id > b.id ? 1 : -1));
      setServiceData(prev => ({ ...prev, [serviceName]: sortedData }))
    } catch (error: any) {
      setError(`Error fetching ${serviceName} data: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }

  const handleReset = async () => {
    setIsLoading(true)
    setError(null)

    try {
      await resetService(currentService, jwtToken!)
      await loadServiceData(currentService)
    } catch (error: any) {
      setError(`Error resetting ${currentService} service: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }

  const handleAdd = async (data: z.infer<typeof formSchema>) => {
    setIsLoading(true)
    setError(null)

    try {
      await createItem(currentService, jwtToken!, data)
      await loadServiceData(currentService)
      setIsDialogOpen(false)
    } catch (error: any) {
      setError(`Error adding ${currentService}: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }

  const handleEdit = (item: any) => {
    setEditingItem(item)
    form.reset(item)
    setIsDialogOpen(true)
  }

  const handleUpdate = async (data: z.infer<typeof formSchema>) => {
    setIsLoading(true)
    setError(null)

    try {
      await updateItem(currentService, jwtToken!, { ...data, id: editingItem!.id })
      await loadServiceData(currentService)
      setIsDialogOpen(false)
      setEditingItem(null)
    } catch (error: any) {
      setError(`Error updating ${currentService}: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }

  const handleDelete = async (id: string) => {
    setIsLoading(true)
    setError(null)

    try {
      await deleteItem(currentService, jwtToken!, id)
      await loadServiceData(currentService)
    } catch (error: any) {
      setError(`Error deleting ${currentService}: ${error.message}`)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <Card className="w-full max-w-4xl mx-auto">
      <CardHeader>
        <CardTitle>Microservices Tester</CardTitle>
      </CardHeader>
      <CardContent>
        {!jwtToken ? (
          <form onSubmit={handleLogin} className="space-y-4 mb-4">
            <div>
              <Label htmlFor="userId">User ID</Label>
              <Input
                id="userId"
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                required
              />
            </div>
            <div>
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? 'Logging in...' : 'Login'}
            </Button>
          </form>
        ) : (
          <>
            <div className="flex justify-between items-center mb-4">
              <Button onClick={handleReset} disabled={isLoading}>Reset</Button>
              <Button onClick={() => setJwtToken(null)} className="mt-4">Logout</Button>
            </div>
            <Tabs value={currentService} onValueChange={setCurrentService}>
              <TabsList>
                {Object.keys(SERVICE_TABLES).map((service) => (
                  <TabsTrigger key={service} value={service}>
                    {service.charAt(0).toUpperCase() + service.slice(1)}
                  </TabsTrigger>
                ))}
              </TabsList>
              {Object.keys(SERVICE_TABLES).map((service) => (
                <TabsContent key={service} value={service}>
                  <ResultsTable
                    service={service}
                    data={serviceData[service] || []}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                  />
                </TabsContent>
              ))}
            </Tabs>
            <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
                <DialogTrigger asChild>
                  <Button onClick={() => { setEditingItem(null); form.reset(); }}>Add</Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>{editingItem ? 'Edit' : 'Add'} {currentService}</DialogTitle>
                    <DialogDescription>Add or edit a {currentService}</DialogDescription>
                  </DialogHeader>
                  <Form {...form}>
                    <form onSubmit={form.handleSubmit(editingItem ? handleUpdate : handleAdd)} className="space-y-4">
                      {SERVICE_TABLES[currentService].columns
                        .filter(column => column.key !== 'id')
                        .map(column => (
                          <FormField
                            key={column.key}
                            control={form.control}
                            name={column.key as any}
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>{column.label}</FormLabel>
                                <FormControl>
                                  <Input
                                    {...field}
                                    type={column.type === 'number' ? 'number' : 'text'}
                                    onChange={(e) => {
                                      const value = column.type === 'number' ? Number(e.target.value) : e.target.value;
                                      field.onChange(value);
                                    }}
                                  />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        ))}
                      <Button type="submit" disabled={isLoading}>
                        {isLoading ? 'Saving...' : 'Save'}
                      </Button>
                    </form>
                  </Form>
                </DialogContent>
              </Dialog>
          </>
        )}
        {error && (
          <Alert variant="destructive" className="mt-4">
            <AlertTitle>Error</AlertTitle>
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}
      </CardContent>
    </Card>
  )
}