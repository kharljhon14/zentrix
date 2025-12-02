import * as React from 'react';
import { Outlet, createRootRoute } from '@tanstack/react-router';
import AppSidebar from '@/components/sidebar/sidebar';
import AppSidebarProvider from '@/providers/sidebar-provider';

export const Route = createRootRoute({
  component: RootComponent
});

function RootComponent() {
  return (
    <React.Fragment>
      <AppSidebarProvider>
        <AppSidebar />
        <div>Hello "__root"!</div>
        <Outlet />
      </AppSidebarProvider>
    </React.Fragment>
  );
}
