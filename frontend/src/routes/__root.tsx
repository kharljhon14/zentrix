import * as React from 'react';
import { Outlet, createRootRoute } from '@tanstack/react-router';
import AppSidebar from '@/components/sidebar/sidebar';
import AppSidebarProvider from '@/providers/sidebar-provider';
import Header from '@/components/header/header';

export const Route = createRootRoute({
  component: RootComponent
});

function RootComponent() {
  return (
    <React.Fragment>
      <Outlet />
    </React.Fragment>
  );
}
