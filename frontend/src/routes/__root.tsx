import QueryProvider from '@/providers/query-provider';
import { createRootRoute, Outlet } from '@tanstack/react-router';
import * as React from 'react';

export const Route = createRootRoute({
  component: RootComponent
});

function RootComponent() {
  return (
    <React.Fragment>
      <QueryProvider>
        <Outlet />
      </QueryProvider>
    </React.Fragment>
  );
}
