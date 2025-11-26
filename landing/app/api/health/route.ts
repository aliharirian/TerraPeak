import { NextResponse } from 'next/server';

/**
 * Health check endpoint for Docker and load balancers
 * Returns 200 OK if the service is healthy
 */
export async function GET() {
  return NextResponse.json(
    {
      status: 'healthy',
      timestamp: new Date().toISOString(),
      service: 'terrapeak-landing',
    },
    { status: 200 }
  );
}

