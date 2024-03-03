// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

import { apiCookieName } from './src/constants';

export async function middleware(request: NextRequest): Promise<NextResponse> {
  if (!request.cookies.has(apiCookieName) && !request.nextUrl.pathname.startsWith('/accept_invitation')) {
    const destParam =
      request.nextUrl.searchParams.get('dest') ??
      encodeURIComponent(`${request.nextUrl.pathname}${request.nextUrl.search}`);
    return NextResponse.redirect(new URL(`/login?dest=${destParam}`, request.url));
  }

  return NextResponse.next();
}

export const config = {
  api: {
    bodyParser: false,
  },
  matcher: ['/(api/v1/.*)', '/(meal_plans/.*)', '/(meals/.*)', '/(settings/.*)'],
};
