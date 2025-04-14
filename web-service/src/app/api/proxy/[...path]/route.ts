import { NextRequest, NextResponse } from 'next/server'

const SERVICES = {
  auth: 'http://auth-service:8080',
  user: 'http://user-service:8082',
  product: 'http://product-service:8081',
  order: 'http://order-service:8083',
  payment: 'http://payment-service:8084',
  notification: 'http://notification-service:8085'
}

export async function GET(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  const resolvedParams = await params
  const path = resolvedParams.path
  const [service, ...rest] = path
  const url = `${SERVICES[service as keyof typeof SERVICES]}/${rest.join('/')}`

  try {
    const response = await fetch(url, {
      headers: {
        'Authorization': request.headers.get('Authorization') || '',
      },
    })

    const contentType = response.headers.get('content-type')
    if (contentType && contentType.includes('application/json')) {
      const data = await response.json()
      return NextResponse.json(data)
    } else {
      const text = await response.text()
      return new NextResponse(text, {
        status: response.status,
        headers: {
          'Content-Type': contentType || 'text/plain',
        },
      })
    }
  } catch (error) {
    console.error('Proxy error:', error)
    return NextResponse.json({ error: 'An error occurred' }, { status: 500 })
  }
}

export async function POST(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  const resolvedParams = await params
  const path = resolvedParams.path
  const [service, ...rest] = path
  const url = `${SERVICES[service as keyof typeof SERVICES]}/${rest.join('/')}`

  try {
    const body = await request.json()
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': request.headers.get('Authorization') || '',
      },
      body: JSON.stringify(body),
    })

    const contentType = response.headers.get('content-type')
    if (contentType && contentType.includes('application/json')) {
      const data = await response.json()
      return NextResponse.json(data)
    } else {
      const text = await response.text()
      return new NextResponse(text, {
        status: response.status,
        headers: {
          'Content-Type': contentType || 'text/plain',
        },
      })
    }
  } catch (error) {
    console.error('Proxy error:', error)
    return NextResponse.json({ error: 'An error occurred' }, { status: 500 })
  }
}

export async function PUT(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  const resolvedParams = await params
  const path = resolvedParams.path
  const [service, ...rest] = path
  const url = `${SERVICES[service as keyof typeof SERVICES]}/${rest.join('/')}`

  try {
    const body = await request.json()
    const response = await fetch(url, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': request.headers.get('Authorization') || '',
      },
      body: JSON.stringify(body),
    })

    const contentType = response.headers.get('content-type')
    if (contentType && contentType.includes('application/json')) {
      const data = await response.json()
      return NextResponse.json(data)
    } else {
      const text = await response.text()
      return new NextResponse(text, {
        status: response.status,
        headers: {
          'Content-Type': contentType || 'text/plain',
        },
      })
    }
  } catch (error) {
    console.error('Proxy error:', error)
    return NextResponse.json({ error: 'An error occurred' }, { status: 500 })
  }
}

export async function DELETE(request: NextRequest, { params }: { params: Promise<{ path: string[] }> }) {
  const resolvedParams = await params
  const path = resolvedParams.path
  const [service, ...rest] = path
  const url = `${SERVICES[service as keyof typeof SERVICES]}/${rest.join('/')}`

  try {
    const response = await fetch(url, {
      method: 'DELETE',
      headers: {
        'Authorization': request.headers.get('Authorization') || '',
      },
    })

    const contentType = response.headers.get('content-type')
    if (contentType && contentType.includes('application/json')) {
      const data = await response.json()
      return NextResponse.json(data)
    } else {
      const text = await response.text()
      return new NextResponse(text, {
        status: response.status,
        headers: {
          'Content-Type': contentType || 'text/plain',
        },
      })
    }
  } catch (error) {
    console.error('Proxy error:', error)
    return NextResponse.json({ error: 'An error occurred' }, { status: 500 })
  }
}